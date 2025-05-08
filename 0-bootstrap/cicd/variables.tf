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

variable "remote_state_bucket" {
  description = "Backend bucket to load Terraform Remote State Data from previous steps."
  type        = string
}

variable "project_deletion_policy" {
  description = "The deletion policy for the project created."
  type        = string
  default     = "PREVENT"
}

variable "bucket_force_destroy" {
  description = "When deleting a bucket, this boolean option will delete all contained objects. If false, Terraform will fail to delete buckets which contain objects."
  type        = bool
  default     = false
}

variable "cicd_config" {
  description = "value"
  type = object({
    repo_type                                   = optional(string, "CLOUDBUILD_CSR")
    cicd_runner_repo                            = optional(string)
    repo_owner                                  = optional(string)
    github_secret_id                            = optional(string)
    github_app_id_secret_id                     = optional(string)
    gitlab_read_authorizer_credential_secret_id = optional(string)
    gitlab_authorizer_credential_secret_id      = optional(string)
    gitlab_webhook_secret_id                    = optional(string)
    gitlab_enterprise_host_uri                  = optional(string)
    gitlab_enterprise_service_directory         = optional(string)
    gitlab_enterprise_ca_certificate            = optional(string)
    secret_project_id                           = optional(string)
    repositories = optional(object({
      seed = object({
        repository_name = optional(string, "gcp-seed")
        repository_url  = string
      }),
      cicd = object({
        repository_name = optional(string, "gcp-cicd")
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
    }), null)
  })
  default = {}
}

variable "cloudbuildv2_repository_config" {
  description = <<-EOT
  Configuration for integrating repositories with Cloud Build v2:
    - repo_type: Specifies the type of repository. Supported types are 'GITHUBv2', 'GITLABv2', and 'CSR'.
    - repositories: A map of repositories to be created. The key must match the exact name of the repository. Each repository is defined by:
        - repository_name: The name of the repository.
        - repository_url: The URL of the repository.
    - github_secret_id: (Optional) The personal access token for GitHub authentication.
    - github_app_id_secret_id: (Optional) The application ID for a GitHub App used for authentication.
    - gitlab_read_authorizer_credential_secret_id: (Optional) The read authorizer credential for GitLab access.
    - gitlab_authorizer_credential_secret_id: (Optional) The authorizer credential for GitLab access.
    - gitlab_webhook_secret_id: (Optional) The secret ID for the GitLab WebHook.
    - gitlab_enterprise_host_uri: (Optional) The URI of the GitLab Enterprise host this connection is for. If not specified, the default value is https://gitlab.com.
    - gitlab_enterprise_service_directory: (Optional) Configuration for using Service Directory to privately connect to a GitLab Enterprise server. This should only be set if the GitLab Enterprise server is hosted on-premises and not reachable by public internet. If this field is left empty, calls to the GitLab Enterprise server will be made over the public internet. Format: projects/{project}/locations/{location}/namespaces/{namespace}/services/{service}.
    - gitlab_enterprise_ca_certificate: (Optional) SSL certificate to use for requests to GitLab Enterprise.
    - secret_project_id: (Optional) The project id where the secret is stored.
  Note: When using GITLABv2, specify `gitlab_read_authorizer_credential` and `gitlab_authorizer_credential` and `gitlab_webhook_secret_id`.
  Note: When using GITHUBv2, specify `github_pat` and `github_app_id`.
  Note: If 'cloudbuildv2_repository_config' variable is not configured, CSR (Cloud Source Repositories) will be used by default.
  EOT
  type = object({
    repo_type = string # Supported values are: GITHUBv2, GITLABv2 and CSR
    # repositories to be created
    repositories = object({
      multitenant = object({
        repository_name = optional(string, "eab-multitenant")
        repository_url  = string
      }),
      applicationfactory = object({
        repository_name = optional(string, "eab-applicationfactory")
        repository_url  = string
      }),
      fleetscope = object({
        repository_name = optional(string, "eab-fleetscope")
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
    secret_project_id                           = optional(string)
  })

  # If cloudbuildv2 is not configured, then auto-creation with CSR will be used
  default = {
    repo_type = "CSR"
    repositories = {
      multitenant = {
        repository_url = ""
      },
      fleetscope = {
        repository_url = ""
      }
      applicationfactory = {
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
        var.cloudbuildv2_repository_config.gitlab_authorizer_credential_secret_id == null &&
        var.cloudbuildv2_repository_config.gitlab_webhook_secret_id == null &&
        var.cloudbuildv2_repository_config.secret_project_id != null
        ) : var.cloudbuildv2_repository_config.repo_type == "GITLABv2" ? (
        var.cloudbuildv2_repository_config.github_secret_id == null &&
        var.cloudbuildv2_repository_config.github_app_id_secret_id == null &&
        var.cloudbuildv2_repository_config.gitlab_read_authorizer_credential_secret_id != null &&
        var.cloudbuildv2_repository_config.gitlab_authorizer_credential_secret_id != null &&
        var.cloudbuildv2_repository_config.gitlab_webhook_secret_id != null &&
        var.cloudbuildv2_repository_config.secret_project_id != null
      ) : var.cloudbuildv2_repository_config.repo_type == "CSR" ? true : false
    )
    error_message = "You must specify a valid repo_type ('GITHUBv2', 'GITLABv2', or 'CSR'). For 'GITHUBv2', all 'github_' prefixed variables must be defined and no 'gitlab_' prefixed variables should be defined. For 'GITLABv2', all 'gitlab_' prefixed variables must be defined and no 'github_' prefixed variables should be defined. Provide the project_id where secrets are hosted on 'secret_project_id'."
  }

}

# variable "cicd_config" {
#   description = "value"
#   type = object({
#     repo_type             = optional(string, "CLOUDBUILD_CSR")
#     tfc_org_name     = optional(string, null)
#     cicd_runner_repo = optional(string, null)
#     repo_owner       = optional(string, null)
#     repositories = optional(object({
#       seed = object({
#         repository_name = optional(string, "gcp-seed")
#         repository_url  = string
#       }),
#       cicd = object({
#         repository_name = optional(string, "gcp-cicd")
#         repository_url  = string
#       }),
#       org = object({
#         repository_name = optional(string, "gcp-org")
#         repository_url  = string
#       }),
#       env = object({
#         repository_name = optional(string, "gcp-environments")
#         repository_url  = string
#       }),
#       net = object({
#         repository_name = optional(string, "gcp-networks")
#         repository_url  = string
#       }),
#       proj = object({
#         repository_name = optional(string, "gcp-projects")
#         repository_url  = string
#       }),
#     }), null)
#   })
#   default = {}
# }

# variable "pat_secret" {
#   description = "value"
#   type = object({
#     github         = optional(string, null)
#     gitlab         = optional(string, null)
#     gitlab_read    = optional(string, null)
#     gitlab_webhook = optional(string, null)
#     tfe            = optional(string, null)
#     tfe_vcs        = optional(string, null)
#   })
#   default = {}
# }
