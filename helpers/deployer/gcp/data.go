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

package gcp

import (
	"fmt"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/gcloud"
	"github.com/mitchellh/go-testing-interface"

	"github.com/terraform-google-modules/terraform-example-foundation/test/integration/testutils"
)

//GetAccessContextManagerPolicyID gets the access context manager policy ID of the organization
func GetAccessContextManagerPolicyID(t testing.TB, ordID string) string {
	filter := fmt.Sprintf("parent:organizations/%s", ordID)
	acmpID := gcloud.Runf(t, "access-context-manager policies list --organization %s --filter %s ", ordID, filter).Array()
	if len(acmpID) == 0 {
		return ""
	}
	return testutils.GetLastSplitElement(acmpID[0].Get("name").String(), "/")
}

// GetOrganizationDomain gets the domain name of the organization
func GetOrganizationDomain(t testing.TB, ordID string) string {
	return gcloud.Runf(t, "organizations describe %s", ordID).Get("displayName").String()
}
