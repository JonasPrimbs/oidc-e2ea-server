SECRETS_DIR="./.secrets/"
PRIVATE_KEY_FILE="${SECRETS_DIR}private.pem"
DB_USERNAME_FILE="${SECRETS_DIR}db_username.txt"
DB_PASSWORD_FILE="${SECRETS_DIR}db_password.txt"
OP_USERNAME_FILE="${SECRETS_DIR}op_username.txt"
OP_PASSWORD_FILE="${SECRETS_DIR}op_password.txt"
RIDT_ENV_FILE="${SECRETS_DIR}ridt.env"
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

# Create DB username secret file if not exists
[ ! -f "$DB_USERNAME_FILE" ] && generate_secret "$DB_USERNAME_FILE" 14

# Create DB password secret file if not exists
[ ! -f "$DB_PASSWORD_FILE" ] && generate_secret "$DB_PASSWORD_FILE" 64

# Create OP username secret file if not exists
[ ! -f "$OP_USERNAME_FILE" ] && generate_secret "$OP_USERNAME_FILE" 14

# Create OP password secret file if not exists
[ ! -f "$OP_PASSWORD_FILE" ] && generate_secret "$OP_PASSWORD_FILE" 64

# Write DB username + password to op.env file
[ ! -f "$OP_ENV_FILE" ] && echo "KC_DB_USERNAME=$(cat $DB_USERNAME_FILE)" > "$OP_ENV_FILE"
                        && echo "KC_DB_PASSWORD=$(cat $DB_PASSWORD_FILE)" >> "$OP_ENV_FILE"
                        && echo "KEYCLOAK_ADMIN=$(cat $OP_USERNAME_FILE)" >> "$OP_ENV_FILE"
                        && echo "KEYCLOAK_ADMIN_PASSWORD=$(cat $OP_PASSWORD_FILE)" >> "$OP_ENV_FILE"

# Write empty key ID to ridt.env file
[ ! -f "$RIDT_ENV_FILE" ] && echo "KID=" > "$RIDT_ENV_FILE"

# Write default hostname and default realm name to .env file
[ ! -f "$RIDT_ENV_FILE" ] && echo "OP_HOST=op.localhost" > $ENV_FILE
                          && echo "REALM_NAME=ridt" >> $ENV_FILE
