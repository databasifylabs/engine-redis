---
name: adding-a-version-line
description: "Add a new Redis version line to this engine repository: the versions.json entry and its pinned tools, any per-line differences in the engine code, and the local build, smoke, conformance and capabilities check. Use when asked to support a new Redis major or minor release line. Do not use for patch bumps inside an existing line (the weekly refresh opens those)."
---

# Adding a Redis version line

Redis 8 takes the newest release of its major (`8`). Redis 7.2 takes patches of 7.2 only, because 7.4 is not open source. A new open-source major gets its own line. The weekly refresh never moves a pin out of its line.

## The versions.json entry

1. Copy the newest `redis-*` entry under `images` to the key `redis-<line>`. Set `line`, `release` (the exact official tag), and `base.tag` to the official tag, keeping the suffix the other lines use.
2. Set `base.digest`. Get it from `docker buildx imagetools inspect <image>:<tag> --format '{{json .Manifest}}' | jq -r .digest`.
3. Set every package version to `"0"`, then run `$IMAGES refresh redis-<line>`. It pins the packages and moves the base, copies, modules and sources to the newest patch of the line. It may also move the shared `go` digest: revert that unless you mean to rebuild every image.
4. Run `$IMAGES check`. It fails when a key isn't `<engine>-<line>`, when a pin left its line, or when `sdk` differs from `go.mod`.

## What differs per line in this engine

The image sets `DATABASIFY_LINE` from the line, and the engine reads it where releases differ:

- The engine code is the `Redis` flavor of [engine-sentinel-family](https://github.com/databasifylabs/engine-sentinel-family). A per-line difference is a pull request there first, then a bump of `github.com/databasifylabs/engine-sentinel-family` in this repository's `go.mod`.
- redis-shake (`sources.redis-shake`) does the migrations. Check that it supports the new release's RDB version.
- Before adding a line, confirm its license is open source.

Add or extend a Go test for each per-line branch you add.

## Check the new line

`$IMAGES` stands for `tools/images.sh` in a clone of [engine-sdk](https://github.com/databasifylabs/engine-sdk) at the SDK version this repository's `go.mod` requires. Run it from this repository's root:

```bash
go vet ./... && go test -race -count=1 ./...
$IMAGES check
$IMAGES build redis-9 --load --tag smoke:redis-9
$IMAGES smoke smoke:redis-9 redis
echo '{}' | docker run --rm -i --network none --entrypoint databasify smoke:redis-9 engine capabilities | jq -S .
echo '{}' | docker run --rm -i --network none --entrypoint databasify smoke:redis-9 config fields | jq -S .
```

- `smoke` runs `images/redis/smoke/check.sh` against the new line, then the SDK's conformance suite.
- Compare `engine capabilities` and `config fields` with the closest existing line. Every difference should be one you meant, such as a password encryption or setting the new release adds.
- Also run the CI scan (`TRIVY_IMAGE` in the SDK's `.github/workflows/engine-images.yml`) on `smoke:redis-9`.
- Put the results in the pull request with the line and OS you ran them on.

## The pull request

Open an issue first, branch from `main`, and open one pull request against `main` with the `versions.json` entry, any per-line code and tests, and the README's Lines section. The app registers the new line in its own repository ([databasifylabs/databasify](https://github.com/databasifylabs/databasify)), after this one merges.
