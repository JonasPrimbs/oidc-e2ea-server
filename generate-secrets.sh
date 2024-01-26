#!/bin/bash

SECRETS_DIR="./.secrets/"
PRIVATE_KEY_FILE="${SECRETS_DIR}private.pem"
DB_USERNAME_FILE="${SECRETS_DIR}db_username.txt"
DB_PASSWORD_FILE="${SECRETS_DIR}db_password.txt"
OP_USERNAME_FILE="${SECRETS_DIR}op_username.txt"
OP_PASSWORD_FILE="${SECRETS_DIR}op_password.txt"
PRIVATE_KEY2_FILE="${SECRETS_DIR}private2.pem"
CERTIFICATE2_FILE="${SECRETS_DIR}certificate2.crt"
DB2_USERNAME_FILE="${SECRETS_DIR}db2_username.txt"
DB2_PASSWORD_FILE="${SECRETS_DIR}db2_password.txt"
OP2_SECRET_KEY="${SECRETS_DIR}op2_secret_key.txt"
ICT_ENV_FILE="${SECRETS_DIR}ict.env"
ICT_ENV2_FILE="${SECRETS_DIR}ict2.env"
OP_ENV_FILE="${SECRETS_DIR}op.env"
ENV_FILE=".env"

# $1 = file name
# $2 = secret length
function generate_secret {
    head /dev/urandom | tr -dc A-Za-z0-9 | head -c $2 > $1
}

# Create secrets directory if not exists
[ ! -d "$SECRETS_DIR" ] && mkdir "$SECRETS_DIR"

# Create private key file if not exists
[ ! -f "$PRIVATE_KEY_FILE" ] && openssl genrsa -out "$PRIVATE_KEY_FILE" 2048

# Create private key 2 file if not exists
[ ! -f "$PRIVATE_KEY2_FILE" ] && openssl genrsa -out "$PRIVATE_KEY2_FILE" 2048

[ ! -f "$CERTIFICATE2_FILE" ] && openssl req -new -x509 -nodes -key "$PRIVATE_KEY2_FILE" -out "$CERTIFICATE2_FILE" -days 365

# Create DB username secret file if not exists
[ ! -f "$DB_USERNAME_FILE" ] && generate_secret "$DB_USERNAME_FILE" 14

# Create DB password secret file if not exists
[ ! -f "$DB_PASSWORD_FILE" ] && generate_secret "$DB_PASSWORD_FILE" 64

# Create OP username secret file if not exists
[ ! -f "$OP_USERNAME_FILE" ] && generate_secret "$OP_USERNAME_FILE" 14

# Create OP password secret file if not exists
[ ! -f "$OP_PASSWORD_FILE" ] && generate_secret "$OP_PASSWORD_FILE" 64

# Create DB username secret file if not exists
[ ! -f "$DB2_USERNAME_FILE" ] && generate_secret "$DB2_USERNAME_FILE" 14

# Create DB password secret file if not exists
[ ! -f "$DB2_PASSWORD_FILE" ] && generate_secret "$DB2_PASSWORD_FILE" 64

# Create secret key for Authentik file if not exists
[ ! -f "$OP2_SECRET_KEY" ] && generate_secret "$OP2_SECRET_KEY" 64

# Write DB username + password to op.env file
if [ ! -f "$OP_ENV_FILE" ]
then
    echo "KC_DB_USERNAME=$(cat $DB_USERNAME_FILE)" > "$OP_ENV_FILE"
    echo "KC_DB_PASSWORD=$(cat $DB_PASSWORD_FILE)" >> "$OP_ENV_FILE"
    echo "KEYCLOAK_ADMIN=$(cat $OP_USERNAME_FILE)" >> "$OP_ENV_FILE"
    echo "KEYCLOAK_ADMIN_PASSWORD=$(cat $OP_PASSWORD_FILE)" >> "$OP_ENV_FILE"
fi

# Write empty key ID to ict.env file
[ ! -f "$ICT_ENV_FILE" ] && echo "KID=" > "$ICT_ENV_FILE"

# Write empty key ID to ict2.env file
[ ! -f "$ICT_ENV2_FILE" ] && echo "KID=" > "$ICT_ENV2_FILE"

# Write default hostname and default realm name to .env file
if [ ! -f "$ENV_FILE" ]
then
    echo "OP_HOST=op.localhost" > $ENV_FILE
    echo "OP2_HOST=op2.localhost" >> $ENV_FILE
fi
