package api

import "errors"

var errFormulaNotFound = errors.New("formula not found")
var errFormulaExists = errors.New("formula already exists")

var errJobNotFound = errors.New("job not found")
var errJobExists = errors.New("job already exists")

var errJSONParseFailure = errors.New("could not parse json")
