package gapi_test

import (
	"testing"

	"markitos-it-svc-golden/internal/domain/model"
	"markitos-it-svc-golden/internal/domain/shared"
	"markitos-it-svc-golden/internal/infrastructure/gapi"
	"markitos-it-svc-golden/internal/testsuite/infrastructure/testdb"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCanSearchWithPattern(t *testing.T) {
	pattern := shared.RandomString(10)
	var ids []string

	for range 5 {
		id := shared.UUIDv4()
		ids = append(ids, id)
		name := pattern + shared.RandomPersonalName()
		golden, _ := model.NewGolden(id, name, shared.RandomString(20), "", "")

		testdb.GetRepository().Create(golden)
	}

	db := testdb.GetDB()
	t.Cleanup(func() {
		db.Delete(&model.Golden{}, "id IN ?", ids)
		db.Delete(&model.GoldenEvent{}, "entity_id IN ? AND entity_name = ?", ids, "golden")
	})

	resp, err := grpcClient.SearchGoldens(ctx, &gapi.SearchGoldensRequest{
		SearchTerm: pattern,
		PageNumber: 1,
		PageSize:   6,
	})

	require.NoError(t, err)
	require.Equal(t, 5, len(resp.Goldens))

	foundCount := 0
	for _, aGolden := range resp.Goldens {
		for _, id := range ids {
			if aGolden.Id == id {
				foundCount++
				break
			}
		}
	}
	require.Equal(t, 5, foundCount)
}

func TestCantSearchWithInvalidOptionalPage(t *testing.T) {
	_, err := grpcClient.SearchGoldens(ctx, &gapi.SearchGoldensRequest{
		PageNumber: -1,
		PageSize:   10,
	})

	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok)
	require.Equal(t, codes.InvalidArgument, st.Code())
}
