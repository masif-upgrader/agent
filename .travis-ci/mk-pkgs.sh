#!/bin/bash

set -e
set -o pipefail
set -x


function mkDeb {
	ARGS=("$@")

	GO_ARCH="${ARGS[0]}"
	DEB_ARCH="${ARGS[1]}"
	GO_ENV=('GOOS=linux' "GOARCH=$GO_ARCH" "${ARGS[@]:2}")

	GO_OUT="$(echo "${GO_ENV[@]}")"
	GO_OUT="${GO_OUT//=/_}"
	GO_OUT="${BIN_CACHE}/${GO_OUT// /-}.bin"

	echo "$DEB_ARCH"

	if ! [ -e "$GO_OUT" ]; then
		for e in "${GO_ENV[@]}"; do
			export $e
		done

		go build -o "$GO_OUT" .
	fi

	cp "$GO_OUT" "pkgroot/usr/sbin/masif-upgrader-agent"

	rm -f pkgpayload.tar

	pushd pkgroot

	tar -cf ../pkgpayload.tar *

	popd

	fpm -s tar -t deb --log debug --verbose --debug \
		-n "$PKG_NAME" \
		-v "$PKG_VERSION" \
		-a "$DEB_ARCH" \
		-m 'Alexander A. Klimov <grandmaster@al2klimov.de>' \
		--description 'The Masif Upgrader agent is a component of Masif Upgrader.
Consult Masif Upgrader'"'"'s manual on its purpose and the agent'"'"'s role in its architecture:
https://github.com/masif-upgrader/manual' \
		--url 'https://github.com/masif-upgrader/agent' \
		-p "${PKG_NAME}-${PKG_VERSION}-${DEB_ARCH}.deb" \
		-d apt -d bash -d systemd --no-auto-depends \
		--config-files /etc/masif-upgrader/agent.ini \
		--after-install packaging/daemon-reload.sh --after-upgrade packaging/daemon-reload.sh --after-remove packaging/daemon-reload.sh \
		pkgpayload.tar
}


export BIN_CACHE="$(mktemp -d)"
export PKG_VERSION="$(git describe)"
export PKG_VERSION="${PKG_VERSION/v/}"
export PKG_NAME="masif-upgrader-agent"

mkdir -p pkgroot/usr/sbin
mkdir -p pkgroot/etc/masif-upgrader
mkdir -p pkgroot/lib/systemd/system

cp packaging/config.ini pkgroot/etc/masif-upgrader/agent.ini
cp packaging/systemd.service pkgroot/lib/systemd/system/masif-upgrader-agent.service


go generate

mkDeb amd64 amd64 GO386=387
mkDeb 386 i386 GO386=387

mkDeb mips mips GOMIPS=softfloat
mkDeb mipsle mipsel GOMIPS=softfloat
mkDeb mips64le mips64el

mkDeb ppc64le ppc64el
mkDeb s390x s390x

mkDeb arm armel GOARM=5
mkDeb arm armhf GOARM=7
mkDeb arm64 arm64

mkDeb arm armv6l GOARM=6
mkDeb arm armv7l GOARM=7
mkDeb arm64 aarch64
