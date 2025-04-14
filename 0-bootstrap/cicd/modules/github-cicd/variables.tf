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

variable "project_id" {
  description = "value"
  type        = string
}

variable "terraform_env_sa" {
  description = "The map of Terraform Service Accounts"
  type        = map(map(string))
}

variable "gcs_bucket_tfstate" {
  description = "Bucket to store Terraform Remote State Data."
  type        = string
}

variable "cicd_config" {
  description = "value"
  type = object({
    type             = optional(string, "CLOUDBUILD_CSR")
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
  description = "A Secret with A fine-grained personal access token for the user or organization. See https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token#creating-a-fine-grained-personal-access-token"
  type        = string
  default     = null
}
