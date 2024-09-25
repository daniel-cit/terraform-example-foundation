<!-- BEGINNING OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| activate\_apis | The api to activate for the GCP project | `list(string)` | `[]` | no |
| application\_name | The name of application where GCP resources relate | `string` | n/a | yes |
| billing\_code | The code that's used to provide chargeback information | `string` | n/a | yes |
| business\_code | The code that describes which business unit owns the project | `string` | `"shared"` | no |
| default\_region | Default region to create resources where applicable. | `string` | `"us-central1"` | no |
| environment | The environment the single project belongs to | `string` | n/a | yes |
| primary\_contact | The primary email contact for the project | `string` | n/a | yes |
| project\_budget | Budget configuration.<br>  budget\_amount: The amount to use as the budget.<br>  alert\_spent\_percents: A list of percentages of the budget to alert on when threshold is exceeded.<br>  alert\_pubsub\_topic: The name of the Cloud Pub/Sub topic where budget related messages will be published, in the form of `projects/{project_id}/topics/{topic_id}`.<br>  alert\_spend\_basis: The type of basis used to determine if spend has passed the threshold. Possible choices are `CURRENT_SPEND` or `FORECASTED_SPEND` (default). | <pre>object({<br>    budget_amount        = optional(number, 1000)<br>    alert_spent_percents = optional(list(number), [1.2])<br>    alert_pubsub_topic   = optional(string, null)<br>    alert_spend_basis    = optional(string, "FORECASTED_SPEND")<br>  })</pre> | `{}` | no |
| remote\_state\_bucket | Backend bucket to load Terraform Remote State Data from previous steps. | `string` | n/a | yes |
| repo\_names | value | `list(string)` | n/a | yes |
| secondary\_contact | The secondary email contact for the project | `string` | `""` | no |
| tfc\_org\_name | Name of the TFC organization | `string` | `""` | no |

## Outputs

| Name | Description |
|------|-------------|
| apply\_triggers\_id | CB apply triggers |
| artifact\_buckets | GCS Buckets to store Cloud Build Artifacts |
| default\_region | Default region to create resources where applicable. |
| enable\_cloudbuild\_deploy | Enable infra deployment using Cloud Build. |
| log\_buckets | GCS Buckets to store Cloud Build logs |
| plan\_triggers\_id | CB plan triggers |
| project\_id | n/a |
| repos | CSRs to store source code |
| state\_buckets | GCS Buckets to store TF state |
| terraform\_service\_accounts | APP Infra Pipeline Terraform Accounts. |

<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
