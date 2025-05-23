package statemachine

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/statemachine"
)

type ReadInstanceByDeviceUidInput struct {
	Uid string // Uid of the device that runs the state machine
}

func NewReadInstanceByDeviceUidInput(deviceUid string) ReadInstanceByDeviceUidInput {
	return ReadInstanceByDeviceUidInput{
		Uid: deviceUid,
	}
}

type ReadInstanceByDeviceUidOutput = statemachine.Instance

var NewReadInstanceByDeviceUidOutputBuilder = statemachine.NewInstanceBuilder

func ReadInstanceByDeviceUid(ctx context.Context, client http.Client, readInp ReadInstanceByDeviceUidInput) (*ReadInstanceByDeviceUidOutput, error) {

	readUrl := url.ReadStateMachineInstance(client.BaseUrl())
	req := client.NewGet(ctx, readUrl)
	req.QueryParams.Add("limit", "1")
	req.QueryParams.Add("q", fmt.Sprintf("objectReference.uid:%s", readInp.Uid))
	req.QueryParams.Add("sort", "lastActiveDate:desc")

	var readRes []ReadInstanceByDeviceUidOutput
	if err := req.Send(&readRes); err != nil {
		return nil, err
	}
	if len(readRes) == 0 {
		return nil, NotFoundError
	}

	if len(readRes) > 1 {
		return nil, MoreThanOneRunningError
	}

	return &readRes[0], nil
}
