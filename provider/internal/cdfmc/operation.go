package cdfmc

import (
	"context"
	"errors"

	sccFwMgrClient "github.com/CiscoDevnet/terraform-provider-sccfm/go-client"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/device/cloudfmc"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func Read(ctx context.Context, resource *Resource, stateData *ResourceModel) error {

	// do read
	readOut, err := resource.client.ReadCloudFmcDevice(ctx)
	if err != nil {
		return err
	}
	cloudFmcSpecificDeviceReadOut, err := resource.client.ReadCloudFmcSpecificDevice(ctx, cloudfmc.NewReadSpecificInput(readOut.Uid))
	if err != nil {
		return err
	}

	// map response to terraform types
	stateData.Id = types.StringValue(readOut.Uid)
	stateData.Name = types.StringValue(readOut.Name)
	stateData.Hostname = types.StringValue(readOut.Host)
	stateData.DomainUuid = types.StringValue(cloudFmcSpecificDeviceReadOut.DomainUid)

	return nil
}

func Create(ctx context.Context, resource *Resource, planData *ResourceModel) error {

	// try reading existing cdFMC first — there can only be one per tenant
	readOut, err := resource.client.ReadCloudFmcDevice(ctx)
	if err == nil {
		tflog.Info(ctx, "cdFMC already exists, adopting into state")
		cloudFmcSpecificDeviceReadOut, readSpecificErr := resource.client.ReadCloudFmcSpecificDevice(ctx, cloudfmc.NewReadSpecificInput(readOut.Uid))
		if readSpecificErr != nil {
			return readSpecificErr
		}
		planData.Id = types.StringValue(readOut.Uid)
		planData.Name = types.StringValue(readOut.Name)
		planData.Hostname = types.StringValue(readOut.Host)
		planData.DomainUuid = types.StringValue(cloudFmcSpecificDeviceReadOut.DomainUid)
		return nil
	}

	// only create if the read returned not found
	if !errors.Is(err, sccFwMgrClient.NotFoundError) {
		return err
	}

	// do create
	createOut, err := resource.client.CreateCloudFmcDevice(ctx, cloudfmc.NewCreateInput())
	if err != nil {
		return err
	}
	cloudFmcSpecificDeviceReadOut, err := resource.client.ReadCloudFmcSpecificDevice(ctx, cloudfmc.NewReadSpecificInput(createOut.Uid))
	if err != nil {
		return err
	}

	// map response to terraform types
	planData.Id = types.StringValue(createOut.Uid)
	planData.Name = types.StringValue(createOut.Name)
	planData.Hostname = types.StringValue(createOut.Host)
	planData.DomainUuid = types.StringValue(cloudFmcSpecificDeviceReadOut.DomainUid)

	return nil
}
