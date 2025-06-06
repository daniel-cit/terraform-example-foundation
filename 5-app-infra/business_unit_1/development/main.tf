/**
 * Copyright 2021 Google LLC
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
  business_unit = "business_unit_1"
  environment   = "development"
}

module "gce_instance" {
  source = "../../modules/env_base"

  environment         = local.environment
  business_unit       = local.business_unit
  project_suffix      = "sample-svpc"
  region              = coalesce(var.instance_region, local.default_region)
  remote_state_bucket = var.remote_state_bucket
}

module "peering_gce_instance" {
  source = "../../modules/env_base"

  environment         = local.environment
  business_unit       = local.business_unit
  project_suffix      = "sample-peering"
  region              = coalesce(var.instance_region, local.default_region)
  remote_state_bucket = var.remote_state_bucket
}
