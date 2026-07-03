package cloudftd_test

import (
	"context"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/device/cloudftd"
	internalHttp "github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	internalTesting "github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/testing"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/ftd/tier"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestReadCloudFtdByUidUsingPublicApi(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testModel := internalTesting.NewRandomModel()
	ftdUid := testModel.FtdUid.String()
	performanceTier := tier.FTDv10

	internalTesting.MockGetOk(url.ReadDeviceUsingPublicApi(testModel.BaseUrl, ftdUid), map[string]any{
		"uid":                ftdUid,
		"name":               "test-ftd",
		"deviceType":         "CDFMC_MANAGED_FTD",
		"state":              "DONE",
		"softwareVersion":    "7.6.0",
		"ftdLicenses":        []string{"BASE", "THREAT"},
		"ftdPerformanceTier": performanceTier,
		"labels": map[string]any{
			"ungroupedLabels": []string{"terraform"},
			"groupedLabels": map[string][]string{
				"environment": {"test"},
			},
		},
		"cdFmcInfo": map[string]string{
			"cliKey": "configure manager add example.test reg-key nat-id example.test",
			"regKey": "reg-key",
			"natId":  "nat-id",
		},
		"fmcAccessPolicy": map[string]string{
			"uid": testModel.FtdAccessPolicyUid.String(),
		},
	})

	output, err := cloudftd.ReadByUidUsingPublicApi(
		context.Background(),
		*internalHttp.MustNewWithConfig(testModel.BaseUrl, "a_valid_token", 0, 0, time.Minute),
		cloudftd.NewReadByUidInput(ftdUid),
	)

	assert.Nil(t, err)
	assert.NotNil(t, output)
	assert.Equal(t, ftdUid, output.Uid)
	assert.Equal(t, "test-ftd", output.Name)
	assert.Equal(t, "CDFMC_MANAGED_FTD", output.DeviceType)
	assert.Equal(t, "DONE", output.State)
	assert.Equal(t, "7.6.0", output.SoftwareVersion)
	assert.Equal(t, "BASE,THREAT", output.Metadata.LicenseCaps)
	assert.Equal(t, &performanceTier, output.Metadata.PerformanceTier)
	assert.Equal(t, testModel.FtdAccessPolicyUid.String(), output.Metadata.AccessPolicyUid)
	assert.Equal(t, "configure manager add example.test reg-key nat-id example.test", output.Metadata.GeneratedCommand)
	assert.Equal(t, "example.test", output.Metadata.CloudManagerDomain)
	assert.Equal(t, "reg-key", output.Metadata.RegKey)
	assert.Equal(t, "nat-id", output.Metadata.NatID)
	assert.Equal(t, []string{"terraform"}, output.Tags.UngroupedTags())
	assert.Equal(t, map[string][]string{"environment": {"test"}}, output.Tags.GroupedTags())
}
