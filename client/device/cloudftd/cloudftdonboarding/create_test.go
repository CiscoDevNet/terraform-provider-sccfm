package cloudftdonboarding_test

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/device/cloudftd/cloudftdonboarding"
	internalHttp "github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/publicapi/transaction/transactiontype"
	internalTesting "github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/testing"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/url"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCloudFtdOnboardingCreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testModel := internalTesting.NewRandomModel()

	ftdOnboardingInput := testModel.FtdOnboardingInput()
	doneTransaction := testModel.CreateDoneTransaction(ftdOnboardingInput.FtdUid, transactiontype.REGISTER_FTD)
	errorTransaction := testModel.CreateErrorTransaction(ftdOnboardingInput.FtdUid, transactiontype.REGISTER_FTD)
	ftdReadOutput := testModel.FtdReadOutput()
	publicApiFtdReadOutput := map[string]any{
		"uid":                ftdReadOutput.Uid,
		"name":               ftdReadOutput.Name,
		"deviceType":         ftdReadOutput.DeviceType,
		"state":              ftdReadOutput.State,
		"softwareVersion":    ftdReadOutput.SoftwareVersion,
		"ftdLicenses":        []string{},
		"ftdPerformanceTier": ftdReadOutput.Metadata.PerformanceTier,
		"labels": map[string]any{
			"ungroupedLabels": []string{},
			"groupedLabels":   map[string][]string{},
		},
		"cdFmcInfo": map[string]string{},
		"fmcAccessPolicy": map[string]string{
			"uid": ftdReadOutput.Metadata.AccessPolicyUid,
		},
	}

	testCases := []struct {
		testName   string
		input      cloudftdonboarding.CreateInput
		setupFunc  func(createInp cloudftdonboarding.CreateInput)
		assertFunc func(output *cloudftdonboarding.CreateOutput, err error, t *testing.T)
	}{
		{
			testName: "successful ftd onboarding",
			input:    cloudftdonboarding.NewCreateInput(ftdOnboardingInput.FtdUid),
			setupFunc: func(createInp cloudftdonboarding.CreateInput) {
				internalTesting.MockPostAccepted(url.RegisterFtd(testModel.BaseUrl), doneTransaction)
				internalTesting.MockGetOk(url.ReadDeviceUsingPublicApi(testModel.BaseUrl, ftdReadOutput.Uid), publicApiFtdReadOutput)
			},
			assertFunc: func(output *cloudftdonboarding.CreateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, ftdReadOutput.Uid, output.Uid)
				assert.Equal(t, ftdReadOutput.Name, output.Name)
				assert.Equal(t, ftdReadOutput.Metadata.AccessPolicyUid, output.Metadata.AccessPolicyUid)
				assert.Equal(t, ftdReadOutput.Metadata.PerformanceTier, output.Metadata.PerformanceTier)
				assert.Equal(t, 1, httpmock.GetCallCountInfo()["GET "+url.ReadDeviceUsingPublicApi(testModel.BaseUrl, ftdReadOutput.Uid)])
				assert.Zero(t, httpmock.GetCallCountInfo()["GET "+url.ReadDevice(testModel.BaseUrl, ftdReadOutput.Uid)])
			},
		},
		{
			testName: "fails ftd onboarding if transaction fails",
			input:    cloudftdonboarding.NewCreateInput(ftdOnboardingInput.FtdUid),
			setupFunc: func(createInp cloudftdonboarding.CreateInput) {
				internalTesting.MockPostError(url.RegisterFtd(testModel.BaseUrl), errorTransaction)
			},
			assertFunc: func(output *cloudftdonboarding.CreateOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
				assert.ErrorContains(t, err, errorTransaction.ErrorMessage)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc(testCase.input)

			output, err := cloudftdonboarding.Create(
				context.Background(),
				*internalHttp.MustNewWithConfig(testModel.BaseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
