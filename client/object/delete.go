package object

import (
	"context"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/url"
)

func Delete(ctx context.Context, client http.Client, deleteInp DeleteInput) error {

	client.Logger.Println("deleting object")

	deleteUrl := url.DeleteObject(client.BaseUrl(), deleteInp.Uid)
	req := client.NewDelete(ctx, deleteUrl)

	if err := req.Send(nil); err != nil {
		return err
	}

	return nil
}
