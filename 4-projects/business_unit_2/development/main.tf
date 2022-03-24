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
  terraform_service_account = data.terraform_remote_state.bootstrap.outputs.terraform_service_account
}

module "env" {
  source = "../../modules/env_base"

  env_code      = "d"
  env           = "development"
  business_code = "bu2"
  business_unit = "business_unit_2"
}
