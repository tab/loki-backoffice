# Generating Certificates for mTLS

For mTLS (mutual TLS), both the server and client need certificates.
The process involves:

- Creating a Certificate Authority (CA)
- Creating server certificates signed by the CA
- Creating client certificates signed by the CA

### Generate the Certificate Authority (CA)

Generate a private key for your CA

```sh
openssl genrsa -out ca.key 4096
openssl req -new -x509 -key ca.key -sha256 -subj "/CN=Loki CA" -out ca.crt -days 3650
```

### Generate the Server Certificate

#### Generate server private key

```sh
openssl genrsa -out server.key 4096
```

#### Create server Certificate Signing Request (CSR)

```sh
openssl req -new -key server.key -out server.csr -config <(
cat <<-EOF
[req]
default_bits = 4096
prompt = no
default_md = sha256
req_extensions = req_ext
distinguished_name = dn

[dn]
CN = loki-server

[req_ext]
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
IP.1 = 127.0.0.1
EOF
)
```

#### Sign the server certificate with CA

```sh
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.pem -days 825 -sha256 -extfile <(
cat <<-EOF
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
IP.1 = 127.0.0.1
EOF
)
```

### Generate the Client Certificate

Generate client private key

```sh
openssl genrsa -out client.key 4096
```

#### Create client Certificate Signing Request (CSR)

```sh
openssl req -new -key client.key -out client.csr -config <(
cat <<-EOF
[req]
default_bits = 4096
prompt = no
default_md = sha256
distinguished_name = dn

[dn]
CN = loki-backoffice
EOF
)
```

#### Sign the client certificate with CA

```sh
openssl x509 -req -in client.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out client.pem -days 825 -sha256
```

### Verify the certificates

```sh
openssl verify -CAfile ca.crt server.pem
```

```sh
openssl verify -CAfile ca.crt client.pem
```
