package object

import (
	"context"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/url"
)

func Update(ctx context.Context, client http.Client, updateInp UpdateInput) (*UpdateOutput, error) {

	client.Logger.Println("updating object")

	updateUrl := url.UpdateObject(client.BaseUrl(), updateInp.Uid)
	req := client.NewPatch(ctx, updateUrl, updateInp)

	var updateOutp UpdateOutput
	if err := req.Send(&updateOutp); err != nil {
		return nil, err
	}

	return &updateOutp, nil
}
