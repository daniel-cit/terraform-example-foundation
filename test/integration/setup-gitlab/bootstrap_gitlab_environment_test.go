// Copyright 2024 Google LLC
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

package bootstrap_gitlab

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/tft"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/terraform-google-modules/terraform-example-foundation/test/integration/testutils"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

// connects to a Google Cloud VM instance using SSH and retrieves the logs from the VM's Startup Script service
func readLogsFromVm(t *testing.T, instanceName string, instanceZone string, instanceProject string) (string, error) {
	args := []string{"compute", "ssh", instanceName, fmt.Sprintf("--zone=%s", instanceZone), fmt.Sprintf("--project=%s", instanceProject), "--command=journalctl -u google-startup-scripts.service -n 20"}
	gcloudCmd := shell.Command{
		Command: "gcloud",
		Args:    args,
	}
	return shell.RunCommandAndGetStdOutE(t, gcloudCmd)
}

func TestValidateStartupScript(t *testing.T) {
	// Retrieve output values from test setup
	setup := tft.NewTFBlueprintTest(t,
		tft.WithTFDir("../../setup"),
	)
	instanceName := setup.GetStringOutput("gitlab_instance_name")
	instanceZone := setup.GetStringOutput("gitlab_instance_zone")
	gitlabSecretProject := setup.GetStringOutput("gitlab_secret_project")
	// Periodically read logs from startup script running on the VM instance
	for count := 0; count < 100; count++ {
		logs, err := readLogsFromVm(t, instanceName, instanceZone, gitlabSecretProject)
		if err != nil {
			t.Fatal(err)
		}

		if strings.Contains(logs, "Finished Google Compute Engine Startup Scripts") {
			if strings.Contains(logs, "exit status 1") {
				t.Fatal("ERROR: Startup Script finished with invalid exit status.")
			}
			break
		}
		time.Sleep(12 * time.Second)
	}
}

func TestBootstrapGitlabVM(t *testing.T) {

	caCert, err := os.ReadFile("/usr/local/share/ca-certificates/gitlab.crt")

	if err != nil {
		t.Fatalf("Failed to read CA certificate: %v", err)
	}

	// Retrieve output values from test setup
	setup := tft.NewTFBlueprintTest(t,
		tft.WithTFDir("../../setup"),
	)

	gitlabSecretProject := setup.GetStringOutput("gitlab_secret_project")
	external_url := setup.GetStringOutput("gitlab_url")
	gitlabPersonalTokenSecretName := setup.GetStringOutput("gitlab_pat_secret_name")

	token, err := testutils.GetSecretFromSecretManager(t, gitlabPersonalTokenSecretName, gitlabSecretProject)
	if err != nil {
		t.Fatal(err)
	}
	git, err := gitlab.NewClient(token, gitlab.WithBaseURL(external_url))
	if err != nil {
		t.Fatal(err)
	}

	repos := map[string]string{
		//repositories
		"seed":         "gcp-seed",
		"cicd":         "gcp-cicd",
		"org":          "gcp-org",
		"env":          "gcp-environments",
		"net":          "gcp-networks",
		"proj":         "gcp-projects",
		"app":          "bu1-example-app", // TODO move to another file to be used in step 4
		"cloudbuilder": "tf-cloudbuilder", // TODO need to be sone how filtered when used in bootstrap step
	}

	repositories := make(map[string]testutils.GitLabRepository)

	for k, repo := range repos {
		p := &gitlab.CreateProjectOptions{
			Name:                 gitlab.Ptr(repo),
			Description:          gitlab.Ptr("Test Repo"),
			InitializeWithReadme: gitlab.Ptr(true),
			Visibility:           gitlab.Ptr(gitlab.PrivateVisibility),
			DefaultBranch:        gitlab.Ptr("master"),
		}
		project, _, err := git.Projects.CreateProject(p)
		if err != nil {
			t.Fatal(err)
		}
		repositories[k] = testutils.GitLabRepository{
			Name: project.Name,
			URL:  project.WebURL,
		}
		t.Log(project.WebURL)
		t.Log(project.Name)
	}

	url := setup.GetStringOutput("gitlab_url")
	gitlabWebhookSecretId := setup.GetStringOutput("gitlab_webhook_secret_id")

	repoConfig := testutils.RepositoryConfig{
		GitlabReadAuthorizerCredentialSecretId: gitlabPersonalTokenSecretName,
		GitlabAuthorizerCredentialSecretId:     gitlabPersonalTokenSecretName,
		GitlabWebhookSecretId:                  gitlabWebhookSecretId,
		GitlabEnterpriseHostUri:                url,
		Repositories:                           repositories,
	}

	err = testutils.WriteGitLabVarConfiguration("../../../0-bootstrap/cicd", "gitlab_repositories.auto.tfvars", repoConfig, caCert)

	if err != nil {
		t.Fatalf("Failed to write repo variables: %v", err)
	}
}
