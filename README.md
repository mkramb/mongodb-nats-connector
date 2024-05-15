# mongodb-nats-connector

## Usage

The minimum number of nodes required to tolerate faults and still reach consensus is three.
Lets start three separate connector instances (in separate terminals):

```
export MONGO_SERVER_URI="mongodb://localhost:27017/test?replicaSet=tilt&directConnection=true"
export MONGO_WATCH_COLLECTIONS="users,movies"
export NATS_SERVER_URL=nats://127.0.0.1:4222

HTTP_PORT=3000 ./mongodb-nats-connector
HTTP_PORT=3001 ./mongodb-nats-connector
HTTP_PORT=3002 ./mongodb-nats-connector
```

## Development

Prerequisite

```
brew install kind
brew install tilt-dev/tap/tilt
brew install tilt-dev/tap/ctlptl

brew install mongosh
brew tap nats-io/nats-tools
brew install nats-io/nats-tools/nats
```

Setting up local cluster:

```
.scripts/cluster_delete.sh
.scripts/cluster_create.sh
```

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

Only start Nats / Mongo server:

```
tilt up -- --only-infra
```
