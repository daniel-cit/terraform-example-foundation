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
	"errors"
	"os"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

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

func ReadTfvars(filename string, val interface{}) error {
	data, diagnostic := hclparse.NewParser().ParseHCLFile(filename)
	if diagnostic.HasErrors() {
		return errors.New(diagnostic.Error())
	}
	decoded := gohcl.DecodeBody(data.Body, nil, val)
	if decoded.HasErrors() {
		return errors.New(decoded.Error())
	}
	return nil
}

func WriteTfvars(filename string, val interface{}) error {
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(val, f.Body())
	return os.WriteFile(filename, f.Bytes(), 0644)
}
