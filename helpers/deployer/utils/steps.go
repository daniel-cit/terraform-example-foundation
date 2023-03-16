// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"fmt"
	"path/filepath"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/git"
	"github.com/gruntwork-io/terratest/modules/terraform"

	"github.com/mitchellh/go-testing-interface"
)

type BootstrapOutputs struct {
	RemoteStateBucket         string
	RemoteStateBucketProjects string
	CICDProject               string
	DefaultRegion             string
}

func GetBootstrapStepOutputs(t testing.TB, options *terraform.Options) BootstrapOutputs {
	return BootstrapOutputs{
		CICDProject:               terraform.Output(t, options, "cloudbuild_project_id"),
		RemoteStateBucket:         terraform.Output(t, options, "gcs_bucket_tfstate"),
		RemoteStateBucketProjects: terraform.Output(t, options, "projects_gcs_bucket_tfstate"),
		DefaultRegion:             terraform.OutputMap(t, options, "common_config")["default_region"],
	}
}

func DeployBootstrapStep(t testing.TB, e State, tfvars GlobalTfvars, options *terraform.Options, checkoutPath, foundationPath string) error {
	repo := "gcp-bootstrap"
	step := "0-bootstrap"

	bootstrapTfvars := BootstrapTfvars{
		OrgID:              tfvars.OrgID,
		DefaultRegion:      tfvars.DefaultRegion,
		BillingAccount:     tfvars.BillingAccount,
		GroupOrgAdmins:     tfvars.GroupOrgAdmins,
		GroupBillingAdmins: tfvars.GroupBillingAdmins,
		ParentFolder:       tfvars.ParentFolder,
		ProjectPrefix:      tfvars.ProjectPrefix,
		FolderPrefix:       tfvars.FolderPrefix,
		BucketForceDestroy: tfvars.BucketForceDestroy,
	}

	err := WriteTfvars(filepath.Join(foundationPath, step, "terraform.tfvars"), bootstrapTfvars)
	if err != nil {
		return err
	}
	// terraform deploy
	initOutput, err := terraform.InitE(t, options)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", initOutput)
	planOutput, err := terraform.PlanE(t, options)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", planOutput)

	applyOutput, err := terraform.ApplyE(t, options)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", applyOutput)

	// read bootstrap outputs
	defaultRegion := terraform.OutputMap(t, options, "common_config")["default_region"]
	cbProjectID := terraform.Output(t, options, "cloudbuild_project_id")
	backendBucket := terraform.Output(t, options, "gcs_bucket_tfstate")
	backendBucketProjects := terraform.Output(t, options, "projects_gcs_bucket_tfstate")

	// replace backend and terraform init migrate
	fmt.Println("## migrate state ##")

	err = CopyFile(filepath.Join(options.TerraformDir, "backend.tf.example"), filepath.Join(options.TerraformDir, "backend.tf"))
	if err != nil {
		return err
	}
	err = ReplaceStringInFile(filepath.Join(options.TerraformDir, "backend.tf"), "UPDATE_ME", backendBucket)
	if err != nil {
		return err
	}

	options.MigrateState = true
	migrate, err := terraform.InitE(t, options)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", migrate)

	// replace all backend files
	files, err := FindFiles(foundationPath, "backend.tf")
	if err != nil {
		return err
	}
	for _, file := range files {
		err = ReplaceStringInFile(file, "UPDATE_ME", backendBucket)
		if err != nil {
			return err
		}
		err = ReplaceStringInFile(file, "UPDATE_PROJECTS_BACKEND", backendBucketProjects)
		if err != nil {
			return err
		}
	}

	fmt.Println("Follow the Cloud Build execution in the following link:")
	fmt.Printf("https://console.cloud.google.com/cloud-build/builds;region=%s?project=%s\n", defaultRegion, cbProjectID)

	// TODO press enter to continue

	// Check if image build was successful.
	// Need to abort if not.
	imageFilter := "source.repoSource.repoName:tf-cloudbuilder"
	imageBuild := GetRunningBuild(t, cbProjectID, defaultRegion, imageFilter)
	if imageBuild != "" {
		status := GetTerminalState(t, cbProjectID, defaultRegion, imageBuild)
		if status != "SUCCESS" {
			fmt.Println("Terraform Image builder Build Failed.")
			//TODO add error retrun
		}
	}

	//prepare policies repo
	gcpPoliciesPath := filepath.Join(checkoutPath, "gcp-policies")
	policiesConf := CloneRepo(t, "gcp-policies", gcpPoliciesPath, cbProjectID)
	policiesBranch := "main"
	err = CheckoutBranch(policiesConf, policiesBranch)
	if err != nil {
		return err
	}
	err = CopyDirectory(filepath.Join(foundationPath, "policy-library"), gcpPoliciesPath)
	if err != nil {
		return err
	}
	err = CommitFiles(policiesConf, "Initialize policy library repo")
	if err != nil {
		return err
	}
	err = PushBranch(policiesConf, policiesBranch)
	if err != nil {
		return err
	}

	//prepare bootstrap repo
	gcpBootstrapPath := filepath.Join(checkoutPath, "gcp-bootstrap")
	bootstrapConf := CloneRepo(t, "gcp-bootstrap", gcpBootstrapPath, cbProjectID)
	err = CheckoutBranch(bootstrapConf, "plan")
	if err != nil {
		return err
	}

	fmt.Println("copy bootstrap files")
	err = copyStepCode(t, bootstrapConf, foundationPath, checkoutPath, repo, step, "envs/shared")
	if err != nil {
		return err
	}

	err = planStep(t, bootstrapConf, cbProjectID, defaultRegion, repo)
	if err != nil {
		return err
	}

	err = applyEnv(t, bootstrapConf, cbProjectID, defaultRegion, repo, "production")
	if err != nil {
		return err
	}

	fmt.Println("end of bootstrap deploy")

	return nil
}

