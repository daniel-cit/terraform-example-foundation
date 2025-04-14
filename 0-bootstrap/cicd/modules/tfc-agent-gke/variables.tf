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

variable "project_id" {
  description = "The Google Cloud Platform project ID to deploy Terraform Cloud agent cluster"
  type        = string
}

variable "region" {
  description = "The GCP region to use when deploying resources"
  type        = string
  default     = "us-central1"
}

variable "zones" {
  description = "The GCP zone to use when deploying resources"
  type        = list(string)
  default     = ["us-central1-a"]
}

variable "nat_bgp_asn" {
  description = "BGP ASN for NAT cloud routes."
  type        = number
  default     = 64514
}

variable "nat_enabled" {
  type    = bool
  default = true
}

variable "nat_num_addresses" {
  type    = number
  default = 2
}

variable "ip_range_pods_name" {
  description = "The secondary IP range to use for pods"
  type        = string
  default     = "ip-range-pods"
}

variable "ip_range_services_name" {
  description = "The secondary IP range to use for services"
  type        = string
  default     = "ip-range-scv"
}

variable "ip_range_pods_cidr" {
  description = "The secondary IP range CIDR to use for pods"
  type        = string
  default     = "192.168.0.0/18"
}

variable "ip_range_services_cider" {
  description = "The secondary IP range CIDR to use for services"
  type        = string
  default     = "192.168.64.0/18"
}

variable "network_name" {
  description = "Name for the VPC network"
  type        = string
  default     = "tfc-agent-network"
}

variable "subnet_ip" {
  description = "IP range for the subnet"
  type        = string
  default     = "10.0.0.0/17"
}

variable "subnet_name" {
  description = "Name for the subnet"
  type        = string
  default     = "tfc-agent-subnet"
}

variable "network_project_id" {
  description = <<-EOF
    The project ID of the shared VPCs host (for shared vpc support).
    If not provided, the project_id is used
  EOF
  type        = string
  default     = ""
}

variable "create_service_account" {
  description = "Set to true to create a new service account, false to use an existing one"
  type        = bool
  default     = true
}

variable "service_account_email" {
  description = "Optional Service Account for the GKE nodes, required if create_service_account is set to false"
  type        = string
  default     = ""
}

variable "service_account_id" {
  description = "Optional Service Account for the GKE nodes, required if create_service_account is set to false"
  type        = string
  default     = ""
}

variable "tfc_agent_k8s_secrets" {
  description = "Name for the k8s secret required to configure TFC agent on GKE"
  type        = string
  default     = "tfc-agent-k8s-secrets"
}

variable "tfc_agent_address" {
  description = "The HTTP or HTTPS address of the Terraform Cloud/Enterprise API"
  type        = string
  default     = "https://app.terraform.io"
}

variable "tfc_agent_single" {
  description = <<-EOF
    Enable single mode. This causes the agent to handle at most one job and
    immediately exit thereafter. Useful for running agents as ephemeral
    containers, VMs, or other isolated contexts with a higher-level scheduler
    or process supervisor.
  EOF
  type        = bool
  default     = false
}

variable "tfc_agent_auto_update" {
  description = "Controls automatic core updates behavior. Acceptable values include disabled, patch, and minor"
  type        = string
  default     = "minor"
}

variable "tfc_agent_name_prefix" {
  description = "This name may be used in the Terraform Cloud user interface to help easily identify the agent"
  type        = string
  default     = "tfc-agent-k8s"
}

variable "tfc_agent_image" {
  description = "The Terraform Cloud agent image to use"
  type        = string
  default     = "hashicorp/tfc-agent:latest"
}

variable "tfc_agent_memory_request" {
  description = "Memory request for the Terraform Cloud agent container"
  type        = string
  default     = "2Gi"
}

variable "tfc_agent_cpu_request" {
  description = "CPU request for the Terraform Cloud agent container"
  type        = string
  default     = "2"
}

variable "tfc_agent_ephemeral_storage" {
  description = "A temporary storage for a container that gets wiped out and lost when the container is stopped or restarted"
  type        = string
  default     = "1Gi"
}

variable "tfc_agent_token" {
  description = "Terraform Cloud agent token. (Organization Settings >> Agents)"
  type        = string
  sensitive   = true
}

variable "tfc_agent_min_replicas" {
  description = "Minimum replicas for the Terraform Cloud agent pod autoscaler"
  type        = string
  default     = "1"
}

variable "firewall_enable_logging" {
  description = "Enable firewall logging"
  type        = bool
  default     = true
}

variable "private_service_connect_ip" {
  default = "10.10.64.5"
  type    = string
}
