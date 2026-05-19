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

func TestDeleteTargets(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	testCases := []struct {
		testName   string
		input      object.DeleteTargetsInput
		setupFunc  func()
		assertFunc func(err error, t *testing.T)
	}{
		{
			testName: "successfully detach targets from object",
			input: object.DeleteTargetsInput{
				Uid:         objectUid,
				TargetUuids: []string{targetId1, targetId2},
			},
			setupFunc: func() {
				httpmock.RegisterResponder(
					"DELETE",
					fmt.Sprintf("/api/rest/v1/objects/%s/targets", objectUid),
					httpmock.NewStringResponder(204, ""),
				)
			},
			assertFunc: func(err error, t *testing.T) {
				assert.Nil(t, err)
			},
		},
		{
			testName: "successfully detach single target",
			input: object.DeleteTargetsInput{
				Uid:         objectUid,
				TargetUuids: []string{targetId1},
			},
			setupFunc: func() {
				httpmock.RegisterResponder(
					"DELETE",
					fmt.Sprintf("/api/rest/v1/objects/%s/targets", objectUid),
					httpmock.NewStringResponder(204, ""),
				)
			},
			assertFunc: func(err error, t *testing.T) {
				assert.Nil(t, err)
			},
		},
		{
			testName: "should error on server error",
			input: object.DeleteTargetsInput{
				Uid:         objectUid,
				TargetUuids: []string{targetId1},
			},
			setupFunc: func() {
				httpmock.RegisterResponder(
					"DELETE",
					fmt.Sprintf("/api/rest/v1/objects/%s/targets", objectUid),
					httpmock.NewStringResponder(500, "internal server error"),
				)
			},
			assertFunc: func(err error, t *testing.T) {
				assert.NotNil(t, err)
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.testName, func(t *testing.T) {
			httpmock.Reset()

			testCase.setupFunc()

			err := object.DeleteTargets(
				context.Background(),
				*http.MustNewWithConfig(baseUrl, "a_valid_token", 0, 0, time.Minute),
				testCase.input,
			)

			testCase.assertFunc(err, t)
		})
	}
}
