package object

import (
	"context"
	"fmt"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/url"
)

type ReadByNameInput struct {
	Name       string
	ObjectType ObjectType
}

func ReadByName(ctx context.Context, client http.Client, inp ReadByNameInput) (*ReadOutput, error) {

	client.Logger.Println("reading object by name")

	readUrl := url.ReadObjectByName(client.BaseUrl())
	req := client.NewGet(ctx, readUrl)
	req.QueryParams.Add("q", fmt.Sprintf("name:%s AND objectType:%s", inp.Name, inp.ObjectType))

	var page struct {
		Count  int          `json:"count"`
		Limit  int          `json:"limit"`
		Offset int          `json:"offset"`
		Items  []ReadOutput `json:"items"`
	}
	if err := req.Send(&page); err != nil {
		return nil, err
	}

	if len(page.Items) == 0 {
		return nil, fmt.Errorf("%w: no object found with name %q and type %s", http.NotFoundError, inp.Name, inp.ObjectType)
	}

	if len(page.Items) > 1 {
		return nil, fmt.Errorf("multiple objects found with name %q and type %s", inp.Name, inp.ObjectType)
	}

	return &page.Items[0], nil
}
