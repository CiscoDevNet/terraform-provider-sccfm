package tenantsettings

import (
	"context"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/settings"
)

func Read(ctx context.Context, client http.Client) (*settings.TenantSettings, error) {
	readUrl := url.ReadTenantSettings(client.BaseUrl())
	req := client.NewGet(ctx, readUrl)

	var settings settings.TenantSettings
	if err := req.Send(&settings); err != nil {
		return nil, err
	}

	return &settings, nil
}
