package object_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/object"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	updatedContent := mustMarshal(object.NetworkContent{Literal: "10.0.0.2"})
	newDescription := "updated description"

	validOutput := object.UpdateOutput{
		Uid:         objectUid,
		Name:        objectName,
		Description: newDescription,
		Value: object.SharedObjectValue{
			ObjectType:     object.NetworkObject,
			DefaultContent: updatedContent,
		},
	}

	testCases := []struct {
		testName   string
		input      object.UpdateInput
		setupFunc  func()
		assertFunc func(output *object.UpdateOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully update object",
			input: object.UpdateInput{
				Uid:         objectUid,
				Name:        objectName,
				Description: &newDescription,
				Value: &object.SharedObjectValue{
					ObjectType:     object.NetworkObject,
					DefaultContent: updatedContent,
				},
			},
			setupFunc: func() {
				httpmock.RegisterResponder(
					"PATCH",
					fmt.Sprintf("/api/rest/v1/objects/%s", objectUid),
					httpmock.NewJsonResponderOrPanic(200, validOutput),
				)
			},
			assertFunc: func(output *object.UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, objectUid, output.Uid)
				assert.Equal(t, newDescription, output.Description)
			},
		},
		{
			testName: "successfully update with overrides added",
			input: func() object.UpdateInput {
				overrideContent := mustMarshal(object.NetworkContent{Literal: "10.0.0.99"})
				return object.UpdateInput{
					Uid:         objectUid,
					Name:        objectName,
					Description: &newDescription,
					Value: &object.SharedObjectValue{
						ObjectType:     object.NetworkObject,
						DefaultContent: updatedContent,
						Overrides: []object.Override{
							{TargetId: targetId1, Content: overrideContent},
						},
					},
				}
			}(),
			setupFunc: func() {
				overrideContent := mustMarshal(object.NetworkContent{Literal: "10.0.0.99"})
				httpmock.RegisterResponder(
					"PATCH",
					fmt.Sprintf("/api/rest/v1/objects/%s", objectUid),
					httpmock.NewJsonResponderOrPanic(200, object.UpdateOutput{
						Uid:         objectUid,
						Name:        objectName,
						Description: newDescription,
						Value: object.SharedObjectValue{
							ObjectType:     object.NetworkObject,
							DefaultContent: updatedContent,
							Overrides: []object.Override{
								{TargetId: targetId1, Content: overrideContent},
							},
						},
						Targets: []object.Target{{Id: targetId1}},
					}),
				)
			},
			assertFunc: func(output *object.UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Len(t, output.Value.Overrides, 1)
				assert.Equal(t, targetId1, output.Value.Overrides[0].TargetId)
				oc, _ := object.UnmarshalContent[object.NetworkContent](output.Value.Overrides[0].Content)
				assert.Equal(t, "10.0.0.99", oc.Literal)
				assert.Len(t, output.Targets, 1)
			},
		},
		{
			testName: "successfully update removing all overrides",
			input: object.UpdateInput{
				Uid:         objectUid,
				Name:        objectName,
				Description: &newDescription,
				Value: &object.SharedObjectValue{
					ObjectType:     object.NetworkObject,
					DefaultContent: updatedContent,
					Overrides:      nil,
				},
			},
			setupFunc: func() {
				httpmock.RegisterResponder(
					"PATCH",
					fmt.Sprintf("/api/rest/v1/objects/%s", objectUid),
					httpmock.NewJsonResponderOrPanic(200, object.UpdateOutput{
						Uid:         objectUid,
						Name:        objectName,
						Description: newDescription,
						Value: object.SharedObjectValue{
							ObjectType:     object.NetworkObject,
							DefaultContent: updatedContent,
						},
					}),
				)
			},
			assertFunc: func(output *object.UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Len(t, output.Value.Overrides, 0)
				assert.Len(t, output.Targets, 0)
			},
		},
		{
			testName: "successfully update only description",
			input: object.UpdateInput{
				Uid:         objectUid,
				Description: &newDescription,
			},
			setupFunc: func() {
				httpmock.RegisterResponder(
					"PATCH",
					fmt.Sprintf("/api/rest/v1/objects/%s", objectUid),
					httpmock.NewJsonResponderOrPanic(200, object.UpdateOutput{
						Uid:         objectUid,
						Name:        objectName,
						Description: newDescription,
						Value: object.SharedObjectValue{
							ObjectType:     object.NetworkObject,
							DefaultContent: mustMarshal(object.NetworkContent{Literal: "10.0.0.1"}),
						},
					}),
				)
			},
			assertFunc: func(output *object.UpdateOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, newDescription, output.Description)
				assert.Equal(t, objectName, output.Name)
			},
		},
		{
			testName: "should error on server error",
			input: object.UpdateInput{
				Uid:  objectUid,
				Name: objectName,
			},
			setupFunc: func() {
				httpmock.RegisterResponder(
					"PATCH",
					fmt.Sprintf("/api/rest/v1/objects/%s", objectUid),
					httpmock.NewStringResponder(500, "internal server error"),
				)
			},
			assertFunc: func(output *object.UpdateOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
		{
			testName: "should error on not found",
			input: object.UpdateInput{
				Uid:  objectUid,
				Name: objectName,
			},
			setupFunc: func() {
				httpmock.RegisterResponder(
					"PATCH",
					fmt.Sprintf("/api/rest/v1/objects/%s", objectUid),
					httpmock.NewStringResponder(404, "not found"),
				)
			},
			assertFunc: func(output *object.UpdateOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := object.Update(
				context.Background(),
				*http.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
