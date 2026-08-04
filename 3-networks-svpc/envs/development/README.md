# 3-networks-svpc/development

The purpose of this step is to set up shared VPCs with default DNS, NAT (optional), Private Service networking, VPC service controls, onprem Dedicated Interconnect, onprem VPN and baseline firewall rules for environment development.

## Prerequisites

1. 0-bootstrap executed successfully.
1. 1-org executed successfully.
1. 2-environments/envs/development executed successfully.
1. 3-networks-svpc/envs/shared executed successfully.
1. Obtain the value for the access_context_manager_policy_id variable. Can be obtained by running `gcloud access-context-manager policies list --organization YOUR_ORGANIZATION_ID --format="value(name)"`.

<!-- BEGINNING OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| access\_context\_manager\_policy\_id | The id of the default Access Context Manager policy created in step `1-org`. Can be obtained by running `gcloud access-context-manager policies list --organization YOUR_ORGANIZATION_ID --format="value(name)"`. | `number` | n/a | yes |
| allow\_additional\_member\_types | Allows use of additional member types: `group`, `principal`, and `principalSet` as members of the perimeter. If true the members will be added in a ingress rules instead of in the access level. | `bool` | `false` | no |
| domain | The DNS name of peering managed zone, for instance 'example.com.'. Must end with a period. | `string` | n/a | yes |
| egress\_policies\_dry\_run\_map | Map of additional egress policies for the dry-run perimeter. Map key is the Terraform state key. | <pre>map(object({<br>    title = optional(string, null)<br>    from = object({<br>      sources = optional(object({<br>        resources     = optional(list(string), [])<br>        access_levels = optional(list(string), [])<br>      }), {}),<br>      identity_type = optional(string, null)<br>      identities    = optional(list(string), null)<br>    })<br>    to = object({<br>      operations = optional(map(object({<br>        methods     = optional(list(string), [])<br>        permissions = optional(list(string), [])<br>      })), {}),<br>      roles              = optional(list(string), null)<br>      resources          = optional(list(string), ["*"])<br>      external_resources = optional(list(string), [])<br>    })<br>  }))</pre> | `{}` | no |
| egress\_policies\_map | Map of additional egress policies for the enforced perimeter. Map key is the Terraform state key. | <pre>map(object({<br>    title = optional(string, null)<br>    from = object({<br>      sources = optional(object({<br>        resources     = optional(list(string), [])<br>        access_levels = optional(list(string), [])<br>      }), {}),<br>      identity_type = optional(string, null)<br>      identities    = optional(list(string), null)<br>    })<br>    to = object({<br>      operations = optional(map(object({<br>        methods     = optional(list(string), [])<br>        permissions = optional(list(string), [])<br>      })), {}),<br>      roles              = optional(list(string), null)<br>      resources          = optional(list(string), ["*"])<br>      external_resources = optional(list(string), [])<br>    })<br>  }))</pre> | `{}` | no |
| ingress\_policies\_dry\_run\_map | Map of additional ingress policies for the dry-run perimeter. Map key is the Terraform state key. | <pre>map(object({<br>    title = optional(string, null)<br>    from = object({<br>      sources = optional(object({<br>        resources     = optional(list(string), [])<br>        access_levels = optional(list(string), [])<br>      }), {}),<br>      identity_type = optional(string, null)<br>      identities    = optional(list(string), null)<br>    })<br>    to = object({<br>      operations = optional(map(object({<br>        methods     = optional(list(string), [])<br>        permissions = optional(list(string), [])<br>      })), {}),<br>      roles     = optional(list(string), null)<br>      resources = optional(list(string), ["*"])<br>    })<br>  }))</pre> | `{}` | no |
| ingress\_policies\_map | Map of additional ingress policies for the enforced perimeter. Map key is the Terraform state key. | <pre>map(object({<br>    title = optional(string, null)<br>    from = object({<br>      sources = optional(object({<br>        resources     = optional(list(string), [])<br>        access_levels = optional(list(string), [])<br>      }), {}),<br>      identity_type = optional(string, null)<br>      identities    = optional(list(string), null)<br>    })<br>    to = object({<br>      operations = optional(map(object({<br>        methods     = optional(list(string), [])<br>        permissions = optional(list(string), [])<br>      })), {}),<br>      roles     = optional(list(string), null)<br>      resources = optional(list(string), ["*"])<br>    })<br>  }))</pre> | `{}` | no |
| perimeter\_additional\_members | The list of additional members to be added to the enforced perimeter access level members list. To be able to see the resources protected by the VPC Service Controls in the restricted perimeter, add your user in this list. Entries must be in the standard GCP form: `user:email@example.com` or `serviceAccount:my-service-account@example.com`. | `list(string)` | `[]` | no |
| perimeter\_additional\_members\_dry\_run | The list of additional members to be added to the dry-run perimeter access level members list. To be able to see the resources protected by the VPC Service Controls in the restricted perimeter, add your user in this list. Entries must be in the standard GCP form: `user:email@example.com` or `serviceAccount:my-service-account@example.com`. | `list(string)` | `[]` | no |
| remote\_state\_bucket | Backend bucket to load Terraform Remote State Data from previous steps. | `string` | n/a | yes |
| tfc\_org\_name | Name of the TFC organization | `string` | `""` | no |
| universe\_domain | The universe domain to use for Google Cloud APIs. This defines the API endpoint boundary for your deployment. The default is 'googleapis.com' for the standard public Google Cloud. Modify this value if you are deploying to isolated environments like Google Distributed Cloud (GDC), Trusted Partner Cloud (TPC), or other sovereign cloud environments. | `string` | `"googleapis.com"` | no |

## Outputs

| Name | Description |
|------|-------------|
| access\_context\_manager\_policy\_id | Access Context Manager Policy ID. |
| access\_level\_name | Access context manager access level name |
| access\_level\_name\_dry\_run | Access context manager access level name for the dry-run perimeter |
| enforce\_vpcsc | Enable the enforced mode for VPC Service Controls. It is not recommended to enable VPC-SC on the first run deploying your foundation. Review [best practices for enabling VPC Service Controls](https://cloud.google.com/vpc-service-controls/docs/enable), then only enforce the perimeter after you have analyzed the access patterns in your dry-run perimeter and created the necessary exceptions for your use cases. |
| network\_name | The name of the VPC being created |
| network\_self\_link | The URI of the VPC being created |
| service\_perimeter\_name | Access context manager service perimeter name |
| shared\_vpc\_host\_project\_id | The shared vpc host project ID |
| subnets\_ips | The IPs and CIDRs of the subnets being created |
| subnets\_names | The names of the subnets being created |
| subnets\_secondary\_ranges | The secondary ranges associated with these subnets |
| subnets\_self\_links | The self-links of subnets being created |

<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
