# OpenID Connect GOes E2EA

This repository provides the proof-of-concept implementation of the end-to-end authentication for an OpenID Connect OpenID Provider server.

The provided implementation is written in [Go](https://golang.org/) and implements a REST endpoint for the Remote ID Token UserInfo endpoint.
Using a reverse proxy in front, it can be mounted to any OpenID Provider implementation.

**Warning:**
Keep in mind that this is the implementation of a research project!
We do not guarantee a secure implementation!
**Do not use this in production!!!**


## Table of Contents

- [Documentation](#documentation)
  - [Architecture](#architecture)
  - [Server Configuration](#server-configuration)
  - [Improve Security](#improve-security)
  - [REST Endpoint](#rest-endpoint)
- [Environment Setup](#environment-setup)
  - [1. Clone Repository](#1-clone-repository)
  - [2. Generate Secrets](#2-generate-secrets)
  - [3. Configure Deployment](#3-configure-deployment)
  - [4. Initial Infrastructure Start](#4-initial-infrastructure-start)
  - [5. Setup OpenID Provider](#5-setup-openid-provider)
  - [6. Configure Deployment Mode](#6-configure-deployment-mode)
  - [7. Restart Infrastructure](#7-restart-infrastructure)
- [Testing](#testing)
  - [Request Token JWT Generation](#request-token-jwt-generation)
  - [Testing with API Specification](#testing-with-api-specification)
  - [Testing with Swagger Editor](#testing-with-swagger-editor)
  - [Testing with Postman](#testing-with-postman)


## Documentation

This section provides an introduction to the architecture and the configuration of the Remote ID Token Endpoint.


### Architecture

The following figure shows the overall architecture how to use the provided Remote ID Token (RIDT) Endpoint with any OpenID Provider implementation.

```
                                                       +---------+                +----------+           +----------+
   /*                                                  |         |       *        |  OpenID  |           |   User   |
  ---------------------------------------------------> |         |--------------->| Provider | <-------> | Database |
                                                       |         |                +----------+           +----------+
                                                       | Reverse |
                                                       |  Proxy  |
                                                       |         |                +----------+
   /realms/test/protocol/openid-connect/userinfo/ridt  |         |       /        |   RIDT   |
  ---------------------------------------------------> |         |--------------->| Endpoint |
                                                       +---------+                +----------+
```

The Docker Compose composition provided [here](./docker-compose.yaml) uses the following implementations:

- Reverse Proxy: [Traefik Proxy](https://traefik.io/traefik/)
- OpenID Provider: [Keycloak](https://www.keycloak.org/)
- User Database: [PostgreSQL](https://www.postgresql.org/)
- Remote ID Token (RIDT) Endpoint: An HTTP endpoint written in [GO](https://go.dev/)


### Server Configuration

This section describes the configuration parameters of the RIDT Endpoint.
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

Signing algorithm for Remote ID Token signatures.

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

This URI is used by the RIDT Endpoint to request the claims of the Remote ID Token from the OpenID Provider.
**Make sure that the running RIDT Endpoint can access the OpenID Provider's Userinfo Endpoint via this URI!**

Example 1:
```bash
USERINFO="https://openid-provider.sample.org/userinfo"
```

Example 2 (Keycloak):
```bash
USERINFO="http://localhost:8080/realms/ridt/protocol/openid-connect/userinfo"
```

Setting this variable is **required**.


#### Issuer Claim

The Remote ID Token's Issuer.

This is the value of the `iss` claim of the issued Remote ID Token.
Typically, this is the public URI of the OpenID Provider where `.well-known/openid-configuration` is added to request the OpenID configuration.

Example 1:
```bash
ISSUER="https://accounts.sample.org/"
```

Example 2:
```bash
ISSUER="http://localhost:8080/realms/ridt"
```

Setting this variable is **required**.


#### Token Validity Period

The Remote ID Token's default validity period in seconds.

Default Value: `3600` (1 hour).

Example:
```bash
DEFAULT_TOKEN_PERIOD=3600
```


#### Maximum Token Validity Period

The Remote ID Token's maximum validity period in seconds.
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


### Improve Security

This section provides some hints how to improve the security of the deployment.


#### 1. Use Containers

The container image of the RIDT endpoint is provided on [DockerHub (external URL)](https://hub.docker.com/r/jonasprimbs/oidc-e2ea-server).

Follow the instructions there to deploy the container or use the Docker Compose file [here](./docker-compose.yaml) to run a predefined composition with [Keycloak (external URL)](https://www.keycloak.org/) as OpenID Provider.


### REST Endpoint

The REST API is described in the OpenAPI format provided [here](./docs/openapi.yaml).


## Environment Setup

This section describes how to setup a test environment locally with Docker Compose.

**WARNING: THIS IS FOR TEST PURPOSES ONLY! DO NOT USE THIS IN PRODUCTION!!!**


### 1. Clone Repository

In your Linux bash, clone this repository to your home directory:

```bash
git clone https://github.com/JonasPrimbs/oidc-e2ea-server.git
```

Now navigate to the cloned directory:

```bash
cd oidc-e2ea-server
```


### 2. Generate Secrets

Execute the following command:

```bash
bash ./generate-secrets.sh
```

This will randomly generate all usernames, passwords, and private keys which are unique for your installation and store them in the new directory `.secrets` in the repository.


### 3. Configure Deployment

Go to the generated `/.env` file and configure the following parameters:

- `OP_HOST=<your-hostname>` the host/domain name of your server. Default is `op.localhost`.
- `REALM_NAME=<your-realm-name>` the name of your preferred Keycloak realm. Default is `ridt`.

For a local deployment, you can leave these settings at default.


### 4. Initial Infrastructure Start

Start up your OpenID Provider for the first time using the following command:

```bash
docker compose up -d op
```

This might take a while to download all related container images.


### 5. Setup OpenID Provider


#### 5.1. Login to Keycloak Admin Console

Open your browser and go to [http://op.localhost/admin/ (external URL)](http://op.localhost/admin/) and *sign in* with the credentials generated in the following files:
- Username: `/.secrets/op_username.txt`
- Password: `/.secrets/op_password.txt`


#### 5.2. Create Realm

Create a new realm called `ridt` as follows:

1. Hover the *Master* realm in the navigation bar and click *Add Realm*.
2. Enter the realm *name* `ridt`.
3. Click *Create*.


#### 5.3. Import Realm Settings

Import the realm settings as follows:

1. In your new test realm, go to *Manage* > *Import*.
2. Click on *Select file* and select the file `/keycloak/realm-export.json`.
3. Click *Import*.


#### 5.4. Import Private Key

Import the private key as follows:

1. In your new test realm, go to *Configure* > *Realm Settings* > *Keys* > *Providers*.
2. On top of the table, open the dropdown menu *Add keystore...* and select *rsa*.
3. Enter the *priority* `101`.
4. As *Private RSA Key*, *Select file* `/.secrets/private.pem`.
5. Click *Save*.
6. In *Configure* > *Realm Settings* > *Keys* > *Active*, copy the *Kid* of the new RSA key with priority `101` and paste it to the file `/.secrets/ridt.env` as value for the parameter `KID`, e.g., `KID=rojPQoDRx_DD-DFs7y45wDLl5T8b9VmX6iQapIK6cRE`.


#### 5.5. Create Test User

Create a new test user as follows:

1. In your new test realm, go to *Manage* > *Users*.
2. On top of the table, click the *Add user* button.
3. Enter a *username*, e.g., `test`.
4. It is recommended that you also enter an *Email* address, a *First Name*, and a *Last Name*.
5. Click *Save*.
6. In your new user, go to the *Credentials* tab.
7. Set a *Password*, re-enter the *Password Confirmation*, and set *Temporary* to `OFF`.
8. Then click *Set Password*.
9. Apply with *Set password*.


### 6. Configure Deployment Mode

This step depends on your intention why you run this deployment.

- **Testing**: Choose this mode if you want to just run the deployment for testing purposes.
- **Development**: Choose this mode if you want to change the implementation of the RIDT endpoint application.


#### 6.1. Test Deployment

*Do this step only if you want to run this deployment for **testing** purposes!*

1. Go to `/docker-compose.yaml`.
2. Uncomment line 61 (`image` attribute in service `ridt`).
3. Comment line 64 to 66 (`build` attribute in service `ridt`).


#### 6.2. Development Deployment

*Do this step only if you want to run this deployment for **development** purposes!*

1. Go to `/docker-compose.yaml`.
2. Make sure that line 61 (`image` attribute in service `ridt`) is commented out.
3. Make sure that line 64 to 66 (`build` attribute in service `ridt`) are not commented.


### 7. Restart Infrastructure

Stop the infrastructure with the following command:

```bash
docker compose down
```

And start it again:

```bash
docker compose up -d
```


## Testing

In the `realm-export.json` are already client applications included to test the API with

- [API Specification](#testing-with-api-specification) or
- [Swagger Editor](#testing-with-swagger-editor) or
- [Postman](#testing-with-postman).


### Request Token JWT Generation

To create a sufficient Token Request JWT, you can go to [JWT.io (external URL)](https://jwt.io) and create one.
You can use the ES256 private key from the [key examples (external repository)](https://github.com/JonasPrimbs/draft-ietf-mla-oidc/tree/main/examples/key-examples.md#private-key) to sign your key or use the [code example (external repository)](https://github.com/JonasPrimbs/draft-ietf-mla-oidc/tree/main/examples/code-examples.md#key-pair-generation-and-export) to generate a new one.
You can also use the Token Request JWT from the [communication examples (external repository)](https://github.com/JonasPrimbs/draft-ietf-mla-oidc/tree/main/examples/communication-examples.md#token-request) as template.

The following JWT can be used as a template:

```jwt
eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCIsImp3ayI6eyJrdHkiOiJFQyIsImNydiI6IlAtMjU2IiwieCI6ImNYUThiZGVOZWVTd2ZMa0h6TWZBVUZySGxMWFpXdkpybW9NMnNDUEdVbmciLCJ5IjoiN0Rwd21Pb0hJbmQwUWNSRVJUS1pBQ2k5YndzYTVnR0tER3hGeG00OEdSQSJ9fQ.eyJpc3MiOiJwb3N0bWFuIiwic3ViIjoiOWJiYWEyZjctNjlhOS00ZWFlLWI2YjgtOTRmYzY2MDExMmZjIiwiYXVkIjoiaHR0cDovL29wLmxvY2FsaG9zdC9yZWFsbXMvcmlkdCIsImlhdCI6MTY1OTM1NTIwNSwibmJmIjoxNjU5MzU1MjA1LCJleHAiOjE2NjkzNTUyMzUsIm5vbmNlIjoiVmpmVTQ2WjV5a0lobjdqSnpxWm9XSytwYXE2M0VLdUgiLCJ0b2tlbl9jbGFpbXMiOiJuYW1lIGVtYWlsIGVtYWlsX3ZlcmlmaWVkIiwidG9rZW5fbGlmZXRpbWUiOjM2MDAsInRva2VuX25vbmNlIjoiQmp4cTI3RlVsQjBYQVcyaWIrWnM2czU3UlFyY21VeEEifQ.BrfJYyrU1bZVWRawXO3Jowic3H84RaIzZDp_e8obviBlLLaq09tAnSUuVGLJ2hw4EIw1enALLtk_F5ZwEMqLlQ
```

It has the following payload (without comments):

```json
{
  "alg": "ES256",
  "typ": "JWT",
  "jwk": {  // The client's public key:
    "kty": "EC",
    "crv": "P-256",
    "x": "cXQ8bdeNeeSwfLkHzMfAUFrHlLXZWvJrmoM2sCPGUng",
    "y": "7DpwmOoHInd0QcRERTKZACi9bwsa5gGKDGxFxm48GRA"
  }
}
```

and the following payload (without comments):

```json
{
  "iss": "postman", // The client ID.
  "sub": "9bbaa2f7-69a9-4eae-b6b8-94fc660112fc",  // The user's unique identifier. In Keycloak, this is a UUID which is displayed in the Users menu.
  "aud": "http://op.localhost/realms/ridt", // The OpenID Provider's URL = issuer of the Remote ID Token.
  "iat": 1659355205,  // Unix timestamp when the token was issued.
  "nbf": 1659355205,  // Unix timestamp when the token becomes valid.
  "exp": 1669355235,  // Unix timestamp when the token expires.
  "nonce": "VjfU46Z5ykIhn7jJzqZoWK+paq63EKuH",  // A random nonce.
  "token_claims": "name email email_verified",  // The requested identity claims for the Remote ID Token.
  "token_lifetime": 3600, // The requested lifetime of the Remote ID Token.
  "token_nonce": "Bjxq27FUlB0XAW2ib+Zs6s57RQrcmUxA" // A random nonce to set into the Remote ID Token.
}
```


### Testing with API Specification

You can test the infrastructure with our API documentation.
This is recommended if you want to play with the API.

Therefore, you must authorize the API documentation as follows:

<details>
  <summary><b>For Public Authorization Server</b></summary>

  1. Open your browser and navigate to the [API documentation (external URL)](https://api.oidc-e2e.primbs.dev/).
  2. Click *Authorize*.
  3. Scroll down to the authorization **oauth2_public**.
  4. Enter the *client_id* `api` and *Select all* scopes.
  5. Click *Authorize* and *Sign In* with your test user.
  6. Click *Close*.
  
</details>
<details>
  <summary><b>For Local Authorization Server</b></summary>

  1. Open your browser and navigate to the [API documentation (external URL)](https://api.oidc-e2e.primbs.dev/).
  2. Click *Authorize*.
  3. Scroll down to the authorization **oauth2_local**.
  4. Enter the *client_id* `api` and *Select all* scopes.
  5. Click *Authorize* and *Sign In* with your test user.
  6. Click *Close*.

</details>

Now you can perform requests to the server as follows:

1. Open the *POST /* Endpoint.
2. Click *Try it out*.
3. Paste a sufficient Token Request JWT to the *Request Body*.
4. Click *Execute* to send the request.


### Testing with Swagger Editor

You can test the infrastructure with Swagger Editor.
This is recommended while editing the API Specification.

Therefore, you must authorize Swagger Editor as follows:

<details>
  <summary><b>For Public Authorization Server</b></summary>

  1. Open your browser and navigate to the [Swagger Editor (external URL)](https://editor.swagger.io/).
  2. Click *Authorize*.
  3. Scroll down to the authorization **oauth2_public**.
  4. Enter the *client_id* `swagger` and *Select all* scopes.
  5. Click *Authorize* and *Sign In* with your test user.
  6. Click *Close*.

</details>
<details>
  <summary><b>For Local Authorization Server</b></summary>

  1. Open your browser and navigate to the [Swagger Editor (external URL)](https://editor.swagger.io/).
  2. Click *Authorize*.
  3. Scroll down to the authorization **oauth2_local**.
  4. Enter the *client_id* `swagger` and *Select all* scopes.
  5. Click *Authorize* and *Sign In* with your test user.
  6. Click *Close*.

</details>

Now you can perform requests to the server as follows:

<details>
  <summary><b>For Public Authorization Server</b></summary>

  1. Make sure that the server starting with URL `https://op.oidc-e2e.primbs.dev/...` is selected.
  2. Open the *POST /* Endpoint.
  3. Click *Try it out*.
  4. Paste a sufficient Token Request JWT to the *Request Body*.
  5. Click *Execute* to send the request.

</details>
<details>
  <summary><b>For Local Authorization Server</b></summary>

  1. Make sure that the server starting with URL `http://op.localhost/...` is selected.
  2. Open the *POST /* Endpoint.
  3. Click *Try it out*.
  4. Paste a sufficient Token Request JWT to the *Request Body*.
  5. Click *Execute* to send the request.

</details>


### Testing with Postman

You can test the infrastructure with Postman.

Therefore, you must authorize Postman as follows:

<details>
  <summary><b>For Public Authorization Server</b></summary>

   1. Open a new Tab and go to the *Authorization* tab.
   2. As *Type*, select `OAuth 2.0`.
   3. In *Configure New Token* > *Configuration Options* insert the following values:
       - *Grant Type*: `Authorization Code (With PKCE)`
       - *Callback URL*: `https://oauth.pstmn.io/v1/callback` and tick *Authorize using browser*.
       - *Auth URL*: `https://op.oidc-e2e.primbs.dev/realms/ridt/protocol/openid-connect/auth`
       - *Access Token URL*: `https://op.oidc-e2e.primbs.dev/realms/ridt/protocol/openid-connect/token`
       - *Client ID*: `postman`
   4. Click *Get New Access Token*.
   5. *Sign In* to your test user account, if requested.
   6. Click *Use Token*.

</details>

<details>
  <summary><b>For Local Authorization Server</b></summary>

   1. Open a new Tab and go to the *Authorization* tab.
   2. As *Type*, select `OAuth 2.0`.
   3. In *Configure New Token* > *Configuration Options* insert the following values:
       - *Grant Type*: `Authorization Code (With PKCE)`
       - *Callback URL*: `https://oauth.pstmn.io/v1/callback` and tick *Authorize using browser*.
       - *Auth URL*: `http://op.localhost/realms/ridt/protocol/openid-connect/auth`
       - *Access Token URL*: `http://op.localhost/realms/ridt/protocol/openid-connect/token`
       - *Client ID*: `postman`
   4. Click *Get New Access Token*.
   5. *Sign In* to your test user account, if requested.
   6. Click *Use Token*.

</details>

Now you can perform requests to the server as follows:

<details>
  <summary><b>For Public Authorization Server</b></summary>

   1. Select the HTTP Method *POST*.
   2. Insert the URL `https://op.oidc-e2e.primbs.dev/realms/ridt/protocol/openid-connect/userinfo/ridt`.
   3. Go to the *Body* tab and insert the Token Request JWT as *raw*.
   4. Click *Send*.

</details>

<details>
  <summary><b>For Local Authorization Server</b></summary>

   1. Select the HTTP Method *POST*.
   2. Insert the URL `http://op.localhost/realms/ridt/protocol/openid-connect/userinfo/ridt`.
   3. Go to the *Body* tab and insert the Token Request JWT as *raw*.
   4. Click *Send*.

</details>
