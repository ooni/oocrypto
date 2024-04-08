#!/bin/bash
set -euo pipefail
upstreamrepo=$(mktemp -d)/upstreamrepo
TAG=$(cat UPSTREAM)
test -d $upstreamrepo || git clone git@github.com:golang/go.git $upstreamrepo
(
	cd $upstreamrepo
	git checkout master
	git pull
	git checkout $TAG
)
diff -ur $upstreamrepo/src/crypto .
