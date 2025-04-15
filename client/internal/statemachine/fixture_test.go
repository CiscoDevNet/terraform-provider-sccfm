package statemachine_test

import "github.com/CiscoDevnet/terraform-provider-scc-firewall-manager/go-client/internal/statemachine"

const (
	baseUrl   = "https://unit-test.cdo.cisco.com"
	deviceUid = "unit-test-device-uid"
)

var (
	validReadStateMachineOutput = statemachine.NewReadInstanceByDeviceUidOutputBuilder().Build()
)
