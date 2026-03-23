#! /bin/sh

SECRET_DIR=".secrets"

# Root CA parameters
CA_ROOT_KEY_PASSWORD_FILE="$SECRET_DIR/ca_root_password.txt"
CA_ROOT_KEY_FILE="$SECRET_DIR/ca_root.key"
CA_ROOT_CERT_FILE="$SECRET_DIR/ca_root.crt"
CA_ROOT_SUBJECT="/C=DE/ST=Baden-Wuerttemberg/L=Tuebingen/O=Uni-Tuebingen/OU=Chair of Communication Networks/CN=OIDC² Root CA"

# Intermediate CA parameters
CA_INTERMEDIATE_KEY_PASSWORD_FILE="$SECRET_DIR/ca_intermediate_password.txt"
CA_INTERMEDIATE_KEY_FILE="$SECRET_DIR/ca_intermediate.key"
CA_INTERMEDIATE_CERT_FILE="$SECRET_DIR/ca_intermediate.crt"
CA_INTERMEDIATE_SUBJECT="/C=DE/ST=Baden-Wuerttemberg/L=Tuebingen/O=Uni-Tuebingen/OU=Chair of Communication Networks/CN=OIDC² Intermediate CA"
CA_DB_PASSWORD_FILE="$SECRET_DIR/ca_db_password.txt"
CA_DB_PASSFILE="$SECRET_DIR/ca_db_passfile.txt"

# OpenID Provider parameters
OP_DB_PASSWORD_FILE="$SECRET_DIR/op_db_password.txt"
OP_PRIVATE_KEY_FILE="$SECRET_DIR/op_private.key"
OP_ENV_FILE="$SECRET_DIR/op.env"

generate_secret() {
  local length=$1
  openssl rand -base64 $((length * 3 / 4)) | tr -d '=+/' | cut -c1-$length
}
generate_rsa_keypair() {
  local key_size=${1:-2048}
  local key_password_file=${2:-"$SECRET_DIR/private_password.txt"}
  local key_file=${3:-"$SECRET_DIR/private.key"}
  local public_key_file=$4

  # Generate RSA private key
  if [ -z "$key_password_file" ] || [ ! -f "$key_password_file" ]; then
    # No password file provided or password file does not exist -> generate unencrypted private key
    openssl genrsa -out "$key_file" $key_size

    # Generate public key if output file is specified
    if [ -n "$public_key_file" ]; then
      openssl rsa -in "$key_file" -pubout -out "$public_key_file"
    fi
  else
    # Password file provided and exists -> generate encrypted private key
    openssl genrsa -aes256 -out "$key_file" -passout file:"$key_password_file" $key_size

    # Generate public key if output file is specified
    if [ -n "$public_key_file" ]; then
      openssl rsa -in "$key_file" -pubout -out "$public_key_file" -passin file:"$key_password_file"
    fi
  fi
}
generate_root_ca() {
  local key_size=${1:-2048}
  local key_password_file=$2
  local key_file=${3:-"$SECRET_DIR/root.key"}
  local cert_lifetime=${4:-365}
  local cert_file=${5:-"$SECRET_DIR/root.crt"}
  local cert_subject=${6:-ROOT_CA_SUBJECT}

  generate_rsa_keypair $key_size "$key_password_file" "$key_file"
  openssl req -x509 -new -nodes -key "$key_file" -sha256 -days "$cert_lifetime" -out "$cert_file" -subj "$cert_subject" -passin file:"$key_password_file"
}
generate_intermediate_ca() {
  local key_size=${1:-2048}
  local key_password_file=$2
  local key_file=${3:-"$SECRET_DIR/intermediate.key"}
  local ca_cert_file=$4
  local ca_key_file=$5
  local ca_key_password_file=$6
  local cert_lifetime=${7:-365}
  local cert_file=${8:-"$SECRET_DIR/intermediate.crt"}
  local cert_subject=${9:-INTERMEDIATE_CA_SUBJECT}
  local ext_file="$SECRET_DIR/intermediate_ext.cnf"

  generate_rsa_keypair $key_size "$key_password_file" "$key_file"
  openssl req -new -key "$key_file" -out "$SECRET_DIR/intermediate.csr" -subj "$cert_subject"  -passin file:"$key_password_file"
  printf "basicConstraints=critical,CA:TRUE,pathlen:0\nkeyUsage=critical,digitalSignature,cRLSign,keyCertSign" > "$ext_file"
  openssl x509 -req -in "$SECRET_DIR/intermediate.csr" -CA "$ca_cert_file" -CAkey "$ca_key_file" -CAcreateserial -out "$cert_file" -days $cert_lifetime -sha256 -passin "file:$ca_key_password_file" -extfile "$ext_file"
  rm "$SECRET_DIR/intermediate.csr"
  rm "$ext_file"
}

mkdir -p "$SECRET_DIR"

# Generate Root CA
echo "Generating Root CA..."
generate_secret 48 > "$CA_ROOT_KEY_PASSWORD_FILE"
generate_root_ca 4096 "$CA_ROOT_KEY_PASSWORD_FILE" "$CA_ROOT_KEY_FILE" 3650 "$CA_ROOT_CERT_FILE" "$CA_ROOT_SUBJECT"
# Generate Intermediate CA
echo "Generating Intermediate CA..."
generate_secret 48 > "$CA_INTERMEDIATE_KEY_PASSWORD_FILE"
generate_intermediate_ca 4096 "$CA_INTERMEDIATE_KEY_PASSWORD_FILE" "$CA_INTERMEDIATE_KEY_FILE" "$CA_ROOT_CERT_FILE" "$CA_ROOT_KEY_FILE" "$CA_ROOT_KEY_PASSWORD_FILE" 365 "$CA_INTERMEDIATE_CERT_FILE" "$CA_INTERMEDIATE_SUBJECT"
# generate_secret 32 > "$CA_PASSWORD_FILE"

# Generate Intermediate CA database secrets
echo "Generating Certificate Authority database secrets..."
generate_secret 32 > "$CA_DB_PASSWORD_FILE"
echo "ca_db:5432:stepca:stepca:$(cat $CA_DB_PASSWORD_FILE)" > "$CA_DB_PASSFILE"

# Generate OpenID Provider password secrets
echo "Generating OpenID Provider password secrets..."
generate_secret 32 > "$OP_DB_PASSWORD_FILE"
echo "KC_DB_PASSWORD=$(cat $OP_DB_PASSWORD_FILE)" >> "$OP_ENV_FILE"
# Generate OpenID Provider ICT signing key pair
echo "Generating OpenID Provider ICT signing key pair..."
generate_rsa_keypair 2048 "" "$OP_PRIVATE_KEY_FILE"
# Generate OpenID Provider's initial admin credentials
echo "Generating OpenID Provider's admin password..."
echo "KC_BOOTSTRAP_ADMIN_PASSWORD=$(generate_secret 24)" >> "$OP_ENV_FILE"

echo "All done. You will find the generated secrets in the '$(pwd)/$SECRET_DIR' directory."
