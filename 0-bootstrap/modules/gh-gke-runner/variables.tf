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
  description = "The project id to deploy Github Runner"
  type        = string
}

variable "cluster_name" {
  default = "gke-b-gh-runner"
  type    = string
}

variable "environment_code" {
  description = "A short form of the folder level resources (environment) within the Google Cloud organization (ex. d)."
  type        = string
}

variable "region" {
  description = "The GCP region to deploy instances into"
  type        = string
}

variable "private_service_connect_ip" {
  description = "Internal IP to be used as the private service connect endpoint."
  type        = string
}

variable "subnet_ip_cidr" {
  description = "IP range for the subnet."
  type        = string
}

variable "master_ipv4_cidr_block" {
  description = "(Beta) The IP range in CIDR notation to use for the hosted master network"
  type        = string
  default     = "10.0.0.0/28"
}

variable "network_tag" {
  description = "Network tag to apply to instances to allow external access."
  type        = string
  default     = "gh-runner-vm"
}

variable "ip_range_pods_cidr" {
  description = "The secondary ip range cidr to use for pods"
  type        = string
}

variable "ip_range_services_cider" {
  description = "The secondary ip range cidr to use for services"
  type        = string
}

variable "firewall_enable_logging" {
  description = "Toggle firewall logging for VPC Firewalls."
  type        = bool
  default     = true
}

variable "nat_enabled" {
  description = "Toggle creation of NAT cloud router."
  type        = bool
  default     = false
}

variable "nat_bgp_asn" {
  description = "BGP ASN for NAT cloud routes."
  type        = number
  default     = 64514
}

variable "nat_num_addresses" {
  description = "Number of external IPs to reserve for Cloud NAT."
  type        = number
  default     = 2
}
