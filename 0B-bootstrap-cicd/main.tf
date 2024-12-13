

locals {
  // Roles required to manage resources in the CI/CD project
  granular_sa_cicd_project = {
    "bootstrap" = [
      "roles/storage.admin",
      "roles/compute.networkAdmin",
      "roles/cloudbuild.builds.editor",
      "roles/cloudbuild.workerPoolOwner",
      "roles/artifactregistry.admin",
      "roles/source.admin",
      "roles/iam.serviceAccountAdmin",
      "roles/workflows.admin",
      "roles/cloudscheduler.admin",
      "roles/resourcemanager.projectDeleter",
      "roles/dns.admin",
      "roles/iam.workloadIdentityPoolAdmin",
    ],
  }

  cicd_project_id = var.cicd_config.type == "CLOUDBUILD_CSR" ? module.cb_csr[0].cloudbuild_project_id : module.cicd_project[0].project_id

}

module "cicd_project" {
  source  = "terraform-google-modules/project-factory/google"
  version = "~> 17.0"

  count   = var.cicd_config.type == "CLOUDBUILD_CSR" ? 0 : 1

  name              = "${local.project_prefix}-b-cicd"
  random_project_id = true
  org_id            = local.org_id
  folder_id         = local.bootstrap_folder
  billing_account   = local.billing_account
  activate_apis = [
    "compute.googleapis.com",
    "admin.googleapis.com",
    "iam.googleapis.com",
    "billingbudgets.googleapis.com",
    "cloudbilling.googleapis.com",
    "serviceusage.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "iamcredentials.googleapis.com",
    "sts.googleapis.com",
    "dns.googleapis.com",
    "secretmanager.googleapis.com",
    "container.googleapis.com",
    "gkeconnect.googleapis.com",
    "gkehub.googleapis.com",
    "connectgateway.googleapis.com"
  ]

  deletion_policy = var.project_deletion_policy
}

module "cb_csr" {
  source = "./modules/cb-csr-cicd"
  count  = var.cicd_config.type == "CLOUDBUILD_CSR" ? 1 : 0

  org_id                      = local.org_id
  billing_account             = local.billing_account
  terraform_env_sa            = local.terraform_env_sa
  bootstrap_folder            = local.bootstrap_folder
  default_region              = local.default_region
  gcs_bucket_tfstate          = local.gcs_bucket_tfstate
  projects_gcs_bucket_tfstate = local.projects_gcs_bucket_tfstate
  group_org_admins            = local.required_groups.group_org_admins
  bucket_prefix               = local.bucket_prefix
  project_prefix              = local.project_prefix
  project_deletion_policy     = var.project_deletion_policy
  bucket_force_destroy        = var.bucket_force_destroy
}

module "github_actions_cicd" {
  source = "./modules/github-cicd"
  count  = var.cicd_config.type == "GITHUB_ACTIONS" ? 1 : 0

  gcs_bucket_tfstate = local.gcs_bucket_tfstate
  project_id         = local.cicd_project_id
  terraform_env_sa   = local.terraform_env_sa
  repos_owner        = var.cicd_config.repo_owner
  repos              = var.cicd_config.repositories
  token              = var.token.github

  depends_on = [module.cicd_project]
}

# module "tfe_cicd" {
#   source = "./modules/tfe-cicd"
#   count  = var.cicd_config.type == "TFE" ? 1 : 0

#   terraform_env_sa   = local.terraform_env_sa
#   project_id         = local.cicd_project_id
#   tfc_org_name       = ""
#   tfc_token          = var.token.tfe
#   vcs_oauth_token_id = var.token.tfe_vcs
#   vcs_repos_owner    = var.cicd_config.repo_owner
#   vcs_repos          = var.cicd_config.repositories

#   depends_on = [module.cicd_project]
# }

# module "gitlab_cicd" {
#   source = "./modules/gitlab-cicd"
#   count  = var.cicd_config.type == "GITLAB_CI" ? 1 : 0

#   gcs_bucket_tfstate = local.gcs_bucket_tfstate
#   project_id         = local.cicd_project_id
#   terraform_env_sa   = local.terraform_env_sa
#   repos_owner        = var.cicd_config.repo_owner
#   cicd_runner_repo   = var.cicd_config.cicd_runner_repo
#   repos              = var.cicd_config.repositories
#   token              = var.token.gitlab

#   depends_on = [module.cicd_project]
# }

module "cicd_project_iam_member" {
  source   = "./modules/parent-iam-member"
  for_each = local.granular_sa_cicd_project

  member      = "serviceAccount:${local.terraform_env_sa[each.key].email}"
  parent_type = "project"
  parent_id   = local.cicd_project_id
  roles       = each.value
}

// When the CI/CD project is created, the Compute Engine
// default service account is disabled but it still has the Editor
// role associated with the service account. This default SA is the
// only member with the editor role.
// This module will remove all editors from the CI/CD project.
module "bootstrap_projects_remove_editor" {
  source = "./modules/parent-iam-remove-role"

  parent_type = "project"
  parent_id   = local.cicd_project_id
  roles       = ["roles/editor"]

  depends_on = [
    module.cicd_project_iam_member
  ]
}
