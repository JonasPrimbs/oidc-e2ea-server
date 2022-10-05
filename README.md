# OpenID Connect GOes E2EA

This repository provides the proof-of-concept implementation of the end-to-end authentication for an OpenID Connect OpenID Provider server.

The provided implementation is written in [Go](https://golang.org/) and implements a REST endpoint for the Remote ID Token UserInfo endpoint.
Using a reverse proxy in front, it can be mounted to any OpenID Provider implementation.

**Warning:**
Keep in mind that this is the implementation of a research project!
We do not guarantee a secure implementation!
**Do not use this in production!!!**


## Documentation

This section provides an introduction to the architecture and the configuration of the Remote ID Token Endpoint.


### Architecture

The following figure shows the overall architecture how to use the provided Remote ID Token (RIDT) Endpoint with any OpenID Provider implementation.

```
                            +---------+                +----------+
                            |         |       *        |  OpenID  |
                            |         |--------------->| Provider |
   ------                   |         |                +----------+
 /        \  localhost:8080 | Reverse |
| Internet |--------------->|  Proxy  |
 \        /                 |         |                                                    +----------+
   ------                   |         | /realms/test/protocol/openid-connect/userinfo/ridt |   RIDT   |
                            |         |--------------------------------------------------->| Endpoint |
                            +---------+                                                    +----------+
```

The Docker Compose composition provided [here](./docker-compose.yaml) uses the following implementations:

- Reverse Proxy: [Traefik Proxy](https://traefik.io/traefik/)
- OpenID Provider: [Keycloak](https://www.keycloak.org/)


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


#### Key ID

The ID of the OpenID Provider's Public Key provided in the `jwks_uri` endpoint.

Example 1:
```bash
KID="abcdef"
```

Example 2:
```bash
KID=1
```


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
USERINFO="http://localhost:8080/realms/test/protocol/openid-connect/userinfo"
```


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
ISSUER="http://localhost:8080/realms/test"
```


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


### REST Endpoint

The REST API is described in the OpenAPI format provided [here](./docs/openapi.yaml).


## Environment Setup

This section describes how to setup the test infrastructure with Docker Compose.

**WARNING: THIS IS FOR TEST PURPOSES ONLY! DO NOT USE THIS IN PRODUCTION!!!**


### 1. Generate Secrets

In your Linux bash, navigate to `/poc` of this cloned repository and run the following command:

```bash
/poc$ bash ./generate-secrets.sh
```

This will randomly generate all usernames, passwords, and private keys which are unique for your installation and store them in the new directory `/poc/.secrets`.


### 2. Initial Infrastructure Start

Then start up your infrastructure for the first time using the following command:

```bash
/poc$ docker-compose up -d
```


### 3. Setup OpenID Provider

Now, open your browser and go to [http://op.localhost/admin/](http://op.localhost/admin/) and *sign in* with the default credentials:
- Username: `admin`
- Password: `admin`

Then, create a new realm called `test` as follows:

1. Hover the *Master* realm in the navigation bar and click *Add Realm*.
2. Enter the realm *name* `test`.
3. Click *Create*.

Then, import the realm settings as follows:

1. In your new test realm, go to *Manage* > *Import*.
2. Click on *Select file* and select the file `/poc/keycloak/realm-export.json`.
3. Click *Import*.

Then, import the private key as follows:

1. In your new test realm, go to *Configure* > *Realm Settings* > *Keys* > *Providers*
2. On top of the table, open the dropdown menu *Add keystore...* and select *rsa*.
3. Enter the *priority* `101`.
4. As *Private RSA Key*, *Select file* `/poc/.secrets/private.pem`.
5. Click *Save*.
6. In *Configure* > *Realm Settings* > *Keys* > *Active*, copy the *Kid* of the new RSA key with priority `101` and paste it to the file `/poc/.secrets/ridt.env` as value for the parameter `KID`, e.g., `KID=rojPQoDRx_DD-DFs7y45wDLl5T8b9VmX6iQapIK6cRE`.

Finally, create a new test user:

1. In your new test realm, go to *Manage* > *Users*
2. On top of the table, click the *Add user* button.
3. Enter a *username*, e.g., `test`.
4. It is recommended that you also enter an *Email* address, a *First Name*, and a *Last Name*.
5. Click *Save*.
6. In your new user, go to the *Credentials* tab.
7. Set a *Password*, re-enter the *Password Confirmation*, and set *Temporary* to `OFF`.
8. Then click *Set Password*.
9. Apply with *Set password*.


### 4. Restart Infrastructure

Now, stop the infrastructure with the following command:

```bash
/poc$ docker-compose down
```

And start it again:

```bash
/poc$ docker-compose up -d
```


## Testing

In the `realm-export.json` are already client applications included to test the API with

- [Swagger Editor](#testing-with-swagger-editor) or
- [Postman](#testing-with-postman).


### Request Token JWT Generation

To create a sufficient Token Request JWT, you can go to [https://jwt.io](https://jwt.io) and create one.
You can use the ES256 private key from the [examples](../examples/key-examples.md#private-key) to sign your key or use the [code example](../examples/code-examples.md#key-pair-generation-and-export) to generate a new one.
You can also use the Token Request JWT from the [communication examples](../examples/communication-examples.md#token-request) as template.

The following JWT can be used as a template:

```jwt
eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCIsImp3ayI6eyJrdHkiOiJFQyIsImNydiI6IlAtMjU2IiwieCI6ImNYUThiZGVOZWVTd2ZMa0h6TWZBVUZySGxMWFpXdkpybW9NMnNDUEdVbmciLCJ5IjoiN0Rwd21Pb0hJbmQwUWNSRVJUS1pBQ2k5YndzYTVnR0tER3hGeG00OEdSQSJ9fQ.eyJpc3MiOiJwb3N0bWFuIiwic3ViIjoiOWJiYWEyZjctNjlhOS00ZWFlLWI2YjgtOTRmYzY2MDExMmZjIiwiYXVkIjoiaHR0cDovL29wLmxvY2FsaG9zdC9yZWFsbXMvdGVzdCIsImlhdCI6MTY1OTM1NTIwNSwibmJmIjoxNjU5MzU1MjA1LCJleHAiOjE2NjkzNTUyMzUsIm5vbmNlIjoiVmpmVTQ2WjV5a0lobjdqSnpxWm9XSytwYXE2M0VLdUgiLCJ0b2tlbl9jbGFpbXMiOiJuYW1lIGVtYWlsIGVtYWlsX3ZlcmlmaWVkIiwidG9rZW5fbGlmZXRpbWUiOjM2MDAsInRva2VuX25vbmNlIjoiQmp4cTI3RlVsQjBYQVcyaWIrWnM2czU3UlFyY21VeEEifQ.VXvKD-ZzrU_ESdFu8sa10GVK-fUvX3IlUGzCYJ27a-S-fdmKD72KmQRtL_91non7fUjWJOZLJrWg4vwKUYqrDA
```

It has the following payload (without comments):

```json
{
  "alg": "ES256",
  "typ": "JWT",
  "jwk": { // The client's public key:
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
  "sub": "9bbaa2f7-69a9-4eae-b6b8-94fc660112fc", // The user's unique identifier. In Keycloak, this is a UUID which is displayed in the Users menu.
  "aud": "http://op.localhost/realms/test", // The OpenID Provider's URL = issuer of the Remote ID Token.
  "iat": 1659355205, // Unix timestamp when the token was issued.
  "nbf": 1659355205, // Unix timestamp when the token becomes valid.
  "exp": 1669355235, // Unix timestamp when the token expires.
  "nonce": "VjfU46Z5ykIhn7jJzqZoWK+paq63EKuH", // A random nonce.
  "token_claims": "name email email_verified", // The requested identity claims for the Remote ID Token.
  "token_lifetime": 3600, // The requested lifetime of the Remote ID Token.
  "token_nonce": "Bjxq27FUlB0XAW2ib+Zs6s57RQrcmUxA" // A random nonce to set into the Remote ID Token.
}
```


### Testing with Swagger Editor

You can test the infrastructure with Swagger Editor.

Therefore, you must authorize Swagger Editor as follows:

1. Open your browser and go to [https://editor.swagger.io](https://editor.swagger.io).
2. Copy the content of the file `/poc/ridt-endpoint/api/swagger.yaml` to the left side of the Swagger Editor.
3. On the right side, click *Authorize*.
4. Enter the *client_id* `swagger` and *Select all* scopes.
6. Click *Authorize* and *Sign In* with your test user.
7. Click *Close*.

Now you can perform requests to the server as follows:

1. Open the *POST /* Endpoint.
2. Click *Try it out*.
3. Paste a sufficient Token Request JWT to the *Request Body*.
4. Click *Execute* to send the request.


### Testing with Postman

You can test the infrastructure with Postman.

Therefore, you must authorize Postman as follows:

1. Open a new Tab and go to the *Authorization* tab.
2. As *Type*, select `OAuth 2.0`.
3. In *Configure New Token* > *Configuration Options* insert the following values:
    - *Grant Type*: `Authorization Code (With PKCE)`
    - *Callback URL*: `https://oauth.pstmn.io/v1/callback` and tick *Authorize using browser*.
    - *Auth URL*: `http://op.localhost/realms/test/protocol/openid-connect/auth`
    - *Access Token URL*: `http://op.localhost/realms/test/protocol/openid-connect/token`
    - *Client ID*: `postman`
4. Click *Get New Access Token*
5. *Sign In* to your test user account, if requested.
6. Click *Use Token*.

Now you can perform requests to the server as follows:

1. Select the HTTP Method *POST*.
2. Insert the URL `http://op.localhost/realms/test/protocol/openid-connect/userinfo/ridt`.
3. Go to the *Body* tab and insert the Token Request JWT as *raw*.
4. Click *Send*.
