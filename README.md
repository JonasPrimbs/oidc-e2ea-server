# Open Identity Certification with OpenID Connect (OIDC²) Middleware
A middleware written in [Go](https://go.dev/) to implement OIDC² as a middleware in front of any OpenID Provider.

**Warning:**
This implementation is a research project!
We do not guarantee a secure implementation!
**Do not use this in production!!!**

## 1. How it works
Use a reverse proxy in front of your OpenID Provider and reroute token requests through the middleware.

The middleware will forward all regular token requests to the OpenID provider and pass through its response (see Figure 1).

```
                   +---------+                +--------------------------------+
                   |         |  /*            |                                |
                   |         | -------------> |         OpenID Provider        |
                   |         | <------------- |                                |
  +--------+       |         |                +--------------------------------+
  |        |  /*   | Reverse |                            ^        |
  | Client | ----> |  Proxy  |                            | /token |
  |        | <---- |         |                            |        V
  +--------+       |         |                +--------------------------------+
                   |         |  /token -> /   |                                |
                   |         | -------------> |        OIDC² Middleware        |
                   |         | <------------- |                                |
                   +---------+                +--------------------------------+
```
**Figure 1**: Regular token request.

If the client performs a token exchange request `requested_token_type=urn:ietf:params:oauth:token-type:ic-token`, the middleware performs the following steps (see Figure 2):

1. It validates the access token provided in `subject_token` using the OpenID Provider's Token Introspection Endpoint (see [RFC 7662](https://datatracker.ietf.org/doc/rfc7662/)).
2. If the access token is valid, it generates an Identity Certification Token using the End User's identity claims provided by the OpenID Provider from the UserInfo Endpoint (see [OIDC Core](https://openid.net/specs/openid-connect-core-1_0.html#UserInfo)).

```
                   +---------+                +--------------------------------+
                   |         |  /*            |                                |
                   |         | -------------> |         OpenID Provider        |
                   |         | <------------- |                                |
  +--------+       |         |                +--------------------------------+
  |        |  /*   | Reverse |                  ^            |   ^           |
  | Client | ----> |  Proxy  |                  | /tokeninfo |   | /userinfo |
  |        | <---- |         |                  |            V   |           V
  +--------+       |         |                +--------------------------------+
                   |         |  /token -> /   |                                |
                   |         | -------------> |        OIDC² Middleware        |
                   |         | <------------- |                                |
                   +---------+                +--------------------------------+
```
**Figure 2**: Identity Certification Token request.

## 2. Setup
The provided [Docker composition](./docker-compose.yaml) contains an example for the following software:

- Reverse Proxy: [Traefik Proxy](https://traefik.io/traefik/)
- OpenID Provider: [Keycloak](https://www.keycloak.org/)
- User Database: [PostgreSQL](https://www.postgresql.org/)

First, copy the file [`example.env`](./example.env`) to `.env` and adjust its [configuration](#configuration) parameters:
```bash
copy example.env .env
nano .env
```

Then, generate your secret files:

```bash
./scripts/generate-secrets.sh
```

If you run the service locally, you should import the generated root certificate from `./.secrets/ca_root.crt` to your browser and you have to configure your local certicate authority:
```bash
./scripts/generate-ca-config.sh
```

Then, use the following command to run the composition in the required profile:
```bash
# Start the composition in detached mode (-d):
docker compose --profile "test" up -d
```

Available profiles are:
- `prod`: Public production mode
- `test`: Local test mode
- `dev`: Local development mode

To stop the composition, use the following command:
```bash
# Stop the composition:
docker compose --profile "test" down
```

See step-by-step guide [here](./docs-dev/setup.md).

## 3. Configuration
Apply the following configuration options in your `.env` file.

### 3.1. Key Configuration
To issue a valid Identity Certification Token, the middleware needs a private signing key whose corresponding public key is listed on the OpenID Provider's `jwks_uri`.

#### 3.1.1. Key File
The PEM-encoded private signing key file path.

Example:
```bash
KEY_FILE="/path/to/private_key.pem"
```

This variable is **required**.

#### 3.1.2. Key ID
The key ID of the public key listed on the OpenID Provider's `jwks_uri`.

Example 1:
```bash
KEY_ID="rojPQoDRx_DD-DFs7y45wDLl5T8b9VmX6iQapIK6cRE"
```

Example 2:
```bash
KEY_ID=1
```

This variable is **required**.

### 3.2. OpenID Provider Configuration
To issue valid Identity Certification Tokens, the middleware must know the OpenID Provider's issuer identifier.
It also needs to know the OpenID Provider's Token Introspection and UserInfo Endpoint to validate token requests and issue correct identity claims.

#### 3.2.1. Issuer Claim
The issuer identifier (`iss` claim) of issued Identity Certification Tokens.

Example 1:
```bash
ISSUER="https://op.example.com/"
```

Example 2 (Keycloak running locally):
```bash
ISSUER="http://localhost:8080/realms/oidc2"
```

This variable is **required**.

#### 3.2.2. Token Introspection Endpoint
The absolute URI of the OpenID Provider's Token Introspection Endpoint to validate provided access tokens.

Example 1:
```bash
INTROSPECTION_ENDPOINT="https://op.example.com/tokeninfo"
```

Example 2 (Keycloak running locally):
```bash
INTROSPECTION_ENDPOINT="http://localhost:8080/realms/ict/protocol/openid-connect/token/introspect"
```

If not provided, its value falls back to `token_introspection_endpoint` value of the discovery document provided in `{ISSUER}/.well-known/openid-configuration`.

#### 3.2.3. UserInfo Endpoint
The absolute URI of the OpenID Provider's UserInfo Endpoint to observe the requesting user's identity claims for the Identity Certification Token:

Example 1:
```bash
USERINFO_ENDPOINT="https://op.example.com/userinfo"
```

Example 2 (Keycloak running locally):
```bash
USERINFO_ENDPOINT="http://localhost:8080/realms/ict/protocol/openid-connect/userinfo"
```

If not provided, its value falls back to the `userinfo_endpoint` value of the discovery document provided in `{ISSUER}/.well-known/openid-configuration`.

#### 3.2.4. Token Endpoint
The absolute URI of the OpenID Provider's Token Endpoint to forward the token request if no Identity Certification Token is requested.

Example 1:
```bash
TOKEN_ENDPOINT="https://op.example.com/token"
```

Example 2 (Keycloak running locally):
```bash
TOKEN_ENDPOINT="http://localhost:8080/realms/ict/protocol/openid-connect/token"
```

If not provided, its value falls back to the `token_endpoint` value of the discovery document provided in `{ISSUER}/.well-known/openid-configuration`.

### 3.3. Token Configuration

#### 3.3.1. Signing Algorithm
The JSON Web Algorithms (JWA) used for signing the Identity Certification Token.

Available options:

- `RS256` for RSASSA-PKCS1-v1_5 using SHA-256
- `RS384` for RSASSA-PKCS1-v1_5 using SHA-384
- `RS512` for RSASSA-PKCS1-v1_5 using SHA-512
- `ES256` for ECDSA using P-256 and SHA-256 (recommended)
- `ES384` for ECDSA using P-384 and SHA-384
- `ES512` for ECDSA using P-521 and SHA-512
- `EdDSA` for Eduard Digital Signing Algorithm using Ed25519 curve

Example:
```bash
ALG="ES256"
```

Default value is `RS256`.

#### 3.3.2. Token Validity Period
The Identity Certification Token's validity period in seconds.

Example:
```bash
TOKEN_PERIOD=300
```

Default value is `300` (5 minutes).

### 3.4. Hosting Configuration

#### 3.4.1. Port
The port to listen to.

Example:
```bash
PORT=8080
```

Default value is `8080`.

### 3.4.2. Hosts
Comma-separated hosts in CIDR format to listen to.

Example (localhost only):
```bash
HOSTS=127.0.0.1,::1
```

Default value is `0.0.0.0` (all IPv4 hosts).

## 4. REST Documentation
The [OpenAPI](https://swagger.io/specification/) documentation of the RESTful API is provided [here](./api/oidc2middleware.yaml).

## 5. Testing
To play around with the API, check out the testing manual [here](./docs-dev/testing.md).
