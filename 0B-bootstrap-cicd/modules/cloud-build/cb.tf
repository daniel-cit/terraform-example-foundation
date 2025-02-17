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
  // terraform version image configuration
  terraform_version = "1.5.7"
  // The version of the terraform docker image to be used in the workspace builds
  docker_tag_version_terraform = "v1"

  bucket_self_link_prefix             = "https://www.googleapis.com/storage/v1/b/"
  default_state_bucket_self_link      = "${local.bucket_self_link_prefix}${var.gcs_bucket_tfstate}"
  gcp_projects_state_bucket_self_link = "${local.bucket_self_link_prefix}${var.projects_gcs_bucket_tfstate}"

  cb_config = {
    "seed" = {
      source       = "gcp-seed",
      state_bucket = local.default_state_bucket_self_link,
    },
    "cicd" = {
      source       = "gcp-cicd",
      state_bucket = local.default_state_bucket_self_link,
    },
    "org" = {
      source       = "gcp-org",
      state_bucket = local.default_state_bucket_self_link,
    },
    "env" = {
      source       = "gcp-environments",
      state_bucket = local.default_state_bucket_self_link,
    },
    "net" = {
      source       = "gcp-networks",
      state_bucket = local.default_state_bucket_self_link,
    },
    "proj" = {
      source       = "gcp-projects",
      state_bucket = local.gcp_projects_state_bucket_self_link,
    },
  }

  cloud_source_repos = [for v in local.cb_config : v.source]
  cloudbuilder_repo  = "tf-cloudbuilder"
  base_cloud_source_repos = [
    "gcp-policies",
    local.cloudbuilder_repo,
  ]
  gar_repository           = split("/", module.tf_cloud_builder.artifact_repo)[length(split("/", module.tf_cloud_builder.artifact_repo)) - 1]
  cloud_builder_trigger_id = element(split("/", module.tf_cloud_builder.cloudbuild_trigger_id), index(split("/", module.tf_cloud_builder.cloudbuild_trigger_id), "triggers") + 1, )

  # If user is not bringing its own repositories, will create CSR
  create_cloud_source_repos = var.cloudbuildv2_repository_config.repo_type == "CSR" ? distinct(concat(local.base_cloud_source_repos, local.cloud_source_repos)) : []
  # Used for backwards compatibility on conditional permission assignment
  cloud_source_repos_granular_sa = var.cloudbuildv2_repository_config.repo_type == "CSR" ? var.terraform_env_sa : {}
  create_cloudbuildv2_connection = var.cloudbuildv2_repository_config.repo_type != "CSR"

  use_csr              = var.cloudbuildv2_repository_config.repo_type == "CSR"
  is_github_connection = var.cloudbuildv2_repository_config.repo_type == "GITHUBv2"
  is_gitlab_connection = var.cloudbuildv2_repository_config.repo_type == "GITLABv2"
  bootstrap_csr_repo_id = length(local.create_cloud_source_repos) == 0 ? "" : (
    contains(keys(module.tf_source.csr_repos), local.cloudbuilder_repo) ? split("/", module.tf_source.csr_repos[local.cloudbuilder_repo].id)[3] : ""
  )

  cicd_project_apis = [
    "serviceusage.googleapis.com",
    "servicenetworking.googleapis.com",
    "compute.googleapis.com",
    "logging.googleapis.com",
    "iam.googleapis.com",
    "admin.googleapis.com",
    "workflows.googleapis.com",
    "artifactregistry.googleapis.com",
    "cloudbuild.googleapis.com",
    "cloudscheduler.googleapis.com",
    "bigquery.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "cloudbilling.googleapis.com",
    "appengine.googleapis.com",
    "storage-api.googleapis.com",
    "billingbudgets.googleapis.com",
    "dns.googleapis.com",
    "secretmanager.googleapis.com"
  ]

  activate_apis = local.use_csr ? distinct(concat(local.cicd_project_apis, ["sourcerepo.googleapis.com"])) : local.cicd_project_apis
}

resource "random_string" "suffix" {
  length  = 4
  special = false
  upper   = false
}

module "tf_source" {
  source  = "terraform-google-modules/bootstrap/google//modules/tf_cloudbuild_source"
  version = "~> 11.0"

  org_id                  = var.org_id
  folder_id               = var.bootstrap_folder
  project_id              = "${var.project_prefix}-b-cicd-${random_string.suffix.result}"
  location                = var.default_region
  billing_account         = var.billing_account
  group_org_admins        = var.group_org_admins
  buckets_force_destroy   = var.bucket_force_destroy
  project_deletion_policy = var.project_deletion_policy
  activate_apis           = local.activate_apis

  cloud_source_repos = local.create_cloud_source_repos

