#!/bin/bash

set -exo pipefail

cd /usr/local/debs

rm -f nothing.deb

fpm -s empty -t deb \
	-n nothing \
	-v "$(perl -e 'print time')" \
	-a all \
	-m 'T.O.N.I. <toni@localhost>' \
	--description Nothing \
	--url 'https://localhost' \
	-p nothing.deb \
	--no-auto-depends

dpkg-scanpackages . /dev/null |gzip >Packages.gz
