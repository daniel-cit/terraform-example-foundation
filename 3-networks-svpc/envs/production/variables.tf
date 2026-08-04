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

variable "target_name_server_addresses" {
  description = "List of IPv4 address of target name servers for the forwarding zone configuration. See https://cloud.google.com/dns/docs/overview#dns-forwarding-zones for details on target name servers in the context of Cloud DNS forwarding zones."
  type        = list(map(any))
  default     = []
}

variable "remote_state_bucket" {
  description = "Backend bucket to load Terraform Remote State Data from previous steps."
  type        = string
}

variable "universe_domain" {
  description = "The universe domain to use for Google Cloud APIs. This defines the API endpoint boundary for your deployment. The default is 'googleapis.com' for the standard public Google Cloud. Modify this value if you are deploying to isolated environments like Google Distributed Cloud (GDC), Trusted Partner Cloud (TPC), or other sovereign cloud environments."
  type        = string
  default     = "googleapis.com"
}

variable "access_context_manager_policy_id" {
  type        = number
  description = "The id of the default Access Context Manager policy created in step `1-org`. Can be obtained by running `gcloud access-context-manager policies list --organization YOUR_ORGANIZATION_ID --format=\"value(name)\"`."
}

variable "domain" {
  type        = string
  description = "The DNS name of peering managed zone, for instance 'example.com.'. Must end with a period."
}

variable "ingress_policies_dry_run_map" {
  description = "Map of additional ingress policies for the dry-run perimeter. Map key is the Terraform state key."
  type = map(object({
    title = optional(string, null)
    from = object({
      sources = optional(object({
        resources     = optional(list(string), [])
        access_levels = optional(list(string), [])
      }), {}),
      identity_type = optional(string, null)
      identities    = optional(list(string), null)
    })
    to = object({
      operations = optional(map(object({
        methods     = optional(list(string), [])
        permissions = optional(list(string), [])
      })), {}),
      roles     = optional(list(string), null)
      resources = optional(list(string), ["*"])
    })
  }))
  default = {}
}

variable "egress_policies_dry_run_map" {
  description = "Map of additional egress policies for the dry-run perimeter. Map key is the Terraform state key."
  type = map(object({
    title = optional(string, null)
    from = object({
      sources = optional(object({
        resources     = optional(list(string), [])
        access_levels = optional(list(string), [])
      }), {}),
      identity_type = optional(string, null)
      identities    = optional(list(string), null)
    })
    to = object({
      operations = optional(map(object({
        methods     = optional(list(string), [])
        permissions = optional(list(string), [])
      })), {}),
      roles              = optional(list(string), null)
      resources          = optional(list(string), ["*"])
      external_resources = optional(list(string), [])
    })
  }))
  default = {}
}

variable "ingress_policies_map" {
  description = "Map of additional ingress policies for the enforced perimeter. Map key is the Terraform state key."
  type = map(object({
    title = optional(string, null)
    from = object({
      sources = optional(object({
        resources     = optional(list(string), [])
        access_levels = optional(list(string), [])
      }), {}),
      identity_type = optional(string, null)
      identities    = optional(list(string), null)
    })
    to = object({
      operations = optional(map(object({
        methods     = optional(list(string), [])
        permissions = optional(list(string), [])
      })), {}),
      roles     = optional(list(string), null)
      resources = optional(list(string), ["*"])
    })
  }))
  default = {}
}

variable "egress_policies_map" {
  description = "Map of additional egress policies for the enforced perimeter. Map key is the Terraform state key."
  type = map(object({
    title = optional(string, null)
    from = object({
      sources = optional(object({
        resources     = optional(list(string), [])
        access_levels = optional(list(string), [])
      }), {}),
      identity_type = optional(string, null)
      identities    = optional(list(string), null)
    })
    to = object({
      operations = optional(map(object({
        methods     = optional(list(string), [])
        permissions = optional(list(string), [])
      })), {}),
      roles              = optional(list(string), null)
      resources          = optional(list(string), ["*"])
      external_resources = optional(list(string), [])
    })
  }))
  default = {}
}

variable "allow_additional_member_types" {
  description = "Allows use of additional member types: `group`, `principal`, and `principalSet` as members of the perimeter. If true the members will be added in a ingress rules instead of in the access level."
  type        = bool
  default     = false
}

variable "perimeter_additional_members" {
  description = "The list of additional members to be added to the enforced perimeter access level members list. To be able to see the resources protected by the VPC Service Controls in the restricted perimeter, add your user in this list. Entries must be in the standard GCP form: `user:email@example.com` or `serviceAccount:my-service-account@example.com`."
  type        = list(string)
  default     = []
}

variable "perimeter_additional_members_dry_run" {
  description = "The list of additional members to be added to the dry-run perimeter access level members list. To be able to see the resources protected by the VPC Service Controls in the restricted perimeter, add your user in this list. Entries must be in the standard GCP form: `user:email@example.com` or `serviceAccount:my-service-account@example.com`."
  type        = list(string)
  default     = []
}

variable "tfc_org_name" {
  description = "Name of the TFC organization"
  type        = string
  default     = ""
}
