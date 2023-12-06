# Private GKE cluster to run GitHub Action Runners

<!-- BEGINNING OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| cluster\_name | n/a | `string` | `"gke-b-gh-runner"` | no |
| environment\_code | A short form of the folder level resources (environment) within the Google Cloud organization (ex. d). | `string` | n/a | yes |
| firewall\_enable\_logging | Toggle firewall logging for VPC Firewalls. | `bool` | `true` | no |
| ip\_range\_pods\_cidr | The secondary ip range cidr to use for pods | `string` | n/a | yes |
| ip\_range\_services\_cider | The secondary ip range cidr to use for services | `string` | n/a | yes |
| master\_ipv4\_cidr\_block | (Beta) The IP range in CIDR notation to use for the hosted master network | `string` | `"10.0.0.0/28"` | no |
| nat\_bgp\_asn | BGP ASN for NAT cloud routes. | `number` | `64514` | no |
| nat\_enabled | Toggle creation of NAT cloud router. | `bool` | `false` | no |
| nat\_num\_addresses | Number of external IPs to reserve for Cloud NAT. | `number` | `2` | no |
| network\_tag | Network tag to apply to instances to allow external access. | `string` | `"gh-runner-vm"` | no |
| private\_service\_connect\_ip | Internal IP to be used as the private service connect endpoint. | `string` | n/a | yes |
| project\_id | The project id to deploy Github Runner | `string` | n/a | yes |
| region | The GCP region to deploy instances into | `string` | n/a | yes |
| subnet\_ip\_cidr | IP range for the subnet. | `string` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| cluster\_host | n/a |
| cluster\_membership\_id | n/a |
| cluster\_name | n/a |

<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
