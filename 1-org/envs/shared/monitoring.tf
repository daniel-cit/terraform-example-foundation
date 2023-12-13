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
  base_projects = {
    seed_project        = data.google_project.seed_project.number
    cicd_project        = data.google_project.cicd_project.number
    audit_logs          = module.org_audit_logs.project_number
    billing_logs        = module.org_billing_logs.project_number
    secrets             = module.org_secrets.project_number
    interconnect        = module.interconnect.project_number
    scc_notifications   = module.scc_notifications.project_number
    dns_hub             = module.dns_hub.project_number
    dev_base            = module.base_restricted_environment_network["development"].base_shared_vpc_project_number
    dev_restricted      = module.base_restricted_environment_network["development"].restricted_shared_vpc_project_number
    non_prod_base       = module.base_restricted_environment_network["non-production"].base_shared_vpc_project_number
    non_prod_restricted = module.base_restricted_environment_network["non-production"].restricted_shared_vpc_project_number
    prod_base           = module.base_restricted_environment_network["production"].base_shared_vpc_project_number
    prod_restricted     = module.base_restricted_environment_network["production"].restricted_shared_vpc_project_number
  }

  monitored_projects = merge(local.base_projects,
    var.enable_hub_and_spoke ?
    {
      base_network_hub       = module.base_network_hub[0].project_number
      restricted_network_hub = module.restricted_network_hub[0].project_number
    }
    :
    {}
  )

  monitored_projects_list = join(", ", [for k, v in local.monitored_projects : "\"projects/${v}\""])
}

resource "random_string" "ending" {
  length  = 4
  special = false
  upper   = false
}

data "google_project" "seed_project" {
  project_id = local.seed_project_id
}

data "google_project" "cicd_project" {
  project_id = local.cicd_project
}

resource "google_monitoring_monitored_project" "projects_monitored" {
  for_each = local.monitored_projects

  metrics_scope = join("", ["locations/global/metricsScopes/", module.org_monitoring.project_id])
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

resource "time_sleep" "wait_for_propagation" {
  create_duration  = "180s"
  destroy_duration = "180s"

  depends_on = [google_logging_metric.logging_metric]
}

//create topic
module "pubsub_monitoring" {
  source  = "terraform-google-modules/pubsub/google"
  version = "~> 6.0"

  topic      = "top-iam-monitoring-${random_string.ending.result}-event"
  project_id = module.org_monitoring.project_id
}

resource "google_monitoring_notification_channel" "pubsub_channel" {
  project      = module.org_monitoring.project_id
  display_name = ""
  type         = "pubsub"

  labels = {
    topic = module.pubsub_monitoring.topic
  }
}

resource "time_sleep" "wait_for_propagation_notification_channel" {
  destroy_duration = "180s"

  depends_on = [google_logging_metric.logging_metric]
}

resource "google_monitoring_alert_policy" "alert_policy" {
  display_name = "GCS Bucket - Set IAM policy Alert"
  combiner     = "AND"
  project      = module.org_monitoring.project_id

  conditions {
    display_name = "Set Bucket IAM policy custom metric"

    condition_threshold {
      filter          = "metric.type=\"logging.googleapis.com/user/set-bucket-iam-policy\" AND resource.type=\"gcs_bucket\""
      duration        = "60s"
      comparison      = "COMPARISON_GT"
      threshold_value = 0

      aggregations {
        alignment_period   = "60s"
        per_series_aligner = "ALIGN_SUM"
      }

      trigger {
        count = 0
      }
    }
  }

  notification_channels = [google_monitoring_notification_channel.pubsub_channel.name]

  depends_on = [
    time_sleep.wait_for_propagation,
    time_sleep.wait_for_propagation_notification_channel,
    google_monitoring_notification_channel.pubsub_channel
  ]
}

resource "google_monitoring_dashboard" "dashboard" {
  project = module.org_monitoring.project_id
  dashboard_json = templatefile("${path.module}/templates/dashboard.json.tftpl",
    {
      alert_policy_name  = google_monitoring_alert_policy.alert_policy.name
      monitored_projects = local.monitored_projects_list
  })
}
