#!/bin/bash

TEMPLATE=$(cat <<EOF
# Instructions to load new secrets
## 1. Add secret names in the workflows_call -> secrets and mark it as required
## 2. Export the secret name with value as JSON in the run section of the "Loading secrets" job

name: Secrets Loader
permissions:
  contents: read

on:
  workflow_call:
    outputs:
      encrypted_secrets:
        description: "Encrypt loaded secrets in base64 JSON format"
        value: \${{ jobs.loading.outputs.encrypted_secrets }}
    secrets:
      APOLLO_KEY:
        required: true
      ## Add addition secrets here

env:
  GHA_GPG_PASSPHRASE: \${{ secrets.GHA_GPG_PASSPHRASE }}

jobs:
  loading:
    name: loading
    runs-on: ubuntu-latest
    outputs:
      encrypted_secrets: \${{ steps.loading.outputs.encrypted_secrets }}
    steps:
      - name: Loading secrets
        id: loading
        run: |
          PLAINTEXT_JSON=\$(cat <<EOM
          {
            "APOLLO_KEY": "\${{ secrets.APOLLO_KEY }}"
            ## Add addition secrets here (With comma separated JSON format)
          }
          EOM
          )
          ENCRYPTED_SECRET=\$(echo "\$PLAINTEXT_JSON" | gpg --symmetric --cipher-algo AES256 --batch --yes --passphrase "\$GHA_GPG_PASSPHRASE" | base64 | tr -d '\n')
          echo "encrypted_secrets=\$ENCRYPTED_SECRET" >> \$GITHUB_OUTPUT

EOF
)
TARGET=".github/workflows/secrets-loader.yaml"

if [ ! -f "$TARGET" ]; then
    echo "$TEMPLATE" > $TARGET
fi
