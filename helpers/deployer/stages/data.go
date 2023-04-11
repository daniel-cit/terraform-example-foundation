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

package stages

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/mitchellh/go-testing-interface"

	"github.com/terraform-google-modules/terraform-example-foundation/helpers/deployer/utils"
)


type CommonConf struct {
	FoundationPath    string
	CheckoutPath      string
	EnableHubAndSpoke bool
	DisablePrompt     bool
	Logger            *logger.Logger
}

type BootstrapOutputs struct {
	RemoteStateBucket         string
	RemoteStateBucketProjects string
	CICDProject               string
	DefaultRegion             string
	NetworkSA                 string
	ProjectsSA                string
	EnvsSA                    string
	OrgSA                     string
}

type InfraPipelineOutputs struct {
	RemoteStateBucket string
	InfraPipeProj     string
	DefaultRegion     string
	TerraformSA       string
	StateBucket       string
}

type ServerAddress struct {
	Ipv4Address    string `cty:"ipv4_address"`
	ForwardingPath string `cty:"forwarding_path"`
}

type GlobalTfvars struct {
	OrgID                                 string          `hcl:"org_id"`
	BillingAccount                        string          `hcl:"billing_account"`
	GroupOrgAdmins                        string          `hcl:"group_org_admins"`
	GroupBillingAdmins                    string          `hcl:"group_billing_admins"`
	BillingDataUsers                      string          `hcl:"billing_data_users"`
	MonitoringWorkspaceUsers              string          `hcl:"monitoring_workspace_users"`
	AuditDataUsers                        string          `hcl:"audit_data_users"`
	DefaultRegion                         string          `hcl:"default_region"`
	ParentFolder                          *string         `hcl:"parent_folder"`
	Domain                                string          `hcl:"domain"`
	DomainsToAllow                        []string        `hcl:"domains_to_allow"`
	EssentialContactsDomains              []string        `hcl:"essential_contacts_domains_to_allow"`
	PerimeterAdditionalMembers            []string        `hcl:"perimeter_additional_members"`
	TargetNameServerAddresses             []ServerAddress `hcl:"target_name_server_addresses"`
	SccNotificationName                   string          `hcl:"scc_notification_name"`
	ProjectPrefix                         *string         `hcl:"project_prefix"`
	FolderPrefix                          *string         `hcl:"folder_prefix"`
	BucketForceDestroy                    *bool           `hcl:"bucket_force_destroy"`
	AuditLogsTableDeleteContentsOnDestroy *bool           `hcl:"audit_logs_table_delete_contents_on_destroy"`
	LogExportStorageForceDestroy          *bool           `hcl:"log_export_storage_force_destroy"`
	EnableHubAndSpoke                     bool            `hcl:"enable_hub_and_spoke"`
	EnableHubAndSpokeTransitivity         bool            `hcl:"enable_hub_and_spoke_transitivity"`
	CreateUniqueTagKey                    bool            `hcl:"create_unique_tag_key"`
	ProjectsKMSLocation                   string          `hcl:"projects_kms_location"`
	ProjectsGCSLocation                   string          `hcl:"projects_gcs_location"`
	CodeCheckoutPath                      string          `hcl:"code_checkout_path"`
	FoundationCodePath                    string          `hcl:"foundation_code_path"`
	ValidatorProjectId                    *string         `hcl:"validator_project_id"`
}

// HasValidatorProj checks if a Validator Project was provided
func (g GlobalTfvars) HasValidatorProj() bool {
	return g.ValidatorProjectId != nil && *g.ValidatorProjectId != ""
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
	DomainsToAllow                        []string `hcl:"domains_to_allow"`
	EssentialContactsDomains              []string `hcl:"essential_contacts_domains_to_allow"`
	BillingDataUsers                      string   `hcl:"billing_data_users"`
	AuditDataUsers                        string   `hcl:"audit_data_users"`
	SccNotificationName                   string   `hcl:"scc_notification_name"`
	RemoteStateBucket                     string   `hcl:"remote_state_bucket"`
	EnableHubAndSpoke                     bool     `hcl:"enable_hub_and_spoke"`
	CreateACMAPolicy                      bool     `hcl:"create_access_context_manager_access_policy"`
	CreateUniqueTagKey                    bool     `hcl:"create_unique_tag_key"`
	AuditLogsTableDeleteContentsOnDestroy *bool    `hcl:"audit_logs_table_delete_contents_on_destroy"`
	LogExportStorageForceDestroy          *bool    `hcl:"log_export_storage_force_destroy"`
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

func GetBootstrapStepOutputs(t testing.TB, foundationPath string) BootstrapOutputs {
	options := &terraform.Options{
		TerraformDir: filepath.Join(foundationPath, "0-bootstrap"),
		Logger:       logger.Discard,
		NoColor:      true,
	}
	return BootstrapOutputs{
		CICDProject:               terraform.Output(t, options, "cloudbuild_project_id"),
		RemoteStateBucket:         terraform.Output(t, options, "gcs_bucket_tfstate"),
		RemoteStateBucketProjects: terraform.Output(t, options, "projects_gcs_bucket_tfstate"),
		DefaultRegion:             terraform.OutputMap(t, options, "common_config")["default_region"],
		NetworkSA:                 terraform.Output(t, options, "networks_step_terraform_service_account_email"),
		ProjectsSA:                terraform.Output(t, options, "projects_step_terraform_service_account_email"),
		EnvsSA:                    terraform.Output(t, options, "environment_step_terraform_service_account_email"),
		OrgSA:                     terraform.Output(t, options, "organization_step_terraform_service_account_email"),
	}
}

func GetInfraPipelineOutputs(t testing.TB, checkoutPath, workspace string) InfraPipelineOutputs {
	options := &terraform.Options{
		TerraformDir: filepath.Join(checkoutPath, "gcp-projects", "business_unit_1", "shared"),
		Logger:       logger.Discard,
		NoColor:      true,
	}
	return InfraPipelineOutputs{
		InfraPipeProj: terraform.Output(t, options, "cloudbuild_project_id"),
		DefaultRegion: terraform.Output(t, options, "default_region"),
		TerraformSA:   terraform.OutputMap(t, options, "terraform_service_accounts")["bu1-example-app"],
		StateBucket:   terraform.OutputMap(t, options, "state_buckets")["bu1-example-app"],
	}
}

func ReadGlobalTfvars(file string) (GlobalTfvars, error) {
	var globalTfvars GlobalTfvars
	if file == "" {
		return globalTfvars, fmt.Errorf("tfvars file is required.")
	}
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		return globalTfvars, fmt.Errorf("tfvars file '%s' does not exits\n", file)
	}
	err = utils.ReadTfvars(file, &globalTfvars)
	if err != nil {
		return globalTfvars, fmt.Errorf("Failed to load tfvars file %s. Error: %s\n", file, err.Error())
	}
	return globalTfvars, nil
}
