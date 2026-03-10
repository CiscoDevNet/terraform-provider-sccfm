package validators_test

import (
	"context"
	"testing"

	"github.com/CiscoDevnet/terraform-provider-sccfm/validators"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestSocketAddressValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		in        types.String
		expErrors int
	}

	testCases := map[string]testCase{
		"unknown-value":            {in: types.StringUnknown(), expErrors: 0},
		"null-value":               {in: types.StringNull(), expErrors: 0},
		"valid-hostname":           {in: types.StringValue("example.com:443"), expErrors: 0},
		"valid-ipv4":               {in: types.StringValue("10.0.0.1:8443"), expErrors: 0},
		"valid-ipv6":               {in: types.StringValue("[2001:db8::1]:443"), expErrors: 0},
		"invalid-missing-port":     {in: types.StringValue("example.com"), expErrors: 1},
		"invalid-empty-host":       {in: types.StringValue(":443"), expErrors: 1},
		"invalid-non-numeric-port": {in: types.StringValue("example.com:https"), expErrors: 1},
		"invalid-port-too-low":     {in: types.StringValue("example.com:0"), expErrors: 1},
		"invalid-port-too-high":    {in: types.StringValue("example.com:65536"), expErrors: 1},
	}

	v := validators.ValidateSocketAddress()

	for name, test := range testCases {
		name, test := name, test

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := validator.StringRequest{
				ConfigValue: test.in,
			}
			res := validator.StringResponse{}

			v.ValidateString(context.Background(), req, &res)

			if got := res.Diagnostics.ErrorsCount(); got != test.expErrors {
				t.Fatalf("expected %d error(s), got %d: %v", test.expErrors, got, res.Diagnostics)
			}
		})
	}
}
