/**
 * Copyright 2021 Google LLC
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
    env_monitoring = module.monitoring_project.project_number
    env_secrets    = module.env_secrets.project_number
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

/******************************************
  Projects for monitoring workspaces
*****************************************/

module "monitoring_project" {
  source  = "terraform-google-modules/project-factory/google"
  version = "~> 14.0"

  random_project_id           = true
  random_project_id_length    = 4
  name                        = "${local.project_prefix}-${var.environment_code}-monitoring"
  org_id                      = local.org_id
  billing_account             = local.billing_account
  folder_id                   = google_folder.env.id
  disable_services_on_destroy = false
  depends_on                  = [time_sleep.wait_60_seconds]
  activate_apis = [
    "logging.googleapis.com",
    "monitoring.googleapis.com",
    "billingbudgets.googleapis.com"
  ]

  labels = {
    environment       = var.env
    application_name  = "env-monitoring"
    billing_code      = "1234"
    primary_contact   = "example1"
    secondary_contact = "example2"
    business_code     = "abcd"
    env_code          = var.environment_code
  }
  budget_alert_pubsub_topic   = var.project_budget.monitoring_alert_pubsub_topic
  budget_alert_spent_percents = var.project_budget.monitoring_alert_spent_percents
  budget_amount               = var.project_budget.monitoring_budget_amount
  budget_alert_spend_basis    = var.project_budget.monitoring_budget_alert_spend_basis
}

