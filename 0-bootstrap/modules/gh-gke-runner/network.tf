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
  vpc_name                = "${var.environment_code}-gh-runner"
  network_name            = "vpc-${local.vpc_name}"
  subnet_name             = "sb-${local.vpc_name}-${var.region}"
  auth_subnetwork         = "sb-${local.vpc_name}-auth-${var.region}"
  ip_range_pods_name      = "rn-${local.vpc_name}-${var.region}-gke-pod"
  ip_range_services_name  = "rn-${local.vpc_name}-${var.region}-gke-svc"
  private_googleapis_cidr = module.private_service_connect.private_service_connect_ip
  subnet_names            = [for subnet_self_link in module.network.subnets_self_links : split("/", subnet_self_link)[length(split("/", subnet_self_link)) - 1]]
}

module "network" {
  source  = "terraform-google-modules/network/google"
  version = "~> 5.2"

  project_id                             = var.project_id
  network_name                           = local.network_name
  delete_default_internet_gateway_routes = true

  routes = [
    {
      name              = "rt-${local.vpc_name}-1000-egress-internet-default"
      description       = "Tag based route through IGW to access internet"
      destination_range = "0.0.0.0/0"
      tags              = var.network_tag
      next_hop_internet = "true"
      priority          = "1000"
    }
  ]

  subnets = [
    {
      description           = "Subnet for GitHub Runner"
      subnet_name           = local.subnet_name
      subnet_ip             = var.subnet_ip
      subnet_region         = var.region
      subnet_private_access = "true"
      subnet_flow_logs      = "true"
    }
  ]

  secondary_ranges = {
    (local.subnet_name) = [
      {
        range_name    = local.ip_range_pods_name
        ip_cidr_range = var.ip_range_pods_cidr
      },
      {
        range_name    = local.ip_range_services_name
        ip_cidr_range = var.ip_range_services_cider
      }
    ]
  }

}

resource "google_dns_policy" "default_policy" {
  project                   = var.project_id
  name                      = "dp-${local.vpc_name}-default-policy"
  enable_inbound_forwarding = true
  enable_logging            = true

  networks {
    network_url = module.network.network_self_link
  }
}

module "private_service_connect" {
  source  = "terraform-google-modules/network/google//modules/private-service-connect"
  version = "~> 5.2"

  project_id                 = var.project_id
  dns_code                   = "dz-${local.vpc_name}"
  network_self_link          = module.network.network_self_link
  private_service_connect_ip = var.private_service_connect_ip
  forwarding_rule_target     = "all-apis"
}

resource "google_compute_firewall" "allow_private_api_egress" {
  name      = "fw-${local.vpc_name}-65430-e-a-allow-google-apis-all-tcp-443"
  network   = module.network.network_name
  project   = var.project_id
  direction = "EGRESS"
  priority  = 65430

  dynamic "log_config" {
    for_each = var.firewall_enable_logging == true ? [{
      metadata = "INCLUDE_ALL_METADATA"
    }] : []

    content {
      metadata = log_config.value.metadata
    }
  }

  allow {
    protocol = "tcp"
    ports    = ["443"]
  }

  destination_ranges = [local.private_googleapis_cidr]

  target_tags = [var.network_tag]
}

/******************************************
  NAT Cloud Router & NAT config
 *****************************************/

resource "google_compute_router" "nat" {
  count = var.nat_enabled ? 1 : 0

  name    = "cr-${local.vpc_name}-${var.region}-nat-router"
  project = var.project_id
  region  = var.region
  network = module.network.network_self_link

  bgp {
    asn = var.nat_bgp_asn
  }
}

resource "google_compute_address" "nat_external_addresses" {
  count = var.nat_enabled ? var.nat_num_addresses : 0

  project = var.project_id
  name    = "ca-${local.vpc_name}-${var.region}-${count.index}"
  region  = var.region
}

resource "google_compute_router_nat" "egress" {
  count = var.nat_enabled ? 1 : 0

  name                               = "rn-${local.vpc_name}-${var.region}-egress"
  project                            = var.project_id
  router                             = google_compute_router.nat[0].name
  region                             = var.region
  nat_ip_allocate_option             = "MANUAL_ONLY"
  nat_ips                            = google_compute_address.nat_external_addresses.*.self_link
  source_subnetwork_ip_ranges_to_nat = "ALL_SUBNETWORKS_ALL_IP_RANGES"

  log_config {
    filter = "TRANSLATIONS_ONLY"
    enable = true
  }
}
