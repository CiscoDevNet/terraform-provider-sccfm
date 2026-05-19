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

func TestRead(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	networkContent := mustMarshal(object.NetworkContent{Literal: "10.0.0.1"})

	validOutput := object.ReadOutput{
		Uid:         objectUid,
		Name:        objectName,
		Description: description,
		Value: object.SharedObjectValue{
			ObjectType:     object.NetworkObject,
			DefaultContent: networkContent,
		},
		Targets: []object.Target{
			{Id: targetId1},
		},
	}

	testCases := []struct {
		testName   string
		uid        string
		setupFunc  func()
		assertFunc func(output *object.ReadOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully read object by uid",
			uid:      objectUid,
			setupFunc: func() {
				httpmock.RegisterResponder(
					"GET",
					fmt.Sprintf("/api/rest/v1/objects/%s", objectUid),
					httpmock.NewJsonResponderOrPanic(200, validOutput),
				)
			},
			assertFunc: func(output *object.ReadOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, objectUid, output.Uid)
				assert.Equal(t, objectName, output.Name)
				assert.Equal(t, object.NetworkObject, output.Value.ObjectType)
				assert.Len(t, output.Targets, 1)
				assert.Equal(t, targetId1, output.Targets[0].Id)
			},
		},
		{
			testName: "successfully read object with overrides and multiple targets",
			uid:      objectUid,
			setupFunc: func() {
				overrideContent := mustMarshal(object.NetworkContent{Literal: "10.0.0.2"})
				outputWithOverrides := object.ReadOutput{
					Uid:         objectUid,
					Name:        objectName,
					Description: description,
					Value: object.SharedObjectValue{
						ObjectType:     object.NetworkObject,
						DefaultContent: networkContent,
						Overrides: []object.Override{
							{TargetId: targetId1, Content: overrideContent},
							{TargetId: targetId2, Content: mustMarshal(object.NetworkContent{Literal: "10.0.0.3"})},
						},
					},
					Targets: []object.Target{
						{Id: targetId1},
						{Id: targetId2},
					},
					ReadOnly:      false,
					ObjectVersion: 3,
				}
				httpmock.RegisterResponder(
					"GET",
					fmt.Sprintf("/api/rest/v1/objects/%s", objectUid),
					httpmock.NewJsonResponderOrPanic(200, outputWithOverrides),
				)
			},
			assertFunc: func(output *object.ReadOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Len(t, output.Value.Overrides, 2)
				assert.Equal(t, targetId1, output.Value.Overrides[0].TargetId)
				assert.Equal(t, targetId2, output.Value.Overrides[1].TargetId)
				assert.Len(t, output.Targets, 2)
				assert.Equal(t, int64(3), output.ObjectVersion)
				assert.False(t, output.ReadOnly)

				oc1, _ := object.UnmarshalContent[object.NetworkContent](output.Value.Overrides[0].Content)
				assert.Equal(t, "10.0.0.2", oc1.Literal)
				oc2, _ := object.UnmarshalContent[object.NetworkContent](output.Value.Overrides[1].Content)
				assert.Equal(t, "10.0.0.3", oc2.Literal)
			},
		},
		{
			testName: "successfully read read-only object",
			uid:      objectUid,
			setupFunc: func() {
				readOnlyOutput := object.ReadOutput{
					Uid:      objectUid,
					Name:     objectName,
					ReadOnly: true,
					Value: object.SharedObjectValue{
						ObjectType:     object.NetworkObject,
						DefaultContent: networkContent,
					},
				}
				httpmock.RegisterResponder(
					"GET",
					fmt.Sprintf("/api/rest/v1/objects/%s", objectUid),
					httpmock.NewJsonResponderOrPanic(200, readOnlyOutput),
				)
			},
			assertFunc: func(output *object.ReadOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.True(t, output.ReadOnly)
			},
		},
		{
			testName: "returns error when object not found",
			uid:      objectUid,
			setupFunc: func() {
				httpmock.RegisterResponder(
					"GET",
					fmt.Sprintf("/api/rest/v1/objects/%s", objectUid),
					httpmock.NewStringResponder(404, ""),
				)
			},
			assertFunc: func(output *object.ReadOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
		{
			testName: "returns error on server error",
			uid:      objectUid,
			setupFunc: func() {
				httpmock.RegisterResponder(
					"GET",
					fmt.Sprintf("/api/rest/v1/objects/%s", objectUid),
					httpmock.NewStringResponder(500, "service error"),
				)
			},
			assertFunc: func(output *object.ReadOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			output, err := object.Read(
				context.Background(),
				*http.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				object.ReadInput{Uid: testCase.uid},
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
