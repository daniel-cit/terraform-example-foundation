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
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/mitchellh/go-testing-interface"

	"github.com/terraform-google-modules/terraform-example-foundation/helpers/deployer/steps"
	"github.com/terraform-google-modules/terraform-example-foundation/helpers/deployer/utils"
	"github.com/terraform-google-modules/terraform-example-foundation/test/integration/testutils"
)

const (
	MaxBuildRetries = 40
)

func DestroyBootstrapStage(t testing.TB, s steps.Steps, c CommonConf) error {
	repo := "gcp-bootstrap"
	step := "0-bootstrap"
	gcpPath := filepath.Join(c.CheckoutPath, repo)

	// remove backend.tf file
	tfDir := filepath.Join(gcpPath, "envs", "shared")
	backendF := filepath.Join(tfDir, "backend.tf")
	exist, err := utils.FileExists(backendF)
	if err != nil {
		return err
	}
	if exist {
		options := &terraform.Options{
			TerraformDir: tfDir,
			Logger:       c.Logger,
			NoColor:      true,
		}
		initOutput, err := terraform.InitE(t, options)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", initOutput)
		err = utils.CopyFile(backendF, filepath.Join(tfDir, "backend.tf.backup"))
		if err != nil {
			return err
		}
		err = os.Remove(backendF)
		if err != nil {
			return err
		}
	}

	err = s.RunDestroyStep("gcp-bootstrap.production", func() error {
		options := &terraform.Options{
			TerraformDir: tfDir,
			Logger:       c.Logger,
			NoColor:      true,
			MigrateState: true,
		}
		conf := utils.CloneCSR(t, repo, gcpPath, "", c.Logger)
		err := conf.CheckoutBranch("production")
		if err != nil {
			return err
		}
		return destroyEnv(t, options, "")
	})
	if err != nil {
		return err
	}
	fmt.Println("end of", step, "destroy")
	return nil
}

func DestroyOrgStage(t testing.TB, s steps.Steps, outputs BootstrapOutputs, c CommonConf) error {
	repo := "gcp-org"
	step := "1-org"
	gcpPath := filepath.Join(c.CheckoutPath, repo)

	err := s.RunDestroyStep("gcp-org.production", func() error {
		options := &terraform.Options{
			TerraformDir: filepath.Join(gcpPath, "envs", "shared"),
			Logger:       c.Logger,
			NoColor:      true,
		}
		conf := utils.CloneCSR(t, repo, gcpPath, outputs.CICDProject, c.Logger)
		err := conf.CheckoutBranch("production")
		if err != nil {
			return err
		}
		return destroyEnv(t, options, outputs.OrgSA)
	})
	if err != nil {
		return err
	}
	fmt.Println("end of", step, "destroy")
	return nil
}

func DestroyEnvStage(t testing.TB, s steps.Steps, outputs BootstrapOutputs, c CommonConf) error {
	repo := "gcp-environments"
	step := "2-environments"

	gcpPath := filepath.Join(c.CheckoutPath, repo)

	for _, e := range []string{"development", "non-production", "production"} {
		err := s.RunDestroyStep(fmt.Sprintf("gcp-environments.%s", e), func() error {
			options := &terraform.Options{
				TerraformDir: filepath.Join(gcpPath, "envs", e),
				Logger:       c.Logger,
				NoColor:      true,
			}
			conf := utils.CloneCSR(t, repo, gcpPath, outputs.CICDProject, c.Logger)
			err := conf.CheckoutBranch(e)
			if err != nil {
				return err
			}
			return destroyEnv(t, options, outputs.EnvsSA)
		})
		if err != nil {
			return err
		}
	}
	fmt.Println("end of", step, "destroy")
	return nil
}

