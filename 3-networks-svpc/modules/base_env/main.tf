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

locals {
  bgp_asn_number = var.enable_partner_interconnect ? "16550" : "64514"
  spoke_group    = "default"
}

data "google_compute_network" "vpc_dns_hub" {
  count = var.environment_code != "p" ? 1 : 0

  name    = "vpc-p-svpc"
  project = local.dns_project_id
}

/******************************************
 Shared VPC
*****************************************/
module "shared_vpc" {
  source = "git::https://github.com/daniel-cit/terraform-google-network.git//modules/foundation/network?ref=ncc-and-peering-changes"

  project_id                 = local.shared_vpc_project_id
  vpc_name                   = "${local.environment_code}-svpc"
  shared_vpc_host            = true
  resource_code              = local.environment_code
  private_service_cidr       = var.private_service_cidr
  private_service_connect_ip = var.private_service_connect_ip
  firewall_enable_logging    = true
  windows_activation_enabled = false


  ncc_hub_config = {
    create_hub   = true
    name         = "ncc-hub-${local.env}"
    description  = "NCC Hub for ${local.env}"
    hub_labels   = { environment = local.env }
    spoke_labels = { type = "hub_vpc" }

    preset_topology              = "MESH"
    spoke_group                  = "default"
    auto_accept_projects_center  = null
    auto_accept_projects_edge    = null
    auto_accept_projects_default = [local.shared_vpc_project_id]
  }

  nat_config = {
    enabled = false
    bgp_asn = local.bgp_asn_number
    regions = [
      {
        name          = local.default_region1
        num_addresses = 2
      },
      {
        name          = local.default_region2
        num_addresses = 2
      }
    ]
  }

  dns_config = {
    type                         = var.environment_code == "p" ? "hub" : "spoke"
    enable_logging               = true
    enable_inbound_forwarding    = var.environment_code == "p"
    onprem_forwarding            = true
    domain                       = var.domain
    target_name_server_addresses = var.target_name_server_addresses
    dns_hub_project_id           = var.environment_code == "p" ? null : local.dns_project_id
    dns_hub_network_name         = var.environment_code == "p" ? null : data.google_compute_network.vpc_dns_hub[0].self_link
  }

  subnets = [
    {
      subnet_name                      = "sb-${var.environment_code}-svpc-${var.default_region1}"
      subnet_ip                        = var.subnet_primary_ranges[var.default_region1]
      subnet_region                    = var.default_region1
      subnet_private_access            = "true"
      subnet_flow_logs                 = true
      subnet_flow_logs_interval        = var.vpc_flow_logs.aggregation_interval
      subnet_flow_logs_sampling        = var.vpc_flow_logs.flow_sampling
      subnet_flow_logs_metadata        = var.vpc_flow_logs.metadata
      subnet_flow_logs_metadata_fields = var.vpc_flow_logs.metadata_fields
      subnet_flow_logs_filter          = var.vpc_flow_logs.filter_expr
      description                      = "First ${var.env} subnet example."
    },
    {
      subnet_name                      = "sb-${var.environment_code}-svpc-${var.default_region2}"
      subnet_ip                        = var.subnet_primary_ranges[var.default_region2]
      subnet_region                    = var.default_region2
      subnet_private_access            = "true"
      subnet_flow_logs                 = true
      subnet_flow_logs_interval        = var.vpc_flow_logs.aggregation_interval
      subnet_flow_logs_sampling        = var.vpc_flow_logs.flow_sampling
      subnet_flow_logs_metadata        = var.vpc_flow_logs.metadata
      subnet_flow_logs_metadata_fields = var.vpc_flow_logs.metadata_fields
      subnet_flow_logs_filter          = var.vpc_flow_logs.filter_expr
      description                      = "Second ${var.env} subnet example."
    },
    {
      subnet_name      = "sb-${var.environment_code}-svpc-${var.default_region1}-proxy"
      subnet_ip        = var.subnet_proxy_ranges[var.default_region1]
      subnet_region    = var.default_region1
      subnet_flow_logs = false
      description      = "First ${var.env} proxy-only subnet example."
      role             = "ACTIVE"
      purpose          = "REGIONAL_MANAGED_PROXY"
    },
    {
      subnet_name      = "sb-${var.environment_code}-svpc-${var.default_region2}-proxy"
      subnet_ip        = var.subnet_proxy_ranges[var.default_region2]
      subnet_region    = var.default_region2
      subnet_flow_logs = false
      description      = "Second ${var.env} proxy-only subnet example."
      role             = "ACTIVE"
      purpose          = "REGIONAL_MANAGED_PROXY"
    }
  ]
  secondary_ranges = {
    "sb-${var.environment_code}-svpc-${var.default_region1}" = var.subnet_secondary_ranges[var.default_region1]
  }
}
