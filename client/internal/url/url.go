// central url management, all urls goes into here
package url

import (
	"fmt"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/devicetype"
)

func ReadDevice(baseUrl string, uid string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/devices/%s", baseUrl, uid)
}

func ReadDeviceByNameAndType(baseUrl string, deviceName string, deviceType devicetype.Type) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/devices?q=name:%s+AND+deviceType:%s", baseUrl, deviceName, deviceType)
}
func ReadAllDevicesByType(baseUrl string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/devices", baseUrl)
}

func CreateDevice(baseUrl string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/devices", baseUrl)
}

func UpdateDevice(baseUrl string, uid string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/devices/%s", baseUrl, uid)
}

func DeleteDevice(baseUrl string, uid string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/devices/%s", baseUrl, uid)
}

func ReadSpecificDevice(baseUrl string, specificUid string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/device/%s/specific-device", baseUrl, specificUid)
}

func ReadAllConnectors(baseUrl string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/proxies", baseUrl)
}

func ReadAsaConfig(baseUrl string, specificUid string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/asa/configs/%s", baseUrl, specificUid)
}

func UpdateAsaConfig(baseUrl string, specificUid string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/asa/configs/%s", baseUrl, specificUid)
}

func ReadConnectorByUid(baseUrl string, connectorUid string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/proxies/%s", baseUrl, connectorUid)
}

func ReadConnectorByName(baseUrl string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/proxies", baseUrl)
}

func CreateConnector(baseUrl string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/proxies", baseUrl)
}

func UpdateConnector(baseUrl string, connectorUid string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/proxies/%s", baseUrl, connectorUid)
}

func DeleteConnector(baseUrl string, connectorUid string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/proxies/%s", baseUrl, connectorUid)
}

func UserToken(baseUrl string) string {
	return fmt.Sprintf("%s/anubis/rest/v1/oauth/token", baseUrl)
}

func ExternalComputeToken(baseUrl string) string {
	return fmt.Sprintf("%s/anubis/rest/v1/oauth/token/external-compute", baseUrl)
}

func ReadSmartLicense(baseUrl string) string {
	return fmt.Sprintf("%s/fmc/api/fmc_platform/v1/license/smartlicenses", baseUrl)
}

func ReadAccessPolicies(baseUrl string, domainUid string) string {
	return fmt.Sprintf("%s/fmc/api/fmc_config/v1/domain/%s/policy/accesspolicies", baseUrl, domainUid)
}

func UpdateSpecificCloudFtd(baseUrl string, ftdSpecificUid string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/firepower/ftds/%s", baseUrl, ftdSpecificUid)
}

func UpdateFmcAppliance(baseUrl string, fmcApplianceUid string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/fmc/appliance/%s", baseUrl, fmcApplianceUid)
}

func ReadStateMachineInstance(baseUrl string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/state-machines/instances", baseUrl)
}

func ReadFmcDomainInfo(fmcHost string) string {
	return fmt.Sprintf("https://%s/api/fmc_platform/v1/info/domain", fmcHost)
}

func ReadFmcDeviceLicenses(baseUrl string) string {
	return fmt.Sprintf("%s/fmc/api/fmc_platform/v1/license/devicelicenses", baseUrl)
}

func UpdateFmcDeviceLicenses(baseUrl string, objectId string) string {
	return fmt.Sprintf("%s/fmc/api/fmc_platform/v1/license/devicelicenses/%s", baseUrl, objectId)
}

func CreateUser(baseUrl string) string {
	return fmt.Sprintf("%s/api/rest/v1/users", baseUrl)
}

func ReadUserByUsername(baseUrl string) string {
	return fmt.Sprintf("%s/api/rest/v1/users", baseUrl)
}

func ReadApiOnlyUserByUsername(baseUrl string) string {
	return fmt.Sprintf("%s/api/rest/v1/users/api-only", baseUrl)
}

func ReadUserByUid(baseUrl string, uid string) string {
	return fmt.Sprintf("%s/api/rest/v1/users/%s", baseUrl, uid)
}

func GenerateApiToken(baseUrl string, userUid string) string {
	return fmt.Sprintf("%s/api/rest/v1/users/%s/apiToken/generate", baseUrl, userUid)
}

func RevokeApiTokenForUser(baseUrl string, userUid string) string {
	return fmt.Sprintf("%s/api/rest/v1/users/%s/apiToken/revoke", baseUrl, userUid)
}

func RevokeApiTokenUsingPublicApi(baseUrl string) string {
	return fmt.Sprintf("%s/api/rest/v1/token/revoke", baseUrl)
}

func ReadAnubisTokenInfo(baseUrl string) string {
	return fmt.Sprintf("%s/anubis/rest/v1/oauth/check_token", baseUrl)
}

func ReadTenantContext(baseUrl string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/common/tenantcontext", baseUrl)
}

func CreateSystemToken(baseUrl string, scope string) string {
	return fmt.Sprintf("%s/anubis/rest/v1/oauth/token/system/%s", baseUrl, scope)
}

func CreateFmcDeviceRecord(baseUrl string, fmcDomainId string) string {
	return fmt.Sprintf("%s/fmc/api/fmc_config/v1/domain/%s/devices/devicerecords", baseUrl, fmcDomainId)
}

func ReadFmcDeviceRecord(baseUrl string, fmcDomainId string, deviceUid string) string {
	return fmt.Sprintf("%s/fmc/api/fmc_config/v1/domain/%s/devices/devicerecords/%s", baseUrl, fmcDomainId, deviceUid)
}

func ReadFmcAllDeviceRecords(baseUrl string, fmcDomainId string) string {
	return fmt.Sprintf("%s/fmc/api/fmc_config/v1/domain/%s/devices/devicerecords", baseUrl, fmcDomainId)
}

