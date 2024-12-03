/**
 * Copyright 2023 Google LLC
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



/* ----------------------------------------
    Specific to github_bootstrap
   ---------------------------------------- */

# Un-comment github_bootstrap and its outputs if you want to use GitHub Actions instead of Cloud Build
# variable "gh_repos" {
#   description = <<EOT
#   Configuration for the GitHub Repositories to be used to deploy the Terraform Example Foundation stages.
#   owner: The owner of the repositories. An user or an organization.
#   bootstrap: The repository to host the code of the bootstrap stage.
#   organization: The repository to host the code of the organization stage.
#   environments: The repository to host the code of the environments stage.
#   networks: The repository to host the code of the networks stage.
#   projects: The repository to host the code of the projects stage.
#   EOT
#   type = object({
#     owner        = string,
#     bootstrap    = string,
#     organization = string,
#     environments = string,
#     networks     = string,
#     projects     = string,
#   })
# }

# variable "gh_token" {
#   description = "A fine-grained personal access token for the user or organization. See https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token#creating-a-fine-grained-personal-access-token"
#   type        = string
#   sensitive   = true
# }

/* ----------------------------------------
    Specific to jenkins_bootstrap module
   ---------------------------------------- */

# # Un-comment the jenkins_bootstrap module and its outputs if you want to use Jenkins instead of Cloud Build
# variable "jenkins_agent_gce_subnetwork_cidr_range" {
#   description = "The subnetwork to which the Jenkins Agent will be connected to (in CIDR range 0.0.0.0/0)"
#   type        = string
# }

# variable "jenkins_agent_gce_private_ip_address" {
#   description = "The private IP Address of the Jenkins Agent. This IP Address must be in the CIDR range of `jenkins_agent_gce_subnetwork_cidr_range` and be reachable through the VPN that exists between on-prem (Jenkins Controller) and GCP (CICD Project, where the Jenkins Agent is located)."
#   type        = string
# }

# variable "jenkins_agent_gce_ssh_pub_key" {
#   description = "SSH public key needed by the Jenkins Agent GCE Instance. The Jenkins Controller holds the SSH private key. The correct format is `'ssh-rsa [KEY_VALUE] [USERNAME]'`"
#   type        = string
# }

# variable "jenkins_agent_sa_email" {
#   description = "Email for Jenkins Agent service account."
#   type        = string
#   default     = "jenkins-agent-gce"
# }

# variable "jenkins_controller_subnetwork_cidr_range" {
#   description = "A list of CIDR IP ranges of the Jenkins Controller in the form ['0.0.0.0/0']. Usually only one IP in the form '0.0.0.0/32'. Needed to create a FW rule that allows communication with the Jenkins Agent GCE Instance."
#   type        = list(string)
# }

# variable "nat_bgp_asn" {
#   type        = number
#   description = "BGP ASN for NAT cloud route. This is needed to allow the Jenkins Agent to download packages and updates from the internet without having an external IP address."
# }

# variable "vpn_shared_secret" {
#   description = "The shared secret used in the VPN"
#   type        = string
# }

# variable "on_prem_vpn_public_ip_address" {
#   description = "The public IP Address of the Jenkins Controller."
#   type        = string
# }

# variable "on_prem_vpn_public_ip_address2" {
#   description = "The second public IP Address of the Jenkins Controller."
#   type        = string
# }

# variable "router_asn" {
#   type        = number
#   description = "BGP ASN for cloud routes."
#   default     = "64515"
# }

# variable "bgp_peer_asn" {
#   type        = number
#   description = "BGP ASN for cloud routes."
# }

# variable "tunnel0_bgp_peer_address" {
#   type        = string
#   description = "BGP session address for tunnel 0"
# }

# variable "tunnel0_bgp_session_range" {
#   type        = string
#   description = "BGP session range for tunnel 0"
# }

# variable "tunnel1_bgp_peer_address" {
#   type        = string
#   description = "BGP session address for tunnel 1"
# }

# variable "tunnel1_bgp_session_range" {
#   type        = string
#   description = "BGP session range for tunnel 1"
# }

/* ----------------------------------------
    Specific to gitlab_bootstrap
   ---------------------------------------- */

# Un-comment gitlab_bootstrap and its outputs if you want to use GitLab Pipelines instead of Cloud Build
# variable "gl_repos" {
#   description = <<EOT
#   Configuration for the GitLab Repositories to be used to deploy the Terraform Example Foundation stages.
#   owner: The owner of the repositories. An user or a group.
#   bootstrap: The repository to host the code of the bootstrap stage.
#   organization: The repository to host the code of the organization stage.
#   environments: The repository to host the code of the environments stage.
#   networks: The repository to host the code of the networks stage.
#   projects: The repository to host the code of the projects stage.
#   cicd_runner: The repository to host the code of docker image used for CI/CD.
#   EOT
#   type = object({
#     owner        = string,
#     bootstrap    = string,
#     organization = string,
#     environments = string,
#     networks     = string,
#     projects     = string,
#     cicd_runner  = string,
#   })
# }

# variable "gitlab_token" {
#   description = <<EOT
#   A GitLab personal access token or group access token.
#   See:
#       https://docs.gitlab.com/ee/user/group/settings/group_access_tokens.html
#       https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html
#   EOT
#   type        = string
#   sensitive   = true
# }

/* ----------------------------------------
    Specific to tfc_bootstrap
   ---------------------------------------- */

# Un-comment tfc_bootstrap and its outputs if you want to use Terraform Cloud instead of Cloud Build
# variable "vcs_repos" {
#   description = <<EOT
#   Configuration for the Terraform Cloud VCS Repositories to be used to deploy the Terraform Example Foundation stages.
#   owner: The owner of the repositories. An user or an organization.
#   bootstrap: The repository to host the code of the bootstrap stage.
#   organization: The repository to host the code of the organization stage.
#   environments: The repository to host the code of the environments stage.
#   networks: The repository to host the code of the networks stage.
#   projects: The repository to host the code of the projects stage.
#   EOT
#   type = object({
#     owner        = string,
#     bootstrap    = string,
#     organization = string,
#     environments = string,
#     networks     = string,
#     projects     = string,
#   })
# }

# variable "tfc_token" {
#   description = " The token used to authenticate with Terraform Cloud. See https://registry.terraform.io/providers/hashicorp/tfe/latest/docs#authentication"
#   type        = string
#   sensitive   = true
# }

# variable "tfc_org_name" {
#   description = "Name of the TFC organization"
#   type        = string
# }

# variable "tfc_terraform_version" {
#   description = "TF version desired for TFC workspaces"
#   type        = string
# }

# variable "vcs_oauth_token_id" {
#   description = "The VCS Connection OAuth Connection Token ID. This is the ID of the connection between TFC and VCS. See https://developer.hashicorp.com/terraform/cloud-docs/vcs#supported-vcs-providers"
#   type        = string
#   sensitive   = true
# }

# variable "tfc_agent_pool_name" {
#   type        = string
#   description = "Terraform Cloud agent pool name to be created"
#   default     = "tfc-agent-gke-simple-pool"
# }

# variable "tfc_agent_pool_token_description" {
#   type        = string
#   description = "Terraform Cloud agent pool token description"
#   default     = "tfc-agent-gke-simple-pool-token"
# }

# variable "enable_tfc_cloud_agents" {
#   type = bool
#   description = "If false TFC will provide remote runners to run the jobs. If true, TFC will use Agents on a private autopilot GKE cluster."
#   default = false
# }
