package networkobject

import (
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/object"
	objectresource "github.com/CiscoDevnet/terraform-provider-sccfm/internal/object"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func NewDataSource() datasource.DataSource {
	return objectresource.NewObjectDataSource(objectresource.ObjectDataSourceConfig[ResourceModel]{
		TypeNameSuffix:      "_network_object",
		MarkdownDescription: "Use this data source to look up an existing network object by name.",
		ObjectType:          object.NetworkObject,
		ExtraSchemaAttributes: func() map[string]schema.Attribute {
			return map[string]schema.Attribute{
				"value": schema.StringAttribute{
					MarkdownDescription: "The value of the network object.",
					Computed:            true,
				},
				"overrides": schema.ListNestedAttribute{
					MarkdownDescription: "List of per-target overrides.",
					Computed:            true,
					NestedObject: schema.NestedAttributeObject{
						Attributes: map[string]schema.Attribute{
							"target_id": schema.StringAttribute{
								Computed: true,
							},
							"value": schema.StringAttribute{
								Computed: true,
							},
						},
					},
				},
			}
		},
		MapReadOutput: mapReadOutput,
	})
}
