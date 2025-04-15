package user_test

import (
	"context"
	"fmt"
	netHttp "net/http"
	"net/url"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-scc-firewall-manager/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-scc-firewall-manager/go-client/model"
	"github.com/CiscoDevnet/terraform-provider-scc-firewall-manager/go-client/user"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestRevokeApiToken(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	userDetails := model.UserDetails{
		Username:    "barack@example.com",
		Uid:         "barack123",
		ApiOnlyUser: false,
		Roles:       []string{"ROLE_ADMIN"},
	}
	userPage := model.UserPage{
		Count: 1,
		Items: []model.UserDetails{userDetails},
	}

	t.Run("Successfully revoke an API token", func(t *testing.T) {
		httpmock.Reset()

		httpmock.RegisterResponder(
			netHttp.MethodGet,
			"/api/rest/v1/users?limit=1&offset=0&q=name%3A"+url.QueryEscape(userDetails.Username),
			httpmock.NewJsonResponderOrPanic(200, userPage),
		)
		httpmock.RegisterResponder(
			netHttp.MethodPost,
			fmt.Sprintf("/api/rest/v1/users/%s/apiToken/revoke", userDetails.Uid),
			httpmock.NewJsonResponderOrPanic(200, nil),
		)

		err := user.RevokeApiToken(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), *user.NewRevokeApiTokenInput(userDetails.Username))
		assert.Nil(t, err, "Error cannot be non-nil")
	})

	t.Run("Should fail if API token revocation failed", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder(
			netHttp.MethodGet,
			"/api/rest/v1/users?q=name:"+userDetails.Username,
			httpmock.NewJsonResponderOrPanic(200, userPage),
		)
		httpmock.RegisterResponder(
			netHttp.MethodPost,
			fmt.Sprintf("/api/rest/v1/users/%s/apiToken/revoke", userDetails.Uid),
			httpmock.NewJsonResponderOrPanic(400, nil),
		)

		err := user.RevokeApiToken(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), *user.NewRevokeApiTokenInput(userDetails.Username))
		assert.NotNil(t, err, "Error should be nil")
	})

	t.Run("Should fail if API token revocation failed because the user could not be found", func(t *testing.T) {
		httpmock.Reset()
		httpmock.RegisterResponder(
			netHttp.MethodGet,
			"/api/rest/v1/users?q=name:"+userDetails.Username,
			httpmock.NewJsonResponderOrPanic(500, nil),
		)

		err := user.RevokeApiToken(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), *user.NewRevokeApiTokenInput(userDetails.Username))
		assert.NotNil(t, err, "Error should be nil")
	})
}
