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
	"strings"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/mitchellh/go-testing-interface"
	"github.com/tidwall/gjson"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/gcloud"
	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/utils"
)

func TerraformVet(t testing.TB, terraformDir, policyPath, project string, c CommonConf) (string, error) {
	var result string
	options := &terraform.Options{
		TerraformDir: terraformDir,
		Logger:       c.Logger,
		NoColor:      true,
		PlanFilePath: filepath.Join(os.TempDir(), "plan.tfplan"),
	}
	_, err := terraform.PlanE(t, options)
	if err != nil {
		return result, err
	}
	jsonPlan, err := terraform.ShowE(t, options)
	if err != nil {
		return result, err
	}
	jsonFile, err := utils.WriteTmpFileWithExtension(jsonPlan, "json")
	defer os.Remove(jsonFile)
	defer os.Remove(options.PlanFilePath)
	if err != nil {
		return result, err
	}
	result, err = gcloud.RunCmdE(t, fmt.Sprintf("beta terraform vet %s --policy-library=%s --project=%s", jsonFile, policyPath, project))
	if err != nil && !(strings.Contains(err.Error(), "Validating resources") && strings.Contains(err.Error(), "done")) {
		return result, err
	}
	if !gjson.Valid(result) {
		return result, fmt.Errorf("Error parsing output, invalid json: %s", result)
	}
	return gjson.Parse(result).String(), nil
}
