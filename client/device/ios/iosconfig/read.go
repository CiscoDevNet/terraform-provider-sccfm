package iosconfig

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/statemachine"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/statemachine/state"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/url"
)

type ReadInput struct {
	SpecificUid string
}

type ReadOutput struct {
	Uid                 string               `json:"uid"`
	State               state.Type           `json:"state"`
	StateMachineDetails statemachine.Details `json:"stateMachineDetails"`
}

func NewReadInput(specificUid string) *ReadInput {
	return &ReadInput{
		SpecificUid: specificUid,
	}
}

func NewReadRequest(ctx context.Context, client http.Client, readReq ReadInput) *http.Request {

	readUrl := url.ReadDevice(client.BaseUrl(), readReq.SpecificUid)

	req := client.NewGet(ctx, readUrl)

	return req
}

func Read(ctx context.Context, client http.Client, readReq ReadInput) (*ReadOutput, error) {

	client.Logger.Println("reading iosconfig")

	req := NewReadRequest(ctx, client, readReq)

	var outp ReadOutput
	err := req.Send(&outp)
	if err != nil {
		return nil, err
	}

	return &outp, nil
}
