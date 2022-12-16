# OpenID Connect GOes E2EA

This repository provides the proof-of-concept implementation of the end-to-end authentication for an OpenID Connect OpenID Provider server.

The provided implementation is written in [Go](https://golang.org/) and implements a REST endpoint for the ID Assertion Token UserInfo endpoint.
Using a reverse proxy in front, it can be mounted to any OpenID Provider implementation.

**Warning:**
Keep in mind that this is the implementation of a research project!
We do not guarantee a secure implementation!
**Do not use this in production!!!**


## Documentation

This section provides an introduction to the architecture and the configuration of the ID Assertion Token Endpoint.


### Architecture

The following figure shows the overall architecture how to use the provided ID Assertion Token (IAT) Endpoint with any OpenID Provider implementation.

```
                                                      +---------+                +----------+           +----------+
   /*                                                 |         |       *        |  OpenID  |           |   User   |
  --------------------------------------------------> |         |--------------->| Provider | <-------> | Database |
                                                      |         |                +----------+           +----------+
                                                      | Reverse |
                                                      |  Proxy  |
                                                      |         |                +----------+
   /realms/test/protocol/openid-connect/userinfo/iat  |         |       /        |   IAT    |
  --------------------------------------------------> |         |--------------->| Endpoint |
                                                      +---------+                +----------+
```

The Docker Compose composition provided [here](./docker-compose.yaml) uses the following implementations:

- Reverse Proxy: [Traefik Proxy](https://traefik.io/traefik/)
- OpenID Provider: [Keycloak](https://www.keycloak.org/)
- User Database: [PostgreSQL](https://www.postgresql.org/)
- ID Assertion Token (IAT) Endpoint: An HTTP endpoint written in [GO](https://go.dev/)


### Server Configuration

This section describes the configuration parameters of the IAT Endpoint.
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

Signing algorithm for ID Assertion Token signatures.

Allowed values are:

- `RS256` for RSASSA-PKCS1-v1_5 using SHA-256
- `RS384` for RSASSA-PKCS1-v1_5 using SHA-384
- `RS512` for RSASSA-PKCS1-v1_5 using SHA-512
- `ES256` for ECDSA using P-256 and SHA-256 (recommended)
- `ES384` for ECDSA using P-384 and SHA-384
- `ES512` for ECDSA using P-521 and SHA-512

Default Value: `ES256`.

Example:
```bash
ALG="ES256"
```


#### Userinfo Endpoint

Absolute URI to the OpenID Provider's Userinfo Endpoint.

This URI is used by the IAT Endpoint to request the claims of the ID Assertion Token from the OpenID Provider.
**Make sure that the running IAT Endpoint can access the OpenID Provider's Userinfo Endpoint via this URI!**

Example 1:
```bash
USERINFO="https://openid-provider.sample.org/userinfo"
```

Example 2 (Keycloak):
```bash
USERINFO="http://localhost:8080/realms/iat/protocol/openid-connect/userinfo"
```

Setting this variable is **required**.


#### Issuer Claim

The ID Assertion Token's Issuer.

This is the value of the `iss` claim of the issued ID Assertion Token.
Typically, this is the public URI of the OpenID Provider where `.well-known/openid-configuration` is added to request the OpenID configuration.

Example 1:
```bash
ISSUER="https://accounts.sample.org/"
```

Example 2:
```bash
ISSUER="http://localhost:8080/realms/iat"
```

Setting this variable is **required**.


#### Token Validity Period

The ID Assertion Token's default validity period in seconds.

Default Value: `3600` (1 hour).

Example:
```bash
DEFAULT_TOKEN_PERIOD=3600
```


#### Maximum Token Validity Period

The ID Assertion Token's maximum validity period in seconds.
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
