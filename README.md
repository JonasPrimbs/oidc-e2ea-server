# OpenID Connect GOes E2EA

This repository provides the proof-of-concept implementation of the end-to-end authentication for an OpenID Connect OpenID Provider server.

The provided implementation is written in [Go](https://golang.org/) and implements a REST endpoint for the Identity Certification Token UserInfo endpoint.
Using a reverse proxy in front, it can be mounted to any OpenID Provider implementation.

**Warning:**
Keep in mind that this is the implementation of a research project!
We do not guarantee a secure implementation!
**Do not use this in production!!!**


## Documentation

This section provides an introduction to the architecture and the configuration of the Identity Certification Token Endpoint.


### Architecture

The following figure shows the overall architecture how to use the provided Identity Certification Token (ICT) Endpoint with any OpenID Provider implementation.

```
                           +---------+                +----------+           +----------+
   /*                      |         |       *        |  OpenID  |           |   User   |
  -----------------------> |         |--------------->| Provider | <-------> | Database |
                           |         |                +----------+           +----------+
                           |         |                  ^
                           | Reverse |                  | /realms/ict/protocol/openid-
                           |  Proxy  |                  |   connect/userinfo
   /realms/ict/protocol/   |         |                +----------+
     openid-connect/ict    |         |       /        |   ICT    |
  -----------------------> |         |--------------->| Endpoint |
                           +---------+                +----------+
```

The Docker Compose composition provided [here](./docker-compose.yaml) uses the following implementations:

- Reverse Proxy: [Traefik Proxy](https://traefik.io/traefik/)
- OpenID Provider: [Keycloak](https://www.keycloak.org/)
- User Database: [PostgreSQL](https://www.postgresql.org/)
- Identity Certification Token (ICT) Endpoint: An HTTP endpoint written in [GO](https://go.dev/)


### Server Configuration

This section describes the configuration parameters of the ICT Endpoint.
They are applied by injecting them as environment variables to the running application.
This can be done by defining the variables in the Docker container or by placing an `.env` file in the execution directory.


#### Key File

Absolute or relative file path to the OpenID Provider's private key file in PEM format.

Example:
```bash
KEY_FILE="/path/to/private_key.pem"
```

Setting this variable is **required**.


#### Key ID

The ID of the OpenID Provider's Public Key provided in the `jwks_uri` endpoint.

Example 1:
```bash
KID="rojPQoDRx_DD-DFs7y45wDLl5T8b9VmX6iQapIK6cRE"
```

Example 2:
```bash
KID=1
```

Setting this variable is **required**.


#### Signing Algorithm

Signing algorithm for Identity Certification Token signatures.

Allowed values are:

- `RS256` for RSASSA-PKCS1-v1_5 using SHA-256
- `RS384` for RSASSA-PKCS1-v1_5 using SHA-384
- `RS512` for RSASSA-PKCS1-v1_5 using SHA-512
- `ES256` for ECDSA using P-256 and SHA-256 (recommended)
- `ES384` for ECDSA using P-384 and SHA-384
- `ES512` for ECDSA using P-521 and SHA-512
- `EdDSA` for Eduard Digital Signing Algorithm using Ed25519 curve

Default Value: `ES256`.

Example:
```bash
ALG="ES256"
```


#### Userinfo Endpoint

Absolute URI to the OpenID Provider's Userinfo Endpoint.

This URI is used by the ICT Endpoint to request the claims of the Identity Certification Token from the OpenID Provider.
**Make sure that the running ICT Endpoint can access the OpenID Provider's Userinfo Endpoint via this URI!**

Example 1:
```bash
USERINFO="https://openid-provider.sample.org/userinfo"
```

Example 2 (Keycloak):
```bash
USERINFO="http://localhost:8080/realms/ict/protocol/openid-connect/userinfo"
```

Setting this variable is **required**.


#### Token Introspection Endpoint

Absolute URI to the OpenID Provider's Token Introspection Endpoint described in [RFC 7662](https://datatracker.ietf.org/doc/html/rfc7662).
This URI is provided on the Discovery Endpoint as attribute `introspection_endpoint`.
**Make sure that the running ICT Endpoint can access the OpenID Provider's Token Introspection Endpoint via this URI!**

Example 1:
```bash
TOKEN_INTROSPECTION="https://openid-provider.sample.org/introspect"
```

Example 2 (Keycloak):
```bash
TOKEN_INTROSPECTION="http://localhost:8080/realms/ict/protocol/openid-connect/token/introspect"
```

Setting this variable is **required**.


#### Token Introspection Host

The hostname in HTTP Host Header when requesting the Token Introspection Endpoint.
If not provided, the hostname from the `TOKEN_INTROSPECTION` URL will be used.

Example:
```bash
TOKEN_INTROSPECTION_HOST="openid-provider.sample.org
```


#### Token Introspection Credentials

The HTTP Authorization Header required for the Token Introspection Endpoint.
Typically a HTTP Basic Authentication Header using the ICT Endpoint's Client Credentials.

Example:
```bash
INTROSPECTION_CREDENTIALS="Basic aWN0X2VuZHBvaW50OlMzY3JldCE="
```
For Client ID `ict_endpoint` and Client Secret `S3cret!`

Setting this variable is **required**, except the OpenID Provider does not requires any authorization for the Token Introspection Endpoint (not recommended).


#### Context Prefix

Prefix of scopes which indicate the granted end-to-end authentication context.

Default Value: `e2e_ctx_`.

Example:
```bash
CONTEXT_PREFIX="ctx_"
```
The scope `ctx_email` will authorize the End-User for the `email` context.


#### Issuer Claim

The Identity Certification Token's Issuer.

This is the value of the `iss` claim of the issued Identity Certification Token.
Typically, this is the public URI of the OpenID Provider where `.well-known/openid-configuration` is added to request the OpenID configuration.

Example 1:
```bash
ISSUER="https://accounts.sample.org/"
```

Example 2:
```bash
ISSUER="http://localhost:8080/realms/ict"
```

Setting this variable is **required**.


#### Token Validity Period

The Identity Certification Token's default validity period in seconds.

Default Value: `3600` (1 hour).

Example:
```bash
DEFAULT_TOKEN_PERIOD=3600
```


#### Maximum Token Validity Period

The Identity Certification Token's maximum validity period in seconds.
If the requested token period is longer than this value, this value is used.

Default Value: `2592000` (30 days).

Example:
```bash
MAX_TOKEN_PERIOD=2592000
```


#### Port

The Port where the endpoint is running on.

Default Value: `8080`.

Example:
```bash
PORT=8080
```


#### Database File

The SQLite database file to store used nonce values in.

Default Value (standalone): `./db.sqlite`
<br>
Default Value (Docker image): `/config/db.sqlite`

Example:
```bash
DB_SQLITE_FILE="/config/db.sqlite"
```


### REST Endpoint

The REST API is described in the OpenAPI format provided [here](./docs/openapi.yaml).


### Environment Setup

To setup a test environment locally, refer to the manual [here](./docs-dev/environment-setup.md).


### Testing

To play around with the API, check out the testing manual [here](./docs-dev/testing.md).


### Security

To improve the security of a deployment, check out the security manual [here](./docs-dev/security.md).
