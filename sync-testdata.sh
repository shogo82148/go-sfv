#!/usr/bin/env bash

set -uex

OWNER=httpwg
REPO=structured-field-tests
SHA=$(gh api --jq '.commit.sha' "repos/$OWNER/$REPO/branches/main")

rm -rf testdata/structured-field-tests
mkdir -p testdata/structured-field-tests
curl -sSL "https://github.com/$OWNER/$REPO/archive/$SHA.tar.gz" | tar xz -C testdata/structured-field-tests --strip=1

git add testdata
git commit -m "sync test cases with https://github.com/$OWNER/$REPO/commit/$SHA"
