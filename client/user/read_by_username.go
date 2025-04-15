package user

import (
	"context"
	"errors"

	"github.com/CiscoDevnet/terraform-provider-scc-firewall-manager/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-scc-firewall-manager/go-client/model"
)

func ReadByUsername(ctx context.Context, client http.Client, readInp ReadByUsernameInput) (*ReadUserOutput, error) {

	readReq := NewReadByUsernameRequest(ctx, client, readInp.Name, readInp.ApiOnlyUser)
	var userPage model.UserPage
	if readErr := readReq.Send(&userPage); readErr != nil {
		return nil, readErr
	}
	if userPage.Count != 1 {
		return nil, errors.New("user not found")
	}
	return &userPage.Items[0], nil
}
