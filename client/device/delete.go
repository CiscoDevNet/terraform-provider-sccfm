package device

import (
	"context"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/url"
)

type DeleteInput struct {
	Uid string `json:"uid"`
}

type DeleteOutput struct {
}

func NewDeleteInput(uid string) *DeleteInput {
	return &DeleteInput{
		Uid: uid,
	}
}

func NewDeleteRequest(ctx context.Context, client http.Client, deleteInp DeleteInput) *http.Request {

	url := url.DeleteDevice(client.BaseUrl(), deleteInp.Uid)

	req := client.NewDelete(ctx, url)

	return req
}

func Delete(ctx context.Context, client http.Client, deleteInp DeleteInput) (*DeleteOutput, error) {

	client.Logger.Println("deleting device")

	req := NewDeleteRequest(ctx, client, deleteInp)

	var outp DeleteOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}
