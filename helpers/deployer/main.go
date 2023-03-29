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
	"strings"
	gotest "testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/mitchellh/go-testing-interface"

	"github.com/terraform-google-modules/terraform-example-foundation/helpers/deployer/gcp"
	"github.com/terraform-google-modules/terraform-example-foundation/helpers/deployer/msg"
	"github.com/terraform-google-modules/terraform-example-foundation/helpers/deployer/stages"
	"github.com/terraform-google-modules/terraform-example-foundation/helpers/deployer/steps"
	"github.com/terraform-google-modules/terraform-example-foundation/helpers/deployer/utils"
)

const (
	replaceME     = "REPLACE_ME"
	exampleDotCom = "example.com"
)

type depCfg struct {
	tfvarsFile    string
	stepsFile     string
	resetStep     string
	quiet         bool
	help          bool
	listSteps     bool
	disablePrompt bool
	validate      bool
}

func parseFlags() depCfg {
	var d depCfg

	flag.StringVar(&d.tfvarsFile, "tfvars_file", "", "Full path to the Terraform .tfvars `file` with the configuration to be used.")
	flag.StringVar(&d.stepsFile, "steps_file", ".steps.json", "Path to the steps `file` to be used to save progress.")
	flag.StringVar(&d.resetStep, "reset_step", "", "Name of a `step` to be reset.")
	flag.BoolVar(&d.quiet, "quiet", false, "If true, additional output is suppressed.")
	flag.BoolVar(&d.help, "help", false, "Prints this help text and exits.")
	flag.BoolVar(&d.listSteps, "list_steps", false, "List the existing steps.")
	flag.BoolVar(&d.disablePrompt, "disable_prompt", false, "Disable interactive prompt.")
	flag.BoolVar(&d.validate, "validate", false, "Validate tfvars file inputs.")

	flag.Parse()
	return d
}

func getLogger(quiet bool) *logger.Logger {
	if quiet {
		return logger.Discard
	}
	return logger.Default
}

func getTfvars(file string) stages.GlobalTfvars {
	if file == "" {
		fmt.Println("# Stopping execution, tfvars file is required.")
		os.Exit(1)
	}
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		fmt.Printf("# Stopping execution, tfvars file '%s' does not exits\n", file)
		os.Exit(1)
	}
	var globalTfvars stages.GlobalTfvars
	err = utils.ReadTfvars(file, &globalTfvars)
	if err != nil {
		fmt.Printf("# Failed to load tfvars file %s. Error: %s\n", file, err.Error())
		os.Exit(1)
	}
	return globalTfvars
}

func validateDirectories(g stages.GlobalTfvars) {
	_, err := os.Stat(g.FoundationCodePath)
	if os.IsNotExist(err) {
		fmt.Printf("# Stopping execution, FoundationCodePath directory '%s' does not exits\n", g.FoundationCodePath)
		os.Exit(1)
	}
	_, err = os.Stat(g.CodeCheckoutPath)
	if os.IsNotExist(err) {
		fmt.Printf("# Stopping execution, CodeCheckoutPath directory '%s' does not exits\n", g.CodeCheckoutPath)
		os.Exit(1)
	}
}

