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
