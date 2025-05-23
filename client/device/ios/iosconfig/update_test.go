package iosconfig_test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/device/ios/iosconfig"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/crypto"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"

	internalHttp "github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/jsonutil"
	"github.com/jarcoal/httpmock"
)

func TestIosConfigUpdate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	username := "unit-test-username"
	password := "not a real password"

	rsaKeyBits := 512
	rsaKey, err := rsa.GenerateKey(rand.Reader, rsaKeyBits)
	assert.Nil(t, err, "could not generate rsa key")

	testCases := []struct {
		testName   string
		input      iosconfig.UpdateInput
		setupFunc  func(input iosconfig.UpdateInput, t *testing.T)
		assertFunc func(output *iosconfig.UpdateOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully updates iOS config",
			input: iosconfig.UpdateInput{
				SpecificUid: iosConfigUid,
				Username:    username,
				Password:    password,
			},

			setupFunc: func(input iosconfig.UpdateInput, t *testing.T) {
				httpmock.RegisterResponder(
					"PUT",
					buildIosConfigPath(iosConfigUid),
					func(r *http.Request) (*http.Response, error) {
						requestBody, err := internalHttp.ReadRequestBody[iosconfig.UpdateBody](r)
						assert.Nil(t, err)

						expectedBody := iosconfig.UpdateBody{
							SmContext: iosconfig.SmContext{
								AcceptCert: true,
							},
							Credentials: fmt.Sprintf(`{"username":"%s","password":"%s"}`, input.Username, input.Password),
						}
						assert.Equal(t, expectedBody, *requestBody)

						return httpmock.NewJsonResponse(200, iosconfig.UpdateOutput{Uid: iosConfigUid})
					},
				)
			},

			assertFunc: func(output *iosconfig.UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, iosconfig.UpdateOutput{Uid: iosConfigUid}, *output)
			},
		},

		{
			testName: "successfully updates iOS config when encrypting credentials",
			input: iosconfig.UpdateInput{
				SpecificUid: iosConfigUid,
				Username:    username,
				Password:    password,
				PublicKey: &model.PublicKey{
					KeyId:      "12341234-1234-1234-1234-123412341234",
					Version:    2,
					EncodedKey: crypto.MustBase64PublicKeyFromRsaKey(rsaKey),
				},
			},

			setupFunc: func(input iosconfig.UpdateInput, t *testing.T) {
				httpmock.RegisterResponder(
					"PUT",
					buildIosConfigPath(iosConfigUid),
					func(r *http.Request) (*http.Response, error) {
						requestBody, err := internalHttp.ReadRequestBody[iosconfig.UpdateBody](r)
						assert.Nil(t, err)

						if !requestBody.SmContext.AcceptCert {
							t.Errorf("expected 'SmContext.AcceptCert' to true got: %t", requestBody.SmContext.AcceptCert)
						}

						credentials, err := jsonutil.UnmarshalStruct[model.Credentials]([]byte(requestBody.Credentials))
						assert.Nil(t, err)

						decryptedUsername := crypto.MustDecryptBase64EncodedPkcs1v15Value(rsaKey, []byte(credentials.Username))
						assert.Equal(t, input.Username, decryptedUsername, `expected decrypted username to equal '%s', got: '%s'`, input.Username, decryptedUsername)

						decryptedPassword := crypto.MustDecryptBase64EncodedPkcs1v15Value(rsaKey, []byte(credentials.Password))
						assert.Equal(t, input.Password, decryptedPassword, `expected decrypted password to equal '%s', got: '%s'`, input.Password, decryptedPassword)

						assert.Equal(t, input.PublicKey.KeyId, credentials.KeyId, "expected keyId to equal '%s', got: '%s'", input.PublicKey.KeyId, credentials.KeyId)

						return httpmock.NewJsonResponse(200, iosconfig.UpdateOutput{Uid: iosConfigUid})
					},
				)
			},

			assertFunc: func(output *iosconfig.UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, iosconfig.UpdateOutput{Uid: iosConfigUid}, *output)
			},
		},

		{
			testName: "returns error when updating iOS config that does not exist",
			input: iosconfig.UpdateInput{
				SpecificUid: iosConfigUid,
				Username:    username,
				Password:    password,
			},

			setupFunc: func(input iosconfig.UpdateInput, t *testing.T) {
				httpmock.RegisterResponder(
					"PUT",
					buildIosConfigPath(input.SpecificUid),
					httpmock.NewStringResponder(404, ""),
				)
			},

			assertFunc: func(output *iosconfig.UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
			},
		},

		{
			testName: "returns error when remote service updating iOS config experiences an issue",
			input: iosconfig.UpdateInput{
				SpecificUid: iosConfigUid,
				Username:    username,
				Password:    password,
			},

			setupFunc: func(input iosconfig.UpdateInput, t *testing.T) {
				httpmock.RegisterResponder(
					"PUT",
					buildIosConfigPath(input.SpecificUid),
					httpmock.NewStringResponder(500, ""),
				)
			},

			assertFunc: func(output *iosconfig.UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, output)
				assert.NotNil(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc(testCase.input, t)

			output, err := iosconfig.Update(
				context.Background(),
				*internalHttp.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
