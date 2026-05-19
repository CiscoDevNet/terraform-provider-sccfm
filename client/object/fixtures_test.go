package object_test

import "encoding/json"

const (
	baseUrl     = "https://unittest.cdo.cisco.com"
	objectUid   = "11111111-1111-1111-1111-111111111111"
	objectName  = "test-network-object"
	description = "test description"
	targetId1   = "aaaa1111-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	targetId2   = "bbbb2222-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
)

func mustMarshal(v interface{}) *json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	raw := json.RawMessage(b)
	return &raw
}
