# mongodb-nats-connector

## Prerequisite

```
brew install kind
brew install tilt-dev/tap/tilt
brew install tilt-dev/tap/ctlptl
```

Setting up local cluster:

```
.scripts/kind-cluster-delete.sh
.scripts/kind-cluster-create.sh
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
tilt down
```

## Tips

Only start Nats server:

```
tilt up -- --only-infra
```

Start 2 separate connector instances (in separate terminals):

```
HTTP_PORT=3000 NATS_SERVER_URL=nats://127.0.0.1:4222 ./mongodb-nats-connector
HTTP_PORT=3001 NATS_SERVER_URL=nats://127.0.0.1:4222 ./mongodb-nats-connector
```
