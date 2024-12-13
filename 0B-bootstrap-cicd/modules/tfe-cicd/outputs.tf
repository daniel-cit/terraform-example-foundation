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

output "tfc_org_name" {
  description = "Name of the TFC organization."
  value       = var.tfc_org_name
}

output "kubernetes_endpoint" {
  description = "The GKE cluster endpoint"
  sensitive   = true
  value       = module.tfc_agent_gke[0].kubernetes_endpoint
}

output "service_account" {
  description = "The default service account used for TFC agent nodes"
  value       = module.tfc_agent_gke[0].service_account
}

output "cluster_name" {
  description = "GKE cluster name"
  value       = module.tfc_agent_gke[0].cluster_name
}

output "hub_cluster_membership_id" {
  value       = module.tfc_agent_gke[0].hub_cluster_membership_id
  description = "The ID of the cluster membership"
}
