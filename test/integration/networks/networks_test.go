// Copyright 2022 Google LLC
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

package networks

import (
	"testing"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/tft"
	"github.com/stretchr/testify/assert"
)

func TestNetworks(t *testing.T) {

	for _, tt := range []struct {
		name  string
		tfDir string
	}{
		{
			name: "development",
			tfDir: "../../../3-networks/envs/development",
		},
		{
			name: "non-production",
			tfDir: "../../../3-networks/envs/non-production",
		},
		{
			name: "production",
			tfDir: "../../../3-networks/envs/production",
		},
	}{
		t.Run(tt.name, func(t *testing.T) {
			networks := tft.NewTFBlueprintTest(t,
				tft.WithTFDir(tt.tfDir),
			)
			networks.DefineVerify(
				func(assert *assert.Assertions) {
					// perform default verification ensuring Terraform reports no additional changes on an applied blueprint
					networks.DefaultVerify(assert)
				})
				networks.Test()
		})

	}
}