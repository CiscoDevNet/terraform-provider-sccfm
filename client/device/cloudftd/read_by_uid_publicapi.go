package cloudftd

import (
	"context"
	"strings"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/device/publicapilabels"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/device/tags"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/ftd/license"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/ftd/tier"
)

type publicApiReadOutput struct {
	Uid                string                `json:"uid"`
	DeviceType         string                `json:"deviceType"`
	Name               string                `json:"name"`
	State              string                `json:"state"`
	SoftwareVersion    string                `json:"softwareVersion"`
	FtdLicenses        []license.Type        `json:"ftdLicenses"`
	FtdPerformanceTier *tier.Type            `json:"ftdPerformanceTier"`
	Labels             publicapilabels.Type  `json:"labels"`
	CdFmcInfo          publicApiCdFmcInfo    `json:"cdFmcInfo"`
	FmcAccessPolicy    publicApiAccessPolicy `json:"fmcAccessPolicy"`
}

type publicApiCdFmcInfo struct {
	CliKey string `json:"cliKey"`
	RegKey string `json:"regKey"`
	NatID  string `json:"natId"`
}

type publicApiAccessPolicy struct {
	Uid string `json:"uid"`
}

func ReadByUidUsingPublicApi(ctx context.Context, client http.Client, readInp ReadByUidInput) (*ReadOutput, error) {

	readUrl := url.ReadDeviceUsingPublicApi(client.BaseUrl(), readInp.Uid)
	req := client.NewGet(ctx, readUrl)

	var publicOutput publicApiReadOutput
	if err := req.Send(&publicOutput); err != nil {
		return nil, err
	}

	return &ReadOutput{
		Uid:             publicOutput.Uid,
		DeviceType:      publicOutput.DeviceType,
		Name:            publicOutput.Name,
		State:           publicOutput.State,
		Tags:            tags.New(publicOutput.Labels.UngroupedLabels, publicOutput.Labels.GroupedLabels),
		SoftwareVersion: publicOutput.SoftwareVersion,
		Metadata: Metadata{
			AccessPolicyUid:    publicOutput.FmcAccessPolicy.Uid,
			CloudManagerDomain: cloudManagerDomainFromCliKey(publicOutput.CdFmcInfo.CliKey),
			GeneratedCommand:   publicOutput.CdFmcInfo.CliKey,
			LicenseCaps:        license.LicensesToString(publicOutput.FtdLicenses),
			NatID:              publicOutput.CdFmcInfo.NatID,
			PerformanceTier:    publicOutput.FtdPerformanceTier,
			RegKey:             publicOutput.CdFmcInfo.RegKey,
		},
	}, nil
}

func cloudManagerDomainFromCliKey(cliKey string) string {
	fields := strings.Fields(cliKey)
	if len(fields) >= 4 && fields[0] == "configure" && fields[1] == "manager" && fields[2] == "add" {
		return fields[3]
	}
	return ""
}