func DeployOrgStep(t testing.TB, e State, tfvars GlobalTfvars, checkoutPath, foundationPath string, outputs BootstrapOutputs) error {
	repo := "gcp-org"
	step := "1-org"

	orgTfvars := OrgTfvars{
		DomainsToAllow:           tfvars.DomainsToAllow,
		EssentialContactsDomains: tfvars.EssentialContactsDomains,
		BillingDataUsers:         tfvars.BillingDataUsers,
		AuditDataUsers:           tfvars.AuditDataUsers,
		SccNotificationName:      tfvars.SccNotificationName,
		RemoteStateBucket:        outputs.RemoteStateBucket,
		EnableHubAndSpoke:        tfvars.EnableHubAndSpoke,
		CreateACMAPolicy:         false,
		CreateUniqueTagKey:       tfvars.CreateUniqueTagKey,
	}

	err := WriteTfvars(filepath.Join(foundationPath, step, "envs", "shared", "terraform.tfvars"), orgTfvars)
	if err != nil {
		return err
	}

	gcpPath := filepath.Join(checkoutPath, repo)
	conf := CloneRepo(t, repo, gcpPath, outputs.CICDProject)
	err = CheckoutBranch(conf, "plan")
	if err != nil {
		return err
	}

	err = copyStepCode(t, conf, foundationPath, checkoutPath, repo, step, "")
	if err != nil {
		return err
	}

	err = planStep(t, conf, outputs.CICDProject, outputs.DefaultRegion, repo)
	if err != nil {
		return err
	}

	err = applyEnv(t, conf, outputs.CICDProject, outputs.DefaultRegion, repo, "production")
	if err != nil {
		return err
	}

	fmt.Println("end of", step, "deploy")
	return nil
}

