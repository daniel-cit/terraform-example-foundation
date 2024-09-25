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
  repo_names = ["bu1-example-app"]
}

module "main" {
  source = "../../modules/base_env"

  environment         = "common"
  repo_names          = local.repo_names
  business_code       = "bu1"
  billing_code        = "1234"
  primary_contact     = "example@example.com"
  secondary_contact   = "example2@example.com"
  application_name    = "app-infra-pipelines"
  remote_state_bucket = var.remote_state_bucket

  activate_apis = [
    "cloudbuild.googleapis.com",
    "sourcerepo.googleapis.com",
    "cloudkms.googleapis.com",
    "iam.googleapis.com",
    "artifactregistry.googleapis.com",
    "cloudresourcemanager.googleapis.com"
  ]
}
