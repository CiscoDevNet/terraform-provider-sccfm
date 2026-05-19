package networkgroup

import (
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/object"
	objectresource "github.com/CiscoDevnet/terraform-provider-sccfm/internal/object"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewDataSource() datasource.DataSource {
	return objectresource.NewObjectDataSource(objectresource.ObjectDataSourceConfig[ResourceModel]{
		TypeNameSuffix:      "_network_group",
		MarkdownDescription: "Use this data source to look up an existing network group by name.",
		ObjectType:          object.NetworkGroup,
		ExtraSchemaAttributes: func() map[string]schema.Attribute {
			return map[string]schema.Attribute{
				"values": schema.ListAttribute{
					MarkdownDescription: "Inline network values included in the group.",
					Computed:            true,
					ElementType:         types.StringType,
				},
				"referenced_object_uids": schema.SetAttribute{
					MarkdownDescription: "Set of UIDs of network objects referenced by this group.",
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
							"values": schema.ListAttribute{
								Computed:    true,
								ElementType: types.StringType,
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
