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

type GCP struct {
	Runf func(t testing.TB, cmd string, args ...interface{}) gjson.Result
}

func NewGCP() GCP {
	return GCP{
		Runf: gcloud.Runf,
	}
}

func (g GCP) GetBuilds(t testing.TB, projectID, region, filter string) []gjson.Result {
	return g.Runf(t, "builds list --project %s --region %s --filter %s", projectID, region, filter).Array()
}

func (g GCP) GetLastBuildStatus(t testing.TB, projectID, region, filter string) string {
	return g.Runf(t, "builds list --project %s --region %s --limit 1 --sort-by ~createTime --filter %s", projectID, region, filter).Array()[0].Get("status").String()
}

func (g GCP) GetBuildState(t testing.TB, projectID, region, buildID string) string {
	return g.Runf(t, "builds describe %s  --project %s --region %s", buildID, projectID, region).Get("status").String()
}

func (g GCP) GetRunningBuildID(t testing.TB, projectID, region, filter string) string {
	time.Sleep(20 * time.Second)
	builds := g.GetBuilds(t, projectID, region, filter)
	for _, build := range builds {
		status := build.Get("status").String()
		if status == "QUEUED" || status == "WORKING" {
			return build.Get("id").String()
		}
	}
	return ""
}

func (g GCP) GetFinalBuildState(t testing.TB, projectID, region, buildID string) string {
	var status string
	fmt.Printf("waiting for build %s execution.\n", buildID)
	status = g.GetBuildState(t, projectID, region, buildID)
	for status != "SUCCESS" && status != "FAILURE" && status != "CANCELLED" {
		fmt.Printf("build status is %s\n", status)
		time.Sleep(20 * time.Second)
		status = g.GetBuildState(t, projectID, region, buildID)
	}
	fmt.Printf("final build status is %s\n", status)
	return status
}

func (g GCP) WaitBuildSuccess(t testing.TB, project, region, repo, failureMsg string) error {
	filter := fmt.Sprintf("source.repoSource.repoName:%s", repo)
	build := g.GetRunningBuildID(t, project, region, filter)
	if build != "" {
		status := g.GetFinalBuildState(t, project, region, build)
		if status != "SUCCESS" {
			return fmt.Errorf("%s\nSee:\nhttps://console.cloud.google.com/cloud-build/builds;region=%s/%s?project=%s\nfor details.\n", failureMsg, region, build, project)
		}
	} else {
		status := g.GetLastBuildStatus(t, project, region, filter)
		if status != "SUCCESS" {
			return fmt.Errorf("%s\nSee:\nhttps://console.cloud.google.com/cloud-build/builds;region=%s/%s?project=%s\nfor details.\n", failureMsg, region, build, project)
		}
	}
	return nil
}
