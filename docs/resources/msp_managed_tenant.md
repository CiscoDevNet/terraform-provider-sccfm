---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sccfm_msp_managed_tenant Resource - sccfm"
subcategory: ""
description: |-
  Provides an MSP managed tenant resource. This allows MSP managed tenants to be created. Note: deleting this resource removes the created tenant from the MSP portal by disassociating the tenant from the MSP portal, but the tenant will continue to exist. To completely delete a tenant, please contact Cisco TAC.
---

# sccfm_msp_managed_tenant (Resource)

Provides an MSP managed tenant resource. This allows MSP managed tenants to be created. Note: deleting this resource removes the created tenant from the MSP portal by disassociating the tenant from the MSP portal, but the tenant will continue to exist. To completely delete a tenant, please contact Cisco TAC.



<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `api_token` (String, Sensitive) API token for an API-only user with super-admin privileges on the tenant. This should be specified only when adding an existing tenant to the MSP portal, and should not be provided if a new tenant is being created (i.e., the `name` and/or `display_name` attributes are specified).
- `display_name` (String) Display name of the tenant. If no display name is specified, the display name will be set to the tenant name. This should be specified only if a new tenant is being created, and should not be provided if an existing tenant is being added to the MSP protal (i.e., the `api_token` attribute is specified).
- `name` (String) Name of the tenant. This should be specified only if a new tenant is being created, and should not be provided if an existing tenant is being added to the MSP protal (i.e., the `api_token` attribute is specified).

### Read-Only

- `generated_name` (String) Actual name of the tenant returned by the API. This auto-generated name will differ from the name entered by the customer.
- `id` (String) Universally unique identifier of the tenant
- `region` (String) SCC Firewall region in which the tenant is created. This is the same region as the region of the MSP portal.
