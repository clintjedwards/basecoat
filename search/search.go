// Package search defines objects that are used to enable search functions within basecoat
package search

import (
	"fmt"
	"strings"
	"time"

	"github.com/blevesearch/bleve"
	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/basecoat/storage"
	"go.uber.org/zap"
)

// searchSyntax is a wrapper for all search terms to improve fuzzy searching
// We use the wildcard query type and then wrap user search terms with *
// meaning search for 0 or more characters
const searchSyntax string = "*%s*"

// Search represents a index that can be used to look up basecoat items
type Search struct {
	formulaIndex map[string]bleve.Index
	jobIndex     map[string]bleve.Index
	store        storage.BoltDB
}

// extendedJob extends a typical job to include the contractor
// so that we can index both at the same time
type compositeJob struct {
	Job        *api.Job
	Contractor *api.Contractor
}

// InitSearch creates a search index object which can then be queried for search results
func InitSearch(store storage.BoltDB) (*Search, error) {
	return &Search{
		formulaIndex: map[string]bleve.Index{},
		jobIndex:     map[string]bleve.Index{},
		store:        store,
	}, nil
}

// BuildIndex will query basecoat's database and populate the search index
// it clears out any current index with fresh data
func (si *Search) BuildIndex(store storage.BoltDB) {
	// TODO: Log how long it took to build the index in prometheus
	start := time.Now()

	accounts, err := store.GetAllAccounts()
	if err != nil {
		zap.S().Fatalw("failed to query database for accounts",
			"error", err)
	}

	for account := range accounts {
		si.populateIndex(account)
	}

	elapsed := time.Since(start)
	zap.S().Infow("compiled index", "time", elapsed)
}

// UpdateFormulaIndex updates an already loaded formula index
func (si *Search) UpdateFormulaIndex(account string, formulaID string) {
	if _, ok := si.formulaIndex[account]; !ok {
		si.formulaIndex[account] = createNewIndex()
	}

	formula, err := si.store.GetFormula(account, formulaID)
	if err != nil {
		zap.S().Errorw("could not get formula from database",
			"account", account, "formulaID", formulaID)
	}

	index := si.formulaIndex[account]
	index.Index(formulaID, formula)
	return
}

// UpdateJobIndex updates an already loaded job index
func (si *Search) UpdateJobIndex(account string, jobID string) {
	if _, ok := si.jobIndex[account]; !ok {
		si.jobIndex[account] = createNewIndex()
	}

	job, err := si.store.GetJob(account, jobID)
	if err != nil {
		zap.S().Errorw("could not get job from database",
			"account", account, "jobID", jobID)
	}

	contractor := &api.Contractor{}
	if job.ContractorId != "" {
		contractor, err = si.store.GetContractor(account, job.ContractorId)
		if err != nil {
			zap.S().Errorw("could not get contractor from database",
				"account", account, "contractorID", job.ContractorId)
		}
	}

	compJob := &compositeJob{
		Job:        job,
		Contractor: contractor,
	}

	index := si.jobIndex[account]
	index.Index(job.Id, compJob)
	return
}

// DeleteFormulaIndex updates an already loaded formula index
func (si *Search) DeleteFormulaIndex(account string, formulaID string) {
	index := si.formulaIndex[account]
	index.Delete(formulaID)
	return
}

// DeleteJobIndex updates an already loaded job index
func (si *Search) DeleteJobIndex(account string, jobID string) {
	index := si.jobIndex[account]
	index.Delete(jobID)
	return
}

// createNewIndex creates a new empty bleve index
func createNewIndex() bleve.Index {
	indexMapping := bleve.NewIndexMapping()
	index, err := bleve.NewMemOnly(indexMapping)
	if err != nil {
		zap.S().Errorw("failed to create search index", "error", err)
		return nil
	}

	return index
}

// newAccountIndex populates a new account with blank indexes;
// will only create if index has not been created yet
func (si *Search) newAccountIndex(account string) {
	si.formulaIndex[account] = createNewIndex()
	si.jobIndex[account] = createNewIndex()
}

// populateIndex queries the database and loads the index for a specific account
func (si *Search) populateIndex(account string) {
	si.newAccountIndex(account)

	// Index all formulas
	formulas, err := si.store.GetAllFormulas(account)
	if err != nil {
		zap.S().Errorw("failed to query database for formulas",
			"error", err,
			"account", account)
	}

	for _, formula := range formulas {
		si.formulaIndex[account].Index(formula.Id, &formula)
	}

	// Index all jobs
	jobs, err := si.store.GetAllJobs(account)
	if err != nil {
		zap.S().Errorw("failed to query database for jobs",
			"error", err,
			"account", account)
	}

	//Get all contractors to be added into job indexes
	contractors, err := si.store.GetAllContractors(account)
	if err != nil {
		zap.S().Errorw("failed to query database for contractors",
			"error", err,
			"account", account)
	}

	for _, job := range jobs {
		si.jobIndex[account].Index(job.Id, &job)

		contractor := &api.Contractor{}
		if contra, ok := contractors[job.ContractorId]; ok {
			contractor = contra
		}

		si.jobIndex[account].Index(job.Id, compositeJob{
			Job:        job,
			Contractor: contractor,
		})
	}

	return
}

// SearchFormulas searches the index for matching terms and then returns formulas which might match
func (si *Search) SearchFormulas(account, searchPhrase string) ([]string, error) {
	if index, ok := si.formulaIndex[account]; ok {

		formulaHits, err := queryIndex(index, strings.ToLower(searchPhrase))
		if err != nil {
			return nil, err
		}

		return formulaHits, nil
	}

	return nil, fmt.Errorf("could not find account: %s", account)
}

// SearchJobs searches the index for matching terms and then returns jobs which might match
func (si *Search) SearchJobs(account, searchPhrase string) ([]string, error) {
	if index, ok := si.jobIndex[account]; ok {

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
	searchPhrase = sanitizeQuery(searchPhrase)

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

//TODO(clintjedwards): Make this better by using better logic for replacing.
// it should be possible to remove the second set used for tracking "done" replacements.
// sanitizeQuery removes bleve syntax special chars. This prevents users from shooting themselves
// in the foot when, for example, searching for something like "hello-world" which according to
// bleve syntax should exclude all results including the word "world".
func sanitizeQuery(query string) string {

	reserved := map[rune]struct{}{
		'+': {},
		'-': {},
		'=': {},
		'&': {},
		'|': {},
		'>': {},
		'<': {},
		'!': {},
		'(': {},
		')': {},
		'{': {},
		'}': {},
		'[': {},
		']': {},
		'^': {},
		'"': {},
		'~': {},
		'*': {},
		'?': {},
		':': {},
		'/': {},
	}

	escaped := map[rune]struct{}{}

	for _, char := range query {
		if _, ok := reserved[char]; !ok {
			continue
		}
		if _, ok := escaped[char]; ok {
			continue
		}

		query = strings.Replace(query, string(char), " ", -1)
		escaped[char] = struct{}{}
	}

	return query
}
