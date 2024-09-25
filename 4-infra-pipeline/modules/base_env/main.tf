/**
 * Copyright 2024 Google LLC
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
  env_code = element(split("", var.environment), 0)
}

module "app_infra_cloudbuild_project" {
  source  = "terraform-google-modules/project-factory/google"
  version = "~> 16.0"
  count   = local.enable_cloudbuild_deploy ? 1 : 0

  random_project_id        = true
  random_project_id_length = 4
  name                     = "${local.project_prefix}-${local.env_code}-${var.business_code}-infra-pipeline"
  org_id                   = local.org_id
  billing_account          = local.billing_account
  folder_id                = local.common_folder_name

  activate_apis = distinct(concat(var.activate_apis, ["billingbudgets.googleapis.com"]))

  labels = {
    environment       = var.environment
    application_name  = "app-infra-pipelines"
    billing_code      = var.billing_code
    primary_contact   = element(split("@", var.primary_contact), 0)
    secondary_contact = element(split("@", var.secondary_contact), 0)
    business_code     = var.business_code
    env_code          = local.env_code
    vpc               = "none"
  }

  budget_alert_pubsub_topic   = var.project_budget.alert_pubsub_topic
  budget_alert_spent_percents = var.project_budget.alert_spent_percents
  budget_amount               = var.project_budget.budget_amount
  budget_alert_spend_basis    = var.project_budget.alert_spend_basis
}


module "infra_pipelines" {
  source = "../infra_pipelines"
  count  = local.enable_cloudbuild_deploy ? 1 : 0

  org_id                      = local.org_id
  cloudbuild_project_id       = module.app_infra_cloudbuild_project[0].project_id
  cloud_builder_artifact_repo = local.cloud_builder_artifact_repo
  remote_tfstate_bucket       = local.projects_remote_bucket_tfstate
  billing_account             = local.billing_account
  default_region              = var.default_region
  app_infra_repos             = var.repo_names
  private_worker_pool_id      = local.cloud_build_private_worker_pool_id
}

/**
 * When a CI/CD other than Cloud Build is used for deployment this resource
 * is created so that terraform validation works.
 * Without this resource, this module creates zero resources
 * and it breaks terraform validation throwing the following error:
 * ERROR: [Terraform plan json does not contain resource_changes key]
 */
resource "null_resource" "no_cloud_build_cicd" {
  count = !local.enable_cloudbuild_deploy ? 1 : 0
}
