package object_test

import (
	"context"
	"testing"
	"time"

	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/object"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestReadByName(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	networkContent := mustMarshal(object.NetworkContent{Literal: "10.0.0.1"})

	validObject := object.ReadOutput{
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
		input      object.ReadByNameInput
		setupFunc  func()
		assertFunc func(output *object.ReadOutput, err error, t *testing.T)
	}{
		{
			testName: "successfully read object by name",
			input: object.ReadByNameInput{
				Name:       objectName,
				ObjectType: object.NetworkObject,
			},
			setupFunc: func() {
				httpmock.RegisterResponder(
					"GET",
					"/api/rest/v1/objects",
					httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
						"count":  1,
						"limit":  50,
						"offset": 0,
						"items":  []object.ReadOutput{validObject},
					}),
				)
			},
			assertFunc: func(output *object.ReadOutput, err error, t *testing.T) {
				assert.Nil(t, err)
				assert.NotNil(t, output)
				assert.Equal(t, objectUid, output.Uid)
				assert.Equal(t, objectName, output.Name)
			},
		},
		{
			testName: "returns error when no object found",
			input: object.ReadByNameInput{
				Name:       "nonexistent",
				ObjectType: object.NetworkObject,
			},
			setupFunc: func() {
				httpmock.RegisterResponder(
					"GET",
					"/api/rest/v1/objects",
					httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
						"count":  0,
						"limit":  50,
						"offset": 0,
						"items":  []object.ReadOutput{},
					}),
				)
			},
			assertFunc: func(output *object.ReadOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
				assert.ErrorContains(t, err, "no object found")
			},
		},
		{
			testName: "returns error when multiple objects found",
			input: object.ReadByNameInput{
				Name:       objectName,
				ObjectType: object.NetworkObject,
			},
			setupFunc: func() {
				httpmock.RegisterResponder(
					"GET",
					"/api/rest/v1/objects",
					httpmock.NewJsonResponderOrPanic(200, map[string]interface{}{
						"count":  2,
						"limit":  50,
						"offset": 0,
						"items":  []object.ReadOutput{validObject, validObject},
					}),
				)
			},
			assertFunc: func(output *object.ReadOutput, err error, t *testing.T) {
				assert.NotNil(t, err)
				assert.Nil(t, output)
				assert.ErrorContains(t, err, "multiple objects found")
			},
		},
		{
			testName: "returns error on server error",
			input: object.ReadByNameInput{
				Name:       objectName,
				ObjectType: object.NetworkObject,
			},
			setupFunc: func() {
				httpmock.RegisterResponder(
					"GET",
					"/api/rest/v1/objects",
					httpmock.NewStringResponder(500, "server error"),
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

			output, err := object.ReadByName(
				context.Background(),
				*http.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(output, err, t)
		})
	}
}
