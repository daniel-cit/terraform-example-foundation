
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

variable "bootstrap_folder" {
  description = "The numerical ID of the Bootstrap folder"
  type        = string
}

variable "terraform_env_sa" {
  description = "The map of Terraform Service Accounts"
  type        = map(map(string))
}

variable "gcs_bucket_tfstate" {
  description = "Bucket used for storing terraform state for Foundations Pipelines in the Seed Project."
  type        = string
}

variable "projects_gcs_bucket_tfstate" {
  description = "Bucket used for storing terraform state for Foundations Pipelines for the step projects in the Seed Project."
  type        = string
}

variable "org_id" {
  description = "GCP Organization ID"
  type        = string
}

variable "billing_account" {
  description = "The ID of the billing account to associate projects with."
  type        = string
}

variable "default_region" {
  description = "Default region to create resources where applicable."
  type        = string
  default     = "us-central1"
}

variable "project_prefix" {
  description = "Name prefix to use for projects created. Should be the same in all steps. Max size is 3 characters."
  type        = string
  default     = "prj"
}

variable "bucket_prefix" {
  description = "Name prefix to use for state bucket created."
  type        = string
  default     = "bkt"
}

variable "bucket_force_destroy" {
  description = "When deleting a bucket, this boolean option will delete all contained objects. If false, Terraform will fail to delete buckets which contain objects."
  type        = bool
  default     = false
}

variable "project_deletion_policy" {
  description = "The deletion policy for the project created."
  type        = string
  default     = "PREVENT"
}

variable "workflow_deletion_protection" {
  description = "Whether Terraform will be prevented from destroying the workflow. When the field is set to true or unset in Terraform state, a `terraform apply` or `terraform destroy` that would delete the workflow will fail. When the field is set to false, deleting the workflow is allowed."
  type        = bool
  default     = true
}

variable "group_org_admins" {
  description = "Google Group for GCP Organization Administrators"
  type        = string
}

variable "cloudbuildv2_repository_config" {
  description = <<-EOT
  Configuration for integrating repositories with Cloud Build v2:
    - repo_type: Specifies the type of repository. Supported types are 'GITHUBv2', 'GITLABv2', and 'CSR'.
    - repositories: A map of repositories to be created. The key must match the exact name of the repository. Each repository is defined by:
        - repository_name: The name of the repository.
        - repository_url: The URL of the repository.
    - github_secret_id: (Optional) The secret ID for GitHub credentials.
    - github_app_id_secret_id: (Optional) The secret ID for the application ID for a GitHub App used for authentication. For app installation, follow this link: https://github.com/apps/google-cloud-build
    - gitlab_read_authorizer_credential_secret_id: (Optional) The secret ID for the GitLab read authorizer credential.
    - gitlab_authorizer_credential_secret_id: (Optional) The secret ID for the GitLab authorizer credential.
    - gitlab_webhook_secret_id: (Optional) The secret ID for the GitLab WebHook.
    - gitlab_enterprise_host_uri: (Optional) The URI of the GitLab Enterprise host this connection is for. If not specified, the default value is https://gitlab.com.
    - gitlab_enterprise_service_directory: (Optional) Configuration for using Service Directory to privately connect to a GitLab Enterprise server. This should only be set if the GitLab Enterprise server is hosted on-premises and not reachable by public internet. If this field is left empty, calls to the GitLab Enterprise server will be made over the public internet. Format: projects/{project}/locations/{location}/namespaces/{namespace}/services/{service}.
    - gitlab_enterprise_ca_certificate: (Optional) SSL certificate to use for requests to GitLab Enterprise.

  Note: When using GITLABv2, specify `gitlab_read_authorizer_credential_secret_id` and `gitlab_authorizer_credential_secret_id`.
  Note: When using GITHUBv2, specify `github_secret_id` and `github_app_id_secret_id`.
  Note: If 'cloudbuildv2' is not configured, CSR (Cloud Source Repositories) will be used by default.
  EOT
  type = object({
    repo_type = string # Supported values are: GITHUBv2, GITLABv2 and CSR
    # repositories to be created
    repositories = object({
      bootstrap = object({
        repository_name = optional(string, "gcp-bootstrap")
        repository_url  = string
      }),
      org = object({
        repository_name = optional(string, "gcp-org")
        repository_url  = string
      }),
      env = object({
        repository_name = optional(string, "gcp-environments")
        repository_url  = string
      }),
      net = object({
        repository_name = optional(string, "gcp-networks")
        repository_url  = string
      }),
      proj = object({
        repository_name = optional(string, "gcp-projects")
        repository_url  = string
      }),
      tf_cloud_builder = object({
        repository_name = optional(string, "tf-cloud-builder")
        repository_url  = string
      }),
    })
    # Credential Config for each repository type
    github_secret_id                            = optional(string)
    github_app_id_secret_id                     = optional(string)
    gitlab_read_authorizer_credential_secret_id = optional(string)
    gitlab_authorizer_credential_secret_id      = optional(string)
    gitlab_webhook_secret_id                    = optional(string)
    gitlab_enterprise_host_uri                  = optional(string)
    gitlab_enterprise_service_directory         = optional(string)
    gitlab_enterprise_ca_certificate            = optional(string)
  })

  # If cloudbuildv2 is not configured, then auto-creation with CSR will be used
  default = {
    repo_type = "CSR"
    repositories = {
      bootstrap = {
        repository_url = ""
      },
      env = {
        repository_url = ""
      }
      net = {
        repository_url = ""
      }
      org = {
        repository_url = ""
      }
      proj = {
        repository_url = ""
      }
      tf_cloud_builder = {
        repository_url = ""
      }
    }
  }
  validation {
    condition = (
      var.cloudbuildv2_repository_config.repo_type == "GITHUBv2" ? (
        var.cloudbuildv2_repository_config.github_secret_id != null &&
        var.cloudbuildv2_repository_config.github_app_id_secret_id != null &&
        var.cloudbuildv2_repository_config.gitlab_read_authorizer_credential_secret_id == null &&
        var.cloudbuildv2_repository_config.gitlab_authorizer_credential_secret_id == null
        ) : var.cloudbuildv2_repository_config.repo_type == "GITLABv2" ? (
        var.cloudbuildv2_repository_config.github_secret_id == null &&
        var.cloudbuildv2_repository_config.github_app_id_secret_id == null &&
        var.cloudbuildv2_repository_config.gitlab_read_authorizer_credential_secret_id != null &&
        var.cloudbuildv2_repository_config.gitlab_authorizer_credential_secret_id != null
      ) : var.cloudbuildv2_repository_config.repo_type == "CSR" ? true : false
    )
    error_message = "You must specify a valid repo_type ('GITHUBv2', 'GITLABv2', or 'CSR'). For 'GITHUBv2', all 'github_' prefixed variables must be defined and no 'gitlab_' prefixed variables should be defined. For 'GITLABv2', all 'gitlab_' prefixed variables must be defined and no 'github_' prefixed variables should be defined."
  }

}
