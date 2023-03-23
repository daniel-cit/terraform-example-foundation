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
	"os"
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"

	"github.com/terraform-google-modules/terraform-example-foundation/helpers/deployer/gcp"
	"github.com/terraform-google-modules/terraform-example-foundation/helpers/deployer/state"
	"github.com/terraform-google-modules/terraform-example-foundation/helpers/deployer/utils"

	"github.com/mitchellh/go-testing-interface"
)

type BootstrapOutputs struct {
	RemoteStateBucket         string
	RemoteStateBucketProjects string
	CICDProject               string
	DefaultRegion             string
	NetworkSA                 string
	ProjectsSA                string
}

type InfraPipelineOutputs struct {
	RemoteStateBucket string
	InfraPipeProj     string
	DefaultRegion     string
}

type ServerAddress struct {
	Ipv4Address    string `cty:"ipv4_address"`
	ForwardingPath string `cty:"forwarding_path"`
}

type GlobalTfvars struct {
	OrgID                         string          `hcl:"org_id"`
	BillingAccount                string          `hcl:"billing_account"`
	GroupOrgAdmins                string          `hcl:"group_org_admins"`
	GroupBillingAdmins            string          `hcl:"group_billing_admins"`
	BillingDataUsers              string          `hcl:"billing_data_users"`
	MonitoringWorkspaceUsers      string          `hcl:"monitoring_workspace_users"`
	AuditDataUsers                string          `hcl:"audit_data_users"`
	DefaultRegion                 string          `hcl:"default_region"`
	ParentFolder                  *string         `hcl:"parent_folder"`
	Domain                        string          `hcl:"domain"`
	DomainsToAllow                []string        `hcl:"domains_to_allow"`
	EssentialContactsDomains      []string        `hcl:"essential_contacts_domains_to_allow"`
	PerimeterAdditionalMembers    []string        `hcl:"perimeter_additional_members"`
	TargetNameServerAddresses     []ServerAddress `hcl:"target_name_server_addresses"`
	SccNotificationName           string          `hcl:"scc_notification_name"`
	ProjectPrefix                 *string         `hcl:"project_prefix"`
	FolderPrefix                  *string         `hcl:"folder_prefix"`
	BucketForceDestroy            *bool           `hcl:"bucket_force_destroy"`
	EnableHubAndSpoke             bool            `hcl:"enable_hub_and_spoke"`
	EnableHubAndSpokeTransitivity bool            `hcl:"enable_hub_and_spoke_transitivity"`
	CreateUniqueTagKey            bool            `hcl:"create_unique_tag_key"`
	ProjectsKMSLocation           string          `hcl:"projects_kms_location"`
	ProjectsGCSLocation           string          `hcl:"projects_gcs_location"`
	CodeCheckoutPath              string          `hcl:"code_checkout_path"`
	FoundationCodePath            string          `hcl:"foundation_code_path"`
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

type NetCommonTfvars struct {
	Domain                        string   `hcl:"domain"`
	PerimeterAdditionalMembers    []string `hcl:"perimeter_additional_members"`
	RemoteStateBucket             string   `hcl:"remote_state_bucket"`
	EnableHubAndSpokeTransitivity *bool    `hcl:"enable_hub_and_spoke_transitivity"`
}

type NetSharedTfvars struct {
	TargetNameServerAddresses []ServerAddress `hcl:"target_name_server_addresses"`
}

type NetAccessContextTfvars struct {
	AccessContextManagerPolicyID string `hcl:"access_context_manager_policy_id"`
}

type ProjCommonTfvars struct {
	RemoteStateBucket string `hcl:"remote_state_bucket"`
}

type ProjSharedTfvars struct {
	DefaultRegion string `hcl:"default_region"`
}

type ProjEnvTfvars struct {
	ProjectsKMSLocation string `hcl:"projects_kms_location"`
	ProjectsGCSLocation string `hcl:"projects_gcs_location"`
}

type AppInfraCommonTfvars struct {
	InstanceRegion    string `hcl:"instance_region"`
	RemoteStateBucket string `hcl:"remote_state_bucket"`
}

func GetBootstrapStepOutputs(t testing.TB, options *terraform.Options) BootstrapOutputs {
	return BootstrapOutputs{
		CICDProject:               terraform.Output(t, options, "cloudbuild_project_id"),
		RemoteStateBucket:         terraform.Output(t, options, "gcs_bucket_tfstate"),
		RemoteStateBucketProjects: terraform.Output(t, options, "projects_gcs_bucket_tfstate"),
		DefaultRegion:             terraform.OutputMap(t, options, "common_config")["default_region"],
		NetworkSA:                 terraform.Output(t, options, "networks_step_terraform_service_account_email"),
		ProjectsSA:                terraform.Output(t, options, "projects_step_terraform_service_account_email"),
	}
}

func GetInfraPipelineOutputs(t testing.TB, options *terraform.Options, workspace string) InfraPipelineOutputs {
	return InfraPipelineOutputs{
		InfraPipeProj: terraform.Output(t, options, "cloudbuild_project_id"),
		DefaultRegion: terraform.Output(t, options, "default_region"),
	}
}

func DeployBootstrapStep(t testing.TB, s state.State, tfvars GlobalTfvars, options *terraform.Options, checkoutPath, foundationPath string, logger *logger.Logger) error {
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
	err = state.RunStepE(s, "gcp-bootstrap.migrate-state", func() error {
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
	err = state.RunStepE(s, "gcp-bootstrap.replace-backend-files", func() error {
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
	err = gcp.WaitBuildSuccess(t, cbProjectID, defaultRegion, "tf-cloudbuilder", "Terraform Image builder Build Failed for tf-cloudbuilder repository.")
	if err != nil {
		return err
	}

	//prepare policies repo
	gcpPoliciesPath := filepath.Join(checkoutPath, "gcp-policies")
	policiesConf := utils.CloneCSR(t, "gcp-policies", gcpPoliciesPath, cbProjectID, logger)
	policiesBranch := "main"

	err = state.RunStepE(s, "gcp-bootstrap.gcp-policies", func() error {
		err = policiesConf.CheckoutBranch(policiesBranch)
		if err != nil {
			return err
		}
		err = utils.CopyDirectory(filepath.Join(foundationPath, "policy-library"), gcpPoliciesPath)
		if err != nil {
			return err
		}
		err = policiesConf.CommitFiles("Initialize policy library repo")
		if err != nil {
			return err
		}
		err = policiesConf.PushBranch(policiesBranch, "origin")
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
	bootstrapConf := utils.CloneCSR(t, "gcp-bootstrap", gcpBootstrapPath, cbProjectID, logger)
	err = bootstrapConf.CheckoutBranch("plan")
	if err != nil {
		return err
	}

	err = state.RunStepE(s, "gcp-bootstrap.copy-code", func() error {
		return copyStepCode(t, bootstrapConf, foundationPath, checkoutPath, repo, step, "envs/shared")
	})
	if err != nil {
		return err
	}

	err = state.RunStepE(s, "gcp-bootstrap.plan", func() error {
		return planStep(t, bootstrapConf, cbProjectID, defaultRegion, repo)
	})
	if err != nil {
		return err
	}

	err = state.RunStepE(s, "gcp-bootstrap.production", func() error {
		return applyEnv(t, bootstrapConf, cbProjectID, defaultRegion, repo, "production")
	})
	if err != nil {
		return err
	}

	fmt.Println("end of bootstrap deploy")

	return nil
}

func DeployOrgStep(t testing.TB, s state.State, tfvars GlobalTfvars, checkoutPath, foundationPath string, outputs BootstrapOutputs, logger *logger.Logger) error {
	repo := "gcp-org"
	step := "1-org"
	createACMAPolicy := gcp.GetAccessContextManagerPolicyID(t, tfvars.OrgID) == ""

	orgTfvars := OrgTfvars{
		DomainsToAllow:           tfvars.DomainsToAllow,
		EssentialContactsDomains: tfvars.EssentialContactsDomains,
		BillingDataUsers:         tfvars.BillingDataUsers,
		AuditDataUsers:           tfvars.AuditDataUsers,
		SccNotificationName:      tfvars.SccNotificationName,
		RemoteStateBucket:        outputs.RemoteStateBucket,
		EnableHubAndSpoke:        tfvars.EnableHubAndSpoke,
		CreateACMAPolicy:         createACMAPolicy,
		CreateUniqueTagKey:       tfvars.CreateUniqueTagKey,
	}

	err := utils.WriteTfvars(filepath.Join(foundationPath, step, "envs", "shared", "terraform.tfvars"), orgTfvars)
	if err != nil {
		return err
	}

	gcpPath := filepath.Join(checkoutPath, repo)
	conf := utils.CloneCSR(t, repo, gcpPath, outputs.CICDProject, logger)
	err = conf.CheckoutBranch("plan")
	if err != nil {
		return err
	}

	err = state.RunStepE(s, "gcp-org.copy-code", func() error {
		return copyStepCode(t, conf, foundationPath, checkoutPath, repo, step, "")
	})
	if err != nil {
		return err
	}

	err = state.RunStepE(s, "gcp-org.plan", func() error {
		return planStep(t, conf, outputs.CICDProject, outputs.DefaultRegion, repo)
	})
	if err != nil {
		return err
	}

	err = state.RunStepE(s, "gcp-org.production", func() error {
		return applyEnv(t, conf, outputs.CICDProject, outputs.DefaultRegion, repo, "production")
	})
	if err != nil {
		return err
	}

	fmt.Println("end of", step, "deploy")
	return nil
}

func DeployEnvStep(t testing.TB, s state.State, tfvars GlobalTfvars, checkoutPath, foundationPath string, outputs BootstrapOutputs, logger *logger.Logger) error {
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
	conf := utils.CloneCSR(t, repo, gcpPath, outputs.CICDProject, logger)
	err = conf.CheckoutBranch("plan")
	if err != nil {
		return err
	}

	err = state.RunStepE(s, "gcp-environments.copy-code", func() error {
		return copyStepCode(t, conf, foundationPath, checkoutPath, repo, step, "")
	})
	if err != nil {
		return err
	}

	err = state.RunStepE(s, "gcp-environments.plan", func() error {
		return planStep(t, conf, outputs.CICDProject, outputs.DefaultRegion, repo)
	})
	if err != nil {
		return err
	}

	err = state.RunStepE(s, "gcp-environments.production", func() error {
		return applyEnv(t, conf, outputs.CICDProject, outputs.DefaultRegion, repo, "production")
	})
	if err != nil {
		return err
	}

	err = state.RunStepE(s, "gcp-environments.non-production", func() error {
		return applyEnv(t, conf, outputs.CICDProject, outputs.DefaultRegion, repo, "non-production")
	})
	if err != nil {
		return err
	}

	err = state.RunStepE(s, "gcp-environments.development", func() error {
		return applyEnv(t, conf, outputs.CICDProject, outputs.DefaultRegion, repo, "development")
	})
	if err != nil {
		return err
	}

	fmt.Println("end of", step, "deploy")
	return nil
}

func DeployNetworksStep(t testing.TB, s state.State, tfvars GlobalTfvars, checkoutPath, foundationPath string, outputs BootstrapOutputs, logger *logger.Logger) error {
	repo := "gcp-networks"
	var step string
	if tfvars.EnableHubAndSpoke {
		step = "3-networks-hub-and-spoke"
	} else {
		step = "3-networks-dual-svpc"
	}

	// shared
	sharedTfvars := NetSharedTfvars{
		TargetNameServerAddresses: tfvars.TargetNameServerAddresses,
	}
	err := utils.WriteTfvars(filepath.Join(foundationPath, step, "shared.auto.tfvars"), sharedTfvars)
	if err != nil {
		return err
	}
	// common
	commonTfvars := NetCommonTfvars{
		Domain:                     tfvars.Domain,
		PerimeterAdditionalMembers: tfvars.PerimeterAdditionalMembers,
		RemoteStateBucket:          outputs.RemoteStateBucket,
	}
	if tfvars.EnableHubAndSpoke {
		commonTfvars.EnableHubAndSpokeTransitivity = &tfvars.EnableHubAndSpokeTransitivity
	}
	err = utils.WriteTfvars(filepath.Join(foundationPath, step, "common.auto.tfvars"), commonTfvars)
	if err != nil {
		return err
	}
	//access_context
	accessContextTfvars := NetAccessContextTfvars{
		AccessContextManagerPolicyID: gcp.GetAccessContextManagerPolicyID(t, tfvars.OrgID),
	}
	err = utils.WriteTfvars(filepath.Join(foundationPath, step, "access_context.auto.tfvars"), accessContextTfvars)
	if err != nil {
		return err
	}

	gcpPath := filepath.Join(checkoutPath, repo)
	conf := utils.CloneCSR(t, repo, gcpPath, outputs.CICDProject, logger)
	err = conf.CheckoutBranch("plan")
	if err != nil {
		return err
	}

	err = state.RunStepE(s, "gcp-networks.copy-code", func() error {
		return copyStepCode(t, conf, foundationPath, checkoutPath, repo, step, "")
	})
	if err != nil {
		return err
	}

	// Apply shared
	options := &terraform.Options{
		TerraformDir: filepath.Join(gcpPath, "envs", "shared"),
		Logger:       logger,
		NoColor:      true,
	}

	err = state.RunStepE(s, "gcp-networks.apply-shared", func() error {
		return applyShared(t, options, outputs.NetworkSA)
	})
	if err != nil {
		return err
	}

	err = state.RunStepE(s, "gcp-networks.plan", func() error {
		return planStep(t, conf, outputs.CICDProject, outputs.DefaultRegion, repo)
	})
	if err != nil {
		return err
	}

	err = state.RunStepE(s, "gcp-networks.production", func() error {
		return applyEnv(t, conf, outputs.CICDProject, outputs.DefaultRegion, repo, "production")
	})
	if err != nil {
		return err
	}

	err = state.RunStepE(s, "gcp-networks.non-production", func() error {
		return applyEnv(t, conf, outputs.CICDProject, outputs.DefaultRegion, repo, "non-production")
	})
	if err != nil {
		return err
	}

	err = state.RunStepE(s, "gcp-networks.development", func() error {
		return applyEnv(t, conf, outputs.CICDProject, outputs.DefaultRegion, repo, "development")
	})
	if err != nil {
		return err
	}

	fmt.Println("end of", step, "deploy")
	return nil
}

func DeployProjectsStep(t testing.TB, s state.State, tfvars GlobalTfvars, checkoutPath, foundationPath string, outputs BootstrapOutputs, logger *logger.Logger) error {
	repo := "gcp-projects"
	step := "4-projects"

	// shared
	sharedTfvars := ProjSharedTfvars{
		DefaultRegion: tfvars.DefaultRegion,
	}
	err := utils.WriteTfvars(filepath.Join(foundationPath, step, "shared.auto.tfvars"), sharedTfvars)
	if err != nil {
		return err
	}
	// common
	commonTfvars := ProjCommonTfvars{
		RemoteStateBucket: outputs.RemoteStateBucket,
	}
	err = utils.WriteTfvars(filepath.Join(foundationPath, step, "common.auto.tfvars"), commonTfvars)
	if err != nil {
		return err
	}
	//for each environment
	envTfvars := ProjEnvTfvars{
		ProjectsKMSLocation: tfvars.ProjectsKMSLocation,
		ProjectsGCSLocation: tfvars.ProjectsGCSLocation,
	}
	for _, envfile := range []string{
		"development.auto.tfvars",
		"non-production.auto.tfvars",
		"production.auto.tfvars"} {
		err = utils.WriteTfvars(filepath.Join(foundationPath, step, envfile), envTfvars)
		if err != nil {
			return err
		}
	}

	gcpPath := filepath.Join(checkoutPath, repo)
	conf := utils.CloneCSR(t, repo, gcpPath, outputs.CICDProject, logger)
	err = conf.CheckoutBranch("plan")
	if err != nil {
		return err
	}

	err = state.RunStepE(s, "gcp-projects.copy-code", func() error {
		return copyStepCode(t, conf, foundationPath, checkoutPath, repo, step, "")
	})
	if err != nil {
		return err
	}

	// Apply shared
	optbu1 := &terraform.Options{
		TerraformDir: filepath.Join(gcpPath, "business_unit_1", "shared"),
		Logger:       logger,
		NoColor:      true,
	}

	err = state.RunStepE(s, "gcp-projects.bu1.apply-shared", func() error {
		return applyShared(t, optbu1, outputs.ProjectsSA)
	})
	if err != nil {
		return err
	}

	optbu2 := &terraform.Options{
		TerraformDir: filepath.Join(gcpPath, "business_unit_2", "shared"),
		Logger:       logger,
		NoColor:      true,
	}

	err = state.RunStepE(s, "gcp-projects.bu2.apply-shared", func() error {
		return applyShared(t, optbu2, outputs.ProjectsSA)
	})
	if err != nil {
		return err
	}

	err = state.RunStepE(s, "gcp-projects.plan", func() error {
		return planStep(t, conf, outputs.CICDProject, outputs.DefaultRegion, repo)
	})
	if err != nil {
		return err
	}

	err = state.RunStepE(s, "gcp-projects.production", func() error {
		return applyEnv(t, conf, outputs.CICDProject, outputs.DefaultRegion, repo, "production")
	})
	if err != nil {
		return err
	}

	err = state.RunStepE(s, "gcp-projects.non-production", func() error {
		return applyEnv(t, conf, outputs.CICDProject, outputs.DefaultRegion, repo, "non-production")
	})
	if err != nil {
		return err
	}

	err = state.RunStepE(s, "gcp-projects.development", func() error {
		return applyEnv(t, conf, outputs.CICDProject, outputs.DefaultRegion, repo, "development")
	})
	if err != nil {
		return err
	}

	fmt.Println("end of", step, "deploy")
	return nil
}

func DeployExampleAppStep(t testing.TB, s state.State, tfvars GlobalTfvars, checkoutPath, foundationPath string, outputs InfraPipelineOutputs, logger *logger.Logger) error {
	repo := "bu1-example-app"
	step := "5-app-infra"

	commonTfvars := AppInfraCommonTfvars{
		InstanceRegion:    tfvars.DefaultRegion,
		RemoteStateBucket: outputs.RemoteStateBucket,
	}
	err := utils.WriteTfvars(filepath.Join(foundationPath, step, "common.auto.tfvars"), commonTfvars)
	if err != nil {
		return err
	}

	//prepare policies repo
	gcpPoliciesPath := filepath.Join(checkoutPath, "gcp-policies-app-infra")
	policiesConf := utils.CloneCSR(t, "gcp-policies", gcpPoliciesPath, outputs.InfraPipeProj, logger)
	policiesBranch := "main"

	err = state.RunStepE(s, "bu1-example-app.gcp-policies-app-infra", func() error {
		err = policiesConf.CheckoutBranch(policiesBranch)
		if err != nil {
			return err
		}
		err = utils.CopyDirectory(filepath.Join(foundationPath, "policy-library"), gcpPoliciesPath)
		if err != nil {
			return err
		}
		err = policiesConf.CommitFiles("Initialize policy library repo")
		if err != nil {
			return err
		}
		err = policiesConf.PushBranch(policiesBranch, "origin")
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	gcpPath := filepath.Join(checkoutPath, repo)
	conf := utils.CloneCSR(t, repo, gcpPath, outputs.InfraPipeProj, logger)
	err = conf.CheckoutBranch("plan")
	if err != nil {
		return err
	}

	err = state.RunStepE(s, "bu1-example-app.copy-code", func() error {
		return copyStepCode(t, conf, foundationPath, checkoutPath, repo, step, "")
	})
	if err != nil {
		return err
	}

	err = state.RunStepE(s, "bu1-example-app.plan", func() error {
		return planStep(t, conf, outputs.InfraPipeProj, outputs.DefaultRegion, repo)
	})
	if err != nil {
		return err
	}

	err = state.RunStepE(s, "bu1-example-app.production", func() error {
		return applyEnv(t, conf, outputs.InfraPipeProj, outputs.DefaultRegion, repo, "production")
	})
	if err != nil {
		return err
	}

	err = state.RunStepE(s, "bu1-example-app.non-production", func() error {
		return applyEnv(t, conf, outputs.InfraPipeProj, outputs.DefaultRegion, repo, "non-production")
	})
	if err != nil {
		return err
	}

	err = state.RunStepE(s, "bu1-example-app.development", func() error {
		return applyEnv(t, conf, outputs.InfraPipeProj, outputs.DefaultRegion, repo, "development")
	})
	if err != nil {
		return err
	}

	fmt.Println("end of", step, "deploy")
	return nil
}

func copyStepCode(t testing.TB, conf utils.GitRepo, foundationPath, checkoutPath, repo, step, customPath string) error {
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

func planStep(t testing.TB, conf utils.GitRepo, project, region, repo string) error {
	err := conf.CommitFiles(fmt.Sprintf("Initialize %s repo", repo))
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = conf.PushBranch("plan", "origin")
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = gcp.WaitBuildSuccess(t, project, region, repo, fmt.Sprintf("Terraform %s plan build Failed.", repo))
	if err != nil {
		return err
	}
	return nil
}

func applyEnv(t testing.TB, conf utils.GitRepo, project, region, repo, environment string) error {
	err := conf.CheckoutBranch(environment)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = conf.PushBranch(environment, "origin")
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = gcp.WaitBuildSuccess(t, project, region, repo, fmt.Sprintf("Terraform %s apply %s build Failed.", repo, environment))
	if err != nil {
		return err
	}
	return nil
}

func applyShared(t testing.TB, options *terraform.Options, serviceAccount string) error {
	err := os.Setenv("GOOGLE_IMPERSONATE_SERVICE_ACCOUNT", serviceAccount)
	if err != nil {
		return err
	}
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
	err = os.Unsetenv("GOOGLE_IMPERSONATE_SERVICE_ACCOUNT")
	if err != nil {
		return err
	}
	return nil
}
