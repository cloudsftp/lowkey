package main

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLowkeyAPI(t *testing.T) {
	res, err := http.Get("http://lowkey-api:6670/hey")
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)

	body, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, "hi :)\n", string(body))
}
