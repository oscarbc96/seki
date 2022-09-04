---
sidebar_position: 2
---

# Install

## Install the pre-compiled binary

### homebrew tap

```bash
brew install 
```

### go install

Not working 
```
The go.mod file for the module providing named packages contains one or
more replace directives. It must not contain directives that would cause
it to be interpreted differently than if it were the main module.
```
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

### Clone

```bash
git clone https://github.com/oscarbc96/seki
cd seki
```

### Install dependencies

```bash
make install
```

### Build

```bash
make build
```
