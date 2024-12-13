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

variable "vcs_repos_owner" {
  description = "The owner of the repositories. An user or an organization."
  type        = string
}

variable "vcs_repos" {
  description = <<EOT
  Configuration for the Terraform Cloud VCS Repositories to be used to deploy the Terraform Example Foundation stages.
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

variable "tfc_token" {
  description = " The token used to authenticate with Terraform Cloud. See https://registry.terraform.io/providers/hashicorp/tfe/latest/docs#authentication"
  type        = string
  sensitive   = true
}

variable "tfc_org_name" {
  description = "Name of the TFC organization"
  type        = string
}

variable "tfc_terraform_version" {
  description = "TF version desired for TFC workspaces"
  type        = string
  default     = "1.5.7"
}

variable "vcs_oauth_token_id" {
  description = "The VCS Connection OAuth Connection Token ID. This is the ID of the connection between TFC and VCS. See https://developer.hashicorp.com/terraform/cloud-docs/vcs#supported-vcs-providers"
  type        = string
  sensitive   = true
}

variable "tfc_agent_pool_name" {
  type        = string
  description = "Terraform Cloud agent pool name to be created"
  default     = "tfc-agent-gke-simple-pool"
}

variable "tfc_agent_pool_token_description" {
  type        = string
  description = "Terraform Cloud agent pool token description"
  default     = "tfc-agent-gke-simple-pool-token"
}

variable "enable_tfc_cloud_agents" {
  type        = bool
  description = "If false TFC will provide remote runners to run the jobs. If true, TFC will use Agents on a private autopilot GKE cluster."
  default     = false
}