func DeployEnvStep(t testing.TB, e State, tfvars GlobalTfvars, checkoutPath, foundationPath string, outputs BootstrapOutputs) error {
	repo := "gcp-environments"
	step := "2-environments"

	envsTfvars := EnvsTfvars{
		MonitoringWorkspaceUsers: tfvars.MonitoringWorkspaceUsers,
		RemoteStateBucket:        outputs.RemoteStateBucket,
	}

	err := WriteTfvars(filepath.Join(foundationPath, step, "terraform.tfvars"), envsTfvars)
	if err != nil {
		return err
	}

	gcpPath := filepath.Join(checkoutPath, repo)
	conf := CloneRepo(t, repo, gcpPath, outputs.CICDProject)
	err = CheckoutBranch(conf, "plan")
	if err != nil {
		return err
	}

	err = RunStepE(e, "gcp-environments.copy-code", func() error {
		return copyStepCode(t, conf, foundationPath, checkoutPath, repo, step, "")
	})
	if err != nil {
		return err
	}

	err = RunStepE(e, "gcp-environments.plan", func() error {
		return planStep(t, conf, outputs.CICDProject, outputs.DefaultRegion, repo)
	})
	if err != nil {
		return err
	}

	err = RunStepE(e, "gcp-environments.production", func() error {
		return applyEnv(t, conf, outputs.CICDProject, outputs.DefaultRegion, repo, "production")
	})
	if err != nil {
		return err
	}

	err = RunStepE(e, "gcp-environments.non-production", func() error {
		return applyEnv(t, conf, outputs.CICDProject, outputs.DefaultRegion, repo, "non-production")
	})
	if err != nil {
		return err
	}

	err = RunStepE(e, "gcp-environments.development", func() error {
		return applyEnv(t, conf, outputs.CICDProject, outputs.DefaultRegion, repo, "development")
	})
	if err != nil {
		return err
	}

	fmt.Println("end of", step, "deploy")
	return nil
}

func copyStepCode(t testing.TB, conf *git.CmdCfg, foundationPath, checkoutPath, repo, step, customPath string) error {
	gcpPath := filepath.Join(checkoutPath, repo)
	targetDir := gcpPath
	if customPath != "" {
		targetDir = filepath.Join(gcpPath, customPath)
	}
	err := CopyDirectory(filepath.Join(foundationPath, step), targetDir)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = CopyFile(filepath.Join(foundationPath, "build/cloudbuild-tf-apply.yaml"), filepath.Join(gcpPath, "cloudbuild-tf-apply.yaml"))
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = CopyFile(filepath.Join(foundationPath, "build/cloudbuild-tf-plan.yaml"), filepath.Join(gcpPath, "cloudbuild-tf-plan.yaml"))
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = CopyFile(filepath.Join(foundationPath, "build/tf-wrapper.sh"), filepath.Join(gcpPath, "tf-wrapper.sh"))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func planStep(t testing.TB, conf *git.CmdCfg, project, region, repo string) error {
	err := CommitFiles(conf, fmt.Sprintf("Initialize %s repo", repo))
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = PushBranch(conf, "plan")
	if err != nil {
		fmt.Println(err)
		return err
	}

	failed := waitBuild(t, conf, project, region, repo, fmt.Sprintf("Terraform %s plan build Failed.", repo))
	if failed {
		return fmt.Errorf("plan build for repository %s failed", repo)
	}
	return nil
}

func applyEnv(t testing.TB, conf *git.CmdCfg, project, region, repo, environment string) error {
	err := CheckoutBranch(conf, environment)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = PushBranch(conf, environment)
	if err != nil {
		fmt.Println(err)
		return err
	}

	failed := waitBuild(t, conf, project, region, repo, fmt.Sprintf("Terraform %s apply %s build Failed.", repo, environment))
	if failed {
		return fmt.Errorf("build for environment %s in repository %s failed", environment, repo)
	}
	return nil
}

func waitBuild(t testing.TB, conf *git.CmdCfg, project, region, repo, failureMsg string) bool {
	filter := fmt.Sprintf("source.repoSource.repoName:%s", repo)
	build := GetRunningBuild(t, project, region, filter)
	if build != "" {
		status := GetTerminalState(t, project, region, build)
		if status != "SUCCESS" {
			fmt.Println(failureMsg)
			return true
		}
	}
	return false
}