  project_labels = {
    environment       = "bootstrap"
    application_name  = "cloudbuild-bootstrap"
    billing_code      = "1234"
    primary_contact   = "example1"
    secondary_contact = "example2"
    business_code     = "shared"
    env_code          = "b"
    vpc               = "none"
  }
}

resource "google_project_service_identity" "workflows_identity" {
  provider = google-beta

  project = module.tf_source.cloudbuild_project_id
  service = "workflows.googleapis.com"

  depends_on = [module.tf_source]
}

module "tf_private_pool" {
  source = "../cb-private-pool"

  project_id = module.tf_source.cloudbuild_project_id

  private_worker_pool = {
    region                   = var.default_region,
    enable_network_peering   = true,
    create_peered_network    = true,
    peered_network_subnet_ip = "10.3.0.0/24"
    peering_address          = "192.168.0.0"
    peering_prefix_length    = 24
  }

  vpn_configuration = {
    enable_vpn = false
  }
}

module "tf_cloud_builder" {
  source  = "terraform-google-modules/bootstrap/google//modules/tf_cloudbuild_builder"
  version = "~> 11.0"

  project_id                   = module.tf_source.cloudbuild_project_id
  dockerfile_repo_uri          = local.use_csr ? module.tf_source.csr_repos[local.cloudbuilder_repo].url : module.git_repo_connection[0].cloud_build_repositories_2nd_gen_repositories["tf_cloud_builder"].id
  use_cloudbuildv2_repository  = !local.use_csr
  dockerfile_repo_type         = local.is_github_connection ? "GITHUB" : (local.is_gitlab_connection ? "UNKNOWN" : "CLOUD_SOURCE_REPOSITORIES")
  gar_repo_location            = var.default_region
  workflow_region              = var.default_region
  terraform_version            = local.terraform_version
  build_timeout                = "1200s"
  cb_logs_bucket_force_destroy = var.bucket_force_destroy
  trigger_location             = var.default_region
  enable_worker_pool           = true
  worker_pool_id               = module.tf_private_pool.private_worker_pool_id
  bucket_name                  = "${var.bucket_prefix}-${module.tf_source.cloudbuild_project_id}-tf-cloudbuilder-build-logs"
  workflow_deletion_protection = var.workflow_deletion_protection
}

module "bootstrap_csr_repo" {
  source  = "terraform-google-modules/gcloud/google"
  version = "~> 3.1"
  count   = local.use_csr ? 1 : 0

  upgrade = false

  create_cmd_entrypoint = "${path.module}/scripts/push-to-repo.sh"
  create_cmd_body       = "${module.tf_source.cloudbuild_project_id} ${local.bootstrap_csr_repo_id} ${path.module}/Dockerfile"
}

data "google_secret_manager_secret_version_access" "github_token" {
  count = local.is_github_connection ? 1 : 0

  secret = var.cloudbuildv2_repository_config.github_secret_id
}

module "bootstrap_github_repo" {
  source  = "terraform-google-modules/gcloud/google"
  version = "~> 3.1"
  count   = local.is_github_connection ? 1 : 0

  upgrade = false

  create_cmd_entrypoint = "${path.module}/scripts/github-bootstrap-builder-repo.sh"
  create_cmd_body       = "${data.google_secret_manager_secret_version_access.github_token[0].secret_data} ${var.cloudbuildv2_repository_config.repositories.tf_cloud_builder.repository_url} ${path.module}/Dockerfile"
}

data "google_secret_manager_secret_version_access" "gitlab_token" {
  count = local.is_gitlab_connection ? 1 : 0

  secret = var.cloudbuildv2_repository_config.gitlab_authorizer_credential_secret_id
}

module "bootstrap_gitlab_repo" {
  source  = "terraform-google-modules/gcloud/google"
  version = "~> 3.1"
  count   = local.is_gitlab_connection ? 1 : 0

  upgrade = false

  create_cmd_entrypoint = "${path.module}/scripts/gitlab-bootstrap-builder-repo.sh"
  create_cmd_body       = "${data.google_secret_manager_secret_version_access.gitlab_token[0].secret_data} ${var.cloudbuildv2_repository_config.repositories.tf_cloud_builder.repository_url} ${path.module}/Dockerfile"
}

resource "time_sleep" "cloud_builder" {
  create_duration = "30s"

  depends_on = [
    module.tf_cloud_builder,
    module.bootstrap_csr_repo,
    module.bootstrap_github_repo,
    module.bootstrap_gitlab_repo
  ]
}

module "build_terraform_image" {
  source  = "terraform-google-modules/gcloud/google"
  version = "~> 3.1"
  upgrade = false

  create_cmd_triggers = {
    "terraform_version" = local.terraform_version
  }

