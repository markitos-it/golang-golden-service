package gapi_test

import (
	"testing"

	"markitos-it-svc-golden/internal/domain/model"
	"markitos-it-svc-golden/internal/infrastructure/gapi"
	"markitos-it-svc-golden/internal/testsuite/infrastructure/testdb"

	"github.com/stretchr/testify/require"
)

func TestGoldenCanListAllResources(t *testing.T) {
	golden1 := createPersistedRandomGolden()
	golden2 := createPersistedRandomGolden()

	db := testdb.GetDB()
	t.Cleanup(func() {
		db.Delete(&model.Golden{}, "id IN ?", []string{golden1.Id, golden2.Id})
		db.Delete(&model.GoldenEvent{}, "entity_id IN ? AND entity_name = ?", []string{golden1.Id, golden2.Id}, "golden")
	})

	resp, err := grpcClient.ListGoldens(ctx, &gapi.ListGoldensRequest{})

	require.NoError(t, err)
	require.NotNil(t, resp.Goldens)

	found1, found2 := false, false
	for _, b := range resp.Goldens {
		if b.Id == golden1.Id {
			found1 = true
		}
		if b.Id == golden2.Id {
			found2 = true
		}
	}
	require.True(t, found1, "First golden not found in response")
	require.True(t, found2, "Second golden not found in response")

	deletePersistedRandomGolden(golden1.Id)
	deletePersistedRandomGolden(golden2.Id)
}
