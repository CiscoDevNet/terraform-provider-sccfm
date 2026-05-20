package object

import "encoding/json"

// ObjectType enumerates the supported unified object types.
type ObjectType string

const (
	NetworkObject ObjectType = "NETWORK_OBJECT"
	UrlObject     ObjectType = "URL_OBJECT"
	ServiceObject ObjectType = "SERVICE_OBJECT"
	NetworkGroup  ObjectType = "NETWORK_GROUP"
	UrlGroup      ObjectType = "URL_GROUP"
	ServiceGroup  ObjectType = "SERVICE_GROUP"
)

// Override represents a per-target content override for an object.
type Override struct {
	TargetId string           `json:"targetId"`
	Content  *json.RawMessage `json:"content"`
}

// SharedObjectValue is the value envelope sent to and received from the object service.
type SharedObjectValue struct {
	ObjectType     ObjectType       `json:"objectType"`
	DefaultContent *json.RawMessage `json:"defaultContent"`
	Overrides      []Override       `json:"overrides,omitempty"`
}

// CreateInput is the input for creating a unified object.
type CreateInput struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Value       SharedObjectValue `json:"value"`
}

// CreateOutput is the output from creating a unified object.
type CreateOutput = ReadOutput

// ReadInput is the input for reading a unified object.
type ReadInput struct {
	Uid string
}

// Target represents an associated target on a unified object (e.g., a CDFMC instance).
type Target struct {
	Id string `json:"id"`
}

// ReadOutput is the response from reading a unified object.
type ReadOutput struct {
	Uid           string            `json:"uid"`
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Value         SharedObjectValue `json:"value"`
	Targets       []Target          `json:"targets"`
	ReadOnly      bool              `json:"readOnly"`
	ObjectVersion int64             `json:"objectVersion"`
}

// UpdateInput is the input for updating a unified object.
type UpdateInput struct {
	Uid         string             `json:"-"`
	Name        string             `json:"name,omitempty"`
	Description *string            `json:"description,omitempty"`
	Value       *SharedObjectValue `json:"value,omitempty"`
}

// UpdateOutput is the output from updating a unified object.
type UpdateOutput = ReadOutput

// DeleteInput is the input for deleting a unified object.
type DeleteInput struct {
	Uid string
}
