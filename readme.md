# Traefik plugin: Vault Auth

## Development

### Prerequisites

* A valid [Traefik Pilot](https://pilot.traefik.io) token for your Traefik instance.
* A running Vault server with LDAP authentication configured

### Docker Compose

```shell
TRAEFIK_PILOT_TOKEN="xxxx" VAULT_ADDR=http://127.0.0.1:8200 VAULT_AUTH_MOUNT_NAME=ldap docker-compose up -d --build 
```

This will launch:

* A Traefik instance with [dashboard](http://traefik.localhost)
* A [`whoami` instance](http://whoami.localhost)

## Installation

Declare it in the Traefik configuration:

```yaml
pilot:
  token: "xxxx"
experimental:
  plugins:
    traefik-vault-ldap-auth:
      moduleName: "github.com/rigrassm/traefik-vault-ldap-auth"
      version: "v0.2.5"
```

### Configuration

```yaml
    middlewares:
      my-traefik-vault-auth:
        plugin:
          traefik-vault-auth:
            customRealm: Use a valid Vault user to authenticate
            vault:
              address: http://127.0.0.1:8200
              mount_name: ldap
              add_token_header: true
```
