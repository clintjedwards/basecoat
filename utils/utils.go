package utils

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/sirupsen/logrus"
)

// SendResponse formats and sends an error message to supplied writer in json format
func SendResponse(w http.ResponseWriter, httpStatusCode int, data interface{}, error bool) error {

	if httpStatusCode != 200 {
		w.WriteHeader(httpStatusCode)
	}

	if error {
		err := json.NewEncoder(w).Encode(struct {
			StatusText string      `json:"status_text"`
			Message    interface{} `json:"message"`
		}{http.StatusText(httpStatusCode), data})
		if err != nil {
			return err
		}
		return nil
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(data)

	return nil
}

// ParseJSON json request into interface
func ParseJSON(rc io.ReadCloser, object interface{}) error {
	decoder := json.NewDecoder(rc)
	err := decoder.Decode(object)
	if err != nil {
		return err
	}
	return nil
}

// StructuredLog allows the application to log to stdout, json formatted,
//  levels accepted are debug, info, warn, error, and fatal
func StructuredLog(level, description string, object interface{}) {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	logger := logrus.WithFields(logrus.Fields{
		"data": object,
	})

	switch level {
	case "debug":
		logger.Debugln(description)
	case "info":
		logger.Infoln(description)
	case "warn":
		logger.Warnln(description)
	case "error":
		logger.Errorln(description)
	case "fatal":
		logger.Fatalln(description)
	default:
		logger.Infoln(description)
	}
}

// RemoveIntFromList removes an element form an array of ints
// does not preserve list order
func RemoveIntFromList(list []int, value int) []int {
	for index, item := range list {
		if item == value {
			list[index] = list[len(list)-1]
			return list[:len(list)-1]
		}
	}

	return list
}

// FindListDifference returns list elements that are in list A
// but not found in B
func FindListDifference(a, b []int) []int {
	m := make(map[int]bool)
	diff := []int{}

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return diff
}

// FindListUpdates is used to compare a new and old version of lists
// it will compare the old version to the new version and return
// which elements have been added or removed from the new version
func FindListUpdates(oldList []int, newList []int) (additions []int, removals []int) {

	removals = FindListDifference(oldList, newList)
	additions = FindListDifference(newList, oldList)

	return additions, removals
}
