package urlobject

import (
	"context"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/object"
	objectresource "github.com/CiscoDevnet/terraform-provider-sccfm/internal/object"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// OverrideModel represents a per-target override for a URL object.
type OverrideModel struct {
	TargetId types.String `tfsdk:"target_id"`
	Url      types.String `tfsdk:"url"`
}

// ResourceModel is the Terraform state model for a URL object.
type ResourceModel struct {
	Id          types.String    `tfsdk:"id"`
	Name        types.String    `tfsdk:"name"`
	Description types.String    `tfsdk:"description"`
	Url         types.String    `tfsdk:"url"`
	Overrides   []OverrideModel `tfsdk:"overrides"`
}

func NewResource() resource.Resource {
	return objectresource.NewObjectResource(objectresource.ObjectResourceConfig[ResourceModel]{
		TypeNameSuffix:      "_url_object",
		MarkdownDescription: "Provides a URL Object resource. Use this resource to create, read, update, and delete URL objects in the Security Cloud Control unified object store.",
		ExtraSchemaAttributes: func() map[string]schema.Attribute {
			return map[string]schema.Attribute{
				"url": schema.StringAttribute{
					MarkdownDescription: "The URL value of the object (e.g., `https://www.example.com`).",
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
							"url": schema.StringAttribute{
								MarkdownDescription: "The override URL value for this target.",
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
		content, err := object.MarshalContent(object.UrlContent{
			Url: o.Url.ValueString(),
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
	content, err := object.MarshalContent(object.UrlContent{
		Url: model.Url.ValueString(),
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
			ObjectType:     object.UrlObject,
			DefaultContent: content,
			Overrides:      overrides,
		},
	}, nil
}

func buildUpdateInput(plan *ResourceModel, state *ResourceModel) (object.UpdateInput, error) {
	content, err := object.MarshalContent(object.UrlContent{
		Url: plan.Url.ValueString(),
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
			ObjectType:     object.UrlObject,
			DefaultContent: content,
			Overrides:      overrides,
		},
	}, nil
}

func mapReadOutput(_ context.Context, output *object.ReadOutput, model *ResourceModel) error {
	model.Id = types.StringValue(output.Uid)
	model.Name = types.StringValue(output.Name)
	model.Description = types.StringValue(output.Description)

	uc, err := object.UnmarshalContent[object.UrlContent](output.Value.DefaultContent)
	if err != nil {
		return err
	}
	if uc != nil {
		model.Url = types.StringValue(uc.Url)
	}

	model.Overrides = nil
	for _, o := range output.Value.Overrides {
		oc, err := object.UnmarshalContent[object.UrlContent](o.Content)
		if err != nil {
			return err
		}
		if oc != nil {
			model.Overrides = append(model.Overrides, OverrideModel{
				TargetId: types.StringValue(o.TargetId),
				Url:      types.StringValue(oc.Url),
			})
		}
	}
	return nil
}
