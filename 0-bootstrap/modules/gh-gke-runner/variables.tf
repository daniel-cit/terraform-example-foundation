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
  type        = string
  description = "The project id to deploy Github Runner"
}

variable "cluster_name" {
  type    = string
  default = "gke-b-gh-runner"
}

variable "environment_code" {
  type        = string
  description = "A short form of the folder level resources (environment) within the Google Cloud organization (ex. d)."
}

variable "region" {
  type        = string
  description = "The GCP region to deploy instances into"
}

variable "private_service_connect_ip" {
  type        = string
  description = "Internal IP to be used as the private service connect endpoint."
}

variable "subnet_ip" {
  type        = string
  description = "IP range for the subnet"
}

variable "auth_subnet_ip" {
  type = string
}

variable "master_ipv4_cidr_block" {
  type = string
}

variable "network_tag" {
  type = string
  description = "value"
  default = "gh-runner-vm"
}

variable "ip_range_pods_cidr" {
  type        = string
  description = "The secondary ip range cidr to use for pods"
}

variable "ip_range_services_cider" {
  type        = string
  description = "The secondary ip range cidr to use for services"
}

variable "firewall_enable_logging" {
  type        = bool
  description = "Toggle firewall logging for VPC Firewalls."
  default     = true
}

variable "nat_enabled" {
  type        = bool
  description = "Toggle creation of NAT cloud router."
  default     = false
}

variable "nat_bgp_asn" {
  type        = number
  description = "BGP ASN for NAT cloud routes."
  default     = 64514
}

variable "nat_num_addresses" {
  type        = number
  description = "Number of external IPs to reserve for Cloud NAT."
  default     = 2
}