func ReadFmcTaskStatus(baseUrl string, fmcDomainUid string, taskId string) string {
	return fmt.Sprintf("%s/fmc/api/fmc_config/v1/domain/%s/job/taskstatuses/%s", baseUrl, fmcDomainUid, taskId)
}

func ReadTenantDetails(baseUrl string) string {
	return fmt.Sprintf("%s/anubis/rest/v1/oauth/check_token", baseUrl)
}

func CreateApplication(baseUrl string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/applications", baseUrl)
}

func ReadApplication(baseUrl string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/applications", baseUrl)
}

func CreateSec(baseUrl string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/estreamers", baseUrl)
}

func ReadSec(baseUrl string, secUid string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/estreamers/%s", baseUrl, secUid)
}

func DeleteSec(baseUrl string, secUid string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/estreamers/%s", baseUrl, secUid)
}

func ReadAllSecs(baseUrl string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/estreamers", baseUrl)
}

func CreateDuoAdminPanel(baseUrl string) string {
	return fmt.Sprintf("%s/api/rest/v1/inventory/devices/duoAdminPanels", baseUrl)
}

func UpdateDuoAdminPanel(baseUrl string, uid string) string {
	return fmt.Sprintf("%s/aegis/rest/v1/services/targets/devices/%s", baseUrl, uid)
}

func CreateAsa(baseUrl string) string {
	return fmt.Sprintf("%s/api/rest/v1/inventory/devices/asas", baseUrl)
}

func CreateFtd(baseUrl string) string {
	return fmt.Sprintf("%s/api/rest/v1/inventory/devices/ftds", baseUrl)
}

func RegisterFtd(baseUrl string) string {
	return fmt.Sprintf("%s/api/rest/v1/inventory/devices/ftds/register", baseUrl)
}

func CreateIos(baseUrl string) string {
	return fmt.Sprintf("%s/api/rest/v1/inventory/devices/ios", baseUrl)
}

func CreateMspManagedTenant(baseUrl string) string {
	return fmt.Sprintf("%s/api/rest/v1/msp/tenants/create", baseUrl)
}

func AddExistingTenantToMspManagedTenant(baseUrl string) string {
	return fmt.Sprintf("%s/api/rest/v1/msp/tenants", baseUrl)
}

func MspManagedTenantByUid(baseUrl string, tenantUid string) string {
	return fmt.Sprintf("%s/api/rest/v1/msp/tenants/%s", baseUrl, tenantUid)
}

func FindMspManagedTenantsByName(baseUrl string, tenantName string) string {
	return fmt.Sprintf("%s/api/rest/v1/msp/tenants?q=name:%s", baseUrl, tenantName)
}

func CreateUsersInMspManagedTenant(baseUrl string, tenantUid string) string {
	return fmt.Sprintf("%s/api/rest/v1/msp/tenants/%s/users", baseUrl, tenantUid)
}

func GetUsersInMspManagedTenant(baseUrl string, tenantUid string, limit int, offset int) string {
	return fmt.Sprintf("%s/api/rest/v1/msp/tenants/%s/users?limit=%d&offset=%d", baseUrl, tenantUid, limit, offset)
}

func GetApiOnlyUsersInMspManagedTenant(baseUrl string, tenantUid string, limit int, offset int) string {
	return fmt.Sprintf("%s/api/rest/v1/msp/tenants/%s/users/api-only?limit=%d&offset=%d", baseUrl, tenantUid, limit, offset)
}

func DeleteUsersInMspManagedTenant(baseUrl string, tenantUid string) string {
	return fmt.Sprintf("%s/api/rest/v1/msp/tenants/%s/users/delete", baseUrl, tenantUid)
}

func GenerateApiTokenForUserInMspManagedTenant(baseUrl string, tenantUid string, userUid string) string {
	return fmt.Sprintf("%s/api/rest/v1/msp/tenants/%s/users/%s/token", baseUrl, tenantUid, userUid)
}

func CreateUserGroupsInMspManagedTenant(baseUrl string, tenantUid string) string {
	return fmt.Sprintf("%s/api/rest/v1/msp/tenants/%s/users/groups", baseUrl, tenantUid)
}

func GetUserGroupsInMspManagedTenant(baseUrl string, tenantUid string, limit int, offset int) string {
	return fmt.Sprintf("%s/api/rest/v1/msp/tenants/%s/users/groups?limit=%d&offset=%d", baseUrl, tenantUid, limit, offset)
}

func DeleteUserGroupsInMspManagedTenant(baseUrl string, tenantUid string) string {
	return fmt.Sprintf("%s/api/rest/v1/msp/tenants/%s/users/groups/delete", baseUrl, tenantUid)
}

func GetCompatibleAsaVersions(baseUrl string, deviceUid string) string {
	return fmt.Sprintf("%s/api/rest/v1/inventory/devices/asas/%s/upgrades/versions", baseUrl, deviceUid)
}

func GetUpgradeAsaUrl(baseUrl string, deviceUid string) string {
	return fmt.Sprintf("%s/api/rest/v1/inventory/devices/asas/%s/upgrades/trigger", baseUrl, deviceUid)
}

func GetFtdUpgradePackagesUrl(baseUrl string, deviceUid string) string {
	return fmt.Sprintf("%s/api/rest/v1/inventory/devices/ftds/%s/upgrades/versions", baseUrl, deviceUid)
}

func GetFtdUpgradeUrl(baseUrl string, deviceUid string) string {
	return fmt.Sprintf("%s/api/rest/v1/inventory/devices/ftds/%s/upgrades/trigger", baseUrl, deviceUid)
}
