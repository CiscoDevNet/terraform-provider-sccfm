package ios

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/device/tags"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/devicetype"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/device"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
)

type ReadInput = device.ReadByUidInput
type ReadOutput struct {
	Uid             string          `json:"uid"`
	Name            string          `json:"name"`
	CreatedDate     int64           `json:"createdDate"`
	LastUpdatedDate int64           `json:"lastUpdatedDate"`
	DeviceType      devicetype.Type `json:"deviceType"`
	ConnectorUid    string          `json:"larUid"`
	ConnectorType   string          `json:"larType"`
	SocketAddress   string          `json:"ipv4"`
	Port            string          `json:"port"`
	Host            string          `json:"host"`
	Tags            tags.Type       `json:"tags"`

	IgnoreCertificate bool `json:"ignoreCertificate"`
}

func NewReadInput(uid string) *ReadInput {
	return device.NewReadByUidInput(uid)
}

func NewReadRequest(ctx context.Context, client http.Client, readInp ReadInput) *http.Request {
	return device.NewReadByUidRequest(ctx, client, readInp)
}

func Read(ctx context.Context, client http.Client, readInp ReadInput) (*ReadOutput, error) {

	client.Logger.Println("reading ios device")

	req := NewReadRequest(ctx, client, readInp)

	var outp ReadOutput
	if err := req.Send(&outp); err != nil {
		return nil, err
	}

	return &outp, nil
}
