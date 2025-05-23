package tenantsettings_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	internalHttp "github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/settings"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/settings/tenantsettings"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestReadTenantSettings(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	baseURL := "https://unit-test.cdo.cisco.com"
	tenantSettings := settings.TenantSettings{
		Uid:                                   uuid.New(),
		ChangeRequestSupportEnabled:           false,
		AutoAcceptDeviceChangesEnabled:        false,
		WebAnalyticsEnabled:                   false,
		ScheduledDeploymentsEnabled:           false,
		DenyCiscoSupportAccessToTenantEnabled: false,
		MultiCloudDefenseEnabled:              false,
		AutoDiscoverOnPremFmcsEnabled:         false,
		ConflictDetectionInterval:             settings.ConflictDetectionIntervalEvery10Minutes,
	}

	testCases := []struct {
		testName   string
		setupFunc  func()
		assertFunc func(output *settings.TenantSettings, err error, t *testing.T)
	}{
		{
			testName: "successfully read tenant settings",
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodGet,
					url.ReadTenantSettings(baseURL),
					httpmock.NewJsonResponderOrPanic(http.StatusOK, tenantSettings),
				)
			},
			assertFunc: func(output *settings.TenantSettings, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, tenantSettings, *output)
			},
		},
		{
			testName: "return error when read tenant settings error",
			setupFunc: func() {
				httpmock.RegisterResponder(
					http.MethodGet,
					url.ReadTenantSettings(baseURL),
					httpmock.NewStringResponder(http.StatusInternalServerError, "internal server error"),
				)
			},
			assertFunc: func(output *settings.TenantSettings, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := tenantsettings.Read(
				context.Background(),
				*internalHttp.MustNewWithConfig(baseURL, "a_valid_token", 0, 0, time.Minute),
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
