package traefik_vault_ldap_auth

import (
	"context"
	"fmt"
	"net/http"
)

// VaultAuth a Traefik Basic Auth middleware backed by Vault LDAP.
type VaultAuth struct {
	next   http.Handler
	name   string
	config *Config
}

// New creates a new VaultAuth plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &VaultAuth{
		next:   next,
		name:   name,
		config: config,
	}, nil
}

// Use traefiks basic auth middleware as a guide
// https://github.com/traefik/traefik/blob/master/pkg/middlewares/auth/basic_auth.go
func (va *VaultAuth) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var vaultToken string
	var err error
	tokenValid := false

	vaultToken = req.Header.Get("X-VAULT-TOKEN")
	// Check if we have a Vault token
	if vaultToken != "" {
		tokenValid, _ = va.config.Vault.validateToken(vaultToken)
		if tokenValid {
			// Renewal Check and Process
			fmt.Print("Do renewal process")
		}
	}

	if !tokenValid {
		user, pass, ok := req.BasicAuth()
		if !ok {
			// No valid 'Authentication: Basic xxxx' header found in request
			rw.Header().Set("WWW-Authenticate", `Basic realm="`+va.config.CustomRealm+`"`)
			http.Error(rw, "Unauthorized.", http.StatusUnauthorized)
			return
		}
		if vaultToken, err = va.config.Vault.login(user, pass); err != nil {
			// Failed to login with provided user/pass
			rw.Header().Set("WWW-Authenticate", `Basic realm="`+va.config.CustomRealm+`"`)
			http.Error(rw, "Unauthorized.", http.StatusUnauthorized)
			return
		}
		if va.config.AddTokenHeader {
			req.Header.Set("X-VAULT-TOKEN", vaultToken)
		}
	}
	va.next.ServeHTTP(rw, req)
}