func validate(t testing.TB, g stages.GlobalTfvars) {

	gcpConf := gcp.NewGCP()
	fmt.Println("")
	fmt.Println("# Validating tfvar file.")
	if gcpConf.HasSccNotification(t, g.OrgID, g.SccNotificationName) {
		fmt.Printf("# Notification '%s' exists in organization '%s'. Chose a different one.\n", g.SccNotificationName, g.OrgID)
		fmt.Printf("# See existing Notifications for organization '%s'.\n", g.OrgID)
		fmt.Printf("# gcloud scc notifications list organizations/%s --filter=\"name:organizations/%s/notificationConfigs/%s\" --format=\"value(name)\"\n", g.OrgID, g.OrgID, g.SccNotificationName)
		fmt.Println("")
	}
	if !g.CreateUniqueTagKey && gcpConf.HasTagKey(t, g.OrgID, "environment") {
		fmt.Printf("# Tag key 'environment' exists in organization '%s'.\n", g.OrgID)
		fmt.Println("# Set variable 'create_unique_tag_key' to 'true' in the tfvar file.")
		fmt.Println("")
	}
	if g.OrgID == replaceME {
		fmt.Println("# Replace value for input 'org_id'")
	}
	if g.BillingAccount == replaceME {
		fmt.Println("# Replace value for input 'billing_account'")
	}
	if g.GroupOrgAdmins == replaceME {
		fmt.Println("# Replace value for input 'group_org_admins'")
	}
	if g.GroupBillingAdmins == replaceME {
		fmt.Println("# Replace value for input 'group_billing_admins'")
	}
	if g.BillingDataUsers == replaceME {
		fmt.Println("# Replace value for input 'billing_data_users'")
	}
	if g.MonitoringWorkspaceUsers == replaceME {
		fmt.Println("# Replace value for input 'monitoring_workspace_users'")
	}
	if g.AuditDataUsers == replaceME {
		fmt.Println("# Replace value for input 'audit_data_users'")
	}
	if strings.Contains(g.Domain, exampleDotCom) {
		fmt.Println("# Replace value 'example.com' for input 'domain'")
	}
	if g.Domain != "" && g.Domain[len(g.Domain)-1:] != "." {
		fmt.Println("# Value for input 'domain' must end with '.'")
	}
	for _, d := range g.DomainsToAllow {
		if strings.Contains(d, exampleDotCom) {
			fmt.Println("# Replace value 'example.com' for input 'domains_to_allow'")
		}
	}
	for _, e := range g.EssentialContactsDomains {
		if strings.Contains(e, exampleDotCom) {
			fmt.Println("# Replace value 'example.com' for input 'essential_contacts_domains_to_allow'")
		}
		if e != "" && e[0:1] != "@" {
			fmt.Printf("# Essential contacts must start with '@': '%s'\n", e)
		}
	}
	for _, p := range g.PerimeterAdditionalMembers {
		if strings.Contains(p, exampleDotCom) {
			fmt.Printf("# Replace value for input 'perimeter_additional_members': '%s'\n", p)
		}
		if strings.Contains(p, "group:") {
			fmt.Printf("# VPC Service Controls does not allow groups in the perimeter: '%s'\n", p)
		}
	}
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
		fmt.Printf("# failed to load state file %s. Error: %s\n", cfg.stepsFile, err.Error())
		os.Exit(2)
	}
	if cfg.listSteps {
		fmt.Println("# Executed steps:")
		e := s.ListSteps()
		if len(e) == 0 {
			fmt.Println("# No steps executed")
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
	globalTfvars := getTfvars(cfg.tfvarsFile)

	// validate inputs
	validateDirectories(globalTfvars)

	// init infra
	gotest.Init()
	t := &testing.RuntimeT{}
	conf := stages.CommonConf{
		FoundationPath: globalTfvars.FoundationCodePath,
		CheckoutPath:   globalTfvars.CodeCheckoutPath,
		DisablePrompt:  cfg.disablePrompt,
		Logger:         getLogger(cfg.quiet),
	}

	if cfg.validate {
		validate(t, globalTfvars)
		return
	}

	// deploy stages
	msg.PrintStageMsg("Deploying 0-bootstrap stage")
	skipInnerBuildMsg := s.IsStepComplete("gcp-bootstrap")
	err = s.RunStep("gcp-bootstrap", func() error {
		return stages.DeployBootstrapStage(t, s, globalTfvars, conf)
	})
	if err != nil {
		fmt.Printf("# Bootstrap step failed. Error: %s\n", err.Error())
		os.Exit(3)
	}

	bo := stages.GetBootstrapStepOutputs(t, conf.FoundationPath)

	if skipInnerBuildMsg {
		msg.PrintBuildMsg(bo.CICDProject, bo.DefaultRegion, conf.DisablePrompt)
	}
	msg.PrintQuotaMsg(bo.ProjectsSA, conf.DisablePrompt)

	msg.PrintStageMsg("Deploying 1-org stage")
	err = s.RunStep("gcp-org", func() error {
		return stages.DeployOrgStage(t, s, globalTfvars, bo, conf)
	})
	if err != nil {
		fmt.Printf("# Org step failed. Error: %s\n", err.Error())
		os.Exit(3)
	}

	msg.PrintStageMsg("Deploying 2-environments stage")
	err = s.RunStep("gcp-environments", func() error {
		return stages.DeployEnvStage(t, s, globalTfvars, bo, conf)
	})
	if err != nil {
		fmt.Printf("# Environments step failed. Error: %s\n", err.Error())
		os.Exit(3)
	}

	msg.PrintStageMsg("Deploying 3-networks stage")
	err = s.RunStep("gcp-networks", func() error {
		return stages.DeployNetworksStage(t, s, globalTfvars, bo, conf)
	})
	if err != nil {
		fmt.Printf("# Networks step failed. Error: %s\n", err.Error())
		os.Exit(3)
	}

	msg.PrintStageMsg("Deploying 4-projects stage")
	msg.ConfirmQuota(bo.ProjectsSA, conf.DisablePrompt)

	err = s.RunStep("gcp-projects", func() error {
		return stages.DeployProjectsStage(t, s, globalTfvars, bo, conf)
	})
	if err != nil {
		fmt.Printf("# Projects step failed. Error: %s\n", err.Error())
		os.Exit(3)
	}

	msg.PrintStageMsg("Deploying 5-app-infra stage")
	io := stages.GetInfraPipelineOutputs(t, conf.CheckoutPath, "bu1-example-app")
	io.RemoteStateBucket = bo.RemoteStateBucketProjects

	msg.PrintBuildMsg(io.InfraPipeProj, io.DefaultRegion, conf.DisablePrompt)

	err = s.RunStep("bu1-example-app", func() error {
		return stages.DeployExampleAppStage(t, s, globalTfvars, io, conf)
	})
	if err != nil {
		fmt.Printf("# Example app step failed. Error: %s\n", err.Error())
		os.Exit(3)
	}
}
