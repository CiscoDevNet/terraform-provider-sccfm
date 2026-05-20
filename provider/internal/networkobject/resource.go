package networkobject

import (
	"context"
	"net"
	"strings"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/object"
	objectresource "github.com/CiscoDevnet/terraform-provider-sccfm/internal/object"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// normalizeRangeLiteral collapses any surrounding whitespace around the dash
// in an IP range so state matches a user's natural "a-b" input. The server
// canonicalizes ranges to "a - b" with spaces; this strips that back out.
// Non-range literals (single IPs, CIDR) are returned as-is.
func normalizeRangeLiteral(s string) string {
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return s
	}
	left, right := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
	if net.ParseIP(left) == nil || net.ParseIP(right) == nil {
		return s
	}
	return left + "-" + right
}

// OverrideModel represents a per-target override for a network object.
type OverrideModel struct {
	TargetId types.String `tfsdk:"target_id"`
	Value    types.String `tfsdk:"value"`
}

// ResourceModel is the Terraform state model for a network object.
type ResourceModel struct {
	Id          types.String    `tfsdk:"id"`
	Name        types.String    `tfsdk:"name"`
	Description types.String    `tfsdk:"description"`
	Value       types.String    `tfsdk:"value"`
	Overrides   []OverrideModel `tfsdk:"overrides"`
}

func NewResource() resource.Resource {
	return objectresource.NewObjectResource(objectresource.ObjectResourceConfig[ResourceModel]{
		TypeNameSuffix:      "_network_object",
		MarkdownDescription: "Provides a Network Object resource. Use this resource to create, read, update, and delete network objects in the Security Cloud Control unified object store.",
		ExtraSchemaAttributes: func() map[string]schema.Attribute {
			return map[string]schema.Attribute{
				"value": schema.StringAttribute{
					MarkdownDescription: "The value of the network object. This can be an IP address (e.g., `10.0.0.1`), a CIDR block (e.g., `10.0.0.0/24`), an IPv6 address (e.g., `a:b:c::1`), or a range (e.g., `10.0.0.1-10.0.0.10`).",
					Required:            true,
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
							"value": schema.StringAttribute{
								MarkdownDescription: "The override network value for this target.",
								Required:            true,
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

func buildOverrides(overrides []OverrideModel) ([]object.Override, error) {
	if len(overrides) == 0 {
		return nil, nil
	}
	result := make([]object.Override, 0, len(overrides))
	for _, o := range overrides {
		content, err := object.MarshalContent(object.NetworkContent{
			Literal: o.Value.ValueString(),
		})
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
	content, err := object.MarshalContent(object.NetworkContent{
		Literal: model.Value.ValueString(),
	})
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
			ObjectType:     object.NetworkObject,
			DefaultContent: content,
			Overrides:      overrides,
		},
	}, nil
}

func buildUpdateInput(plan *ResourceModel, state *ResourceModel) (object.UpdateInput, error) {
	content, err := object.MarshalContent(object.NetworkContent{
		Literal: plan.Value.ValueString(),
	})
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
			ObjectType:     object.NetworkObject,
			DefaultContent: content,
			Overrides:      overrides,
		},
	}, nil
}

func mapReadOutput(_ context.Context, output *object.ReadOutput, model *ResourceModel) error {
	model.Id = types.StringValue(output.Uid)
	model.Name = types.StringValue(output.Name)
	model.Description = types.StringValue(output.Description)

	nc, err := object.UnmarshalContent[object.NetworkContent](output.Value.DefaultContent)
	if err != nil {
		return err
	}
	if nc != nil {
		model.Value = types.StringValue(normalizeRangeLiteral(nc.Literal))
	}

	model.Overrides = nil
	for _, o := range output.Value.Overrides {
		oc, err := object.UnmarshalContent[object.NetworkContent](o.Content)
		if err != nil {
			return err
		}
		if oc != nil {
			model.Overrides = append(model.Overrides, OverrideModel{
				TargetId: types.StringValue(o.TargetId),
				Value:    types.StringValue(normalizeRangeLiteral(oc.Literal)),
			})
		}
	}
	return nil
}
