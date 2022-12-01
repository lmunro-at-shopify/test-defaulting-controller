#!/bin/bash
set -eu pipefail

echo "Generating local Certs"
key_dir="certs/"
serviceName=webhook-service
namespace=system
mkdir -p $key_dir

# Generate the CA cert and private key
openssl req -nodes -new -days 358000 -x509 -keyout $key_dir/ca.key -out $key_dir/ca.crt -subj "/CN=Admission Controller Webhook CA"

# Generate the private key for the webhook server
openssl genrsa -out $key_dir/tls.key 2048

# Generate a Certificate Signing Request (CSR) for the private key, and sign it with the private key of the CA.
openssl req -new -days 358000 -key $key_dir/tls.key -subj "/CN=$serviceName.$namespace.svc" \
    | openssl x509 -req -days 358000 -CA $key_dir/ca.crt -CAkey $key_dir/ca.key -CAcreateserial -out $key_dir/tls.crt \
    -extensions SAN \
    -extfile <( printf "[SAN]\nsubjectAltName=DNS:$serviceName.$namespace.svc" )

echo "Local Certs Generated"
