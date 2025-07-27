package repository_test

import (
	"testing"

	_ "github.com/lib/pq"

	"ftrack/repository"
	"ftrack/tests"
)

func TestInsertRow(t *testing.T) {
	tdb := tests.GetTestDB(t)

	repo := repository.NewRepository(tdb.DB)

	t.Run("ins row into valid table", func(t *testing.T) {
		newRow := map[string]any{
			"name":   "deadlift",
			"weight": 200,
			"ieninreps":   10,
		}

		err := repo.InsertRow(tests.TABLE1, &newRow)
		if err != nil {
			t.Errorf("Err: %v", err)
		}
	})
}
