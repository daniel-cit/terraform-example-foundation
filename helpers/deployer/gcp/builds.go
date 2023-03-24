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
	"time"

	"github.com/GoogleCloudPlatform/cloud-foundation-toolkit/infra/blueprint-test/pkg/gcloud"
	"github.com/mitchellh/go-testing-interface"
	"github.com/tidwall/gjson"
)

func GetBuilds(t testing.TB, projectID, region, filter string) []gjson.Result {
	gcOps := gcloud.WithCommonArgs([]string{"--project", projectID, "--region", region, "--filter", filter, "--format", "json"})
	return gcloud.Run(t, "builds list", gcOps).Array()
}

func GetRunningBuild(t testing.TB, projectID string, region string, filter string) string {
	time.Sleep(20 * time.Second)
	builds := GetBuilds(t, projectID, region, filter)
	for _, build := range builds {
		status := build.Get("status").String()
		if status == "QUEUED" || status == "WORKING" {
			return build.Get("id").String()
		}
	}
	return ""
}

func GetLastBuildStatus(t testing.TB, projectID string, region string, filter string) string {
	gcOps := gcloud.WithCommonArgs([]string{"--project", projectID, "--region", region, "--limit", "1", "--sort-by", "~createTime", "--filter", filter, "--format", "json"})
	build := gcloud.Run(t, "builds list", gcOps).Array()[0]
	return build.Get("status").String()
}

func GetBuildState(t testing.TB, projectID, region, buildID string) string {
	gcOps := gcloud.WithCommonArgs([]string{buildID, "--project", projectID, "--region", region, "--format", "json(status)"})
	return gcloud.Run(t, "builds describe", gcOps).Get("status").String()
}

func GetTerminalState(t testing.TB, projectID string, region string, buildID string) string {
	var status string
	fmt.Printf("waiting for build %s execution.\n", buildID)
	status = GetBuildState(t, projectID, region, buildID)
	for status != "SUCCESS" && status != "FAILURE" && status != "CANCELLED" {
		fmt.Printf("build status is %s\n", status)
		time.Sleep(20 * time.Second)
		status = GetBuildState(t, projectID, region, buildID)
	}
	fmt.Printf("final build status is %s\n", status)
	return status
}

func WaitBuildSuccess(t testing.TB, project, region, repo, failureMsg string) error {
	filter := fmt.Sprintf("source.repoSource.repoName:%s", repo)
	build := GetRunningBuild(t, project, region, filter)
	if build != "" {
		status := GetTerminalState(t, project, region, build)
		if status != "SUCCESS" {
			return fmt.Errorf("%s\nSee:\nhttps://console.cloud.google.com/cloud-build/builds;region=%s/%s?project=%s\nfor details.\n", failureMsg, region, build, project)
		}
	} else {
		status := GetLastBuildStatus(t, project, region, filter)
		if status != "SUCCESS" {
			return fmt.Errorf("%s\nSee:\nhttps://console.cloud.google.com/cloud-build/builds;region=%s/%s?project=%s\nfor details.\n", failureMsg, region, build, project)
		}
	}
	return nil
}
