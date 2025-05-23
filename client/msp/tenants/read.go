package tenants

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/url"
)

func ReadByUid(ctx context.Context, client http.Client, readInp ReadByUidInput) (*MspTenantOutput, error) {
	client.Logger.Println("reading tenant by UID " + readInp.Uid)

	readUrl := url.MspManagedTenantByUid(client.BaseUrl(), readInp.Uid)
	req := client.NewGet(ctx, readUrl)

	var outp MspTenantOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}

func ReadByName(ctx context.Context, client http.Client, readInp ReadByNameInput) (*MspTenantsOutput, error) {
	client.Logger.Println("reading tenant by name " + readInp.Name)
	findByNameUrl := url.FindMspManagedTenantsByName(client.BaseUrl(), readInp.Name)
	req := client.NewGet(ctx, findByNameUrl)

	var outp MspTenantsOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}
