package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLowkeyAPI(t *testing.T) {
	res, err := http.Get("http://lowkey-api:6670/hey")
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, res.StatusCode)

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, "hi :)\n", string(body))
}

type ExtensionResponse struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
	Successful bool   `json:"successful"`
}

func TestDevServer(t *testing.T) {
	body := bytes.NewBufferString("{}")

	response, err := http.Post(
		"http://local-dev:8080/internal/calls/extensionAddedToContext",
		"application/json",
		body,
	)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, response.StatusCode)

	responseBody, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	log.Println(string(responseBody))

	var extensionResponse ExtensionResponse
	err = json.Unmarshal(responseBody, &extensionResponse)
	require.NoError(t, err)

	require.Equal(
		t,
		ExtensionResponse{
			Successful: true,
			StatusCode: http.StatusOK,
			Body:       "Ok",
		},
		extensionResponse,
	)
}
