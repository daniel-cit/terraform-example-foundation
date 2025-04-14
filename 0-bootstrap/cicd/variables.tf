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
    type             = optional(string, "CLOUDBUILD_CSR")
    tfc_org_name     = optional(string, null)
    cicd_runner_repo = optional(string, null)
    repo_owner       = optional(string, null)
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

variable "pat_secret" {
  description = "value"
  type = object({
    github         = optional(string, null)
    gitlab         = optional(string, null)
    gitlab_read    = optional(string, null)
    gitlab_webhook = optional(string, null)
    tfe            = optional(string, null)
    tfe_vcs        = optional(string, null)
  })
  default = {}
}
