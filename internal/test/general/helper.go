package general_helpers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	bookmark_handler "github.com/bookmark-project-learn/bookmark-service/internal/handler/bookmark"
)

func MakeJSONRequest(method, target string, body interface{}) *http.Request {
	var bodyBytes []byte
	switch typedBody := body.(type) {
	case nil:
		bodyBytes = nil
	case string:
		bodyBytes = []byte(typedBody)
	default:
		bodyBytes, _ = json.Marshal(typedBody)
	}

	req := httptest.NewRequest(method, target, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func MakeMultipartRequest(target string, hasfile bool, content, fileContentType string) (*http.Request, error) {
	bodyBytes, writer, err := bookmark_handler.CreateMultipartWriterWithContentString(hasfile, fileContentType, content)
	if err != nil {
		return nil, err
	}
	req := httptest.NewRequest("POST", target, bodyBytes)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}
