package traefik_vault_auth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Vault info

func (v *Vault) login(user string, password string) (string, error) {
	client := &http.Client{}
	authURL := v.authURL(user)
	resp, err := client.PostForm(authURL, url.Values{
		"password": {password},
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("HTTP Response Status:", resp.StatusCode, http.StatusText(resp.StatusCode))
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("Authentication request send to %s failed: status code %d", authURL, resp.StatusCode)
	}
	ar, err := parseAuthResponse(body)
	if err != nil {
		return "", err
	}
	return ar.Auth.ClientToken, nil
}

func (v *Vault) validateToken(token string) (bool, error) {
	var tokenInfo TokenInfo
	reqURL := fmt.Sprintf("%s/%s", v.Address, "/auth/token/lookup-self")
	req, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return false, err
	}
	req.Header.Add("X-VAULT-TOKEN", token)
	resp, err := v.Client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	if err := json.Unmarshal(body, &tokenInfo); err != nil {
		return false, err
	}
	tokenValid := true
	return tokenValid, nil
}

func (v *Vault) renewToken(token string) (bool, error) {
	reqURL := fmt.Sprintf("%s/%s", v.Address, "/auth/token/renew-self")
	fmt.Print(reqURL)
	tokenValid := true
	return tokenValid, nil
}

func (v *Vault) authURL(username string) string {
	return fmt.Sprintf("%s/auth/%s/login/%s", v.Address, v.MountName, username)
}

func parseAuthResponse(body []byte) (AuthResponse, error) {
	var authResponse AuthResponse
	if err := json.Unmarshal(body, &authResponse); err != nil {
		return AuthResponse{}, err
	}
	return authResponse, nil
}

type AuthResponse struct {
	LeaseID       string `json:"lease_id"`
	Renewable     bool   `json:"renewable"`
	LeaseDuration int    `json:"lease_duration"`
	Auth          struct {
		ClientToken string   `json:"client_token"`
		Policies    []string `json:"policies"`
		Metadata    struct {
			Username string `json:"username"`
		} `json:"metadata"`
		LeaseDuration int  `json:"lease_duration"`
		Renewable     bool `json:"renewable"`
	} `json:"auth"`
}

type TokenInfo struct {
	Data struct {
		Accessor         string   `json:"accessor"`
		CreationTime     int      `json:"creation_time"`
		CreationTTL      int      `json:"creation_ttl"`
		DisplayName      string   `json:"display_name"`
		EntityID         string   `json:"entity_id"`
		ExpireTime       string   `json:"expire_time"`
		ExplicitMaxTTL   int      `json:"explicit_max_ttl"`
		ID               string   `json:"id"`
		IdentityPolicies []string `json:"identity_policies"`
		IssueTime        string   `json:"issue_time"`
		Meta             struct {
			Username string `json:"username"`
		} `json:"meta"`
		NumUses   int      `json:"num_uses"`
		Orphan    bool     `json:"orphan"`
		Path      string   `json:"path"`
		Policies  []string `json:"policies"`
		Renewable bool     `json:"renewable"`
		TTL       int      `json:"ttl"`
	} `json:"data"`
}
