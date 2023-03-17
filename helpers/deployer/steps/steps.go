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

package steps

import (
	"fmt"
	"path/filepath"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/git"
	"github.com/gruntwork-io/terratest/modules/terraform"

	"github.com/terraform-google-modules/terraform-example-foundation/helpers/deployer/utils"

	"github.com/mitchellh/go-testing-interface"
)

type BootstrapOutputs struct {
	RemoteStateBucket         string
	RemoteStateBucketProjects string
	CICDProject               string
	DefaultRegion             string
}

type ServerAddress struct {
	Ipv4Address    string `cty:"ipv4_address"`
	ForwardingPath string `cty:"forwarding_path"`
}

type GlobalTfvars struct {
	OrgID                     string           `hcl:"org_id"`
	BillingAccount            string           `hcl:"billing_account"`
	GroupOrgAdmins            string           `hcl:"group_org_admins"`
	GroupBillingAdmins        string           `hcl:"group_billing_admins"`
	BillingDataUsers          string           `hcl:"billing_data_users"`
	MonitoringWorkspaceUsers  string           `hcl:"monitoring_workspace_users"`
	AuditDataUsers            string           `hcl:"audit_data_users"`
	DefaultRegion             string           `hcl:"default_region"`
	ParentFolder              *string          `hcl:"parent_folder"`
	DomainsToAllow            []string         `hcl:"domains_to_allow"`
	EssentialContactsDomains  []string         `hcl:"essential_contacts_domains_to_allow"`
	TargetNameServerAddresses *[]ServerAddress `hcl:"target_name_server_addresses"`
	SccNotificationName       string           `hcl:"scc_notification_name"`
	ProjectPrefix             *string          `hcl:"project_prefix"`
	FolderPrefix              *string          `hcl:"folder_prefix"`
	BucketForceDestroy        *bool            `hcl:"bucket_force_destroy"`
	EnableHubAndSpoke         bool             `hcl:"enable_hub_and_spoke"`
	CreateUniqueTagKey        bool             `hcl:"create_unique_tag_key"`
	CodeCheckoutPath          string           `hcl:"code_checkout_path"`
	FoundationCodePath        string           `hcl:"foundation_code_path"`
}

type BootstrapTfvars struct {
	OrgID              string  `hcl:"org_id"`
	BillingAccount     string  `hcl:"billing_account"`
	GroupOrgAdmins     string  `hcl:"group_org_admins"`
	GroupBillingAdmins string  `hcl:"group_billing_admins"`
	DefaultRegion      string  `hcl:"default_region"`
	ParentFolder       *string `hcl:"parent_folder"`
	ProjectPrefix      *string `hcl:"project_prefix"`
	FolderPrefix       *string `hcl:"folder_prefix"`
	BucketForceDestroy *bool   `hcl:"bucket_force_destroy"`
}

type OrgTfvars struct {
	DomainsToAllow           []string `hcl:"domains_to_allow"`
	EssentialContactsDomains []string `hcl:"essential_contacts_domains_to_allow"`
	BillingDataUsers         string   `hcl:"billing_data_users"`
	AuditDataUsers           string   `hcl:"audit_data_users"`
	SccNotificationName      string   `hcl:"scc_notification_name"`
	RemoteStateBucket        string   `hcl:"remote_state_bucket"`
	EnableHubAndSpoke        bool     `hcl:"enable_hub_and_spoke"`
	CreateACMAPolicy         bool     `hcl:"create_access_context_manager_access_policy"`
	CreateUniqueTagKey       bool     `hcl:"create_unique_tag_key"`
}

type EnvsTfvars struct {
	MonitoringWorkspaceUsers string `hcl:"monitoring_workspace_users"`
	RemoteStateBucket        string `hcl:"remote_state_bucket"`
}

func GetBootstrapStepOutputs(t testing.TB, options *terraform.Options) BootstrapOutputs {
	return BootstrapOutputs{
		CICDProject:               terraform.Output(t, options, "cloudbuild_project_id"),
		RemoteStateBucket:         terraform.Output(t, options, "gcs_bucket_tfstate"),
		RemoteStateBucketProjects: terraform.Output(t, options, "projects_gcs_bucket_tfstate"),
		DefaultRegion:             terraform.OutputMap(t, options, "common_config")["default_region"],
	}
}

