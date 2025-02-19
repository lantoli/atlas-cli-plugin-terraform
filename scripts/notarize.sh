#!/usr/bin/env bash

set -Eeou pipefail

echo "GRS_CONFIG_USER1_USERNAME=${ARTIFACTORY_SIGN_USER}" >> "signing-envfile"
echo "GRS_CONFIG_USER1_PASSWORD=${ARTIFACTORY_SIGN_PASSWORD}" >> "signing-envfile"

if [[ -f "${artifact:?}" ]]; then
  echo "notarizing package ${artifact}"

  docker run \
    --env-file=signing-envfile \
    --rm -v "$(pwd)":"$(pwd)" -w "$(pwd)" \
    "${ARTIFACTORY_REGISTRY}/release-tools-container-registry-local/garasign-gpg" \
    /bin/bash -c "gpgloader && gpg --yes -v --armor -o ${artifact}.sig --detach-sign ${artifact}"
fi

echo "Signing of ${artifact} completed."

