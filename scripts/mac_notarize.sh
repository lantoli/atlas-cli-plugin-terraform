#!/usr/bin/env bash

set -Eeou pipefail

if [[ -f "./dist/macos_darwin_amd64_v1/binary" && -f "./dist/macos_darwin_arm64_v8.0/binary" && ! -f "./dist/mac_bin_signed.zip" ]]; then
	echo "notarizing macOs binaries"

	curl "${NOTARY_SERVICE_URL}" --output macos-notary.zip
	unzip -u macos-notary.zip
	chmod 755 ./linux_amd64/macnotary

	zip -r ./dist/mac_bin.zip ./dist/macos_darwin_amd64_v1/binary ./dist/macos_darwin_arm64_v8.0/binary # The Notarization Service takes an archive as input
	./linux_amd64/macnotary \
		-f ./dist/mac_bin.zip \
		-m notarizeAndSign -u https://dev.macos-notary.build.10gen.cc/api \
		-b com.mongodb.atlas-cli-plugin-terraform \
		-o ./dist/mac_bin_signed.zip

	echo "replacing original files"
	unzip -oj ./dist/mac_bin_signed.zip dist/macos_darwin_amd64_v1/binary -d ./dist/macos_darwin_amd64_v1/
	unzip -oj ./dist/mac_bin_signed.zip dist/macos_darwin_arm64_v8.0/binary -d ./dist/macos_darwin_arm64_v8.0/
fi
