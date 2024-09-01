#!/bin/bash

# Generate JWT secret key
jwt_secret=$(openssl rand -hex 32)
echo "JWT_SECRET=$jwt_secret" >> .env

# Generate CSRF auth key
csrf_auth_key=$(openssl rand -hex 32)
echo "CSRF_AUTH_KEY=$csrf_auth_key" >> .env

# Generate self-signed certificate for development
openssl req -x509 -newkey rsa:4096 -keyout key.pem -out cert.pem -days 365 -nodes -subj "/CN=localhost"

echo "Security setup complete. Self-signed certificates generated for development."