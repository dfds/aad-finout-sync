package finout

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"go.dfds.cloud/aad-finout-sync/internal/util"
	"io"
	"net/http"
	"strings"
	"time"
)

type AuthClientSecret struct {
	config Config
}

type AuthUser struct {
	username    string
	password    string
	totpUrl     string
	tokenClient *util.TokenClient
}

func AuthUserMethod(username string, password string, totpVal *string) *AuthUser {
	var url string
	if totpVal != nil {
		url = *totpVal
	} else {
		url = ""
	}
	method := &AuthUser{
		username:    username,
		password:    password,
		totpUrl:     url,
		tokenClient: nil,
	}
	method.tokenClient = util.NewTokenClient(method.getNewToken)
	return method
}

func AuthClientSecretMethod(conf Config) *AuthClientSecret {
	method := &AuthClientSecret{
		config: conf,
	}
	return method
}

func (a *AuthUser) AcceptedEndpoint(val string) bool {
	return strings.Contains(val, AUTH_API_ENDPOINT)
}

func (a *AuthUser) PrepareHttpRequest(h *http.Request) error {
	if a.tokenClient.Token.IsExpired() {
		err := a.Refresh()
		if err != nil {
			return err
		}
	}

	h.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.tokenClient.Token.GetToken()))

	return nil
}

func (a *AuthUser) Refresh() error {
	err := a.tokenClient.RefreshAuth()
	return err
}

func (a *AuthUser) getNewToken() (*util.RefreshAuthResponse, error) {
	userLoginResp, err := a.login()
	if err != nil {
		return nil, err
	}

	var payload *util.RefreshAuthResponse

	if userLoginResp.MfaRequired {
		mfaResp, err := a.verifyMfa(userLoginResp.MfaToken)
		if err != nil {
			return nil, err
		}
		payload = &util.RefreshAuthResponse{
			TokenType:    "Bearer",
			ExpiresIn:    int64(mfaResp.ExpiresIn),
			ExtExpiresIn: int64(mfaResp.ExpiresIn),
			AccessToken:  mfaResp.AccessToken,
		}
	} else {
		payload = &util.RefreshAuthResponse{
			TokenType:    "Bearer",
			ExpiresIn:    int64(userLoginResp.ExpiresIn),
			ExtExpiresIn: int64(userLoginResp.ExpiresIn),
			AccessToken:  userLoginResp.AccessToken,
		}
	}

	return payload, nil
}

func (a *AuthUser) verifyMfa(mfaToken string) (*VerifyMfaResponse, error) {
	client := http.DefaultClient

	otpCode, err := a.generateOTP(a.totpUrl)
	if err != nil {
		return nil, err
	}

	payload := VerifyMfaRequest{
		MfaToken:       mfaToken,
		Value:          otpCode,
		RememberDevice: false,
	}

	serialisedPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/frontegg/identity/resources/auth/v1/user/mfa/verify", AUTH_API_ENDPOINT), bytes.NewBuffer(serialisedPayload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "aad-finout-sync - github.com/dfds/aad-finout-sync")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var mfaResponse *VerifyMfaResponse

	err = json.Unmarshal(data, &mfaResponse)
	if err != nil {
		return nil, err
	}

	return mfaResponse, nil
}

func (a *AuthUser) login() (*UserLoginResponse, error) {
	client := http.DefaultClient
	payload := UserLoginRequest{
		Email:           a.username,
		Password:        a.password,
		InvitationToken: "",
	}
	serialisedPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/frontegg/identity/resources/auth/v1/user", AUTH_API_ENDPOINT), bytes.NewBuffer(serialisedPayload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "aad-finout-sync - github.com/dfds/aad-finout-sync")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userResponse *UserLoginResponse

	err = json.Unmarshal(data, &userResponse)
	if err != nil {
		return nil, err
	}

	return userResponse, nil
}

func (a *AuthUser) generateOTP(url string) (string, error) {
	key, err := otp.NewKeyFromURL(url)
	if err != nil {
		return "", err
	}

	otpCode, err := totp.GenerateCode(key.Secret(), time.Now())
	if err != nil {
		return "", err
	}

	return otpCode, nil
}

func (a *AuthClientSecret) PrepareHttpRequest(h *http.Request) error {
	h.Header.Set("x-finout-client-id", a.config.ClientId)
	h.Header.Set("x-finout-secret-key", a.config.ClientSecret)

	return nil
}

func (a *AuthClientSecret) Refresh() error {
	return nil
}

func (a *AuthClientSecret) AcceptedEndpoint(val string) bool {
	return strings.Contains(val, APP_API_ENDPOINT)
}
