package duoadminpanel

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-scc-firewall-manager/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-scc-firewall-manager/go-client/internal/http"
)

type DeleteInput struct {
	Uid string
}

type DeleteOutput = device.DeleteOutput

func Delete(ctx context.Context, client http.Client, deleteInp DeleteInput) (*DeleteOutput, error) {

	client.Logger.Println("deleting duo admin panel")

	return device.Delete(ctx, client, *device.NewDeleteInput(deleteInp.Uid))
}
