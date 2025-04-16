package examples

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-scc-firewall-manager/go-client/internal/http"
)

type CreateInput struct {
	Uid string
}

type CreateOutput struct {
}

func Create(ctx context.Context, client http.Client, createInp CreateInput) (*CreateOutput, error) {

	// TODO

	return nil, nil
}
