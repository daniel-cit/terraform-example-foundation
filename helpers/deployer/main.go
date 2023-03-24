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

	"github.com/terraform-google-modules/terraform-example-foundation/helpers/deployer/stages"
	"github.com/terraform-google-modules/terraform-example-foundation/helpers/deployer/steps"
	"github.com/terraform-google-modules/terraform-example-foundation/helpers/deployer/utils"
)

type depCfg struct {
	tfvarsFile string
	stepsFile  string
	resetStep  string
	quiet      bool
	help       bool
	listSteps  bool
}

func parseFlags() depCfg {
	var d depCfg

	flag.StringVar(&d.tfvarsFile, "tfvars_file", "", "Full path to the Terraform .tfvars `file` with the complete configuration to be used.")
	flag.StringVar(&d.stepsFile, "steps_file", ".steps.json", "Path to the stepes `file` to be used to save progress.")
	flag.StringVar(&d.resetStep, "reset_step", "", "Name of a step to be reset.")
	flag.BoolVar(&d.quiet, "quiet", false, "If true, additional output is suppressed.")
	flag.BoolVar(&d.help, "help", false, "Prints this help text and exits.")
	flag.BoolVar(&d.listSteps, "list_steps", false, "List the existing steps.")

	flag.Parse()
	return d
}

func getLogger(c depCfg) *logger.Logger {
	if c.quiet {
		return logger.Discard
	}
	return logger.Default
}

func getTfvars(c depCfg) stages.GlobalTfvars {
	if c.tfvarsFile == "" {
		fmt.Println("stopping execution, tfvars file is required.")
		os.Exit(1)
	}
	_, err := os.Stat(c.tfvarsFile)
	if os.IsNotExist(err) {
		fmt.Printf("stopping execution, tfvars file '%s' does not exits\n", c.tfvarsFile)
		os.Exit(1)
	}
	var globalTfvars stages.GlobalTfvars
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

	s, err := steps.LoadSteps(cfg.stepsFile)
	if err != nil {
		fmt.Printf("failed to load state file %s. Error: %s\n", cfg.stepsFile, err.Error())
		os.Exit(2)
	}
	if cfg.listSteps {
		fmt.Println("Executed steps:")
		e := s.ListSteps()
		if len(e) == 0 {
			fmt.Println("No steps executed")
			return
		}
		for _, step := range e {
			fmt.Println(step)
		}
		return
	}
	if cfg.resetStep != "" {
		s.ResetStep(cfg.resetStep)
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

	// deploy stages
	err = s.RunStep("gcp-bootstrap", func() error {
		return stages.DeployBootstrapStage(t, s, globalTfvars, codeCheckoutPath, foundationCodePath, logger)
	})
	if err != nil {
		fmt.Printf("Bootstrap step failed. Error: %s\n", err.Error())
		os.Exit(3)
	}

	bootstrapOutputs := stages.GetBootstrapStepOutputs(t, foundationCodePath)

	// TODO tell the user about the form for asking additional quota projects for the service account of step 4
	// TODO put the link in the output

	err = s.RunStep("gcp-org", func() error {
		return stages.DeployOrgStage(t, s, globalTfvars, codeCheckoutPath, foundationCodePath, bootstrapOutputs, logger)
	})
	if err != nil {
		fmt.Printf("Org step failed. Error: %s\n", err.Error())
		os.Exit(3)
	}

	err = s.RunStep("gcp-environments", func() error {
		return stages.DeployEnvStage(t, s, globalTfvars, codeCheckoutPath, foundationCodePath, bootstrapOutputs, logger)
	})
	if err != nil {
		fmt.Printf("Environments step failed. Error: %s\n", err.Error())
		os.Exit(3)
	}

	err = s.RunStep("gcp-networks", func() error {
		return stages.DeployNetworksStage(t, s, globalTfvars, codeCheckoutPath, foundationCodePath, bootstrapOutputs, logger)
	})
	if err != nil {
		fmt.Printf("Networks step failed. Error: %s\n", err.Error())
		os.Exit(3)
	}

	err = s.RunStep("gcp-projects", func() error {
		return stages.DeployProjectsStage(t, s, globalTfvars, codeCheckoutPath, foundationCodePath, bootstrapOutputs, logger)
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
	infraPipelineOutputs := stages.GetInfraPipelineOutputs(t, InfraPipelineOptions, "bu1-example-app")
	infraPipelineOutputs.RemoteStateBucket = bootstrapOutputs.RemoteStateBucketProjects
	err = s.RunStep("bu1-example-app", func() error {
		return stages.DeployExampleAppStage(t, s, globalTfvars, codeCheckoutPath, foundationCodePath, infraPipelineOutputs, logger)
	})
	if err != nil {
		fmt.Printf("Example app step failed. Error: %s\n", err.Error())
		os.Exit(3)
	}

	// TODO ask about the answer of thr request form to asking additional quota projects for the service account of step 4
}
