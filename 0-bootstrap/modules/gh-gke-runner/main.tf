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

data "google_project" "project" {
  project_id = var.project_id
}

module "runner_cluster" {
  source  = "terraform-google-modules/kubernetes-engine/google//modules/beta-private-cluster/"
  version = "~> 28.0"

  project_id                      = var.project_id
  name                            = var.cluster_name
  regional                        = true
  region                          = var.region
  network                         = module.network.network_name
  subnetwork                      = local.subnet_names[index(module.network.subnets_names, local.subnet_name)]
  ip_range_pods                   = local.ip_range_pods_name
  ip_range_services               = local.ip_range_services_name
  master_ipv4_cidr_block          = var.master_ipv4_cidr_block
  release_channel                 = "REGULAR"
  enable_vertical_pod_autoscaling = true
  enable_private_endpoint         = true
  enable_private_nodes            = true
  create_service_account          = true
  add_cluster_firewall_rules      = false #TODO
  node_pools_tags = {
    all               = [var.network_tag]
    default-node-pool = []
  }

  master_authorized_networks = [
    {
      cidr_block   = var.auth_subnet_ip
      display_name = "VPC"
    },
  ]
}

module "fleet_membership" {
  source  = "terraform-google-modules/kubernetes-engine/google//modules/fleet-membership"
  version = "~> 28.0"

  project_id   = var.project_id
  location     = var.region
  cluster_name = module.runner_cluster.name
}
