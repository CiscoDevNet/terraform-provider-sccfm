package serviceobject

import (
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/object"
	objectresource "github.com/CiscoDevnet/terraform-provider-sccfm/internal/object"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func NewDataSource() datasource.DataSource {
	return objectresource.NewObjectDataSource(objectresource.ObjectDataSourceConfig[ResourceModel]{
		TypeNameSuffix:      "_service_object",
		MarkdownDescription: "Use this data source to look up an existing service object by name.",
		ObjectType:          object.ServiceObject,
		ExtraSchemaAttributes: func() map[string]schema.Attribute {
			return map[string]schema.Attribute{
				"protocol": schema.StringAttribute{
					MarkdownDescription: "The protocol of the service object.",
					Computed:            true,
				},
				"value": schema.StringAttribute{
					MarkdownDescription: "The port value of the service object.",
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
							"protocol": schema.StringAttribute{
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
