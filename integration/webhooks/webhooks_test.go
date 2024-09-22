package webhooks

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

type WebhooksTestSuite struct {
	suite.Suite
}

func TestWebhooks(t *testing.T) {
	suite.Run(t, new(WebhooksTestSuite))
}

type ExtensionResponse struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
	Successful bool   `json:"successful"`
}

func (s *WebhooksTestSuite) TestBla() {
	body := bytes.NewBufferString("{}")

	response, err := http.Post(
		"http://local-dev:8080/internal/calls/extensionAddedToContext",
		"application/json",
		body,
	)
	s.NoError(err)

	s.Equal(http.StatusOK, response.StatusCode)

	responseBody, err := io.ReadAll(response.Body)
	s.NoError(err)

	var extensionResponse ExtensionResponse
	err = json.Unmarshal(responseBody, &extensionResponse)
	s.NoError(err)

	s.Equal(
		ExtensionResponse{
			Successful: true,
			StatusCode: http.StatusOK,
			Body:       "Ok",
		},
		extensionResponse,
	)
}
