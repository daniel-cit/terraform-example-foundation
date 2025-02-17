output "csr_repos" {
  description = "List of Cloud Source Repos created by the module, linked to Cloud Build triggers."
  value       = var.cicd_config.type == "CLOUDBUILD_CSR" ? module.cb_csr[0].csr_repos : []
}


output "cicd_project_id" {
  description = "Project where the CI/CD infrastructure for GitHub Action resides."
  value       = local.cicd_project_id
}