func DestroyNetworksStage(t testing.TB, s steps.Steps, outputs BootstrapOutputs, c CommonConf) error {
	repo := "gcp-networks"
	var step string
	if c.EnableHubAndSpoke {
		step = "3-networks-hub-and-spoke"
	} else {
		step = "3-networks-dual-svpc"
	}
	gcpPath := filepath.Join(c.CheckoutPath, repo)

	for _, e := range []string{"development", "non-production", "production"} {
		err := s.RunDestroyStep(fmt.Sprintf("gcp-networks.%s", e), func() error {
			options := &terraform.Options{
				TerraformDir:             filepath.Join(gcpPath, "envs", e),
				Logger:                   c.Logger,
				NoColor:                  true,
				RetryableTerraformErrors: testutils.RetryableTransientErrors,
				MaxRetries:               2,
				TimeBetweenRetries:       2 * time.Minute,
			}
			conf := utils.CloneCSR(t, repo, gcpPath, outputs.CICDProject, c.Logger)
			err := conf.CheckoutBranch(e)
			if err != nil {
				return err
			}
			return destroyEnv(t, options, outputs.NetworkSA)
		})
		if err != nil {
			return err
		}
	}
	err := s.RunDestroyStep("gcp-networks.apply-shared", func() error {
		options := &terraform.Options{
			TerraformDir: filepath.Join(gcpPath, "envs", "shared"),
			Logger:       c.Logger,
			NoColor:      true,
			RetryableTerraformErrors: testutils.RetryableTransientErrors,
			MaxRetries:               2,
			TimeBetweenRetries:       2 * time.Minute,
		}
		conf := utils.CloneCSR(t, repo, gcpPath, outputs.CICDProject, c.Logger)
		err := conf.CheckoutBranch("production")
		if err != nil {
			return err
		}
		return destroyEnv(t, options, outputs.NetworkSA)
	})
	if err != nil {
		return err
	}
	fmt.Println("end of", step, "destroy")
	return nil
}

func DestroyProjectsStage(t testing.TB, s steps.Steps, outputs BootstrapOutputs, c CommonConf) error {
	repo := "gcp-projects"
	step := "4-projects"
	gcpPath := filepath.Join(c.CheckoutPath, repo)

	for _, e := range []string{"development", "non-production", "production"} {
		err := s.RunDestroyStep(fmt.Sprintf("gcp-projects.%s", e), func() error {
			for _, u := range []string{"business_unit_1", "business_unit_2"} {
				options := &terraform.Options{
					TerraformDir: filepath.Join(gcpPath, u, e),
					Logger:       c.Logger,
					NoColor:      true,
				}
				conf := utils.CloneCSR(t, repo, gcpPath, outputs.CICDProject, c.Logger)
				err := conf.CheckoutBranch(e)
				if err != nil {
					return err
				}
				err = destroyEnv(t, options, outputs.ProjectsSA)
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	for _, u := range []string{"business_unit_1", "business_unit_2"} {
		err := s.RunDestroyStep(fmt.Sprintf("gcp-projects.%s.apply-shared", u), func() error {
			options := &terraform.Options{
				TerraformDir: filepath.Join(gcpPath, u, "shared"),
				Logger:       c.Logger,
				NoColor:      true,
			}
			conf := utils.CloneCSR(t, repo, gcpPath, outputs.CICDProject, c.Logger)
			err := conf.CheckoutBranch("production")
			if err != nil {
				return err
			}
			return destroyEnv(t, options, outputs.ProjectsSA)
		})
		if err != nil {
			return err
		}
	}

	fmt.Println("end of", step, "destroy")
	return nil
}

func DestroyExampleAppStage(t testing.TB, s steps.Steps, outputs InfraPipelineOutputs, c CommonConf) error {
	repo := "bu1-example-app"
	step := "5-app-infra"
	gcpPath := filepath.Join(c.CheckoutPath, repo)

	for _, e := range []string{"development", "non-production", "production"} {
		err := s.RunDestroyStep(fmt.Sprintf("bu1-example-app.%s", e), func() error {
			options := &terraform.Options{
				TerraformDir: filepath.Join(gcpPath, "business_unit_1", e),
				Logger:       c.Logger,
				NoColor:      true,
			}
			conf := utils.CloneCSR(t, repo, gcpPath, outputs.InfraPipeProj, c.Logger)
			err := conf.CheckoutBranch(e)
			err = utils.ReplaceStringInFile(filepath.Join(options.TerraformDir, "backend.tf"), "UPDATE_APP_INFRA_BUCKET", outputs.StateBucket)
			if err != nil {
				return err
			}
			if err != nil {
				return err
			}
			return destroyEnv(t, options, outputs.TerraformSA)
		})
		if err != nil {
			return err
		}
	}

	fmt.Println("end of", step, "destroy")
	return nil
}

func destroyEnv(t testing.TB, options *terraform.Options, serviceAccount string) error {
	var err error

	if serviceAccount != "" {
		err = os.Setenv("GOOGLE_IMPERSONATE_SERVICE_ACCOUNT", serviceAccount)
		if err != nil {
			return err
		}
	}

	initOutput, err := terraform.InitE(t, options)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", initOutput)

	destroyOutput, err := terraform.DestroyE(t, options)
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", destroyOutput)

	if serviceAccount != "" {
		err = os.Unsetenv("GOOGLE_IMPERSONATE_SERVICE_ACCOUNT")
		if err != nil {
			return err
		}
	}
	return nil
}
