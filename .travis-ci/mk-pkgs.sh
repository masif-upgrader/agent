#!/bin/bash

set -e
set -o pipefail
set -x


function mkDeb {
	ARGS=("$@")

	LSBDISTID="${ARGS[0]}"
	DEB_ARCH="${ARGS[1]}"
	GO_ARCH="${ARGS[2]}"
	GO_ENV=('GOOS=linux' "GOARCH=$GO_ARCH" "${ARGS[@]:3}")

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

	SOURCE_DATE_EPOCH="$(git log -1 --format=%ct)" fpm -s tar -t deb --log debug --verbose --debug \
		-n "$PKG_NAME" \
		-v "$PKG_VERSION" \
		-a "$DEB_ARCH" \
		-m 'Alexander A. Klimov <grandmaster@al2klimov.de>' \
		--description 'The Masif Upgrader agent is a component of Masif Upgrader.
Consult Masif Upgrader'"'"'s manual on its purpose and the agent'"'"'s role in its architecture:
https://github.com/masif-upgrader/manual' \
		--url 'https://github.com/masif-upgrader/agent' \
		-p "${PKG_NAME}-${PKG_VERSION}-${LSBDISTID}-${DEB_ARCH}.deb" \
		-d adduser -d apt -d bash -d systemd --no-auto-depends \
		--config-files /etc/masif-upgrader/agent.ini \
		--after-install packaging/on-install.sh \
		--after-upgrade packaging/on-install.sh \
		--after-remove packaging/on-remove.sh \
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


#     LSBDISTID DEB_ARCH GO_ARCH  GO_ENV

mkDeb Debian    amd64    amd64    GO386=387
mkDeb Debian    i386     386      GO386=387

mkDeb Debian    mips     mips     GOMIPS=softfloat
mkDeb Debian    mipsel   mipsle   GOMIPS=softfloat
mkDeb Debian    mips64el mips64le

mkDeb Debian    ppc64el  ppc64le
mkDeb Debian    s390x    s390x

mkDeb Debian    armel    arm      GOARM=5
mkDeb Debian    armhf    arm      GOARM=7
mkDeb Debian    arm64    arm64

mkDeb Raspbian  armhf    arm      GOARM=6
