/**
 * Copyright 2022 Google LLC
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

  subnet_primary_ranges = {
    (local.default_region1) = "10.8.0.0/18"
    (local.default_region2) = "10.9.0.0/18"
  }
  subnet_proxy_ranges = {
    (local.default_region1) = "10.26.0.0/23"
    (local.default_region2) = "10.27.0.0/23"
  }

  group_configuration = var.enable_hub_and_spoke_transitivity ? {
    topology = "MESH",
    group    = "default"
    center   = null,
    edge     = null,
    default = [
      local.net_hub_project_id,
      local.development_net_project_id,
      local.nonproduction_net_project_id,
      local.production_net_project_id
    ]
    } : {
    topology = "STAR",
    group    = "center"
    center   = [local.net_hub_project_id],
    edge = [
      local.development_net_project_id,
      local.nonproduction_net_project_id,
      local.production_net_project_id
    ],
    default = null
  }
}

/******************************************
  Shared Network VPC
*****************************************/

module "shared_vpc" {
  source = "git::https://github.com/daniel-cit/terraform-google-network.git//modules/foundation/network?ref=ncc-and-peering-changes"

  project_id                 = local.net_hub_project_id
  vpc_name                   = "${local.environment_code}-svpc-hub"
  shared_vpc_host            = true
  resource_code              = local.environment_code
  private_service_connect_ip = "10.17.0.5"
  firewall_enable_logging    = var.hub_firewall_enable_logging
  windows_activation_enabled = var.hub_windows_activation_enabled


  ncc_hub_config = {
    create_hub   = true
    name         = "ncc-hub-${local.env}"
    description  = "NCC Hub for ${local.env}"
    hub_labels   = { environment = local.env }
    spoke_labels = { type = "hub_vpc" }

    preset_topology              = local.group_configuration["topology"]
    spoke_group                  = local.group_configuration["group"]
    auto_accept_projects_center  = local.group_configuration["center"]
    auto_accept_projects_edge    = local.group_configuration["edge"]
    auto_accept_projects_default = local.group_configuration["default"]
  }

  nat_config = {
    enabled = var.hub_nat_enabled
    bgp_asn = local.bgp_asn_number
    regions = [
      {
        name          = local.default_region1
        num_addresses = var.hub_nat_num_addresses_region1
      },
      {
        name          = local.default_region2
        num_addresses = var.hub_nat_num_addresses_region2
      }
    ]
  }

  dns_config = {
    type                         = "hub"
    enable_logging               = var.hub_dns_enable_logging
    enable_inbound_forwarding    = var.hub_dns_enable_inbound_forwarding
    onprem_forwarding            = true
    domain                       = var.domain
    target_name_server_addresses = var.target_name_server_addresses
  }

  subnets = [
    {
      subnet_name                      = "sb-c-svpc-hub-${local.default_region1}"
      subnet_ip                        = local.subnet_primary_ranges[local.default_region1]
      subnet_region                    = local.default_region1
      subnet_private_access            = "true"
      subnet_flow_logs                 = var.vpc_flow_logs.enable_logging
      subnet_flow_logs_interval        = var.vpc_flow_logs.aggregation_interval
      subnet_flow_logs_sampling        = var.vpc_flow_logs.flow_sampling
      subnet_flow_logs_metadata        = var.vpc_flow_logs.metadata
      subnet_flow_logs_metadata_fields = var.vpc_flow_logs.metadata_fields
      subnet_flow_logs_filter          = var.vpc_flow_logs.filter_expr
      description                      = "Network hub subnet for ${local.default_region1}"
    },
    {
      subnet_name                      = "sb-c-svpc-hub-${local.default_region2}"
      subnet_ip                        = local.subnet_primary_ranges[local.default_region2]
      subnet_region                    = local.default_region2
      subnet_private_access            = "true"
      subnet_flow_logs                 = var.vpc_flow_logs.enable_logging
      subnet_flow_logs_interval        = var.vpc_flow_logs.aggregation_interval
      subnet_flow_logs_sampling        = var.vpc_flow_logs.flow_sampling
      subnet_flow_logs_metadata        = var.vpc_flow_logs.metadata
      subnet_flow_logs_metadata_fields = var.vpc_flow_logs.metadata_fields
      subnet_flow_logs_filter          = var.vpc_flow_logs.filter_expr
      description                      = "Network hub subnet for ${local.default_region2}"
    },
    {
      subnet_name      = "sb-c-svpc-hub-${local.default_region1}-proxy"
      subnet_ip        = local.subnet_proxy_ranges[local.default_region1]
      subnet_region    = local.default_region1
      subnet_flow_logs = false
      description      = "Network hub proxy-only subnet for ${local.default_region1}"
      role             = "ACTIVE"
      purpose          = "REGIONAL_MANAGED_PROXY"
    },
    {
      subnet_name      = "sb-c-svpc-hub-${local.default_region2}-proxy"
      subnet_ip        = local.subnet_proxy_ranges[local.default_region2]
      subnet_region    = local.default_region2
      subnet_flow_logs = false
      description      = "Network hub proxy-only subnet for ${local.default_region2}"
      role             = "ACTIVE"
      purpose          = "REGIONAL_MANAGED_PROXY"
    }
  ]
  secondary_ranges = {}
}
