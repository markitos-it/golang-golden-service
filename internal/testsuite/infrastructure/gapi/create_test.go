package gapi_test

import (
	"testing"

	model "markitos-it-svc-golden/internal/domain/model"
	"markitos-it-svc-golden/internal/infrastructure/gapi"
	"markitos-it-svc-golden/internal/testsuite/infrastructure/testdb"
	internal_test "markitos-it-svc-golden/internal/testsuite/internal"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGoldenCanCreate(t *testing.T) {
	golden := internal_test.NewRandomOnlyNameGolden()
	resp, err := grpcClient.CreateGolden(ctx, &gapi.CreateGoldenRequest{
		Name: golden.Name,
	})

	require.NoError(t, err)
	require.NotEmpty(t, resp.Id)
	require.Equal(t, golden.Name, resp.Name)

	db := testdb.GetDB()
	t.Cleanup(func() {
		db.Delete(&model.Golden{}, "id = ?", resp.Id)
		db.Delete(&model.GoldenEvent{}, "entity_id = ? AND entity_name = ? AND name = ?", resp.Id, "golden", "golden-created")
	})
}

func TestGoldenCantCreateWithoutName(t *testing.T) {
	_, err := grpcClient.CreateGolden(ctx, &gapi.CreateGoldenRequest{
		Name: "",
	})

	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, st.Code())
}

func TestGoldenCantCreateWithoutValidName(t *testing.T) {
	_, err := grpcClient.CreateGolden(ctx, &gapi.CreateGoldenRequest{
		Name: "!!!!!invalid!!!name!!!",
	})

	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, st.Code())
}
