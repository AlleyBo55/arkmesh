#!/bin/sh
set -eu

if grep -nE -- '--|—|–' README.md; then
  echo 'README contains forbidden dash style' >&2
  exit 1
fi

unformatted="$(gofmt -l cmd internal)"
if [ -n "$unformatted" ]; then
  printf 'Unformatted Go files:\n%s\n' "$unformatted" >&2
  exit 1
fi

git diff --check
go test -count=1 ./...
go test -race -count=1 ./...
go vet ./...
go build ./...
go test -run='^$' -fuzz='^FuzzCapsuleEnvelopeDecoders$' -fuzztime=3s ./internal/capsule
go test -run='^$' -fuzz='^FuzzIdentityDecoders$' -fuzztime=3s ./internal/identity
go test -run='^$' -fuzz='^FuzzRejectDuplicateKeys$' -fuzztime=3s ./internal/strictjson