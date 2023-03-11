package search

import (
	"fmt"
	"os"
	"testing"

	"github.com/clintjedwards/basecoat/internal/storage"
	"github.com/lithammer/shortuuid"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func ptr[T any](v T) *T {
	return &v
}

var testInfo = struct {
	search       *Search
	databasePath string
	storage      storage.DB
	contractorID string
	formula1ID   string
	formula2ID   string
	job1ID       string
	job2ID       string
}{}

func setup() {
	databasePath := fmt.Sprintf("/tmp/basecoat%s.db", shortuuid.New()[0:7])

	storage, err := storage.New(databasePath, 100)
	if err != nil {
		log.Fatal().Err(err).Msg("could not init storage")
	}

	searchIndex, err := InitSearch(storage)
	if err != nil {
		log.Fatal().Err(err).Msg("could not init storage")
	}

	testInfo.search = searchIndex
	testInfo.databasePath = databasePath
	testInfo.storage = storage

	populateDB()
	searchIndex.BuildIndex(storage)
}

func populateDB() {
	err := testInfo.storage.InsertAccount(testInfo.storage.DB, &storage.Account{
		ID:   "test",
		Name: "test",
	})
	if err != nil {
		panic(err)
	}
	// contractorID, err := testInfo.storage.AddContractor("test", &api.Contractor{
	// 	Company: "testcontractor",
	// })

	err = testInfo.storage.InsertFormula(testInfo.storage.DB, &storage.Formula{
		Account: "test",
		ID:      "formula_1",
		Name:    "testunique",
	})
	if err != nil {
		panic(err)
	}
	err = testInfo.storage.InsertFormula(testInfo.storage.DB, &storage.Formula{
		Account: "test",
		ID:      "formula_2",
		Name:    "test-name",
	})
	if err != nil {
		panic(err)
	}
	// job1ID, err := testInfo.storage.AddJob("test", &api.Job{
	// 	Name:         "testjob1",
	// 	ContractorId: contractorID,
	// })
	// job2ID, err := testInfo.storage.AddJob("test", &api.Job{
	// 	Name: "testjob2",
	// })
	// if err != nil {
	// 	return err
	// }

	// testInfo.contractorID = contractorID
	testInfo.formula1ID = "formula_1"
	testInfo.formula2ID = "formula_2"
	// testInfo.job1ID = job1ID
	// testInfo.job2ID = job2ID
}

func resetIndex() {
	testInfo.search.BuildIndex(testInfo.storage)
}

func TestMain(m *testing.M) {
	setup()
	exitVal := m.Run()
	teardown()
	os.Exit(exitVal)
}

func teardown() {
	os.Remove(testInfo.databasePath)
}

func TestSearchFormulas(t *testing.T) {
	results, err := testInfo.search.SearchFormulas("test", "formula")
	require.NoError(t, err)
	require.NotEmpty(t, results)
	require.Contains(t, results, testInfo.formula1ID)
}

func TestSearchFormulasPartialDashed(t *testing.T) {
	results, err := testInfo.search.SearchFormulas("test", "name")
	require.NoError(t, err)
	require.NotEmpty(t, results)
	require.Contains(t, results, testInfo.formula2ID)
	require.NotContains(t, results, testInfo.formula1ID)
}

func TestSearchFormulasQueryDashed(t *testing.T) {
	results, err := testInfo.search.SearchFormulas("test", `test-name`)
	require.NoError(t, err)
	require.NotEmpty(t, results)
	require.Contains(t, results, testInfo.formula2ID)
	require.NotContains(t, results, testInfo.formula1ID)
}

// func TestSearchJobs(t *testing.T) {
// 	results, err := testInfo.search.SearchJobs("test", "job")
// 	require.NoError(t, err)
// 	require.NotEmpty(t, results)
// 	require.Contains(t, results, testInfo.job1ID)
// 	require.Contains(t, results, testInfo.job2ID)

// 	results, err = testInfo.search.SearchJobs("test", "contract")
// 	require.NoError(t, err)
// 	require.NotEmpty(t, results)
// 	require.Contains(t, results, testInfo.job1ID)
// }

func TestUpdateFormulaIndex(t *testing.T) {
	results, err := testInfo.search.SearchFormulas("test", "formula")
	require.NoError(t, err)
	require.NotEmpty(t, results)

	err = testInfo.storage.UpdateFormula(testInfo.storage.DB, "test", testInfo.formula1ID, storage.UpdatableFormulaFields{
		Name: ptr("testupdate"),
	})
	require.NoError(t, err)

	testInfo.search.UpdateFormulaIndex("test", testInfo.formula1ID)

	results, err = testInfo.search.SearchFormulas("test", "unique")
	require.NoError(t, err)
	require.Empty(t, results)

	results, err = testInfo.search.SearchFormulas("test", "update")
	require.NoError(t, err)
	require.NotEmpty(t, results)
	require.Contains(t, results, testInfo.formula1ID)
}

// func TestUpdateJobIndex(t *testing.T) {
// 	results, err := testInfo.search.SearchJobs("test", "job")
// 	require.NoError(t, err)
// 	require.Len(t, results, 2)

// 	err = testInfo.storage.UpdateJob("test", testInfo.job1ID, &api.Job{
// 		Id:   testInfo.job1ID,
// 		Name: "testdifferent",
// 	})
// 	require.NoError(t, err)

// 	testInfo.search.UpdateJobIndex("test", testInfo.job1ID)

// 	results, err = testInfo.search.SearchJobs("test", "job")
// 	require.NoError(t, err)
// 	require.Len(t, results, 1)

// 	results, err = testInfo.search.SearchJobs("test", "diff")
// 	require.NoError(t, err)
// 	require.NotEmpty(t, results)
// 	require.Contains(t, results, testInfo.job1ID)
// }

func TestDeleteFormulaIndex(t *testing.T) {
	results, err := testInfo.search.SearchFormulas("test", "update")
	require.NoError(t, err)
	require.NotEmpty(t, results)
	require.Contains(t, results, testInfo.formula1ID)

	testInfo.search.DeleteFormulaIndex("test", testInfo.formula1ID)

	results, err = testInfo.search.SearchFormulas("test", "update")
	require.NoError(t, err)
	require.Empty(t, results)
}

// func TestDeleteJobIndex(t *testing.T) {
// 	results, err := testInfo.search.SearchJobs("test", "diff")
// 	require.NoError(t, err)
// 	require.NotEmpty(t, results)
// 	require.Contains(t, results, testInfo.job1ID)

// 	testInfo.search.DeleteJobIndex("test", testInfo.job1ID)

// 	results, err = testInfo.search.SearchJobs("test", "diff")
// 	require.NoError(t, err)
// 	require.Empty(t, results)
// }

func TestSanitizeQueryString(t *testing.T) {
	tests := map[string]struct {
		query string
		want  string
	}{
		"simple term (no change)":       {"helloworld", `helloworld`},
		"compound phrase (no change)":   {"hello world", `hello world`},
		"dashed phrase":                 {"hello-world", `hello world`},
		"multiple changes of same char": {"hello-world-its-me", `hello world its me`},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := sanitizeQuery(tc.query)
			if got != tc.want {
				t.Errorf("want %q, got %s", tc.want, got)
			}
		})
	}
}
