package utils

import (
	"net/http"
	"net/http/httputil"

	"github.com/sirupsen/logrus"
)

func DumpHTTPRequest(req *http.Request, description string) {
	body := false

	if req.Method == http.MethodPost {
		body = true
	}

	dump, err := httputil.DumpRequest(req, body)
	if err != nil {
		logrus.Fatal(err)
	}

	// if there's no body, we don't need extra \n\n
	if !body {
		logrus.Debug(description, "\n\n", string(dump))
		return
	}

	logrus.Debug(description, "\n\n", string(dump), "\n\n")
}

func DumpHTTPResponse(res *http.Response, description string) {
	dump, err := httputil.DumpResponse(res, true)
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Debug(description, "\n\n", string(dump), "\n\n")
}
