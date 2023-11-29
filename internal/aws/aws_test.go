package aws

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateHttpClientWithoutKeepAlive(t *testing.T) {
	httpClient := CreateHttpClientWithoutKeepAlive()
	assert.NotNil(t, httpClient)
}

func TestCreateScimClient(t *testing.T) {
	token := "replaceme"
	endpoint := "https://dummy.awsapps.com/start"
	sc := CreateScimClient(endpoint, token)

	assert.NotNil(t, sc)
	assert.NotNil(t, sc.http)
	assert.Equal(t, sc.token, token)
	assert.Equal(t, sc.endpoint, endpoint)
}

func TestScimClient_prepareHttpRequest(t *testing.T) {
	sc := CreateScimClient("", "dummy")
	req, err := http.NewRequest("GET", "http://dummy", nil)
	assert.NoError(t, err)

	err = sc.prepareHttpRequest(req)
	assert.NoError(t, err)
	assert.Equal(t, req.Header.Get("User-Agent"), "aad-finout-sync - github.com/dfds/aad-finout-sync")
	assert.Equal(t, req.Header.Get("Authorization"), "Bearer dummy")

}
