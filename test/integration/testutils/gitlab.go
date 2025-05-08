package testutils

import (
	"os"
	fp "path/filepath"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func writeGitLabTfvars(filename string, val interface{}) error {
	f := hclwrite.NewEmptyFile()
	gohcl.EncodeIntoBody(val, f.Body())
	return os.WriteFile(filename, f.Bytes(), 0644)
}

type GitLabRepository struct {
	Name string `cty:"repository_name"`
	URL  string `cty:"repository_url"`
}

type RepositoryConfig struct {
	GitlabReadAuthorizerCredentialSecretId string                      `cty:"gitlab_read_authorizer_credential_secret_id"`
	GitlabAuthorizerCredentialSecretId     string                      `cty:"gitlab_authorizer_credential_secret_id"`
	GitlabWebhookSecretId                  string                      `cty:"gitlab_webhook_secret_id"`
	GitlabEnterpriseHostUri                string                      `cty:"gitlab_enterprise_host_uri"`
	GitlabEnterpriseServiceDirectory       string                      `cty:"gitlab_enterprise_service_directory"`
	GitlabEnterpriseCaCertificate          string                      `cty:"gitlab_enterprise_ca_certificate"`
	Repositories                           map[string]GitLabRepository `cty:"repositories"`
}

type GitLabRepoConfigAutoTfvar struct {
	RepoConfig RepositoryConfig `hcl:"cloudbuildv2_repository_config"`
}

func WriteGitLabVarConfiguration(filepath, filename string, repoConfig RepositoryConfig, certificate []byte) error {

	repoConfig.GitlabEnterpriseCaCertificate = "REPLACE_WITH_SSL_CERT"

	autoTfvars := GitLabRepoConfigAutoTfvar{
		RepoConfig: repoConfig,
	}

	err := writeGitLabTfvars(fp.Join(filepath, filename), autoTfvars)
	if err != nil {
		return err
	}
	return ReplacePatternInTfVars(filename, "\"REPLACE_WITH_SSL_CERT\"", string(gitLabCertificate(certificate)), filepath )
}

func gitLabCertificate(certificate []byte) []byte {
	replacement := []byte("<<EOF\n")
	replacement = append(replacement, certificate...)
	return append(replacement, []byte("EOF\n")...)
}
