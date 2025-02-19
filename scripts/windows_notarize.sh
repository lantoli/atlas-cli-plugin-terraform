#!/usr/bin/env bash

set -Eeou pipefail

echo "TEMP signing ${artifact}"

if [[ -f "${artifact:?}" ]]; then
	echo "signing ${artifact}"

  docker run \
    -e GRS_CONFIG_USER1_USERNAME="${ARTIFACTORY_SIGN_USER}" \
    -e GRS_CONFIG_USER1_PASSWORD="${ARTIFACTORY_SIGN_PASSWORD}" \
		--rm -v "$(pwd)":"$(pwd)" -w "$(pwd)" \
    "${ARTIFACTORY_REGISTRY}/release-tools-container-registry-local/garasign-jsign" \
		/bin/bash -c "jsign --tsaurl http://timestamp.digicert.com -a ${AUTHENTICODE_KEY_NAME} \"${artifact}\""
fi