  create_cmd_body = "beta builds triggers run  ${local.cloud_builder_trigger_id} --branch main --region ${var.default_region} --project ${module.tf_source.cloudbuild_project_id}"

  module_depends_on = [
    time_sleep.cloud_builder,
  ]
}

module "git_repo_connection" {
  source  = "terraform-google-modules/bootstrap/google//modules/cloudbuild_repo_connection"
  version = "~> 11.0"
  count   = local.use_csr ? 0 : 1

  project_id = module.tf_source.cloudbuild_project_id

  connection_config = {
    connection_type = var.cloudbuildv2_repository_config.repo_type
    gitlab_authorizer_credential_secret_id      = var.cloudbuildv2_repository_config.gitlab_authorizer_credential_secret_id
    gitlab_read_authorizer_credential_secret_id = var.cloudbuildv2_repository_config.gitlab_read_authorizer_credential_secret_id
    gitlab_webhook_secret_id                    = var.cloudbuildv2_repository_config.gitlab_webhook_secret_id
    github_secret_id                            = var.cloudbuildv2_repository_config.github_secret_id
    github_app_id_secret_id                     = var.cloudbuildv2_repository_config.github_app_id_secret_id
    gitlab_enterprise_host_uri                  = var.cloudbuildv2_repository_config.gitlab_enterprise_host_uri
    gitlab_enterprise_service_directory         = var.cloudbuildv2_repository_config.gitlab_enterprise_service_directory
    gitlab_enterprise_ca_certificate            = var.cloudbuildv2_repository_config.gitlab_enterprise_ca_certificate
  }

  # TODO is the tf-builder part of this list or not ?
  cloud_build_repositories = var.cloudbuildv2_repository_config.repositories

}

module "tf_workspace" {
  source   = "terraform-google-modules/bootstrap/google//modules/tf_cloudbuild_workspace"
  version  = "~> 11.0"
  for_each = var.terraform_env_sa

  project_id                = module.tf_source.cloudbuild_project_id
  location                  = var.default_region
  trigger_location          = var.default_region
  enable_worker_pool        = true
  worker_pool_id            = module.tf_private_pool.private_worker_pool_id
  state_bucket_self_link    = local.cb_config[each.key].state_bucket
  log_bucket_name           = "${var.bucket_prefix}-${module.tf_source.cloudbuild_project_id}-${local.cb_config[each.key].source}-build-logs"
  artifacts_bucket_name     = "${var.bucket_prefix}-${module.tf_source.cloudbuild_project_id}-${local.cb_config[each.key].source}-build-artifacts"
  cloudbuild_plan_filename  = "cloudbuild-tf-plan.yaml"
  cloudbuild_apply_filename = "cloudbuild-tf-apply.yaml"
  tf_repo_uri               = local.use_csr ? module.tf_source.csr_repos[local.cb_config[each.key].source].url : module.git_repo_connection[0].cloud_build_repositories_2nd_gen_repositories[each.key].id
  tf_repo_type              = local.use_csr ? "CLOUD_SOURCE_REPOSITORIES" : "CLOUDBUILD_V2_REPOSITORY"
  cloudbuild_sa             = var.terraform_env_sa[each.key].id
  create_cloudbuild_sa      = false
  diff_sa_project           = true
  create_state_bucket       = false
  buckets_force_destroy     = var.bucket_force_destroy

  substitutions = {
    "_ORG_ID"                       = var.org_id
    "_BILLING_ID"                   = var.billing_account
    "_GAR_REGION"                   = var.default_region
    "_GAR_PROJECT_ID"               = module.tf_source.cloudbuild_project_id
    "_GAR_REPOSITORY"               = local.gar_repository
    "_DOCKER_TAG_VERSION_TERRAFORM" = local.docker_tag_version_terraform
  }

  tf_apply_branches = ["development", "nonproduction", "production"]

  depends_on = [
    module.tf_source,
    module.tf_cloud_builder,
  ]

}

resource "google_artifact_registry_repository_iam_member" "terraform_sa_artifact_registry_reader" {
  for_each = var.terraform_env_sa

  project    = module.tf_source.cloudbuild_project_id
  location   = var.default_region
  repository = local.gar_repository
  role       = "roles/artifactregistry.reader"
  member     = "serviceAccount:${var.terraform_env_sa[each.key].email}"
}

resource "google_sourcerepo_repository_iam_member" "member" {
  for_each = local.cloud_source_repos_granular_sa

  project    = module.tf_source.cloudbuild_project_id
  repository = module.tf_source.csr_repos["gcp-policies"].name
  role       = "roles/viewer"
  member     = "serviceAccount:${var.terraform_env_sa[each.key].email}"
}
