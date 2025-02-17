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
  description = "value"
  type        = map(map(string))
}

variable "gcs_bucket_tfstate" {
  description = "Bucket to store Terraform Remote State Data."
  type        = string
}

variable "repos_owner" {
  description = "The owner of the repositories. An user or a group."
  type        = string
}

variable "cicd_runner_repo" {
  description = "The repository to host the code of docker image used for CI/CD."
  type        = string
}

variable "repos" {
  description = <<EOT
  Configuration for the GitLab Repositories to be used to deploy the Terraform Example Foundation stages.
  bootstrap: The repository to host the code of the bootstrap stage.
  organization: The repository to host the code of the organization stage.
  environments: The repository to host the code of the environments stage.
  networks: The repository to host the code of the networks stage.
  projects: The repository to host the code of the projects stage.
  EOT
  type = object({
    bootstrap    = string,
    organization = string,
    environments = string,
    networks     = string,
    projects     = string,
  })
}

variable "pat_secret" {
  description = <<EOT
  A Secret hosting A GitLab personal access token or group access token.
  See:
      https://docs.gitlab.com/ee/user/group/settings/group_access_tokens.html
      https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html
  EOT
  type        = string
}
