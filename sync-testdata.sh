#!/usr/bin/env bash

set -uex

OWNER=httpwg
REPO=structured-field-tests
SHA=$(gh api --jq '.commit.sha' "repos/$OWNER/$REPO/branches/main")

rm -rf testdata
mkdir testdata
curl -sSL "https://github.com/$OWNER/$REPO/archive/$SHA.tar.gz" | tar xz -C testdata --strip=1

git add testdata
git commit -m "sync test cases with https://github.com/$OWNER/$REPO/commit/$SHA"
