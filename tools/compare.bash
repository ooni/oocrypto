#!/bin/bash
set -euo pipefail
upstreamrepo=./upstreamrepo
TAG=$(cat UPSTREAM)
test -d $upstreamrepo || git clone git@github.com:golang/go.git $upstreamrepo
(
	cd $upstreamrepo
	git checkout master
	git pull
	git checkout $TAG
)
for file in $(cd $upstreamrepo/src/crypto && find . -type f -name \*.go); do
	git diff --no-index $upstreamrepo/src/crypto/$file $file || true
done
