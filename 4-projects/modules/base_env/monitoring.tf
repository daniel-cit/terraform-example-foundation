/**
 * Copyright 2023 Google LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

locals {
  monitored_projects = {
    base_project       = module.base_shared_vpc_project.project_number
    restricted_project = module.restricted_shared_vpc_project.project_number
    floating           = module.floating_project.project_number
    peering            = module.peering_project.project_number
    cmek               = module.env_secrets_project.project_number
  }

}

resource "google_monitoring_monitored_project" "projects_monitored" {
  for_each = local.monitored_projects

  metrics_scope = join("", ["locations/global/metricsScopes/", local.org_monitoring_project_id])
  name          = each.value
}

resource "time_sleep" "wait_for_propagation_adding_monitoring_project" {
  create_duration = "60s"

  depends_on = [google_monitoring_monitored_project.projects_monitored]
}

resource "google_logging_metric" "logging_metric" {
  for_each = local.monitored_projects

  name    = "set-bucket-iam-policy"
  project = each.value
  filter  = "resource.type=gcs_bucket AND protoPayload.methodName=\"storage.setIamPermissions\""

  metric_descriptor {
    metric_kind = "DELTA"
    value_type  = "INT64"
  }

  depends_on = [time_sleep.wait_for_propagation_adding_monitoring_project]
}
