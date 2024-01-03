# mongodb-nats-connector

## Prerequisite

```
brew install kind
brew install tilt-dev/tap/tilt
brew install tilt-dev/tap/ctlptl
```

Setting up local cluster:

```
./scripts/kind-cluster-delete.sh
./scripts/kind-cluster-create.sh
```

## Usage

Running scripts:

```
task compile
task execute
```

Running services:

```
tilt up
```