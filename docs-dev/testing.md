# Testing

In the `realm-export.json` are already client applications included to test the API with

- [API Specification](#testing-with-api-specification) or
- [Swagger Editor](#testing-with-swagger-editor) or
- [Postman](#testing-with-postman).


## Request Token JWT Generation

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
  "aud": "http://op.localhost/realms/iat", // The OpenID Provider's URL = issuer of the ID Assertion Token.
  "iat": 1659355205,  // Unix timestamp when the token was issued.
  "nbf": 1659355205,  // Unix timestamp when the token becomes valid.
  "exp": 1669355235,  // Unix timestamp when the token expires.
  "nonce": "VjfU46Z5ykIhn7jJzqZoWK+paq63EKuH",  // A random nonce.
  "token_claims": "name email email_verified",  // The requested identity claims for the ID Assertion Token.
  "token_lifetime": 3600, // The requested lifetime of the ID Assertion Token.
  "token_nonce": "Bjxq27FUlB0XAW2ib+Zs6s57RQrcmUxA" // A random nonce to set into the ID Assertion Token.
}
```


## Testing with API Specification

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


## Testing with Swagger Editor

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


## Testing with Postman

You can test the infrastructure with Postman.

Therefore, you must authorize Postman as follows:

<details>
  <summary><b>For Public Authorization Server</b></summary>

   1. Open a new Tab and go to the *Authorization* tab.
   2. As *Type*, select `OAuth 2.0`.
   3. In *Configure New Token* > *Configuration Options* insert the following values:
       - *Grant Type*: `Authorization Code (With PKCE)`
       - *Callback URL*: `https://oauth.pstmn.io/v1/callback` and tick *Authorize using browser*.
       - *Auth URL*: `https://op.oidc-e2e.primbs.dev/realms/iat/protocol/openid-connect/auth`
       - *Access Token URL*: `https://op.oidc-e2e.primbs.dev/realms/iat/protocol/openid-connect/token`
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
       - *Auth URL*: `http://op.localhost/realms/iat/protocol/openid-connect/auth`
       - *Access Token URL*: `http://op.localhost/realms/iat/protocol/openid-connect/token`
       - *Client ID*: `postman`
   4. Click *Get New Access Token*.
   5. *Sign In* to your test user account, if requested.
   6. Click *Use Token*.

</details>

Now you can perform requests to the server as follows:

<details>
  <summary><b>For Public Authorization Server</b></summary>

   1. Select the HTTP Method *POST*.
   2. Insert the URL `https://op.oidc-e2e.primbs.dev/realms/iat/protocol/openid-connect/userinfo/iat`.
   3. Go to the *Body* tab and insert the Token Request JWT as *raw*.
   4. Click *Send*.

</details>

<details>
  <summary><b>For Local Authorization Server</b></summary>

   1. Select the HTTP Method *POST*.
   2. Insert the URL `http://op.localhost/realms/iat/protocol/openid-connect/userinfo/iat`.
   3. Go to the *Body* tab and insert the Token Request JWT as *raw*.
   4. Click *Send*.

</details>
