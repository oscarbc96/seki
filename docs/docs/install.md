---
sidebar_position: 2
---

# Install

## Install the pre-compiled binary

### Manually

Download the pre-compiled binaries from the [releases page](https://github.com/oscarbc96/seki/releases) and copy it to the desired location.

### Using go

```bash
go install github.com/oscarbc96/seki@latest
```

## Run with docker

```bash
docker run --rm \
  -v $PWD:/my-repo \
  ghcr.io/oscarbc96/seki /my-repo
```

## Compiling from source

```bash
git clone https://github.com/oscarbc96/seki
cd seki
make build
```

## Continuous Integration

### Github Actions

```yaml
steps:
  - uses: oscarbc96/setup-seki
```

See configuration at [https://github.com/oscarbc96/setup-seki](https://github.com/oscarbc96/setup-seki)
