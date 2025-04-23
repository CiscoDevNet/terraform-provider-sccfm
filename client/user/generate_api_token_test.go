package user_test

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/user"
	netHttp "net/http"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestGenerateApiToken(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	userToGenTokenFor := model.UserDetails{
		Username:    "lbj-api-only-user",
		Uid:         "donald-duck",
		ApiOnlyUser: true,
		Roles:       []string{"ROLE_SUPER_ADMIN"},
	}
	t.Run("Successfully generate an API token", func(t *testing.T) {
		apiTokenResponse := user.ApiTokenResponse{ApiToken: "jwt-token"}
		httpmock.RegisterResponder(
			netHttp.MethodGet,
			fmt.Sprintf("/api/rest/v1/users/api-only?limit=1&offset=0&q=name%%3A%s", userToGenTokenFor.Username),
			httpmock.NewJsonResponderOrPanic(200, model.UserPage{Count: 1, Items: []model.UserDetails{userToGenTokenFor}}))
		httpmock.RegisterResponder(
			netHttp.MethodPost,
			fmt.Sprintf("/api/rest/v1/users/%s/apiToken/generate", userToGenTokenFor.Uid),
			httpmock.NewJsonResponderOrPanic(200, apiTokenResponse))
		actual, err := user.GenerateApiToken(context.Background(),
			*http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute),
			*user.NewGenerateApiTokenInput(userToGenTokenFor.Username))
		assert.NotNil(t, actual, "API token response should not be nil")
		assert.Equal(t, apiTokenResponse, *actual)
		assert.Nil(t, err, "Error cannot be non-nil")
	})

	t.Run("Should fail if user not found", func(t *testing.T) {
		httpmock.RegisterResponder(
			netHttp.MethodGet,
			fmt.Sprintf("/api/rest/v1/users/api-only?limit=1&offset=0&q=name%%3A%s", userToGenTokenFor.Username),
			httpmock.NewJsonResponderOrPanic(500, nil))

		actual, err := user.GenerateApiToken(context.Background(),
			*http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute),
			*user.NewGenerateApiTokenInput(userToGenTokenFor.Username))
		assert.Nil(t, actual, "API token response should be nil")
		assert.NotNil(t, err, "Error cannot be nil")
	})

	t.Run("Should fail if token generation failed", func(t *testing.T) {
		httpmock.RegisterResponder(
			netHttp.MethodGet,
			fmt.Sprintf("/api/rest/v1/users/api-only?limit=1&offset=0&q=name%%3A%s", userToGenTokenFor.Username),
			httpmock.NewJsonResponderOrPanic(200, model.UserPage{Count: 1, Items: []model.UserDetails{userToGenTokenFor}}))
		httpmock.RegisterResponder(
			netHttp.MethodPost,
			fmt.Sprintf("/api/rest/v1/users/%s/apiToken/generate", userToGenTokenFor.Uid),
			httpmock.NewJsonResponderOrPanic(500, nil))

		actual, err := user.GenerateApiToken(context.Background(),
			*http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute),
			*user.NewGenerateApiTokenInput(userToGenTokenFor.Username))
		assert.Nil(t, actual, "API token response should be nil")
		assert.NotNil(t, err, "Error cannot be nil")
	})

}
