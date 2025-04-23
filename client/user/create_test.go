package user_test

import (
	"context"
	"fmt"
	netHttp "net/http"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/user"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("successfully create user", func(t *testing.T) {
		httpmock.Reset()
		expected := model.UserDetails{
			Username:    "george@example.com",
			Uid:         "donald-duck",
			ApiOnlyUser: false,
			Roles:       []string{"ROLE_SUPER_ADMIN"},
		}

		httpmock.RegisterResponder(
			netHttp.MethodPost,
			fmt.Sprintf("/api/rest/v1/users"),
			httpmock.NewJsonResponderOrPanic(200, expected),
		)

		firstName := "George"
		lastName := "Washington"
		actual, err := user.Create(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), *user.NewCreateUserInput(expected.Username, expected.Roles[0], expected.ApiOnlyUser, &firstName, &lastName))

		assert.NotNil(t, actual, "User details returned must not be nil")
		assert.Equal(t, expected, *actual, "Actual user details do not match expected")
		assert.Nil(t, err, "error should be nil")
	})

	t.Run("should error if failed to create user", func(t *testing.T) {
		username := "bill_api"
		userRoles := []string{"ROLE_READ_ONLY"}
		httpmock.RegisterResponder(
			netHttp.MethodPost,
			fmt.Sprintf("/api/rest/v1/users"),
			httpmock.NewJsonResponderOrPanic(500, nil),
		)

		actual, err := user.Create(context.Background(), *http.MustNewWithConfig(baseUrl, "valid_token", 0, 0, time.Minute), *user.NewCreateUserInput(username, userRoles[0], true, nil, nil))
		assert.Nil(t, actual, "Expected actual user not to be created")
		assert.NotNil(t, err, "Expected error")
	})
}
