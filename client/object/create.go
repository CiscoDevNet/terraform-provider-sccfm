package object

import (
	"context"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/url"
)

func Create(ctx context.Context, client http.Client, createInp CreateInput) (*CreateOutput, error) {

	client.Logger.Println("creating object")

	createUrl := url.CreateObject(client.BaseUrl())
	req := client.NewPost(ctx, createUrl, createInp)

	var createOutp CreateOutput
	if err := req.Send(&createOutp); err != nil {
		return nil, err
	}

	return &createOutp, nil
}
