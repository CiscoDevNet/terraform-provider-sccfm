package validators

import (
	"context"
	"net"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = socketAddressValidator{}

type socketAddressValidator struct{}

func (v socketAddressValidator) Description(ctx context.Context) string {
	return v.MarkdownDescription(ctx)
}

func (v socketAddressValidator) MarkdownDescription(_ context.Context) string {
	return "value must contain a host and a port (1–65535) in the format `host:port` or `[ipv6]:port`"
}

func (v socketAddressValidator) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}

	value := request.ConfigValue.ValueString()

	// use Go's standard library to split and validate "host:port" or "[ipv6]:port"
	host, port, err := net.SplitHostPort(value)
	if err != nil || host == "" {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			request.ConfigValue.String(),
		))
		return
	}

	// ensure the port is a valid integer within the allowed TCP/UDP port range.
	p, err := strconv.Atoi(port)
	if err != nil || p < 1 || p > 65535 {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueMatchDiagnostic(
			request.Path,
			v.Description(ctx),
			request.ConfigValue.String(),
		))
	}
}

// ValidateSocketAddress checks that the given socket address is valid.
func ValidateSocketAddress() validator.String {
	return socketAddressValidator{}
}
