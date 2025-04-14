# 0-bootstrap

This repo is part of a multi-part guide that shows how to configure and deploy
the example.com reference architecture described in
[Google Cloud security foundations guide](https://cloud.google.com/architecture/security-foundations). The following table lists the stages of this deployment.

<table>
<tbody>
<tr>
<td>0A-bootstrap seed (this file)</td>
<td>Bootstraps a Google Cloud organization, creating all the required resources
and permissions to start using the Cloud Foundation Toolkit (CFT).</td>
</tr>
<tr>
<td><a href="../0B-bootstrap-cicd"><span style="white-space: nowrap;">0B-bootstrap CI/CD</span></a></td>
<td>This step configures a <a href="../docs/GLOSSARY.md#foundation-cicd-pipeline">CI/CD pipeline</a>
for foundations code in subsequent stages.</td>
</tr>
<tr>
<td><a href="../1-org">1-org</a></td>
<td>Sets up top-level shared folders, networking projects, and
organization-level logging, and sets baseline security settings through
organizational policy.</td>
</tr>
<tr>
<td><a href="../2-environments"><span style="white-space: nowrap;">2-environments</span></a></td>
<td>Sets up development, nonproduction, and production environments within the
Google Cloud organization that you've created.</td>
</tr>
<tr>
<td><a href="../3-networks-dual-svpc">3-networks-dual-svpc</a></td>
<td>Sets up base and restricted shared VPCs with default DNS, NAT (optional),
Private Service networking, VPC service controls, on-premises Dedicated
Interconnect, and baseline firewall rules for each environment. It also sets
up the global DNS hub.</td>
</tr>
<tr>
<td><a href="../3-networks-hub-and-spoke">3-networks-hub-and-spoke</a></td>
<td>Sets up base and restricted shared VPCs with all the default configuration
found on step 3-networks-dual-svpc, but here the architecture will be based on the
Hub and Spoke network model. It also sets up the global DNS hub.</td>
</tr>
</tr>
<tr>
<td><a href="../4-projects">4-projects</a></td>
<td>Set up a folder structure, projects, and application infrastructure pipeline for applications,
 which are connected as service projects to the shared VPC created in the previous stage.</td>
</tr>
<tr>
<td><a href="../5-app-infra">5-app-infra</a></td>
<td>Deploy a <a href="https://cloud.google.com/compute/">Compute Engine</a> instance in one of the business unit projects using the infra pipeline setup in 4-projects.</td>
</tr>
</tbody>
</table>

For an overview of the architecture and the parts, see the
[terraform-example-foundation README](https://github.com/terraform-google-modules/terraform-example-foundation)
file.

## Purpose

The purpose of this step is to bootstrap a Google Cloud organization, creating all the required resources and permissions to start using the Cloud Foundation Toolkit (CFT).

## Intended usage and support

This repository is intended as an example to be forked, tweaked, and maintained in the user's own version-control system; the modules within this repository are not intended for use as remote references.
Though this blueprint can help accelerate your foundation design and build, we assume that you have the engineering skills and teams to deploy and customize your own foundation based on your own requirements.

We will support:
 - Code is semantically valid, pinned to known good versions, and passes terraform validate and lint checks
 - All PR to this repo must pass integration tests to deploy all resources into a test environment before being merged
 - Feature requests about ease of use of the code, or feature requests that generally apply to all users, are welcome

We will not support:
 - In-place upgrades from a foundation deployed with an earlier version to a more recent version, even for minor version changes, might not be feasible. Repository maintainers do not have visibility to what resources a user deploys on top of their foundation or how the foundation was customized in deployment, so we make no guarantee about avoiding breaking changes.
 - Feature requests that are specific to a single user's requirement and not representative of general best practices

## Prerequisites

To run the commands described in this document, install the following:

- [Google Cloud SDK](https://cloud.google.com/sdk/install) version 393.0.0 or later
- [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git) version 2.28.0 or later
- [Terraform](https://www.terraform.io/downloads.html) version 1.5.7
- [jq](https://jqlang.github.io/jq/download/) version 1.6.0 or later

**Note:** Make sure that you use the same version of [Terraform](https://www.terraform.io/downloads.html) throughout this series. Otherwise, you might experience Terraform state snapshot lock errors.

Version 1.5.7 is the last version before the license model change. To use a later version of Terraform, ensure that the Terraform version used in the Operational System to manually execute part of the steps in `3-networks` and `4-projects` is the same version configured in the following code

- 0-bootstrap/modules/jenkins-agent/variables.tf
   ```
   default     = "1.5.7"
   ```

- 0-bootstrap/cb.tf
   ```
   terraform_version = "1.5.7"
   ```

- scripts/validate-requirements.sh
   ```
   TF_VERSION="1.5.7"
   ```

- build/github-tf-apply.yaml
   ```
   terraform_version: '1.5.7'
   ```

- github-tf-pull-request.yaml

   ```
   terraform_version: "1.5.7"
   ```

- 0-bootstrap/Dockerfile
   ```
   ARG TERRAFORM_VERSION=1.5.7
   ```

Also make sure that you've done the following:

1. Set up a Google Cloud
   [organization](https://cloud.google.com/resource-manager/docs/creating-managing-organization).
1. Set up a Google Cloud
   [billing account](https://cloud.google.com/billing/docs/how-to/manage-billing-account).
1. Create Cloud Identity or Google Workspace groups as defined in [groups for access control](https://cloud.google.com/architecture/security-foundations/authentication-authorization#groups_for_access_control).
Set the variables in **terraform.tfvars** (`groups` block) to use the specific group names you create.
1. For the user who will run the procedures in this document, grant the following roles:
   - The `roles/resourcemanager.organizationAdmin` role on the Google Cloud organization.
   - The `roles/orgpolicy.policyAdmin` role on the Google Cloud organization.
   - The `roles/resourcemanager.projectCreator` role on the Google Cloud organization.
   - The `roles/billing.admin` role on the billing account.
   - The `roles/resourcemanager.folderCreator` role.
   - The `roles/securitycenter.admin` role.

     ```bash
     # example:
     gcloud organizations add-iam-policy-binding ${ORG_ID}  --member=user:$SUPER_ADMIN_EMAIL --role=roles/securitycenter.admin --quiet > /dev/null 1>&1
     ```
1. Enable the following additional services on your current bootstrap project:
     ```bash
     gcloud services enable cloudresourcemanager.googleapis.com
     gcloud services enable cloudbilling.googleapis.com
     gcloud services enable iam.googleapis.com
     gcloud services enable cloudkms.googleapis.com
     gcloud services enable servicenetworking.googleapis.com
     ```

### Optional - Automatic creation of Google Cloud Identity groups

In the foundation, Google Cloud Identity groups are used for [authentication and access management](https://cloud.google.com/architecture/security-foundations/authentication-authorization) .

To enable automatic creation of the [groups](https://cloud.google.com/architecture/security-foundations/authentication-authorization#groups_for_access_control), complete the following actions:

- Have an existing project for Cloud Identity API billing.
- Enable the Cloud Identity API (`cloudidentity.googleapis.com`) on the billing project.
- Grant role `roles/serviceusage.serviceUsageConsumer` to the user running Terraform on the billing project.
- Change the field `groups.create_required_groups` to **true** to create the required groups.
- Change the field `groups.create_optional_groups` to **true** and fill the `groups.optional_groups` with the emails to be created.

### Troubleshooting

See [troubleshooting](../docs/TROUBLESHOOTING.md) if you run into issues during this step.



## Deploying

1. Clone [terraform-example-foundation](https://github.com/terraform-google-modules/terraform-example-foundation) into your local environment and navigate to the `0A-bootstrap` folder.

   ```bash
   git clone https://github.com/terraform-google-modules/terraform-example-foundation.git

   cd terraform-example-foundation/0A-bootstrap-seed
   ```

1. Rename `terraform.example.tfvars` to `terraform.tfvars` and update the file with values from your environment:

   ```bash
   mv terraform.example.tfvars terraform.tfvars
   ```

1. Use the helper script [validate-requirements.sh](../scripts/validate-requirements.sh) to validate your environment:

   ```bash
   ../scripts/validate-requirements.sh -o <ORGANIZATION_ID> -b <BILLING_ACCOUNT_ID> -u <END_USER_EMAIL>
   ```

   **Note:** The script is not able to validate if the user is in a Cloud Identity or Google Workspace group with the required roles.

1. Run `terraform init` and `terraform plan` and review the output.

   ```bash
   terraform init
   terraform plan -input=false -out bootstrapA.tfplan
   ```

1. To  validate your policies, run `gcloud beta terraform vet`. For installation instructions, see [Install Google Cloud CLI](https://cloud.google.com/docs/terraform/policy-validation/validate-policies#install).

1. Run the following commands and check for violations:

   ```bash
   export VET_PROJECT_ID=A-VALID-PROJECT-ID
   terraform show -json bootstrapA.tfplan > bootstrapA.json
   gcloud beta terraform vet bootstrapA.json --policy-library="../policy-library" --project ${VET_PROJECT_ID}
   ```

   *`A-VALID-PROJECT-ID`* must be an existing project you have access to. This is necessary because `gcloud beta terraform vet` needs to link resources to a valid Google Cloud Platform project.

1. Run `terraform apply`.

   ```bash
   terraform apply bootstrapA.tfplan
   ```

1. Run `terraform output` to get the email address of the terraform service accounts that will be used to run manual steps for `shared` environments in steps `3-networks-dual-svpc`, `3-networks-hub-and-spoke`, and `4-projects` and the state bucket that will be used by step 4-projects.

   ```bash
   export network_step_sa=$(terraform output -raw networks_step_terraform_service_account_email)
   export projects_step_sa=$(terraform output -raw projects_step_terraform_service_account_email)
   export projects_gcs_bucket_tfstate=$(terraform output -raw projects_gcs_bucket_tfstate)

   echo "network step service account = ${network_step_sa}"
   echo "projects step service account = ${projects_step_sa}"
   echo "projects gcs bucket tfstate = ${projects_gcs_bucket_tfstate}"
   ```

1. Copy the backend and update `backend.tf` with the name of your Google Cloud Storage bucket for Terraform's state. Also update the `backend.tf` of all steps.

   ```bash
   export backend_bucket=$(terraform output -raw gcs_bucket_tfstate)
   echo "backend_bucket = ${backend_bucket}"

   export backend_bucket_projects=$(terraform output -raw projects_gcs_bucket_tfstate)
   echo "backend_bucket_projects = ${backend_bucket_projects}"

   cp backend.tf.example backend.tf
   cd ..

   for i in `find . -name 'backend.tf'`; do sed -i'' -e "s/UPDATE_ME/${backend_bucket}/" $i; done
   for i in `find . -name 'backend.tf'`; do sed -i'' -e "s/UPDATE_PROJECTS_BACKEND/${backend_bucket_projects}/" $i; done

   cd 0A-bootstrap-seed
   ```

1. Re-run `terraform init`. When you're prompted, agree to copy Terraform state to Cloud Storage.

   ```bash
   terraform init
   ```

1. Continue with the instructions in the [0B-bootstrap-cicd](../0B-bootstrap-cicd/README.md) step.

**Note 1:** The stages after `0-bootstrap` use `terraform_remote_state` data source to read common configuration like the organization ID from the output of the `0-bootstrap` stage. They will [fail](../docs/TROUBLESHOOTING.md#error-unsupported-attribute) if the state is not copied to the Cloud Storage bucket.

**Note 2:** After the deploy, even if you did not receive the project quota error described in the [Troubleshooting guide](../docs/TROUBLESHOOTING.md#project-quota-exceeded), we recommend that you request 50 additional projects for the **projects step service account** created in this step.

## Running Terraform locally

If you deploy using Cloud Build, the bucket information is replaced in the state
backends as part of the build process when the build is executed by Cloud Build.
If you want to execute Terraform locally, you need to add your Cloud
Storage bucket to the `backend.tf` files.
Each step has instructions for this change.

<!-- BEGINNING OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| billing\_account | The ID of the billing account to associate projects with. | `string` | n/a | yes |
| bucket\_force\_destroy | When deleting a bucket, this boolean option will delete all contained objects. If false, Terraform will fail to delete buckets which contain objects. | `bool` | `false` | no |
| bucket\_prefix | Name prefix to use for state bucket created. | `string` | `"bkt"` | no |
| bucket\_tfstate\_kms\_force\_destroy | When deleting a bucket, this boolean option will delete the KMS keys used for the Terraform state bucket. | `bool` | `false` | no |
| default\_region | Default region to create resources where applicable. | `string` | `"us-central1"` | no |
| default\_region\_2 | Secondary default region to create resources where applicable. | `string` | `"us-west1"` | no |
| default\_region\_gcs | Case-Sensitive default region to create gcs resources where applicable. | `string` | `"US"` | no |
| default\_region\_kms | Secondary default region to create kms resources where applicable. | `string` | `"us"` | no |
| folder\_deletion\_protection | Prevent Terraform from destroying or recreating the folder. | `string` | `true` | no |
| folder\_prefix | Name prefix to use for folders created. Should be the same in all steps. | `string` | `"fldr"` | no |
| groups | Contain the details of the Groups to be created. | <pre>object({<br>    create_required_groups = optional(bool, false)<br>    create_optional_groups = optional(bool, false)<br>    billing_project        = optional(string, null)<br>    required_groups = object({<br>      group_org_admins     = string<br>      group_billing_admins = string<br>      billing_data_users   = string<br>      audit_data_users     = string<br>    })<br>    optional_groups = optional(object({<br>      gcp_security_reviewer    = optional(string, "")<br>      gcp_network_viewer       = optional(string, "")<br>      gcp_scc_admin            = optional(string, "")<br>      gcp_global_secrets_admin = optional(string, "")<br>      gcp_kms_admin            = optional(string, "")<br>    }), {})<br>  })</pre> | n/a | yes |
| initial\_group\_config | Define the group configuration when it is initialized. Valid values are: WITH\_INITIAL\_OWNER, EMPTY and INITIAL\_GROUP\_CONFIG\_UNSPECIFIED. | `string` | `"WITH_INITIAL_OWNER"` | no |
| org\_id | GCP Organization ID | `string` | n/a | yes |
| org\_policy\_admin\_role | Additional Org Policy Admin role for admin group. You can use this for testing purposes. | `bool` | `false` | no |
| parent\_folder | Optional - for an organization with existing projects or for development/validation. It will place all the example foundation resources under the provided folder instead of the root organization. The value is the numeric folder ID. The folder must already exist. | `string` | `""` | no |
| project\_deletion\_policy | The deletion policy for the project created. | `string` | `"PREVENT"` | no |
| project\_prefix | Name prefix to use for projects created. Should be the same in all steps. Max size is 3 characters. | `string` | `"prj"` | no |

## Outputs

| Name | Description |
|------|-------------|
| bootstrap\_cicd\_step\_terraform\_service\_account\_email | Bootstrap Step Terraform Account |
| bootstrap\_folder | The ID of the Bootstrap folder |
| bootstrap\_seed\_step\_terraform\_service\_account\_email | Bootstrap Step Terraform Account |
| common\_config | Common configuration data to be used in other steps. |
| environment\_step\_terraform\_service\_account\_email | Environment Step Terraform Account |
| gcs\_bucket\_tfstate | Bucket used for storing terraform state for Foundations Pipelines in Seed Project. |
| networks\_step\_terraform\_service\_account\_email | Networks Step Terraform Account |
| optional\_groups | List of Google Groups created that are optional to the Example Foundation steps. |
| organization\_step\_terraform\_service\_account\_email | Organization Step Terraform Account |
| projects\_gcs\_bucket\_tfstate | Bucket used for storing terraform state for stage 4-projects foundations pipelines in seed project. |
| projects\_step\_terraform\_service\_account\_email | Projects Step Terraform Account |
| required\_groups | List of Google Groups created that are required by the Example Foundation steps. |
| seed\_project\_id | Project where service accounts and core APIs will be enabled. |
| terraform\_env\_sa | Map of Terraform Accounts |

<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
