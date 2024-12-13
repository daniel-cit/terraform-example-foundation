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

output "jenkins_agent_gce_instance_id" {
  description = "Jenkins Agent GCE Instance id."
  value       = module.jenkins_bootstrap.jenkins_agent_gce_instance_id
}

output "jenkins_agent_vpc_id" {
  description = "Jenkins Agent VPC name."
  value       = module.jenkins_bootstrap.jenkins_agent_vpc_id
}

output "projects_gcs_bucket_tfstate" {
  description = "Bucket used for storing terraform state for stage 4-projects foundations pipelines in seed project."
  value       = module.seed_bootstrap.gcs_bucket_tfstate
}

output "jenkins_agent_sa_email" {
  description = "Email for privileged custom service account for Jenkins Agent GCE instance."
  value       = module.jenkins_bootstrap.jenkins_agent_sa_email
}

output "jenkins_agent_sa_name" {
  description = "Fully qualified name for privileged custom service account for Jenkins Agent GCE instance."
  value       = module.jenkins_bootstrap.jenkins_agent_sa_name
}

output "gcs_bucket_jenkins_artifacts" {
  description = "Bucket used to store Jenkins artifacts in Jenkins project."
  value       = module.jenkins_bootstrap.gcs_bucket_jenkins_artifacts
}
