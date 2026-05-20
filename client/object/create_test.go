package object_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/object"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestCreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	networkContent := mustMarshal(object.NetworkContent{Literal: "10.0.0.1"})
	urlContent := mustMarshal(object.UrlContent{Url: "https://www.example.com"})
	serviceContent := mustMarshal(object.ServiceContent{Protocol: "TCP", ServiceValue: &object.ServiceValueContent{Literal: "443"}})
	overrideContent := mustMarshal(object.NetworkContent{Literal: "10.0.0.2"})

	validNetworkOutput := object.CreateOutput{
		Uid:         objectUid,
		Name:        objectName,
		Description: description,
		Value: object.SharedObjectValue{
			ObjectType:     object.NetworkObject,
			DefaultContent: networkContent,
		},
	}

	testCases := []struct {
		testName   string
		input      object.CreateInput
		setupFunc  func()
		assertFunc func(output *object.CreateOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully create network object",
			input: object.CreateInput{
				Name:        objectName,
				Description: description,
				Value: object.SharedObjectValue{
					ObjectType:     object.NetworkObject,
					DefaultContent: networkContent,
				},
			},
			setupFunc: func() {
				httpmock.RegisterResponder(
					"POST",
					"/api/rest/v1/objects",
					httpmock.NewJsonResponderOrPanic(200, validNetworkOutput),
				)
			},
			assertFunc: func(output *object.CreateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, objectUid, output.Uid)
				assert.Equal(t, objectName, output.Name)
				assert.Equal(t, description, output.Description)
				assert.Equal(t, object.NetworkObject, output.Value.ObjectType)
			},
		},
		{
			testName: "successfully create url object",
			input: object.CreateInput{
				Name: "test-url-object",
				Value: object.SharedObjectValue{
					ObjectType:     object.UrlObject,
					DefaultContent: urlContent,
				},
			},
			setupFunc: func() {
				httpmock.RegisterResponder(
					"POST",
					"/api/rest/v1/objects",
					httpmock.NewJsonResponderOrPanic(200, object.CreateOutput{
						Uid:  objectUid,
						Name: "test-url-object",
						Value: object.SharedObjectValue{
							ObjectType:     object.UrlObject,
							DefaultContent: urlContent,
						},
					}),
				)
			},
			assertFunc: func(output *object.CreateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, object.UrlObject, output.Value.ObjectType)
				uc, _ := object.UnmarshalContent[object.UrlContent](output.Value.DefaultContent)
				assert.Equal(t, "https://www.example.com", uc.Url)
			},
		},
		{
			testName: "successfully create service object",
			input: object.CreateInput{
				Name: "test-service-object",
				Value: object.SharedObjectValue{
					ObjectType:     object.ServiceObject,
					DefaultContent: serviceContent,
				},
			},
			setupFunc: func() {
				httpmock.RegisterResponder(
					"POST",
					"/api/rest/v1/objects",
					httpmock.NewJsonResponderOrPanic(200, object.CreateOutput{
						Uid:  objectUid,
						Name: "test-service-object",
						Value: object.SharedObjectValue{
							ObjectType:     object.ServiceObject,
							DefaultContent: serviceContent,
						},
					}),
				)
			},
			assertFunc: func(output *object.CreateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, object.ServiceObject, output.Value.ObjectType)
				sc, _ := object.UnmarshalContent[object.ServiceContent](output.Value.DefaultContent)
				assert.Equal(t, "TCP", sc.Protocol)
				assert.Equal(t, "443", sc.ServiceValue.Literal)
			},
		},
		{
			testName: "successfully create network object with overrides",
			input: object.CreateInput{
				Name:        objectName,
				Description: description,
				Value: object.SharedObjectValue{
					ObjectType:     object.NetworkObject,
					DefaultContent: networkContent,
					Overrides: []object.Override{
						{TargetId: targetId1, Content: overrideContent},
					},
				},
			},
			setupFunc: func() {
				outputWithOverrides := object.CreateOutput{
					Uid:         objectUid,
					Name:        objectName,
					Description: description,
					Value: object.SharedObjectValue{
						ObjectType:     object.NetworkObject,
						DefaultContent: networkContent,
						Overrides: []object.Override{
							{TargetId: targetId1, Content: overrideContent},
						},
					},
					Targets: []object.Target{{Id: targetId1}},
				}
				httpmock.RegisterResponder(
					"POST",
					"/api/rest/v1/objects",
					httpmock.NewJsonResponderOrPanic(200, outputWithOverrides),
				)
			},
			assertFunc: func(output *object.CreateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, objectUid, output.Uid)
				assert.Len(t, output.Value.Overrides, 1)
				assert.Equal(t, targetId1, output.Value.Overrides[0].TargetId)
				oc, _ := object.UnmarshalContent[object.NetworkContent](output.Value.Overrides[0].Content)
				assert.Equal(t, "10.0.0.2", oc.Literal)
				assert.Len(t, output.Targets, 1)
			},
		},
		{
			testName: "successfully create network group with literals and refs",
			input: func() object.CreateInput {
				ncRaw, _ := json.Marshal(object.NetworkContent{Literal: "10.0.0.0/24"})
				gc := object.GroupContent{
					Literals:             []json.RawMessage{json.RawMessage(ncRaw)},
					ReferencedObjectUids: []string{"ref-uid-1"},
				}
				gcContent := mustMarshal(gc)
				return object.CreateInput{
					Name: "test-network-group",
					Value: object.SharedObjectValue{
						ObjectType:     object.NetworkGroup,
						DefaultContent: gcContent,
					},
				}
			}(),
			setupFunc: func() {
				ncRaw, _ := json.Marshal(object.NetworkContent{Literal: "10.0.0.0/24"})
				gcContent := mustMarshal(object.GroupContent{
					Literals:             []json.RawMessage{json.RawMessage(ncRaw)},
					ReferencedObjectUids: []string{"ref-uid-1"},
				})
				httpmock.RegisterResponder(
					"POST",
					"/api/rest/v1/objects",
					httpmock.NewJsonResponderOrPanic(200, object.CreateOutput{
						Uid:  objectUid,
						Name: "test-network-group",
						Value: object.SharedObjectValue{
							ObjectType:     object.NetworkGroup,
							DefaultContent: gcContent,
						},
					}),
				)
			},
			assertFunc: func(output *object.CreateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, object.NetworkGroup, output.Value.ObjectType)
				gc, _ := object.UnmarshalContent[object.GroupContent](output.Value.DefaultContent)
				assert.Len(t, gc.Literals, 1)
				assert.Equal(t, []string{"ref-uid-1"}, gc.ReferencedObjectUids)
			},
		},
		{
			testName: "should error on server error",
			input: object.CreateInput{
				Name: objectName,
				Value: object.SharedObjectValue{
					ObjectType:     object.NetworkObject,
					DefaultContent: networkContent,
				},
			},
			setupFunc: func() {
				httpmock.RegisterResponder(
					"POST",
					"/api/rest/v1/objects",
					httpmock.NewStringResponder(500, "internal server error"),
				)
			},
			assertFunc: func(output *object.CreateOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
		{
			testName: "should error on bad request",
			input: object.CreateInput{
				Name: "",
				Value: object.SharedObjectValue{
					ObjectType: object.NetworkObject,
				},
			},
			setupFunc: func() {
				httpmock.RegisterResponder(
					"POST",
					"/api/rest/v1/objects",
					httpmock.NewStringResponder(400, "name is required"),
				)
			},
			assertFunc: func(output *object.CreateOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := object.Create(
				context.Background(),
				*http.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
