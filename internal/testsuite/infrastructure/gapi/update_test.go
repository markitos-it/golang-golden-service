package gapi_test

import (
	"testing"

	model "markitos-it-svc-golden/internal/domain/model"
	"markitos-it-svc-golden/internal/domain/shared"
	"markitos-it-svc-golden/internal/infrastructure/gapi"
	"markitos-it-svc-golden/internal/testsuite/infrastructure/testdb"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCanUpdateAGolden(t *testing.T) {
	golden := createPersistedRandomGolden()
	updatedName := golden.Name + " UPDATED"

	db := testdb.GetDB()
	t.Cleanup(func() {
		db.Delete(&model.Golden{}, "id = ?", golden.Id)
		db.Delete(&model.GoldenEvent{}, "entity_id = ? AND entity_name = ?", golden.Id, "golden")
	})

	resp, err := grpcClient.UpdateGolden(ctx, &gapi.UpdateGoldenRequest{
		Id:   golden.Id,
		Name: updatedName,
	})

	require.NoError(t, err)
	require.Equal(t, golden.Id, resp.Updated)

	getResp, err := grpcClient.GetGolden(ctx, &gapi.GetGoldenRequest{Id: golden.Id})
	require.NoError(t, err)
	require.Equal(t, updatedName, getResp.Name)
}

func TestCantUpdateANonExistingGolden(t *testing.T) {
	randomId := shared.UUIDv4()
	_, err := grpcClient.UpdateGolden(ctx, &gapi.UpdateGoldenRequest{
		Id:   randomId,
		Name: shared.RandomPersonalName(),
	})

	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.NotFound, st.Code())
}

func TestCantUpdateAnInvalidGoldenId(t *testing.T) {
	_, err := grpcClient.UpdateGolden(ctx, &gapi.UpdateGoldenRequest{
		Id:   "an-invalid-id-format",
		Name: shared.RandomPersonalName(),
	})

	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, st.Code())
}
