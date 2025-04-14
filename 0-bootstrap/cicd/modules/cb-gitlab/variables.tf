
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

variable "group_org_admins" {
  description = "Google Group for GCP Organization Administrators"
  type        = string
}

variable "pat_secret" {
  description = "value"
  type        = string
}

variable "read_pat_secret" {
  description = "value"
  type        = string
}

variable "webhook_secret" {
  description = "value"
  type        = string
}

variable "workflow_deletion_protection" {
  description = "Whether Terraform will be prevented from destroying the workflow. When the field is set to true or unset in Terraform state, a `terraform apply` or `terraform destroy` that would delete the workflow will fail. When the field is set to false, deleting the workflow is allowed."
  type        = bool
  default     = true
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
