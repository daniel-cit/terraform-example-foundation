output "csr_repos" {
  description = "List of Cloud Source Repos created by the module, linked to Cloud Build triggers."
  value = { for k, v in module.tf_source.csr_repos : k => {
    "id"      = v.id,
    "name"    = v.name,
    "project" = v.project,
    "url"     = v.url,
    }
  }
}


output "projects_gcs_bucket_tfstate" {
  description = "Bucket used for storing terraform state for stage 4-projects foundations pipelines in seed project."
  value       = module.gcp_projects_state_bucket.bucket.name
}
