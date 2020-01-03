package search

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/basecoat/storage"
	"github.com/clintjedwards/toolkit/random"
	"github.com/stretchr/testify/require"
)

var testInfo = struct {
	search       *Search
	databasePath string
	storage      storage.BoltDB
	contractorID string
	formulaID    string
	job1ID       string
	job2ID       string
}{}

func setup() {

	databasePath := fmt.Sprintf("/tmp/basecoat%s.db", random.GenerateRandString(4))

	storage, err := storage.NewBoltDB(databasePath, 4)
	if err != nil {
		log.Fatal(err)
	}
	searchIndex, err := InitSearch(storage)
	if err != nil {
		log.Fatal(err)
	}

	testInfo.search = searchIndex
	testInfo.databasePath = databasePath
	testInfo.storage = storage

	populateDB()
	searchIndex.BuildIndex(storage)
}

func populateDB() error {
	testInfo.storage.CreateAccount("test", "test")
	contractorID, err := testInfo.storage.AddContractor("test", &api.Contractor{
		Company: "testcontractor",
	})
	formulaID, err := testInfo.storage.AddFormula("test", &api.Formula{
		Name: "testformula",
	})
	job1ID, err := testInfo.storage.AddJob("test", &api.Job{
		Name:         "testjob1",
		ContractorId: contractorID,
	})
	job2ID, err := testInfo.storage.AddJob("test", &api.Job{
		Name: "testjob2",
	})
	if err != nil {
		return err
	}

	testInfo.contractorID = contractorID
	testInfo.formulaID = formulaID
	testInfo.job1ID = job1ID
	testInfo.job2ID = job2ID

	return nil
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
	require.Contains(t, results, testInfo.formulaID)
}

func TestSearchJobs(t *testing.T) {
	results, err := testInfo.search.SearchJobs("test", "job")
	require.NoError(t, err)
	require.NotEmpty(t, results)
	require.Contains(t, results, testInfo.job1ID)
	require.Contains(t, results, testInfo.job2ID)

	results, err = testInfo.search.SearchJobs("test", "contract")
	require.NoError(t, err)
	require.NotEmpty(t, results)
	require.Contains(t, results, testInfo.job1ID)
}

func TestUpdateFormulaIndex(t *testing.T) {
	results, err := testInfo.search.SearchFormulas("test", "formula")
	require.NoError(t, err)
	require.NotEmpty(t, results)

	err = testInfo.storage.UpdateFormula("test", testInfo.formulaID, &api.Formula{
		Id:   testInfo.formulaID,
		Name: "testupdate",
	})

	testInfo.search.UpdateFormulaIndex("test", testInfo.formulaID)

	results, err = testInfo.search.SearchFormulas("test", "formula")
	require.NoError(t, err)
	require.Empty(t, results)

	results, err = testInfo.search.SearchFormulas("test", "update")
	require.NoError(t, err)
	require.NotEmpty(t, results)
	require.Contains(t, results, testInfo.formulaID)
}

func TestUpdateJobIndex(t *testing.T) {
	results, err := testInfo.search.SearchJobs("test", "job")
	require.NoError(t, err)
	require.Len(t, results, 2)

	err = testInfo.storage.UpdateJob("test", testInfo.job1ID, &api.Job{
		Id:   testInfo.job1ID,
		Name: "testdifferent",
	})
	require.NoError(t, err)

	testInfo.search.UpdateJobIndex("test", testInfo.job1ID)

	results, err = testInfo.search.SearchJobs("test", "job")
	require.NoError(t, err)
	require.Len(t, results, 1)

	results, err = testInfo.search.SearchJobs("test", "diff")
	require.NoError(t, err)
	require.NotEmpty(t, results)
	require.Contains(t, results, testInfo.job1ID)
}

func TestDeleteFormulaIndex(t *testing.T) {
	results, err := testInfo.search.SearchFormulas("test", "update")
	require.NoError(t, err)
	require.NotEmpty(t, results)
	require.Contains(t, results, testInfo.formulaID)

	testInfo.search.DeleteFormulaIndex("test", testInfo.formulaID)

	results, err = testInfo.search.SearchFormulas("test", "update")
	require.NoError(t, err)
	require.Empty(t, results)
}

func TestDeleteJobIndex(t *testing.T) {

	results, err := testInfo.search.SearchJobs("test", "diff")
	require.NoError(t, err)
	require.NotEmpty(t, results)
	require.Contains(t, results, testInfo.job1ID)

	testInfo.search.DeleteJobIndex("test", testInfo.job1ID)

	results, err = testInfo.search.SearchJobs("test", "diff")
	require.NoError(t, err)
	require.Empty(t, results)

}
