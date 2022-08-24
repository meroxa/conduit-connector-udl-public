package connector

import (
	"context"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/golang/mock/gomock"
	"testing"

	"github.com/meroxa/conduit-connector-udl/udl/mock"
)

func TestWrite(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	mockClient := mock.NewMockClientInterface(mockCtrl)
	testDestination := Destination{client: mockClient}

	ctx := context.Background()
	mockClient.EXPECT().FiledropUdlElsetPostId(ctx, nil).Return(nil, nil).Times(1)

	testDestination.Write(ctx, []sdk.Record{
		{
			Key:      sdk.RawData("100"),
			Position: []byte(`1`),
			Payload: sdk.Change{
				After: sdk.RawData(`{"idOnOrbit":"100"}`),
			},
			Operation: sdk.OperationCreate,
		},
	})
}
