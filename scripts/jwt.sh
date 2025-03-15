#!/bin/sh

# Check if the environment variables JWT_PRIVATE and JWT_PUBLIC are set
if [ -z "$JWT_PRIVATE" ] || [ -z "$JWT_PUBLIC" ]; then
  echo "Error: JWT_PRIVATE or JWT_PUBLIC environment variables are not set."
  exit 1
fi

# Decode the private key and save it to the privateKey.pem file
echo "$JWT_PRIVATE" | base64 -d > internal/jwt/privateKey.pem

if [ $? -ne 0 ]; then
  echo "Error: Failed to decode JWT_PRIVATE."
  exit 1
fi
echo "Decoded JWT_PRIVATE and saved to internal/jwt/privateKey.pem."

# Decode the public key and save it to the publicKey.pem file
echo "$JWT_PUBLIC" | base64 -d > internal/jwt/publicKey.pem

if [ $? -ne 0 ]; then
  echo "Error: Failed to decode JWT_PUBLIC."
  exit 1
fi
echo "Decoded JWT_PUBLIC and saved to internal/jwt/publicKey.pem."

# Success
echo "Keys successfully decoded and saved."