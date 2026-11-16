---
name: adding-a-version-line
description: "Add a new Redis version line to this engine repository: the versions.json entry and its pinned tools, then the local build, smoke test, runtime conformance and scan. Use when asked to support a new Redis major or minor release line. Do not use for patch bumps inside an existing line (the weekly refresh opens those)."
---

# Adding a Redis version line

Redis 8 takes the newest release of its major (`8`). Redis 7.2 takes patches of 7.2 only, because 7.4 is not open source. A new open-source major gets its own line. The weekly refresh never moves a pin out of its line.

## The versions.json entry

1. Copy the newest `redis-*` entry under `images` to the key `redis-<line>`. Set `line`, `release` (the exact official tag), and `base.tag` to the official tag, keeping the suffix the other lines use.
2. Set `base.digest`. Get it from `docker buildx imagetools inspect <image>:<tag> --format '{{json .Manifest}}' | jq -r .digest`.
3. Set every package version to `"0"`, then run `$IMAGES refresh redis-<line>`. It pins the packages and moves the base, copies, modules and sources to the newest patch of the line. It may also move the shared `go` digest (the image used to build tools from source): revert that unless you mean to rebuild every image.
4. Run `$IMAGES check`. It fails when a key isn't `<engine>-<line>`, when a pin left its line, when `sdk` is not a 1.x release, or when the Dockerfile's `AS sdk` stage is not the pinned `engine-sdk` image.

## What differs per line

- **In this image:** only the pins. Check that redis-shake (`sources`) support the new release.
- **In the app:** engine behaviour per line lives in the Redis driver (`app/Engines/Drivers/Redis`, on the shared `Sentinel` base in [databasifylabs/databasify](https://github.com/databasifylabs/databasify)): `facts($line)`, the settings catalog, and any `$line` branch in its Blade scripts. The app also registers the new line. That pull request comes after this one merges.

## Check the new line

`$IMAGES` stands for `tools/images.sh` in a clone of [engine-sdk](https://github.com/databasifylabs/engine-sdk) at the tag the `sdk` pin names (`v1.0.0`). Run it from this repository's root:

```bash
$IMAGES check
$IMAGES build redis-<line> --load --tag smoke:redis-<line>
$IMAGES smoke smoke:redis-<line> redis
```

- `smoke` runs `images/redis/smoke/check.sh` against the new line, then the SDK's runtime conformance suite.
- Also run the CI scan (`TRIVY_IMAGE` in the SDK's `.github/workflows/engine-images.yml`) on `smoke:redis-<line>`.
- Put the results in the pull request with the line and OS you ran them on.

## The pull request

Open an issue first, branch from `main`, and open one pull request against `main` with the `versions.json` entry and the README's Lines section. The app registers the new line in its own repository ([databasifylabs/databasify](https://github.com/databasifylabs/databasify)), after this one merges.
