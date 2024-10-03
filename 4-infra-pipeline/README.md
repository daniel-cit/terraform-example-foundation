# 4-infra-pipeline

This repo is part of a multi-part guide that shows how to configure and deploy
the example.com reference architecture described in
[Google Cloud security foundations guide](https://cloud.google.com/architecture/security-foundations). The following table lists the parts of the guide.

<table>
<tbody>
<tr>
<td><a href="../0-bootstrap">0-bootstrap</a></td>
<td>Bootstraps a Google Cloud organization, creating all the required resources
and permissions to start using the Cloud Foundation Toolkit (CFT). This
step also configures a <a href="../docs/GLOSSARY.md#foundation-cicd-pipeline">CI/CD Pipeline</a> for foundations code in subsequent
stages.</td>
</tr>
<tr>
<td><a href="../1-org">1-org</a></td>
<td>Sets up top level shared folders, networking projects, and
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
Hub and Spoke network model. It also sets up the global DNS hub</td>
</tr>
<tr>
<td>4-infra-pipeline (this file)</td>
<td>Sets up a application infrastructure pipeline for applications.</td>
</tr>
<tr>
<td><a href="../5-projects">5-projects</a></td>
<td>Sets up a folder structure and projects which are connected as service projects to the shared VPC created in the previous stage.</td>
</tr>
<tr>
<td><a href="../6-project-infra">6-project-infra</a></td>
<td>Deploy a simple <a href="https://cloud.google.com/compute/">Compute Engine</a> instance in one of the business unit projects using the infra pipeline set up in 4-infra-pipeline.</td>
</tr>
</tbody>
</table>

For an overview of the architecture and the parts, see the
[terraform-example-foundation README](https://github.com/terraform-google-modules/terraform-example-foundation).

## Purpose

The purpose of this step is to set up the infrastructure pipelines for applications that will be deployed in the service projects to the shared VPC created in the previous stage.

For each business unit, a shared `infra-pipeline` project is created along with Cloud Build triggers, CSRs for application infrastructure code and Google Cloud Storage buckets for state storage.

This step follows the same [conventions](https://github.com/terraform-google-modules/terraform-example-foundation#branching-strategy) as the Foundation pipeline deployed in [0-bootstrap](https://github.com/terraform-google-modules/terraform-example-foundation/blob/master/0-bootstrap/README.md).
A custom [workspace](https://github.com/terraform-google-modules/terraform-google-bootstrap/blob/master/modules/tf_cloudbuild_workspace/README.md) (`bu1-example-app`) is created by this pipeline and necessary roles are granted to the Terraform Service Account of this workspace by enabling variable `sa_roles` as shown in this [example](https://github.com/terraform-google-modules/terraform-example-foundation/blob/master/5-projects/modules/base_env/example_base_shared_vpc_project.tf).

This pipeline is utilized to deploy resources in projects across development/nonproduction/production in step [6-project-infra](../6-project-infra/README.md).
Other Workspaces can also be created to isolate deployments if needed.

## Prerequisites

1. 0-bootstrap executed successfully.
1. 1-org executed successfully.
1. 2-environments executed successfully.
1. 3-networks executed successfully.

### Troubleshooting

Please refer to [troubleshooting](../docs/TROUBLESHOOTING.md) if you run into issues during this step.

## Usage

**Note:** If you are using MacOS, replace `cp -RT` with `cp -R` in the relevant
commands. The `-T` flag is needed for Linux, but causes problems for MacOS.

### Deploying with Cloud Build

1. Clone the `gcp-infra-pipeline` repo based on the Terraform output from the `0-bootstrap` step.
Clone the repo at the same level of the `terraform-example-foundation` folder, the following instructions assume this layout.
Run `terraform output cloudbuild_project_id` in the `0-bootstrap` folder to get the Cloud Build Project ID.

   ```bash
   export CLOUD_BUILD_PROJECT_ID=$(terraform -chdir="terraform-example-foundation/0-bootstrap/" output -raw cloudbuild_project_id)
   echo ${CLOUD_BUILD_PROJECT_ID}

   gcloud source repos clone gcp-infra-pipeline --project=${CLOUD_BUILD_PROJECT_ID}
   ```

1. Change to the freshly cloned repo, change to the non-main branch and copy contents of foundation to new repo.

   ```bash
   cd gcp-infra-pipeline
   git checkout -b plan

   cp -RT ../terraform-example-foundation/4-infra-pipeline/ .
   cp ../terraform-example-foundation/build/cloudbuild-tf-* .
   cp ../terraform-example-foundation/build/tf-wrapper.sh .
   chmod 755 ./tf-wrapper.sh
   ```

1. Rename `terraform.example.tfvars` files to `terraform.tfvars`.

   ```bash
   mv ./business_unit_1/shared/terraform.example.tfvars ./business_unit_1/shared/terraform.tfvars
   ```

1. See the shared folder [README.md](./business_unit_1/shared/README.md) file for additional information on the values in the `terraform.tfvars` file.

1. Use `terraform output` to get the backend bucket value from 0-bootstrap output.

   ```bash
   export remote_state_bucket=$(terraform -chdir="../terraform-example-foundation/0-bootstrap/" output -raw gcs_bucket_tfstate)
   echo "remote_state_bucket = ${remote_state_bucket}"

   sed -i'' -e "s/REMOTE_STATE_BUCKET/${remote_state_bucket}/" ./business_unit_1/shared/terraform.tfvars
   ```

1. (Optional) If you want additional subfolders for separate business units or entities, make additional copies of the folder `business_unit_1` and modify any values that vary across business unit like `business_code`, `business_unit`, or `subnet_ip_range`.

For example, to create a new business unit similar to business_unit_1, run the following:

   ```bash
   #copy the business_unit_1 folder and it's contents to a new folder business_unit_2
   cp -r  business_unit_1 business_unit_2

   # search all files under the folder `business_unit_2` and replace strings for business_unit_1 with strings for business_unit_2
   grep -rl bu1 business_unit_2/ | xargs sed -i 's/bu1/bu2/g'
   grep -rl business_unit_1 business_unit_2/ | xargs sed -i 's/business_unit_1/business_unit_2/g'
   ```

1. Commit changes.

   ```bash
   git add .
   git commit -m 'Initialize infra-pipeline repo'
   ```

1. Push your plan branch to trigger a plan for all environments. Because the
   _plan_ branch is not a [named environment branch](../docs/FAQ.md#what-is-a-named-branch), pushing your _plan_
   branch triggers _terraform plan_ but not _terraform apply_. Review the plan output in your Cloud Build project https://console.cloud.google.com/cloud-build/builds;region=DEFAULT_REGION?project=YOUR_CLOUD_BUILD_PROJECT_ID

   ```bash
   git push --set-upstream origin plan
   ```

1. Merge changes to production. Because this is a [named environment branch](../docs/FAQ.md#what-is-a-named-branch),
   pushing to this branch triggers both _terraform plan_ and _terraform apply_. Review the apply output in your Cloud Build project. https://console.cloud.google.com/cloud-build/builds;region=DEFAULT_REGION?project=YOUR_CLOUD_BUILD_PROJECT_ID

   ```bash
   git checkout -b production
   git push origin production
   ```

1. You can now move to the instructions in the [6-project-infra](../6-project-infra/README.md) step.

### Deploying with Jenkins

See `0-bootstrap` [README-Jenkins.md](../0-bootstrap/README-Jenkins.md#deploying-step-4).

### Deploying with GitHub Actions

See `0-bootstrap` [README-GitHub.md](../0-bootstrap/README-GitHub.md#deploying-step-4-infra-pipeline).

### Run Terraform locally

1. The next instructions assume that you are at the same level of the `terraform-example-foundation` folder. Change into `4-infra-pipeline` folder, copy the Terraform wrapper script and ensure it can be executed.

   ```bash
   cd terraform-example-foundation/4-infra-pipeline
   cp ../build/tf-wrapper.sh .
   chmod 755 ./tf-wrapper.sh
   ```

1. Rename `terraform.example.tfvars` files to `terraform.tfvars`.

   ```bash
   mv ./business_unit_1/shared/terraform.example.tfvars ./business_unit_1/shared/terraform.tfvars
   ```

1. See the shared folder [README.md](./business_unit_1/shared/README.md) file for additional information on the values in the `terraform.tfvars` file.
   Use `terraform output` to get the remote state bucket (the backend bucket used by previous steps) value from `0-bootstrap` output.

   ```bash
   export remote_state_bucket=$(terraform -chdir="../0-bootstrap/" output -raw gcs_bucket_tfstate)
   echo "remote_state_bucket = ${remote_state_bucket}"

   sed -i'' -e "s/REMOTE_STATE_BUCKET/${remote_state_bucket}/" ./business_unit_1/shared/terraform.tfvars
   ```

We will now deploy to our `shared` environment using the `tf-wrapper.sh` script.

To use the `validate` option of the `tf-wrapper.sh` script, please follow the [instructions](https://cloud.google.com/docs/terraform/policy-validation/validate-policies#install) to install the terraform-tools component.

1. Use `terraform output` to get the Cloud Build Project ID and the environment step Terraform Service Account from 0-bootstrap output. An environment variable `GOOGLE_IMPERSONATE_SERVICE_ACCOUNT` will be set using the Terraform Service Account to enable impersonation.

   ```bash
   export CLOUD_BUILD_PROJECT_ID=$(terraform -chdir="../0-bootstrap/" output -raw cloudbuild_project_id)
   echo ${CLOUD_BUILD_PROJECT_ID}

   export GOOGLE_IMPERSONATE_SERVICE_ACCOUNT=$(terraform -chdir="../0-bootstrap/" output -raw projects_step_terraform_service_account_email)
   echo ${GOOGLE_IMPERSONATE_SERVICE_ACCOUNT}
   ```

1. (Optional) If you want additional subfolders for separate business units or entities, make additional copies of the folder `business_unit_1` and modify any values that vary across business unit like `business_code`, `business_unit`, or `subnet_ip_range`.

For example, to create a new business unit similar to business_unit_1, run the following:

   ```bash
   #copy the business_unit_1 folder and it's contents to a new folder business_unit_2
   cp -r  business_unit_1 business_unit_2

   # search all files under the folder `business_unit_2` and replace strings for business_unit_1 with strings for business_unit_2
   grep -rl bu1 business_unit_2/ | xargs sed -i 's/bu1/bu2/g'
   grep -rl business_unit_1 business_unit_2/ | xargs sed -i 's/business_unit_1/business_unit_2/g'
   ```


1. Run `init` and `plan` and review output for environment shared.

   ```bash
   ./tf-wrapper.sh init shared
   ./tf-wrapper.sh plan shared
   ```

1. Run `validate` and check for violations.

   ```bash
   ./tf-wrapper.sh validate shared $(pwd)/../policy-library ${CLOUD_BUILD_PROJECT_ID}
   ```

1. Run `apply` shared.

   ```bash
   ./tf-wrapper.sh apply shared
   ```

If you received any errors or made any changes to the Terraform config or `terraform.tfvars`, you must re-run `./tf-wrapper.sh plan <env>` before running `./tf-wrapper.sh apply <env>`.

Before executing the next stages, unset the `GOOGLE_IMPERSONATE_SERVICE_ACCOUNT` environment variable.

```bash
unset GOOGLE_IMPERSONATE_SERVICE_ACCOUNT
```
