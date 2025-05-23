---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sccfm_ios_device Data Source - sccfm"
subcategory: ""
description: |-
  IOS data source
---

# sccfm_ios_device (Data Source)

IOS data source



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The human-readable name of the device. This is the name displayed on the CDO Inventory page. Device names are unique across a CDO tenant.

### Read-Only

- `connector_name` (String) The name of the Secure Device Connector (SDC) that is used by CDO to communicate with the device.
- `grouped_labels` (Map of Set of String) The grouped labels applied to the device. Labels are used to group devices in CDO. Refer to the [SCC Firewall Manager documentation](https://docs.manage.security.cisco.com/t-applying-labels-to-devices-and-objects.html#!c-labels-and-filtering.html) for details on how labels are used in CDO.
- `host` (String) The host used to connect to the device.
- `id` (String) Universally unique identifier of the device.
- `ignore_certificate` (Boolean) This attribute indicates whether certificates were ignored when onboarding this device.
- `labels` (List of String) The labels applied to the device. Labels are used to group devices in CDO. Refer to the [SCC Firewall Manager documentation](https://docs.manage.security.cisco.com/t-applying-labels-to-devices-and-objects.html#!c-labels-and-filtering.html) for details on how labels are used in CDO.
- `port` (Number) The port used to connect to the device.
- `socket_address` (String) The address of the device to onboard, specified in the format `host:port`.
