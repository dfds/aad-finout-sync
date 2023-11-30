package ssu

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.dfds.cloud/aad-finout-sync/internal/util"
)

var capsvcTestConfig = Config{
	Host:         "http://dummy",
	TenantId:     "00000-00000-00000-00000-00000",
	ClientId:     "dummy",
	ClientSecret: "secretdummy",
	Scope:        "openid,profile",
}

func TestClient_prepareHttpRequest(t *testing.T) {
	capsvc := NewSsuClient(capsvcTestConfig)
	req, err := http.NewRequest("GET", "http://dummy", nil)
	assert.NoError(t, err)

	capsvc.tokenClient = util.NewTokenClient(func() (*util.RefreshAuthResponse, error) {
		return &util.RefreshAuthResponse{
			TokenType:    "",
			ExpiresIn:    time.Now().Add(time.Minute * 100).Unix(),
			ExtExpiresIn: time.Now().Add(time.Minute * 100).Unix(),
			AccessToken:  "dummy",
		}, nil
	})

	err = capsvc.prepareHttpRequest(req)
	assert.NoError(t, err)
	assert.Equal(t, req.Header.Get("User-Agent"), "aad-finout-sync - github.com/dfds/aad-finout-sync")
	assert.Equal(t, req.Header.Get("Authorization"), "Bearer dummy")

}

func TestNewCapSvcClient(t *testing.T) {
	capsvc := NewSsuClient(capsvcTestConfig)
	assert.NotNil(t, capsvc)
}
