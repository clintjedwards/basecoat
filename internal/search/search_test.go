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
		ID:   "test_account",
		Name: "test_account",
	})
	if err != nil {
		panic(err)
	}

	err = testInfo.storage.InsertContractor(testInfo.storage.DB, &storage.Contractor{
		Account: "test_account",
		ID:      "test_contractor",
		Company: "test_contractor",
	})
	if err != nil {
		panic(err)
	}

	err = testInfo.storage.InsertFormula(testInfo.storage.DB, &storage.Formula{
		Account: "test_account",
		ID:      "formula_1",
		Name:    "testunique",
	})
	if err != nil {
		panic(err)
	}
	err = testInfo.storage.InsertFormula(testInfo.storage.DB, &storage.Formula{
		Account: "test_account",
		ID:      "formula_2",
		Name:    "test-name",
	})
	if err != nil {
		panic(err)
	}

	err = testInfo.storage.InsertJob(testInfo.storage.DB, &storage.Job{
		Account:    "test_account",
		ID:         "test_job1",
		Name:       "testjob1",
		Contractor: "test_contractor",
	})
	if err != nil {
		panic(err)
	}
	err = testInfo.storage.InsertJob(testInfo.storage.DB, &storage.Job{
		Account:    "test_account",
		ID:         "test_job2",
		Name:       "testjob2",
		Contractor: "test_contractor",
	})
	if err != nil {
		panic(err)
	}

	testInfo.formula1ID = "formula_1"
	testInfo.formula2ID = "formula_2"
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
	results, err := testInfo.search.SearchFormulas("test_account", "formula")
	require.NoError(t, err)
	require.NotEmpty(t, results)
	require.Contains(t, results, testInfo.formula1ID)
}

func TestSearchFormulasPartialDashed(t *testing.T) {
	results, err := testInfo.search.SearchFormulas("test_account", "name")
	require.NoError(t, err)
	require.NotEmpty(t, results)
	require.Contains(t, results, testInfo.formula2ID)
	require.NotContains(t, results, testInfo.formula1ID)
}

func TestSearchFormulasQueryDashed(t *testing.T) {
	results, err := testInfo.search.SearchFormulas("test_account", `test-name`)
	require.NoError(t, err)
	require.NotEmpty(t, results)
	require.Contains(t, results, testInfo.formula2ID)
	require.NotContains(t, results, testInfo.formula1ID)
}

func TestSearchJobs(t *testing.T) {
	results, err := testInfo.search.SearchJobs("test_account", "job")
	require.NoError(t, err)
	require.NotEmpty(t, results)
	require.Contains(t, results, "test_job1")
	require.Contains(t, results, "test_job2")
}

func TestUpdateFormulaIndex(t *testing.T) {
	results, err := testInfo.search.SearchFormulas("test_account", "formula")
	require.NoError(t, err)
	require.NotEmpty(t, results)

	err = testInfo.storage.UpdateFormula(testInfo.storage.DB, "test_account", testInfo.formula1ID, storage.UpdatableFormulaFields{
		Name: ptr("testupdate"),
	})
	require.NoError(t, err)

	testInfo.search.UpdateFormulaIndex("test_account", testInfo.formula1ID)

	results, err = testInfo.search.SearchFormulas("test_account", "unique")
	require.NoError(t, err)
	require.Empty(t, results)

	results, err = testInfo.search.SearchFormulas("test_account", "update")
	require.NoError(t, err)
	require.NotEmpty(t, results)
	require.Contains(t, results, testInfo.formula1ID)
}

func TestUpdateJobIndex(t *testing.T) {
	results, err := testInfo.search.SearchJobs("test_account", "job")
	require.NoError(t, err)
	require.Len(t, results, 2)

	err = testInfo.storage.UpdateJob(testInfo.storage.DB, "test_account", "test_job1", storage.UpdatableJobFields{
		Name: ptr("testdifferent"),
	})
	require.NoError(t, err)

	testInfo.search.UpdateJobIndex("test_account", "test_job1")

	results, err = testInfo.search.SearchJobs("test_account", "job")
	require.NoError(t, err)
	require.Len(t, results, 2)

	results, err = testInfo.search.SearchJobs("test_account", "diff")
	require.NoError(t, err)
	require.NotEmpty(t, results)
	require.Contains(t, results, "test_job1")
}

func TestDeleteFormulaIndex(t *testing.T) {
	results, err := testInfo.search.SearchFormulas("test_account", "update")
	require.NoError(t, err)
	require.NotEmpty(t, results)
	require.Contains(t, results, testInfo.formula1ID)

	testInfo.search.DeleteFormulaIndex("test_account", testInfo.formula1ID)

	results, err = testInfo.search.SearchFormulas("test_account", "update")
	require.NoError(t, err)
	require.Empty(t, results)
}

func TestDeleteJobIndex(t *testing.T) {
	results, err := testInfo.search.SearchJobs("test_account", "diff")
	require.NoError(t, err)
	require.NotEmpty(t, results)
	require.Contains(t, results, "test_job1")

	testInfo.search.DeleteJobIndex("test_account", "test_job1")

	results, err = testInfo.search.SearchJobs("test_account", "diff")
	require.NoError(t, err)
	require.Empty(t, results)
}

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
