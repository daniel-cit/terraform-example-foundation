/**
 * Copyright 2022 Google LLC
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
  gh_config = {
    "seed" = var.repos.seed,
    "cicd" = var.repos.cicd,
    "org"  = var.repos.organization,
    "env"  = var.repos.environments,
    "net"  = var.repos.networks,
    "proj" = var.repos.projects,
  }

  sa_mapping = {
    for k, v in local.gh_config : k => {
      sa_name   = var.terraform_env_sa[k].name
      attribute = "attribute.repository/${var.repos_owner}/${v}"
    }
  }

  common_secrets = {
    "PROJECT_ID" : var.project_id,
    "WIF_PROVIDER_NAME" : module.gh_oidc.provider_name,
    "TF_BACKEND" : var.gcs_bucket_tfstate,
    "TF_VAR_token" : var.token,
  }

  secrets_list = flatten([
    for k, v in local.gh_config : [
      for secret, plaintext in local.common_secrets : {
        config          = k
        secret_name     = secret
        plaintext_value = plaintext
        repository      = v
      }
    ]
  ])

  sa_secrets = [for k, v in local.gh_config : {
    config          = k
    secret_name     = "SERVICE_ACCOUNT_EMAIL"
    plaintext_value = var.terraform_env_sa[k]["mail"]
    repository      = v
    }
  ]

  gh_secrets = { for v in concat(local.sa_secrets, local.secrets_list) : "${v.config}.${v.secret_name}" => v }

}

module "gh_oidc" {
  source  = "terraform-google-modules/github-actions-runners/google//modules/gh-oidc"
  version = "~> 4.0"

  project_id  = var.project_id
  pool_id     = "foundation-pool"
  provider_id = "foundation-gh-provider"
  sa_mapping  = local.sa_mapping
}

resource "github_actions_secret" "secrets" {
  for_each = local.gh_secrets

  repository      = each.value.repository
  secret_name     = each.value.secret_name
  plaintext_value = each.value.plaintext_value
}

resource "google_service_account_iam_member" "self_impersonate" {
  for_each = var.terraform_env_sa

  service_account_id = each.value["id"]
  role               = "roles/iam.serviceAccountTokenCreator"
  member             = "serviceAccount:${each.value["email"]}"
}
