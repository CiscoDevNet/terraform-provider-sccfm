package servicegroup

import (
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/object"
	objectresource "github.com/CiscoDevnet/terraform-provider-sccfm/internal/object"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var dsServiceValueSchema = schema.NestedAttributeObject{
	Attributes: map[string]schema.Attribute{
		"protocol": schema.StringAttribute{
			MarkdownDescription: "The protocol (e.g., `TCP`, `UDP`, `ICMP`).",
			Computed:            true,
		},
		"value": schema.StringAttribute{
			MarkdownDescription: "The port value (e.g., `80`, `8000-8080`).",
			Computed:            true,
		},
	},
}

func NewDataSource() datasource.DataSource {
	return objectresource.NewObjectDataSource(objectresource.ObjectDataSourceConfig[ResourceModel]{
		TypeNameSuffix:      "_service_group",
		MarkdownDescription: "Use this data source to look up an existing service group by name.",
		ObjectType:          object.ServiceGroup,
		ExtraSchemaAttributes: func() map[string]schema.Attribute {
			return map[string]schema.Attribute{
				"values": schema.ListNestedAttribute{
					MarkdownDescription: "Inline service values included in the group.",
					Computed:            true,
					NestedObject:        dsServiceValueSchema,
				},
				"referenced_object_uids": schema.SetAttribute{
					MarkdownDescription: "Set of UIDs of service objects referenced by this group.",
					Computed:            true,
					ElementType:         types.StringType,
				},
				"overrides": schema.ListNestedAttribute{
					MarkdownDescription: "List of per-target overrides.",
					Computed:            true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"target_id": schema.StringAttribute{
								Computed: true,
							},
							"values": schema.ListNestedAttribute{
								Computed:     true,
								NestedObject: dsServiceValueSchema,
							},
							"referenced_object_uids": schema.SetAttribute{
								Computed:    true,
								ElementType: types.StringType,
							},
						},
					},
				},
			}
		},
		MapReadOutput: mapReadOutput,
	})
}
