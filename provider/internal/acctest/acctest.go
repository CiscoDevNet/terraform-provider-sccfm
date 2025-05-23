// ci user: terraform-provider-cdo@lockhart.io
package acctest

import (
	"fmt"
	"os"
	"testing"

	"github.com/CiscoDevnet/terraform-provider-sccfm/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

const (
	apiTokenEnvName       = "ACC_TEST_CISCO_CDO_API_TOKEN"
	mspApiTokenEnvName    = "ACC_TEST_CISCO_CDO_MSP_API_TOKEN"
	apiTokenSecretName    = "staging-terraform-provider-cdo-acceptance-test-api-token"
	mspApiTokenSecretName = "staging-terraform-provider-cdo-acceptance-test-api-token"
)

var cdoSecretManager = NewCdoSecretManager("us-west-2")

func GetApiToken() (string, error) {
	tokenFromEnv, ok := os.LookupEnv(apiTokenEnvName)
	if ok {
		return tokenFromEnv, nil
	}

	tokenFromSecretManager, err := cdoSecretManager.getCurrentSecretValue(apiTokenSecretName)
	if err == nil {
		return tokenFromSecretManager, nil
	}

	return "", fmt.Errorf("failed to retrieve api token from environment variable and secret manager.\nenvironment variable name=%s\nsecret manager secret token name=%s\nplease set one of them.\ncause=%v", apiTokenEnvName, apiTokenSecretName, err)
}

func GetMspApiToken() (string, error) {
	tokenFromEnv, ok := os.LookupEnv(mspApiTokenEnvName)
	if ok {
		return tokenFromEnv, nil
	}

	tokenFromSecretManager, err := cdoSecretManager.getCurrentSecretValue(mspApiTokenEnvName)
	if err == nil {
		return tokenFromSecretManager, nil
	}

	return "", fmt.Errorf("failed to retrieve api token from environment variable and secret manager.\nenvironment variable name=%s\nsecret manager secret token name=%s\nplease set one of them.\ncause=%v", mspApiTokenEnvName, mspApiTokenSecretName, err)
}

func PreCheckFunc(t *testing.T) func() {
	return func() {
		_, err := GetApiToken()
		_, mspErr := GetMspApiToken()
		if err != nil {
			t.Fatalf("Precheck failed, cause=%v", err)
		}
		if mspErr != nil {
			t.Fatalf("Precheck failed, cause=%v", mspErr)
		}
	}
}

func ProviderConfig() string {
	token, err := GetApiToken()
	if err != nil {
		panic(fmt.Errorf("failed to retrieve api token, cause=%w", err))
	}

	return fmt.Sprintf(`
	provider "sccfm" {
		api_token = "%s"
		base_url = "https://ci.manage.security.cisco.com"
	}
	// New line
	`, token)
}

func MspProviderConfig() string {
	mspToken, err := GetMspApiToken()
	if err != nil {
		panic(fmt.Errorf("failed to retrieve api token, cause=%w", err))
	}

	return fmt.Sprintf(`
	provider "sccfm" {
		api_token = "%s"
		base_url = "https://ci.manage.security.cisco.com"
	}
	// New line
	`, mspToken)
}

// ProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var ProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"sccfm": providerserver.NewProtocol6WithError(provider.New("test")()),
}
