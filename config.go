package traefik_vault_auth

import "net/http"

// Config for the plugin configuration.
type Config struct {
	Vault          Vault  `yaml:"vault"`
	CustomRealm    string `yaml:"customRealm"` // CustomRealm can be used to personalize Basic Auth window message
	AddTokenHeader bool   `yaml:"add_token_header"`
}

type Vault struct {
	Address   string `yaml:"url"`
	MountName string `yaml:"mount_name,omitempty"`
	Client    *http.Client
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	c := &Config{}
	return c.addMissingFields()
}

func (c *Config) addMissingFields() *Config {
	if c.CustomRealm == "" {
		c.CustomRealm = "Vault login credentials"
	}
	return c
}
