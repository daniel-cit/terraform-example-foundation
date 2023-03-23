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
	"flag"
	"fmt"
	"os"
	"path/filepath"
	gotest "testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/mitchellh/go-testing-interface"

	"github.com/terraform-google-modules/terraform-example-foundation/helpers/deployer/state"
	"github.com/terraform-google-modules/terraform-example-foundation/helpers/deployer/steps"
	"github.com/terraform-google-modules/terraform-example-foundation/helpers/deployer/utils"
)

type depCfg struct {
	tfvarsFile string
	stateFile  string
	quiet      bool
	help       bool
}

func parseFlags() depCfg {
	var d depCfg

	flag.StringVar(&d.tfvarsFile, "tfvars_file", "", "Full path to the Terraform .tfvars `file` with the complete configuration to be used.")
	flag.StringVar(&d.stateFile, "state_file", ".state.json", "Path to the state `file` to be used to save progress.")
	flag.BoolVar(&d.quiet, "quiet", false, "If true, additional output is suppressed.")
	flag.BoolVar(&d.help, "help", false, "Prints this help text and exits.")

	flag.Parse()
	return d
}

func getLogger(c depCfg) *logger.Logger {
	if c.quiet {
		return logger.Discard
	}
	return logger.Default
}

func getTfvars(c depCfg) steps.GlobalTfvars {
	if c.tfvarsFile == "" {
		fmt.Println("stopping execution, tfvars file is required.")
		os.Exit(1)
	}
	_, err := os.Stat(c.tfvarsFile)
	if os.IsNotExist(err) {
		fmt.Printf("stopping execution, tfvars file '%s' does not exits\n", c.tfvarsFile)
		os.Exit(1)
	}
	var globalTfvars steps.GlobalTfvars
	err = utils.ReadTfvars(c.tfvarsFile, &globalTfvars)
	if err != nil {
		fmt.Printf("failed to load tfvars file %s. Error: %s\n", c.tfvarsFile, err.Error())
		os.Exit(1)
	}
	return globalTfvars
}

func main() {

	cfg := parseFlags()
	if cfg.help {
		fmt.Println("Deploys the Terraform Example Foundation")
		flag.PrintDefaults()
		return
	}

	// load tfvars
	globalTfvars := getTfvars(cfg)

	// init infra
	gotest.Init()
	t := &testing.RuntimeT{}
	foundationCodePath := globalTfvars.FoundationCodePath
	codeCheckoutPath := globalTfvars.CodeCheckoutPath
	logger := getLogger(cfg)

	// load state
	s, err := state.LoadState(cfg.stateFile)
	if err != nil {
		fmt.Printf("failed to load state file %s. Error: %s\n", cfg.stateFile, err.Error())
		os.Exit(2)
	}

	// deploy foundation
	bootstrapOptions := &terraform.Options{
		TerraformDir: filepath.Join(foundationCodePath, "0-bootstrap"),
		Logger:       logger,
		NoColor:      true,
	}

	err = state.RunStepE(s, "gcp-bootstrap", func() error {
		return steps.DeployBootstrapStep(t, s, globalTfvars, bootstrapOptions, codeCheckoutPath, foundationCodePath)
	})
	if err != nil {
		fmt.Printf("Bootstrap step failed. Error: %s\n", err.Error())
		os.Exit(3)
	}

	bootstrapOutputs := steps.GetBootstrapStepOutputs(t, bootstrapOptions)

	// TODO tell the user about the form for asking additional quota projects for the service account of step 4
	// TODO put the link in the output

	err = state.RunStepE(s, "gcp-org", func() error {
		return steps.DeployOrgStep(t, s, globalTfvars, codeCheckoutPath, foundationCodePath, bootstrapOutputs)
	})
	if err != nil {
		fmt.Printf("Org step failed. Error: %s\n", err.Error())
		os.Exit(3)
	}

	err = state.RunStepE(s, "gcp-environments", func() error {
		return steps.DeployEnvStep(t, s, globalTfvars, codeCheckoutPath, foundationCodePath, bootstrapOutputs)
	})
	if err != nil {
		fmt.Printf("Environments step failed. Error: %s\n", err.Error())
		os.Exit(3)
	}

	err = state.RunStepE(s, "gcp-networks", func() error {
		return steps.DeployNetworksStep(t, s, globalTfvars, codeCheckoutPath, foundationCodePath, bootstrapOutputs, logger)
	})
	if err != nil {
		fmt.Printf("Networks step failed. Error: %s\n", err.Error())
		os.Exit(3)
	}

	err = state.RunStepE(s, "gcp-projects", func() error {
		return steps.DeployProjectsStep(t, s, globalTfvars, codeCheckoutPath, foundationCodePath, bootstrapOutputs, logger)
	})
	if err != nil {
		fmt.Printf("Projects step failed. Error: %s\n", err.Error())
		os.Exit(3)
	}
	//5-app-infra
	InfraPipelineOptions := &terraform.Options{
		TerraformDir: filepath.Join(codeCheckoutPath, "gcp-projects", "business_unit_1", "shared"),
		Logger:       logger,
		NoColor:      true,
	}
	infraPipelineOutputs := steps.GetInfraPipelineOutputs(t, InfraPipelineOptions, "bu1-example-app")
	infraPipelineOutputs.RemoteStateBucket = bootstrapOutputs.RemoteStateBucketProjects
	err = state.RunStepE(s, "bu1-example-app", func() error {
		return steps.DeployExampleAppStep(t, s, globalTfvars, codeCheckoutPath, foundationCodePath, infraPipelineOutputs)
	})
	if err != nil {
		fmt.Printf("Example app step failed. Error: %s\n", err.Error())
		os.Exit(3)
	}

	// TODO ask about the answer of thr request form to asking additional quota projects for the service account of step 4
}
