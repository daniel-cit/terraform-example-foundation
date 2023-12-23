# Cloud Asset Inventory Notification
Uses Google Cloud Asset Inventory to create a feed of IAM Policy change events, then process them to detect when a roles (from a preset list) is given to a member (service account, user or group). Then generates a SCC Finding with the member, role, resource where it was granted and the time that was granted.

## Usage

```hcl
module "secure_log_notification" {
  source = "terraform-google-modules/terraform-example-foundation/google//1-org/modules/log-monitoring"

  org_id         = <ORG ID>
  project_id     = <PROJECT ID>
  region         = <REGION>
  encryption_key = <CMEK KEY>
  labels         = <LABELS>
}
```

<!-- BEGINNING OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| enable\_cmek | The KMS Key to Encrypt Artifact Registry repository, Cloud Storage Bucket and Pub/Sub. | `bool` | `false` | no |
| encryption\_key | The KMS Key to Encrypt Artifact Registry repository, Cloud Storage Bucket and Pub/Sub. | `string` | `null` | no |
| labels | Labels to be assigned to resources. | `map(any)` | `{}` | no |
| location | Default location to create resources where applicable. | `string` | `"us-central1"` | no |
| org\_id | GCP Organization ID | `string` | n/a | yes |
| project\_id | The Project ID where the resources will be created | `string` | n/a | yes |
| pubsub\_topic\_id | The Pub/Sub topic to listen to for alert notifications. | `string` | n/a | yes |
| random\_suffix | Adds a suffix of 4 random characters to the created resources names. | `bool` | `true` | no |

## Outputs

| Name | Description |
|------|-------------|
| artifact\_registry\_name | Artifact Registry Repo to store the Cloud Function image. |
| bucket\_name | Storage bucket where the source code is. |
| function\_uri | URI of the Cloud Function. |
| scc\_source | SCC Findings Source. |

<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->

## Requirements

### Software

The following dependencies must be available:

* [Terraform](https://www.terraform.io/downloads.html) >= 1.3
* [Terraform Provider for GCP](https://github.com/terraform-providers/terraform-provider-google) < 5.0

### APIs

A project with the following APIs enabled must be used to host the resources of this module:

* Project
  * Google Cloud Key Management Service: `cloudkms.googleapis.com`
  * Cloud Resource Manager API: `cloudresourcemanager.googleapis.com`
  * Cloud Functions API: `cloudfunctions.googleapis.com`
  * Cloud Build API: `cloudbuild.googleapis.com`
  * Cloud Asset API`cloudasset.googleapis.com`
  * Clouod Pub/Sub API: `pubsub.googleapis.com`
  * Identity and Access Management (IAM) API: `iam.googleapis.com`
  * Cloud Billing API: `cloudbilling.googleapis.com`
