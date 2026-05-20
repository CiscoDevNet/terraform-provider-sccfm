package object

import "encoding/json"

// NetworkContent represents the defaultContent for a NETWORK_OBJECT.
type NetworkContent struct {
	Literal string `json:"literal"`
}

// UrlContent represents the defaultContent for a URL_OBJECT.
type UrlContent struct {
	Url string `json:"url"`
}

// ServiceValueContent represents the port/value portion of a service object (maps to PortsValue in Java).
type ServiceValueContent struct {
	Literal string `json:"literal"`
}

// ServiceContent represents the defaultContent for a SERVICE_OBJECT.
type ServiceContent struct {
	Protocol     string               `json:"protocol"`
	ServiceValue *ServiceValueContent `json:"serviceValue,omitempty"`
}

// GroupContent represents the defaultContent for any group type (NETWORK_GROUP, URL_GROUP, SERVICE_GROUP).
type GroupContent struct {
	Literals             []json.RawMessage `json:"literals,omitempty"`
	ReferencedObjectUids []string          `json:"referencedObjectUids,omitempty"`
}

// MarshalContent marshals a typed content struct into a *json.RawMessage suitable for SharedObjectValue.DefaultContent.
func MarshalContent(content interface{}) (*json.RawMessage, error) {
	b, err := json.Marshal(content)
	if err != nil {
		return nil, err
	}
	raw := json.RawMessage(b)
	return &raw, nil
}

// UnmarshalContent unmarshals a *json.RawMessage into a typed content struct.
func UnmarshalContent[T any](raw *json.RawMessage) (*T, error) {
	if raw == nil {
		return nil, nil
	}
	var content T
	if err := json.Unmarshal(*raw, &content); err != nil {
		return nil, err
	}
	return &content, nil
}
