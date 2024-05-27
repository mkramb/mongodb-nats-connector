# mongodb-nats-connector

## Functionality

The connector seamlessly synchronizes data between MongoDB and NATS JetStream, thereby offloading error management, retry logic, and duplicate message handling from the services. This functionality is underpinned by the utilization of MongoDB resume tokens. Each change event in MongoDB contains an `_id` field that functions as a resume token. By setting the resumeAfter parameter with a specific resume token value, MongoDB can continue the change stream from the exact event marked by that token.

Upon processing a change event, the connector persists the corresponding resume token into a designated collection. In the event of a restart, the connector queries this collection to retrieve the most recent token, enabling it to resume the change stream from the last processed event. This ensures reliable and continuous data synchronization without loss or duplication of change events.

Lastly, the connector employs the Raft consensus algorithm to eliminate single points of failure in production environments. This allows for multiple instances of the connector to run concurrently, ensuring high availability and fault tolerance. Only the elected master instance performs the processing, thereby maintaining consistency and reliability in data synchronization.

## Example Usage

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

## Local Development

Prerequisite:

```
brew install kind
brew install tilt-dev/tap/tilt
brew install tilt-dev/tap/ctlptl
brew install golangci-lint
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
