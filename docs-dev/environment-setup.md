# Environment Setup

This section describes how to setup a test environment locally with Docker Compose.

**WARNING: THIS IS FOR TEST PURPOSES ONLY! DO NOT USE THIS IN PRODUCTION!!!**


## 1. Clone Repository

In your Linux bash, clone this repository to your home directory:

```bash
git clone https://github.com/JonasPrimbs/oidc-e2ea-server.git
```

Now navigate to the cloned directory:

```bash
cd oidc-e2ea-server
```


## 2. Generate Secrets

Execute the following command:

```bash
bash ./generate-secrets.sh
```

This will randomly generate all usernames, passwords, and private keys which are unique for your installation and store them in the new directory `.secrets/` and a `.env` file in the repository.


## 3. Configure Deployment

Go to the generated `/.env` file and configure the following parameters:

- `OP_HOST=<your-hostname>` the host/domain name of your server. Default is `op.localhost`.

For a local deployment, you can leave these settings at default.


## 4. Initial Infrastructure Start

Start up your OpenID Provider for the first time using the following command:

```bash
docker compose up -d op
```

This might take a while to download all related container images.


## 5. Setup OpenID Provider

This section describes how to setup the Keycloak OpenID Provider to make it ready to issue ID Assertion Tokens.


### 5.1. Login to Keycloak Admin Console

Open your browser and go to `http://<your-hostname>/admin/` where `<your-hostname>` is your configured hostname.
By default, this is [http://op.localhost/admin](http://op.localhost/admin).
Then *sign in* with the credentials generated in the following files:
- Username: `/.secrets/op_username.txt`
- Password: `/.secrets/op_password.txt`

*If you experience a* **Bad Gateway** *error, wait for up to one minute until you Keycloak instance is ready!*


### 5.2. Switch to Realm

On the top left, click the dropdown menu and select the realm `ict`:

![Screenshot of how to switch to realm 'ict'](./images/switch_realm.png)


### 5.3. Import Private Key

Import the generated private key as follows:

 1. Go to *Configure* > *Realm settings* > *Keys* > *Providers*.
 2. In *Add provider*, select the option *rsa*.

    ![Screenshot of how to add an RSA key provider](./images/add_rsa.png)

 3. In field *Private RSA Key*, select *Browse...* and select the generate `private.pem` private key file in the `/.secrets/` directory of the cloned repository.
 4. Click *Save* to store the changes.

    ![Screenshot of how to save the RSA key provider](./images/save_rsa.png)


### 5.4. Configure Private Key

 1. Go to the file `/.secrets/ict.env`.
 2. Copy the *Kid* of your newly generated key of *Type* `rsa` from *Configure* > *Realm settings* > *Keys* > *Key list*.

    ![Screenshot of the Kid of the RSA key](./images/rsa_kid.png)

 3. Paste the copied *Kid* parameter to the `/.secrets/ict.env` file as value for the key `KID`, e.g.:

```bash
KID=GFSKUd9yi3LiQhT6HKuU4IOymufp_OIIlG8DmGa8hvs
```


### 5.5. Create Test User

Create a new test user as follows:

 1. Go to *Manage* > *Users* > *User list* > *Create new user*.

    ![Screenshot of how to create a new user](./images/create_user.png)

 2. Insert at least a *Username*.
 3. *Create* the user.

    ![Screenshot of how to save the new user](./images/save_user.png)

 4. In the tab *Credentials*, click *Set password*.

    ![Screenshot of how to set a password for the new user](./images/set_password.png)

 5. Insert a *Password*, repeat it in *Password confirmation*, and set *Temporary* to `off`.
    Then click *Save*.

    ![Screenshot of how to save the password for the new user](./images/save_password.png)

 6. Confirm the dialog by clicking *Save password*.


### 5.6. Configure ICT Endpoint

To introspect the Access Token from the Authorization Server, the ICT Endpoint must be registered at the Authorization Server as follows:

 1. Go to *Clients* > *Client list* > *ict_endpoint* > Credentials.
 2. *Regenerate* the Client secret and copy it to clipboard.
 3. Open a [HTTP Basic Authentication Header Generator](https://www.blitter.se/utils/basic-authentication-header-generator/).
 4. As Username, insert `ict_endpoint`.
 5. As Password, paste the Client secret.
 6. Generate the Basic Auth header and copy the header value (e.g., `Basic aWN0X2VuZHBvaW50OjhjOHY2aGRhZ3c5ZXRTOFVMYVdVZ1dhT2ZUNWpKTzNa`).
 7. Paste this value to the [`/.env`](../.env) file in the `ICT_CREDENTIALS` variable
 
 Example:
 ```bash
 ICT_CREDENTIALS="Basic aWN0X2VuZHBvaW50OjhjOHY2aGRhZ3c5ZXRTOFVMYVdVZ1dhT2ZUNWpKTzNa"`
```


## 6. Configure Deployment Mode

This step depends on your intention why you run this deployment.

- **Testing**: Choose this mode if you want to just run the deployment for testing purposes.
- **Development**: Choose this mode if you want to change the implementation of the ICT endpoint application.


### 6.1. Test Deployment

*Do this step only if you want to run this deployment for **testing** purposes!*

1. Go to `/docker-compose.yaml`.
2. Uncomment line 65 (`image` attribute in service `ict`).
3. Comment line 68 to 70 (`build` attribute in service `ict`).


### 6.2. Development Deployment

*Do this step only if you want to run this deployment for **development** purposes!*

1. Go to `/docker-compose.yaml`.
2. Comment line 65 (`image` attribute in service `ict`).
3. Uncomment line 68 to 70 (`build` attribute in service `ict`).


## 7. Restart Infrastructure

Stop the infrastructure with the following command:

```bash
docker compose down
```

And start it again:

```bash
docker compose up -d
```

## Help

Default username in Authentik is `akadmin`.