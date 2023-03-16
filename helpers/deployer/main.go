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

package main

import (
	"fmt"
	"flag"
	"log"
	"os"
	"path/filepath"
	gotest "testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/mitchellh/go-testing-interface"

	"github.com/terraform-google-modules/terraform-example-foundation/helpers/deployer/utils"
)

var (
	tfvarsFile = flag.String("tfvars", "", "Full path to the tfvars file with the complete configuration to be used")
	verbose    = flag.Bool("verbose", false, "If true, additional output is logged.")
	help       = flag.Bool("help", false, "If true, prints help text and exits.")
)

func main() {

	flag.Parse()

	if *tfvarsFile == "" {
		fmt.Println("stopping execution, tfvars file is required.")
		os.Exit(1)
	}
	_, err := os.Stat(*tfvarsFile)
	if os.IsNotExist(err) {
		fmt.Println("stopping execution, tfvars file", *tfvarsFile, "does not exits")
		os.Exit(2)
	}

	var globalTfvars utils.GlobalTfvars
	err = utils.ReadTfvars(*tfvarsFile, &globalTfvars)
	if err != nil {
		log.Fatal(err)
	}
	// TODO: warn destroy flag not set


	foundationCodePath := globalTfvars.FoundationCodePath
	codeCheckoutPath := globalTfvars.CodeCheckoutPath

	logger := logger.Default
	bootstrapOptions := &terraform.Options{
		TerraformDir: filepath.Join(foundationCodePath, "0-bootstrap"),
		Logger:       logger,
		NoColor:      true,
	}

	gotest.Init()
	t := &testing.RuntimeT{}

	bootstrapOutputs := utils.DeployBootstrapStep(t, globalTfvars, bootstrapOptions, codeCheckoutPath, foundationCodePath)

	utils.DeployOrgStep(t, globalTfvars, codeCheckoutPath, foundationCodePath, bootstrapOutputs)

	utils.DeployEnvStep(t, globalTfvars, codeCheckoutPath, foundationCodePath, bootstrapOutputs)
}
