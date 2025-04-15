package user

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-scc-firewall-manager/go-client/internal/http"
)

func GenerateApiToken(ctx context.Context, client http.Client, generateApiTokenInp GenerateApiTokenInput) (*ApiTokenResponse, error) {
	client.Logger.Println(fmt.Sprintf("Generating API token for user %s", generateApiTokenInp.Name))
	user, err := ReadByUsername(ctx, client, ReadByUsernameInput{
		Name:        generateApiTokenInp.Name,
		ApiOnlyUser: true,
	})
	if err != nil {
		return nil, err
	}
	client.Logger.Printf("Found user %s with uid %s\n", user.Username, user.Uid)

	req := NewGenerateApiTokenRequest(ctx, client, user.Uid)
	var apiToken ApiTokenResponse
	if err := req.Send(&apiToken); err != nil {
		return nil, err
	}

	return &apiToken, nil
}
