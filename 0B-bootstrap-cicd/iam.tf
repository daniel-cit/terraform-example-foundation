
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
}

module "cicd_project_iam_member" {
  source   = "./modules/parent-iam-member"
  for_each = local.granular_sa_cicd_project

  member      = "serviceAccount:${local.terraform_env_sa[each.key].email}"
  parent_type = "project"
  parent_id   = local.cicd_project_id
  roles       = each.value
}
