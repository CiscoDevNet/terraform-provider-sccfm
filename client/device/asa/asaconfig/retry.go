package asaconfig

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/retry"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/internal/statemachine"
	"github.com/CiscoDevnet/terraform-provider-sccfm/go-client/model/statemachine/state"
)

func UntilStateDone(ctx context.Context, client http.Client, specificUid string) retry.Func {

	// create asa config read request
	readReq := NewReadRequest(ctx, client, *NewReadInput(
		specificUid,
	))

	var readOutp ReadOutput

	return func() (bool, error) {
		err := readReq.Send(&readOutp)
		if err != nil {
			return false, err
		}

		client.Logger.Printf("asa config state=%s\n", readOutp.State)
		if readOutp.State == state.DONE {
			return true, nil
		}
		if readOutp.State == state.ERROR {
			return false, statemachine.NewWorkflowErrorFromDetails(readOutp.StateMachineDetails)
		}
		if readOutp.State == state.BAD_CREDENTIALS {
			return false, statemachine.NewWorkflowErrorf("Bad Credentials")
		}
		if readOutp.State == state.PRE_WAIT_FOR_USER_TO_UPDATE_CREDS {
			return false, statemachine.NewWorkflowErrorf("Bad Credentials")
		}
		return false, nil
	}
}
