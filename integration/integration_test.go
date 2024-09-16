package main

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegration(t *testing.T) {
	res, err := http.Get("http://lowkey-api:6670/hey")
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, res.StatusCode)

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	require.Equal(t, "hi :)\n", string(body))

	require.True(t, true)
}