func DeployBootstrapStep(t testing.TB, e utils.State, tfvars GlobalTfvars, options *terraform.Options, checkoutPath, foundationPath string) error {
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

	err := utils.WriteTfvars(filepath.Join(foundationPath, step, "terraform.tfvars"), bootstrapTfvars)
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
	err = utils.RunStepE(e, "gcp-bootstrap.migrate-state", func() error {
		options.MigrateState = true
		err = utils.CopyFile(filepath.Join(options.TerraformDir, "backend.tf.example"), filepath.Join(options.TerraformDir, "backend.tf"))
		if err != nil {
			return err
		}
		err = utils.ReplaceStringInFile(filepath.Join(options.TerraformDir, "backend.tf"), "UPDATE_ME", backendBucket)
		if err != nil {
			return err
		}
		migrate, err := terraform.InitE(t, options)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", migrate)
		return nil
	})
	if err != nil {
		return err
	}

	// replace all backend files
	err = utils.RunStepE(e, "gcp-bootstrap.replace-backend-files", func() error {
		files, err := utils.FindFiles(foundationPath, "backend.tf")
		if err != nil {
			return err
		}
		for _, file := range files {
			err = utils.ReplaceStringInFile(file, "UPDATE_ME", backendBucket)
			if err != nil {
				return err
			}
			err = utils.ReplaceStringInFile(file, "UPDATE_PROJECTS_BACKEND", backendBucketProjects)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	fmt.Println("Follow the Cloud Build execution in the following link:")
	fmt.Printf("https://console.cloud.google.com/cloud-build/builds;region=%s?project=%s\n", defaultRegion, cbProjectID)

	// TODO press enter to continue

	// Check if image build was successful.
	err = waitBuildSuccess(t, cbProjectID, defaultRegion, "tf-cloudbuilder", "Terraform Image builder Build Failed for tf-cloudbuilder repository.")
	if err != nil {
		return err
	}

	//prepare policies repo
	gcpPoliciesPath := filepath.Join(checkoutPath, "gcp-policies")
	policiesConf := utils.CloneRepo(t, "gcp-policies", gcpPoliciesPath, cbProjectID)
	policiesBranch := "main"

	err = utils.RunStepE(e, "gcp-bootstrap.gcp-policies", func() error {
		err = utils.CheckoutBranch(policiesConf, policiesBranch)
		if err != nil {
			return err
		}
		err = utils.CopyDirectory(filepath.Join(foundationPath, "policy-library"), gcpPoliciesPath)
		if err != nil {
			return err
		}
		err = utils.CommitFiles(policiesConf, "Initialize policy library repo")
		if err != nil {
			return err
		}
		err = utils.PushBranch(policiesConf, policiesBranch)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	//prepare bootstrap repo
	gcpBootstrapPath := filepath.Join(checkoutPath, "gcp-bootstrap")
	bootstrapConf := utils.CloneRepo(t, "gcp-bootstrap", gcpBootstrapPath, cbProjectID)
	err = utils.CheckoutBranch(bootstrapConf, "plan")
	if err != nil {
		return err
	}

	err = utils.RunStepE(e, "gcp-bootstrap.copy-code", func() error {
		return copyStepCode(t, bootstrapConf, foundationPath, checkoutPath, repo, step, "envs/shared")
	})
	if err != nil {
		return err
	}

	err = utils.RunStepE(e, "gcp-bootstrap.plan", func() error {
		return planStep(t, bootstrapConf, cbProjectID, defaultRegion, repo)
	})
	if err != nil {
		return err
	}

	err = utils.RunStepE(e, "gcp-bootstrap.production", func() error {
		return applyEnv(t, bootstrapConf, cbProjectID, defaultRegion, repo, "production")
	})
	if err != nil {
		return err
	}

	fmt.Println("end of bootstrap deploy")

	return nil
}

func DeployOrgStep(t testing.TB, e utils.State, tfvars GlobalTfvars, checkoutPath, foundationPath string, outputs BootstrapOutputs) error {
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

	err := utils.WriteTfvars(filepath.Join(foundationPath, step, "envs", "shared", "terraform.tfvars"), orgTfvars)
	if err != nil {
		return err
	}

	gcpPath := filepath.Join(checkoutPath, repo)
	conf := utils.CloneRepo(t, repo, gcpPath, outputs.CICDProject)
	err = utils.CheckoutBranch(conf, "plan")
	if err != nil {
		return err
	}

	err = utils.RunStepE(e, "gcp-org.copy-code", func() error {
		return copyStepCode(t, conf, foundationPath, checkoutPath, repo, step, "")
	})
	if err != nil {
		return err
	}

	err = utils.RunStepE(e, "gcp-org.plan", func() error {
		return planStep(t, conf, outputs.CICDProject, outputs.DefaultRegion, repo)
	})
	if err != nil {
		return err
	}

	err = utils.RunStepE(e, "gcp-org.production", func() error {
		return applyEnv(t, conf, outputs.CICDProject, outputs.DefaultRegion, repo, "production")
	})
	if err != nil {
		return err
	}

	fmt.Println("end of", step, "deploy")
	return nil
}

func DeployEnvStep(t testing.TB, e utils.State, tfvars GlobalTfvars, checkoutPath, foundationPath string, outputs BootstrapOutputs) error {
	repo := "gcp-environments"
	step := "2-environments"

	envsTfvars := EnvsTfvars{
		MonitoringWorkspaceUsers: tfvars.MonitoringWorkspaceUsers,
		RemoteStateBucket:        outputs.RemoteStateBucket,
	}

	err := utils.WriteTfvars(filepath.Join(foundationPath, step, "terraform.tfvars"), envsTfvars)
	if err != nil {
		return err
	}

	gcpPath := filepath.Join(checkoutPath, repo)
	conf := utils.CloneRepo(t, repo, gcpPath, outputs.CICDProject)
	err = utils.CheckoutBranch(conf, "plan")
	if err != nil {
		return err
	}

	err = utils.RunStepE(e, "gcp-environments.copy-code", func() error {
		return copyStepCode(t, conf, foundationPath, checkoutPath, repo, step, "")
	})
	if err != nil {
		return err
	}

	err = utils.RunStepE(e, "gcp-environments.plan", func() error {
		return planStep(t, conf, outputs.CICDProject, outputs.DefaultRegion, repo)
	})
	if err != nil {
		return err
	}

	err = utils.RunStepE(e, "gcp-environments.production", func() error {
		return applyEnv(t, conf, outputs.CICDProject, outputs.DefaultRegion, repo, "production")
	})
	if err != nil {
		return err
	}

	err = utils.RunStepE(e, "gcp-environments.non-production", func() error {
		return applyEnv(t, conf, outputs.CICDProject, outputs.DefaultRegion, repo, "non-production")
	})
	if err != nil {
		return err
	}

	err = utils.RunStepE(e, "gcp-environments.development", func() error {
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
	err := utils.CopyDirectory(filepath.Join(foundationPath, step), targetDir)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = utils.CopyFile(filepath.Join(foundationPath, "build/cloudbuild-tf-apply.yaml"), filepath.Join(gcpPath, "cloudbuild-tf-apply.yaml"))
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = utils.CopyFile(filepath.Join(foundationPath, "build/cloudbuild-tf-plan.yaml"), filepath.Join(gcpPath, "cloudbuild-tf-plan.yaml"))
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = utils.CopyFile(filepath.Join(foundationPath, "build/tf-wrapper.sh"), filepath.Join(gcpPath, "tf-wrapper.sh"))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func planStep(t testing.TB, conf *git.CmdCfg, project, region, repo string) error {
	err := utils.CommitFiles(conf, fmt.Sprintf("Initialize %s repo", repo))
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = utils.PushBranch(conf, "plan")
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = waitBuildSuccess(t, project, region, repo, fmt.Sprintf("Terraform %s plan build Failed.", repo))
	if err != nil {
		return err
	}
	return nil
}

func applyEnv(t testing.TB, conf *git.CmdCfg, project, region, repo, environment string) error {
	err := utils.CheckoutBranch(conf, environment)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = utils.PushBranch(conf, environment)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = waitBuildSuccess(t, project, region, repo, fmt.Sprintf("Terraform %s apply %s build Failed.", repo, environment))
	if err != nil {
		return err
	}
	return nil
}

func waitBuildSuccess(t testing.TB, project, region, repo, failureMsg string) error {
	filter := fmt.Sprintf("source.repoSource.repoName:%s", repo)
	build := utils.GetRunningBuild(t, project, region, filter)
	if build != "" {
		// TODO add message to fully identify build
		status := utils.GetTerminalState(t, project, region, build)
		if status != "SUCCESS" {
			return fmt.Errorf("%s\nSee:\nhttps://console.cloud.google.com/cloud-build/builds;region=%s/%s?project=%s\nfor details.\n", failureMsg, region, build, project)
		}
	} else {
		status := utils.GetLastBuildStatus(t, project, region, filter)
		if status != "SUCCESS" {
			return fmt.Errorf("%s\nSee:\nhttps://console.cloud.google.com/cloud-build/builds;region=%s/%s?project=%s\nfor details.\n", failureMsg, region, build, project)
		}
	}
	return nil
}
