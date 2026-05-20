package serviceobject

import (
	"context"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/object"
	objectresource "github.com/CiscoDevnet/terraform-provider-sccfm/internal/object"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// OverrideModel represents a per-target override for a service object.
type OverrideModel struct {
	TargetId types.String `tfsdk:"target_id"`
	Protocol types.String `tfsdk:"protocol"`
	Value    types.String `tfsdk:"value"`
}

// ResourceModel is the Terraform state model for a service object.
type ResourceModel struct {
	Id          types.String    `tfsdk:"id"`
	Name        types.String    `tfsdk:"name"`
	Description types.String    `tfsdk:"description"`
	Protocol    types.String    `tfsdk:"protocol"`
	Value       types.String    `tfsdk:"value"`
	Overrides   []OverrideModel `tfsdk:"overrides"`
}

func NewResource() resource.Resource {
	return objectresource.NewObjectResource(objectresource.ObjectResourceConfig[ResourceModel]{
		TypeNameSuffix:      "_service_object",
		MarkdownDescription: "Provides a Service Object resource. Use this resource to create, read, update, and delete service objects in the Security Cloud Control unified object store.",
		ExtraSchemaAttributes: func() map[string]schema.Attribute {
			return map[string]schema.Attribute{
				"protocol": schema.StringAttribute{
					MarkdownDescription: "The protocol of the service object (e.g., `TCP`, `UDP`, `ICMP`).",
					Required:            true,
				},
				"value": schema.StringAttribute{
					MarkdownDescription: "The port value of the service object (e.g., `80`, `8000-8080`). Required for TCP/UDP protocols.",
					Optional:            true,
					Computed:            true,
				},
				"overrides": schema.ListNestedAttribute{
					MarkdownDescription: "List of per-target overrides. Each override replaces the default content for its target.",
					Optional:            true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"target_id": schema.StringAttribute{
								MarkdownDescription: "The ID of the target (device, service, or shared policy).",
								Required:            true,
							},
							"protocol": schema.StringAttribute{
								MarkdownDescription: "The override protocol for this target.",
								Required:            true,
							},
							"value": schema.StringAttribute{
								MarkdownDescription: "The override port value for this target.",
								Optional:            true,
							},
						},
					},
				},
			}
		},
		BuildCreateInput: buildCreateInput,
		BuildUpdateInput: buildUpdateInput,
		MapReadOutput:    mapReadOutput,
		GetId:            func(m *ResourceModel) string { return m.Id.ValueString() },
	})
}

func buildServiceContent(protocol string, value string) *object.ServiceContent {
	sc := &object.ServiceContent{
		Protocol: protocol,
	}
	if value != "" {
		sc.ServiceValue = &object.ServiceValueContent{
			Literal: value,
		}
	}
	return sc
}

func buildOverrides(overrides []OverrideModel) ([]object.Override, error) {
	if len(overrides) == 0 {
		return nil, nil
	}
	result := make([]object.Override, 0, len(overrides))
	for _, o := range overrides {
		sc := buildServiceContent(o.Protocol.ValueString(), o.Value.ValueString())
		content, err := object.MarshalContent(sc)
		if err != nil {
			return nil, err
		}
		result = append(result, object.Override{
			TargetId: o.TargetId.ValueString(),
			Content:  content,
		})
	}
	return result, nil
}

func buildCreateInput(model *ResourceModel) (object.CreateInput, error) {
	sc := buildServiceContent(model.Protocol.ValueString(), model.Value.ValueString())
	content, err := object.MarshalContent(sc)
	if err != nil {
		return object.CreateInput{}, err
	}
	overrides, err := buildOverrides(model.Overrides)
	if err != nil {
		return object.CreateInput{}, err
	}
	return object.CreateInput{
		Name:        model.Name.ValueString(),
		Description: model.Description.ValueString(),
		Value: object.SharedObjectValue{
			ObjectType:     object.ServiceObject,
			DefaultContent: content,
			Overrides:      overrides,
		},
	}, nil
}

func buildUpdateInput(plan *ResourceModel, state *ResourceModel) (object.UpdateInput, error) {
	sc := buildServiceContent(plan.Protocol.ValueString(), plan.Value.ValueString())
	content, err := object.MarshalContent(sc)
	if err != nil {
		return object.UpdateInput{}, err
	}
	overrides, err := buildOverrides(plan.Overrides)
	if err != nil {
		return object.UpdateInput{}, err
	}
	description := plan.Description.ValueString()
	return object.UpdateInput{
		Uid:         state.Id.ValueString(),
		Name:        plan.Name.ValueString(),
		Description: &description,
		Value: &object.SharedObjectValue{
			ObjectType:     object.ServiceObject,
			DefaultContent: content,
			Overrides:      overrides,
		},
	}, nil
}

func mapReadOutput(_ context.Context, output *object.ReadOutput, model *ResourceModel) error {
	model.Id = types.StringValue(output.Uid)
	model.Name = types.StringValue(output.Name)
	model.Description = types.StringValue(output.Description)

	sc, err := object.UnmarshalContent[object.ServiceContent](output.Value.DefaultContent)
	if err != nil {
		return err
	}
	if sc != nil {
		model.Protocol = types.StringValue(sc.Protocol)
		if sc.ServiceValue != nil {
			model.Value = types.StringValue(sc.ServiceValue.Literal)
		} else {
			model.Value = types.StringNull()
		}
	}

	model.Overrides = nil
	for _, o := range output.Value.Overrides {
		oc, err := object.UnmarshalContent[object.ServiceContent](o.Content)
		if err != nil {
			return err
		}
		if oc != nil {
			om := OverrideModel{
				TargetId: types.StringValue(o.TargetId),
				Protocol: types.StringValue(oc.Protocol),
			}
			if oc.ServiceValue != nil {
				om.Value = types.StringValue(oc.ServiceValue.Literal)
			}
			model.Overrides = append(model.Overrides, om)
		}
	}
	return nil
}
