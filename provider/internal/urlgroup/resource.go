package urlgroup

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/object"
	objectresource "github.com/CiscoDevnet/terraform-provider-sccfm/internal/object"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// OverrideModel represents a per-target override for a URL group.
type OverrideModel struct {
	TargetId             types.String   `tfsdk:"target_id"`
	Values               []types.String `tfsdk:"values"`
	ReferencedObjectUids []types.String `tfsdk:"referenced_object_uids"`
}

// ResourceModel is the Terraform state model for a URL group.
type ResourceModel struct {
	Id                   types.String    `tfsdk:"id"`
	Name                 types.String    `tfsdk:"name"`
	Description          types.String    `tfsdk:"description"`
	Values               types.List      `tfsdk:"values"`
	ReferencedObjectUids types.Set       `tfsdk:"referenced_object_uids"`
	Overrides            []OverrideModel `tfsdk:"overrides"`
}

func NewResource() resource.Resource {
	return objectresource.NewObjectResource(objectresource.ObjectResourceConfig[ResourceModel]{
		TypeNameSuffix:      "_url_group",
		MarkdownDescription: "Provides a URL Group resource. Use this resource to create, read, update, and delete URL groups in the Security Cloud Control unified object store. A URL group can contain inline URL values and/or references to existing URL objects.",
		ExtraSchemaAttributes: func() map[string]schema.Attribute {
			return map[string]schema.Attribute{
				"values": schema.ListAttribute{
					MarkdownDescription: "Inline URL values included in the group.",
					Optional:            true,
					ElementType:         types.StringType,
				},
				"referenced_object_uids": schema.SetAttribute{
					MarkdownDescription: "Set of UIDs of URL objects referenced by this group.",
					Optional:            true,
					ElementType:         types.StringType,
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
							"values": schema.ListAttribute{
								MarkdownDescription: "Inline URL values for this target override.",
								Optional:            true,
								ElementType:         types.StringType,
							},
							"referenced_object_uids": schema.SetAttribute{
								MarkdownDescription: "Set of UIDs of URL objects for this target override.",
								Optional:            true,
								ElementType:         types.StringType,
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

func marshalUrlLiterals(values []types.String) ([]json.RawMessage, error) {
	var literals []json.RawMessage
	for _, v := range values {
		raw, err := json.Marshal(object.UrlContent{Url: v.ValueString()})
		if err != nil {
			return nil, err
		}
		literals = append(literals, json.RawMessage(raw))
	}
	return literals, nil
}

func buildGroupContent(model *ResourceModel) (*object.GroupContent, error) {
	gc := &object.GroupContent{}

	if !model.Values.IsNull() && !model.Values.IsUnknown() {
		var vals []types.String
		for _, v := range model.Values.Elements() {
			vals = append(vals, v.(types.String))
		}
		literals, err := marshalUrlLiterals(vals)
		if err != nil {
			return nil, err
		}
		gc.Literals = literals
	}

	if !model.ReferencedObjectUids.IsNull() && !model.ReferencedObjectUids.IsUnknown() {
		for _, v := range model.ReferencedObjectUids.Elements() {
			gc.ReferencedObjectUids = append(gc.ReferencedObjectUids, v.(types.String).ValueString())
		}
	}

	return gc, nil
}

func buildOverrides(overrides []OverrideModel) ([]object.Override, error) {
	if len(overrides) == 0 {
		return nil, nil
	}
	result := make([]object.Override, 0, len(overrides))
	for _, o := range overrides {
		gc := &object.GroupContent{}
		literals, err := marshalUrlLiterals(o.Values)
		if err != nil {
			return nil, err
		}
		gc.Literals = literals
		for _, uid := range o.ReferencedObjectUids {
			gc.ReferencedObjectUids = append(gc.ReferencedObjectUids, uid.ValueString())
		}
		content, err := object.MarshalContent(gc)
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
	gc, err := buildGroupContent(model)
	if err != nil {
		return object.CreateInput{}, err
	}
	content, err := object.MarshalContent(gc)
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
			ObjectType:     object.UrlGroup,
			DefaultContent: content,
			Overrides:      overrides,
		},
	}, nil
}

func buildUpdateInput(plan *ResourceModel, state *ResourceModel) (object.UpdateInput, error) {
	gc, err := buildGroupContent(plan)
	if err != nil {
		return object.UpdateInput{}, err
	}
	content, err := object.MarshalContent(gc)
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
			ObjectType:     object.UrlGroup,
			DefaultContent: content,
			Overrides:      overrides,
		},
	}, nil
}

func unmarshalUrlLiterals(literals []json.RawMessage) ([]string, error) {
	var values []string
	for _, raw := range literals {
		var uc object.UrlContent
		if err := json.Unmarshal(raw, &uc); err != nil {
			return nil, err
		}
		values = append(values, uc.Url)
	}
	return values, nil
}

func mapReadOutput(ctx context.Context, output *object.ReadOutput, model *ResourceModel) error {
	model.Id = types.StringValue(output.Uid)
	model.Name = types.StringValue(output.Name)
	model.Description = types.StringValue(output.Description)

	gc, err := object.UnmarshalContent[object.GroupContent](output.Value.DefaultContent)
	if err != nil {
		return err
	}
	if gc == nil {
		model.Values = types.ListNull(types.StringType)
		model.ReferencedObjectUids = types.SetNull(types.StringType)
		return nil
	}

	values, err := unmarshalUrlLiterals(gc.Literals)
	if err != nil {
		return err
	}
	if len(values) > 0 {
		listVal, diags := types.ListValueFrom(ctx, types.StringType, values)
		if diags.HasError() {
			return fmt.Errorf("failed to convert values: %s", diags.Errors())
		}
		model.Values = listVal
	} else {
		model.Values = types.ListNull(types.StringType)
	}

	if len(gc.ReferencedObjectUids) > 0 {
		setVal, diags := types.SetValueFrom(ctx, types.StringType, gc.ReferencedObjectUids)
		if diags.HasError() {
			return fmt.Errorf("failed to convert referenced_object_uids: %s", diags.Errors())
		}
		model.ReferencedObjectUids = setVal
	} else {
		model.ReferencedObjectUids = types.SetNull(types.StringType)
	}

	model.Overrides = nil
	for _, o := range output.Value.Overrides {
		ogc, err := object.UnmarshalContent[object.GroupContent](o.Content)
		if err != nil {
			return err
		}
		if ogc == nil {
			continue
		}
		om := OverrideModel{
			TargetId: types.StringValue(o.TargetId),
		}
		oValues, err := unmarshalUrlLiterals(ogc.Literals)
		if err != nil {
			return err
		}
		for _, v := range oValues {
			om.Values = append(om.Values, types.StringValue(v))
		}
		for _, uid := range ogc.ReferencedObjectUids {
			om.ReferencedObjectUids = append(om.ReferencedObjectUids, types.StringValue(uid))
		}
		model.Overrides = append(model.Overrides, om)
	}

	return nil
}
