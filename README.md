# mongodb-nats-connector

## Usage

The minimum number of nodes required to tolerate faults and still reach consensus is three.
Lets start three separate connector instances (in separate terminals):

```
export MONGO_URI="mongodb://localhost:27017/test?replicaSet=tilt&directConnection=true"
export MONGO_WATCH_COLLECTIONS="users,movies"
export NATS_SERVER_URL=nats://127.0.0.1:4222
export NATS_STREAM_NAME=cs

HTTP_PORT=3000 ./mongodb-nats-connector
HTTP_PORT=3001 ./mongodb-nats-connector
HTTP_PORT=3002 ./mongodb-nats-connector
```

To disable raft server, which is useful for local development:

```
# by default set to 3
export RAFT_CLUSTER_SIZE=1
```

## Development

Prerequisite:

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

Running locally:

```
tilt up
tilt down
```

To only start Nats / Mongo server:

```
tilt up -- --only-infra
```
