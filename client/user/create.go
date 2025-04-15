package user

import (
	"context"
	"fmt"

	"github.com/CiscoDevnet/terraform-provider-scc-firewall-manager/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-scc-firewall-manager/go-client/model"
)

func Create(ctx context.Context, client http.Client, createInp CreateUserInput) (*CreateUserOutput, error) {
	client.Logger.Println(fmt.Sprintf("Creating user %s", createInp.Username))
	req := NewCreateRequest(ctx, client, createInp)

	var outp model.UserDetails
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}
