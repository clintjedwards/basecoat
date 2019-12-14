// Package search defines objects that are used to enable search functions within basecoat
package search

import (
	"fmt"
	"strings"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/basecoat/storage"
	"github.com/clintjedwards/basecoat/utils"
)

// searchSyntax is a wrapper for all search terms to improve fuzzy searching
// We use the wildcard query type and then wrap user search terms with *
// meaning search for 0 or more characters
const searchSyntax string = "*%s*"

// Search represents a index that can be used to look up basecoat items
type Search struct {
	formulaIndex map[string]bleve.Index
	jobIndex     map[string]bleve.Index
}

// InitSearch creates a search index object which can then be queried for search results
func InitSearch() (*Search, error) {

	return &Search{
		formulaIndex: map[string]bleve.Index{},
		jobIndex:     map[string]bleve.Index{},
	}, nil
}

// BuildIndex will query basecoat's database and populate the search index
func (searchIndex *Search) BuildIndex() {
	// Log how long it took to build the index in prometheus
	start := time.Now()

	storage, err := storage.InitStorage()
	if err != nil {
		utils.Log().Fatalw("failed to initialize storage",
			"error", err)
	}

	users, err := storage.GetAllUsers()
	if err != nil {
		utils.Log().Fatalw("failed to query database for accounts",
			"error", err)
	}

	for account := range users {
		populateIndex(account, searchIndex)
	}

	elapsed := time.Since(start)
	utils.Log().Infow("compiled index", "time", elapsed)
	return
}

// UpdateFormulaIndex updates an already loaded formula index
func (searchIndex *Search) UpdateFormulaIndex(account string, formula api.Formula) {
	if _, ok := searchIndex.formulaIndex[account]; !ok {
		searchIndex.formulaIndex[account] = createNewIndex()
	}
	index := searchIndex.formulaIndex[account]
	index.Index(formula.Id, formula)
	return
}

// UpdateJobIndex updates an already loaded job index
func (searchIndex *Search) UpdateJobIndex(account string, job api.Job) {
	if _, ok := searchIndex.jobIndex[account]; !ok {
		searchIndex.jobIndex[account] = createNewIndex()
	}
	index := searchIndex.jobIndex[account]
	index.Index(job.Id, job)
	return
}

// DeleteFormulaIndex updates an already loaded formula index
func (searchIndex *Search) DeleteFormulaIndex(account string, formulaID string) {
	index := searchIndex.formulaIndex[account]
	index.Index(formulaID, nil)
	return
}

// DeleteJobIndex updates an already loaded job index
func (searchIndex *Search) DeleteJobIndex(account string, jobID string) {
	index := searchIndex.jobIndex[account]
	index.Index(jobID, nil)
	return
}

// createNewIndex creates a new empty bleve index
func createNewIndex() bleve.Index {
	indexMapping := bleve.NewIndexMapping()
	index, err := bleve.NewMemOnly(indexMapping)
	if err != nil {
		utils.Log().Errorw("failed to create search index", "error", err)
		return nil
	}

	return index
}

// newAccountIndex populates a new account with blank indexes;
// will only create if index has not been created yet
func newAccountIndex(account string, searchIndex *Search) {

	formulaIndex := createNewIndex()
	jobIndex := createNewIndex()

	if _, ok := searchIndex.formulaIndex[account]; !ok {
		searchIndex.formulaIndex[account] = formulaIndex
	}

	if _, ok := searchIndex.jobIndex[account]; !ok {
		searchIndex.jobIndex[account] = jobIndex
	}
}

// populateIndex queries the database and loads the index for a specific account
func populateIndex(account string, searchIndex *Search) {
	storage, err := storage.InitStorage()
	if err != nil {
		utils.Log().Fatalw("failed to initialize storage", "error", err)
	}

	newAccountIndex(account, searchIndex)

	// Index all formulas
	formulas, err := storage.GetAllFormulas(account)
	if err != nil {
		utils.Log().Errorw("failed to query database for formulas",
			"error", err,
			"account", account)
	}

	for _, formula := range formulas {
		searchIndex.formulaIndex[account].Index(formula.Id, &formula)
	}

	// Index all jobs
	jobs, err := storage.GetAllJobs(account)
	if err != nil {
		utils.Log().Errorw("failed to query database for jobs",
			"error", err,
			"account", account)
	}

	for _, job := range jobs {
		searchIndex.jobIndex[account].Index(job.Id, &job)
	}

	return
}

// SearchFormulas searches the index for matching terms and then returns formulas which might match
func (searchIndex *Search) SearchFormulas(account, searchPhrase string) ([]string, error) {
	if index, ok := searchIndex.formulaIndex[account]; ok {

		formulaHits, err := queryIndex(index, strings.ToLower(searchPhrase))
		if err != nil {
			return nil, err
		}

		return formulaHits, nil
	}

	return nil, fmt.Errorf("could not find account: %s", account)
}

// SearchJobs searches the index for matching terms and then returns jobs which might match
func (searchIndex *Search) SearchJobs(account, searchPhrase string) ([]string, error) {
	if index, ok := searchIndex.jobIndex[account]; ok {

		jobHits, err := queryIndex(index, strings.ToLower(searchPhrase))
		if err != nil {
			return nil, err
		}

		return jobHits, err
	}

	return nil, fmt.Errorf("could not find account: %s", account)
}

// queryIndex runs the actual search query against the index.
// It uses the boolean query which is a type of query builder
// The search phrase given is separated into separate search terms, made into a wildcard query
// and then passed to the boolean query. The boolean query checks that all terms are found in any hits
// it returns.
// Example: "hello world" is searched as .*hello.* .*world.* and only when both are present in a document
// will it be present in the results
func queryIndex(index bleve.Index, searchPhrase string) ([]string, error) {
	queryBuilder := bleve.NewBooleanQuery()

	for _, term := range strings.Split(searchPhrase, " ") {
		query := bleve.NewWildcardQuery(fmt.Sprintf(searchSyntax, term))
		queryBuilder.AddMust(query)
	}

	searchRequest := bleve.NewSearchRequest(queryBuilder)
	searchResult, err := index.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	var matchingIDs []string
	for _, result := range searchResult.Hits {
		matchingIDs = append(matchingIDs, result.ID)
	}

	return matchingIDs, nil
}
