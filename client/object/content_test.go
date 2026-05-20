package object_test

import (
	"encoding/json"
	"testing"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/object"
	"github.com/stretchr/testify/assert"
)

func TestMarshalContent_NetworkContent(t *testing.T) {
	nc := object.NetworkContent{Literal: "10.0.0.1"}
	raw, err := object.MarshalContent(nc)
	assert.Nil(t, err)
	assert.NotNil(t, raw)

	var decoded object.NetworkContent
	err = json.Unmarshal(*raw, &decoded)
	assert.Nil(t, err)
	assert.Equal(t, "10.0.0.1", decoded.Literal)
}

func TestMarshalContent_UrlContent(t *testing.T) {
	uc := object.UrlContent{Url: "https://www.example.com"}
	raw, err := object.MarshalContent(uc)
	assert.Nil(t, err)
	assert.NotNil(t, raw)

	var decoded object.UrlContent
	err = json.Unmarshal(*raw, &decoded)
	assert.Nil(t, err)
	assert.Equal(t, "https://www.example.com", decoded.Url)
}

func TestMarshalContent_ServiceContent(t *testing.T) {
	sc := object.ServiceContent{
		Protocol: "TCP",
		ServiceValue: &object.ServiceValueContent{
			Literal: "80",
		},
	}
	raw, err := object.MarshalContent(sc)
	assert.Nil(t, err)
	assert.NotNil(t, raw)

	var decoded object.ServiceContent
	err = json.Unmarshal(*raw, &decoded)
	assert.Nil(t, err)
	assert.Equal(t, "TCP", decoded.Protocol)
	assert.NotNil(t, decoded.ServiceValue)
	assert.Equal(t, "80", decoded.ServiceValue.Literal)
}

func TestMarshalContent_ServiceContent_NoPort(t *testing.T) {
	sc := object.ServiceContent{
		Protocol: "ICMP",
	}
	raw, err := object.MarshalContent(sc)
	assert.Nil(t, err)
	assert.NotNil(t, raw)

	var decoded object.ServiceContent
	err = json.Unmarshal(*raw, &decoded)
	assert.Nil(t, err)
	assert.Equal(t, "ICMP", decoded.Protocol)
	assert.Nil(t, decoded.ServiceValue)
}

func TestMarshalContent_GroupContent(t *testing.T) {
	ncRaw, _ := json.Marshal(object.NetworkContent{Literal: "10.0.0.1"})
	gc := object.GroupContent{
		Literals:             []json.RawMessage{json.RawMessage(ncRaw)},
		ReferencedObjectUids: []string{"uid-1", "uid-2"},
	}
	raw, err := object.MarshalContent(gc)
	assert.Nil(t, err)
	assert.NotNil(t, raw)

	var decoded object.GroupContent
	err = json.Unmarshal(*raw, &decoded)
	assert.Nil(t, err)
	assert.Len(t, decoded.Literals, 1)
	assert.Equal(t, []string{"uid-1", "uid-2"}, decoded.ReferencedObjectUids)
}

func TestUnmarshalContent_NilInput(t *testing.T) {
	result, err := object.UnmarshalContent[object.NetworkContent](nil)
	assert.Nil(t, err)
	assert.Nil(t, result)
}

func TestUnmarshalContent_NetworkContent(t *testing.T) {
	raw := json.RawMessage(`{"literal":"192.168.1.0/24"}`)
	result, err := object.UnmarshalContent[object.NetworkContent](&raw)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "192.168.1.0/24", result.Literal)
}

func TestUnmarshalContent_UrlContent(t *testing.T) {
	raw := json.RawMessage(`{"url":"https://example.org"}`)
	result, err := object.UnmarshalContent[object.UrlContent](&raw)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "https://example.org", result.Url)
}

func TestUnmarshalContent_InvalidJson(t *testing.T) {
	raw := json.RawMessage(`{not valid json}`)
	result, err := object.UnmarshalContent[object.NetworkContent](&raw)
	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func TestMarshalUnmarshalRoundTrip(t *testing.T) {
	original := object.ServiceContent{
		Protocol: "UDP",
		ServiceValue: &object.ServiceValueContent{
			Literal: "53",
		},
	}

	raw, err := object.MarshalContent(original)
	assert.Nil(t, err)

	result, err := object.UnmarshalContent[object.ServiceContent](raw)
	assert.Nil(t, err)
	assert.Equal(t, original.Protocol, result.Protocol)
	assert.Equal(t, original.ServiceValue.Literal, result.ServiceValue.Literal)
}
