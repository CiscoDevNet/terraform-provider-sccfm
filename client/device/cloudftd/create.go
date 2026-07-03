package cloudftd

/**
* The cloud FTD corresponds to the CDO Device Type FTDC, which is an FTD managed by a cdFMC.
 */
import (
	"context"
	"fmt"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/device/cloudfmc"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/device/cloudfmc/fmcplatform"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/cdo"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/publicapi"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/url"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/cloudfmc/accesspolicies"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/device/publicapilabels"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/devicetype"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/ftd/license"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/ftd/tier"
)

type CreateInput struct {
	Name             string
	AccessPolicyName string
	PerformanceTier  *tier.Type // ignored if it is physical device
	Virtual          bool
	Licenses         *[]license.Type
	Labels           publicapilabels.Type
}

type CreateOutput struct {
	Uid      string               `json:"uid"`
	Name     string               `json:"name"`
	Metadata Metadata             `json:"metadata,omitempty"`
	State    string               `json:"state"`
	Labels   publicapilabels.Type `json:"labels"`
}

func FromDeviceReadOutput(readOutput *ReadOutput) *CreateOutput {
	if readOutput == nil {
		return nil
	}

	return &CreateOutput{
		Uid:      readOutput.Uid,
		Name:     readOutput.Name,
		Metadata: readOutput.Metadata,
		State:    readOutput.State,
		Labels:   publicapilabels.New(readOutput.Tags.UngroupedTags(), readOutput.Tags.GroupedTags()),
	}
}

func NewCreateInput(
	name string,
	accessPolicyName string,
	performanceTier *tier.Type,
	virtual bool,
	licenses *[]license.Type,
	labels publicapilabels.Type,
) CreateInput {
	return CreateInput{
		Name:             name,
		AccessPolicyName: accessPolicyName,
		PerformanceTier:  performanceTier,
		Virtual:          virtual,
		Licenses:         licenses,
		Labels:           labels,
	}
}

type createRequestBody struct {
	Name               string               `json:"name"`
	DeviceType         devicetype.Type      `json:"deviceType"`
	FmcAccessPolicyUid string               `json:"fmcAccessPolicyUid"`
	PerformanceTier    *tier.Type           `json:"performanceTier"`
	Virtual            bool                 `json:"virtual"`
	Licenses           *[]license.Type      `json:"licenses"`
	Labels             publicapilabels.Type `json:"labels"`
}

func Create(ctx context.Context, client http.Client, createInp CreateInput) (*CreateOutput, error) {

	client.Logger.Println("creating cloud ftd")

	createUrl := url.CreateFtd(client.BaseUrl())

	selectedPolicy, err := readPolicyUidFromPolicyName(ctx, client, createInp.AccessPolicyName)
	if err != nil {
		return nil, err
	}

	transaction, err := publicapi.TriggerTransaction(
		ctx,
		client,
		createUrl,
		createRequestBody{
			DeviceType:         "CDFMC_MANAGED_FTD",
			Name:               createInp.Name,
			FmcAccessPolicyUid: selectedPolicy.Id,
			PerformanceTier:    createInp.PerformanceTier,
			Virtual:            createInp.Virtual,
			Labels:             createInp.Labels,
			Licenses:           createInp.Licenses,
		},
	)
	if err != nil {
		return nil, err
	}
	createdFtdUid := transaction.EntityUid

	transaction, err = publicapi.WaitForTransactionToFinishWithDefaults(
		ctx,
		client,
		transaction,
		"Waiting for Cloud FTD to onboard...",
	)
	if err != nil {
		_ = deleteUsingPublicApi(ctx, client, createdFtdUid)
		return nil, err
	}

	cloudFtdReadOutput, err := ReadByUidUsingPublicApi(ctx, client, NewReadByUidInput(transaction.EntityUid))
	if err != nil {
		return nil, err
	}
	// The public inventory model exposes the access policy UID but not its
	// display name. Preserve the name supplied by the Terraform configuration.
	cloudFtdReadOutput.Metadata.AccessPolicyName = createInp.AccessPolicyName

	return FromDeviceReadOutput(cloudFtdReadOutput), nil
}

func deleteUsingPublicApi(ctx context.Context, client http.Client, uid string) error {
	if uid == "" {
		return nil
	}

	transaction, err := publicapi.TriggerTransaction(ctx, client, url.DeleteFtdUsingPublicApi(client.BaseUrl(), uid), nil)
	if err != nil {
		return err
	}

	_, err = publicapi.WaitForTransactionToFinishWithDefaults(ctx, client, transaction, "Waiting for Cloud FTD cleanup...")
	return err
}

func readPolicyUidFromPolicyName(ctx context.Context, client http.Client, accessPolicyName string) (accesspolicies.Item, error) {
	readFmcDomainRes, err := fmcplatform.ReadFmcDomainInfoUsingPublicApi(ctx, client)
	if err != nil {
		return accesspolicies.Item{}, err
	}
	if len(readFmcDomainRes.Items) == 0 {
		return accesspolicies.Item{}, fmt.Errorf("%w: fmc domain info not found", http.NotFoundError)
	}

	accessPoliciesRes, err := cloudfmc.ReadAccessPoliciesUsingPublicApi(
		ctx,
		client,
		cloudfmc.NewReadAccessPoliciesUsingPublicApiInput(readFmcDomainRes.Items[0].Uuid, 1000), // 1000 is what CDO UI uses
	)
	if err != nil {
		return accesspolicies.Item{}, err
	}

	selectedPolicy, ok := accessPoliciesRes.Find(accessPolicyName)
	if !ok {
		return accesspolicies.Item{}, fmt.Errorf(
			`access policy: "%s" not found, available policies: %s. In rare cases where you have more than 1000 access policies, please raise an issue at: %s`,
			accessPolicyName,
			accessPoliciesRes.Items,
			cdo.TerraformProviderCDOIssuesUrl,
		)
	}
	return selectedPolicy, nil
}
