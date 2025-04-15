package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-scc-firewall-manager/go-client/internal/url"

	"github.com/CiscoDevnet/terraform-provider-scc-firewall-manager/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-scc-firewall-manager/go-client/model"
)

func RevokeApiToken(ctx context.Context, client http.Client, revokeInp RevokeApiTokenInput) error {
	client.Logger.Println(fmt.Sprintf("Revoking API token for %s", revokeInp.Name))

	// 1. Find the user
	readReq := NewReadByUsernameRequest(ctx, client, revokeInp.Name, revokeInp.ApiOnlyUser)
	var userPage model.UserPage
	if readErr := readReq.Send(&userPage); readErr != nil {
		return readErr
	}
	if userPage.Count != 1 {
		return errors.New("User not found")
	}

	// 2. Revoke the API token by ID for the user
	revokeReq := client.NewPost(ctx, url.RevokeApiTokenForUser(client.BaseUrl(), userPage.Items[0].Uid), nil)
	var revokeOutput interface{}
	if revokeErr := revokeReq.Send(&revokeOutput); revokeErr != nil {
		return revokeErr
	}

	return nil
}
