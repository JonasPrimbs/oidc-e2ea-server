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

### 1.1. ICT and DPoP (strict policy)

For **ICT token exchange** (`grant_type=urn:ietf:params:oauth:grant-type:token-exchange` and `requested_token_type=urn:ietf:params:oauth:token-type:ic_token`, or omitted `requested_token_type`), the middleware enforces:

1. **Subject token** must be a **DPoP-bound access token** (`cnf.jkt` in introspection or in the JWT payload).
2. The client must send a valid **`DPoP` proof** for the ICT POST (`htm=POST`, `htu` = middleware token URL).
3. The proof’s public key must match the subject token’s `cnf.jkt`.

Bearer-only subject tokens are **rejected** for ICT. The issued ICT always includes `cnf.jkt` carrying forward the binding.

All **other** token requests (authorization code, refresh, non-ICT token exchange, etc.) are **proxied to Keycloak** as-is; DPoP is optional on those forwarded calls.

## 2. Setup

There are two ways to run the middleware locally:

| Setup | Compose file | OpenID Provider | Typical use |
|-------|--------------|-----------------|-------------|
| **Full stack** | [`docker-compose.yaml`](./docker-compose.yaml) | Local Keycloak + Traefik + internal CA | Integration tests, Swagger/Postman against `op.localhost` |
| **Middleware only** | [`compose-dev.yaml`](./compose-dev.yaml) | External / deployed Keycloak (e.g. `https://sso.koala.primbs.dev/realms/koala`) | Koala client dev, ICT exchange against a real realm |

### 2.1. Full stack (`docker-compose.yaml`)

The composition includes:

- Reverse Proxy: [Traefik Proxy](https://traefik.io/traefik/)
- OpenID Provider: [Keycloak](https://www.keycloak.org/)
- User Database: [PostgreSQL](https://www.postgresql.org/)

First, copy the file [`example.env`](./example.env`) to `.env` and adjust its [configuration](#configuration) parameters:
```bash
# Local test stack (internal CA, op.localhost):
cp example-test.env .env

# Or production-style stack:
# cp example-prod.env .env
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
docker compose --profile test up -d   # local test
docker compose --profile dev up -d    # local dev (builds middleware from source)
docker compose --profile prod up -d   # production-style
```

To stop the composition, use the following command:
```bash
# Stop the composition:
docker compose --profile test down
```

See the step-by-step guide for the full stack [here](./docs-dev/setup.md).

### 2.2. Middleware only (`compose-dev.yaml`)

Use this when Keycloak already runs elsewhere and you only need the OIDC^2 middleware on your machine (e.g. Koala client with `LOCAL_OIDC2_ICT_ENDPOINT=http://localhost:8082`).

1. Copy [`example-dev.env`](./example-dev.env) to `.env` and set `ISSUER`, the three endpoint URLs, `KEY_ID`, `INTROSPECTION_CREDENTIALS`, and `ALLOWED_ORIGINS` for your realm.
2. Start the middleware:

```bash
docker compose -f compose-dev.yaml up --build -d oidc2middleware
```

The service listens on **host port 8082** (`8082:8080`). POST token requests to `http://localhost:8082/`.

The optional `app` service in the same file is a [Docker Dev Environment](https://docs.docker.com/dev/) shell (`dev.Dockerfile`); start it only if you need that IDE container:

```bash
docker compose -f compose-dev.yaml up -d app
```

See [docs-dev/setup.md](./docs-dev/setup.md#middleware-only-external-keycloak) for details.

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

#### 3.3.3. Authorized-party context mapping (`ATH_CTX`)

When set, the middleware adds a `ctx` claim to issued ICTs based on the subject token's authorized party (`azp` from token introspection, falling back to `client_id`).

The value is a JSON object mapping `azp` strings to arrays of context strings.

If `ATH_CTX` is configured and the subject token's authorized party has no entry, the ICT exchange is rejected.

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

#### 3.4.3. Introspection credentials
HTTP `Authorization` header value for the Token Introspection request (e.g. Keycloak confidential client).

Example:

```bash
INTROSPECTION_CREDENTIALS="Basic <base64(client_id:client_secret)>"
```

#### 3.4.4. Allowed origins (CORS)
Comma-separated list of `Origin` header values allowed for browser clients.

Example:

```bash
ALLOWED_ORIGINS=https://koala-local,https://koala.primbs.dev
```

Optional. Native clients usually do not send `Origin`; this matters mainly for web clients.

#### 3.4.5. DPoP proof max age
Maximum age in seconds for a DPoP proof `iat` claim when the subject token is DPoP-bound.

Example:

```bash
DPOP_MAX_AGE=300
```

Default is `300`.

## 4. REST Documentation
The [OpenAPI](https://swagger.io/specification/) documentation of the RESTful API is provided [here](./api/oidc2middleware.yaml).

## 5. Testing
To play around with the API, check out the testing manual [here](./docs-dev/testing.md).
