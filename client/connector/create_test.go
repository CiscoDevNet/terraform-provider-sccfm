package connector_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/connector"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/user"
	"github.com/stretchr/testify/assert"

	"github.com/jarcoal/httpmock"
)

func TestCreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	validCreateRequestOutput := connector.NewCreateResponseBuilder().
		Uid(connectorUid).
		TenantUid(tenantUid).
		Name(connectorName).
		ServiceConnectivityState(serviceConnectivityState).
		State(state).
		Status(status).
		Build()

	validCreateOutput := connector.NewCreateOutputBuilder().
		CreateRequestOutput(validCreateRequestOutput).
		BootstrapData(bootstrapData).
		Build()

	validUserToken := user.NewGetTokenOutputBuilder().
		AccessToken(accessToken).
		RefreshToken(refreshToken).
		TenantUid(tenantUid).
		TenantName(tenantName).
		Scope(scope).
		TokenType(tokenType).
		Build()

	validQueueReadyOutput := connector.NewConnectorOutputBuilder().
		AsOnPremConnector().
		WithUid(validCreateRequestOutput.Uid).
		WithName(validCreateRequestOutput.Name).
		WithTenantUid(validCreateRequestOutput.TenantUid).
		WithCommunicationReady(true).
		Build()

	testCases := []struct {
		testName   string
		sdcName    string
		setupFunc  func()
		assertFunc func(output *connector.CreateOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully create SDC",
			sdcName:  connectorName,

			setupFunc: func() {
				httpmock.RegisterResponder(
					"POST",
					"/aegis/rest/v1/services/targets/proxies",
					httpmock.NewJsonResponderOrPanic(200, validCreateRequestOutput),
				)
				httpmock.RegisterResponder(
					"POST",
					"/anubis/rest/v1/oauth/token/external-compute",
					httpmock.NewJsonResponderOrPanic(200, validUserToken),
				)
				httpmock.RegisterResponder(
					"GET",
					fmt.Sprintf("/aegis/rest/v1/services/targets/proxies/%s", validQueueReadyOutput.Uid),
					httpmock.NewJsonResponderOrPanic(200, validQueueReadyOutput),
				)
			},

			assertFunc: func(output *connector.CreateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, validCreateOutput, *output)
			},
		},
		{
			testName: "should error if failed to create proxy",
			sdcName:  connectorName,

			setupFunc: func() {
				httpmock.RegisterResponder(
					"POST",
					"/aegis/rest/v1/services/targets/proxies",
					httpmock.NewJsonResponderOrPanic(500, "test error"),
				)
				httpmock.RegisterResponder(
					"POST",
					"/anubis/rest/v1/oauth/token",
					httpmock.NewJsonResponderOrPanic(200, validUserToken),
				)
			},

			assertFunc: func(output *connector.CreateOutput, err error, t *testing.T) {
				assert.NotNil(t, err, "error should not be nil")
				assert.Equal(t, output, &connector.CreateOutput{}, "output should be zero value")
			},
		},
		{
			testName: "should error if failed to retrieve user token",
			sdcName:  connectorName,

			setupFunc: func() {
				httpmock.RegisterResponder(
					"POST",
					"/aegis/rest/v1/services/targets/proxies",
					httpmock.NewJsonResponderOrPanic(200, validCreateRequestOutput),
				)
				httpmock.RegisterResponder(
					"POST",
					"/anubis/rest/v1/oauth/token/external-compute",
					httpmock.NewJsonResponderOrPanic(500, nil),
				)
			},

			assertFunc: func(output *connector.CreateOutput, err error, t *testing.T) {
				assert.NotNil(t, err, "error should not be nil")
				assert.Equal(t, output, &connector.CreateOutput{}, "output should be zero value")
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := connector.Create(
				context.Background(),
				*http.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				*connector.NewCreateInput(testCase.sdcName),
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
