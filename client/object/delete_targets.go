package object

import (
	"context"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/url"
)

type DeleteTargetsInput struct {
	Uid         string
	TargetUuids []string
}

func DeleteTargets(ctx context.Context, client http.Client, inp DeleteTargetsInput) error {

	client.Logger.Println("detaching targets from object")

	deleteUrl := url.DeleteObjectTargets(client.BaseUrl(), inp.Uid)
	req := client.NewDelete(ctx, deleteUrl)
	for _, t := range inp.TargetUuids {
		req.QueryParams.Add("targetUuids", t)
	}

	if err := req.Send(nil); err != nil {
		return err
	}

	return nil
}
