# mongodb-nats-connector

The connector seamlessly synchronizes data between MongoDB and NATS JetStream, thereby offloading error management, retry logic, and duplicate message handling from the services. This functionality is underpinned by the utilization of MongoDB resume tokens. Each change event in MongoDB contains an \_id field that functions as a resume token. By setting the resumeAfter parameter with a specific resume token value, MongoDB can continue the change stream from the exact event marked by that token. Upon processing a change event, the connector persists the corresponding resume token into a designated collection. In the event of a restart, the connector queries this collection to retrieve the most recent token, enabling it to resume the change stream from the last processed event. This ensures reliable and continuous data synchronization without loss or duplication of change events.

## Functionality

The connector still faces a few challenges:

- What if the connector crashes before publishing the message to NATS and persisting the resume token?
- What if it fails to publish the message to NATS, perhaps due to unavailability?
- What if it successfully publishes the message to NATS, but fails to persist the resume token?

In each scenario, the connector retries the operation, restarting from the previous resume token. While the first two scenarios do not present significant issues, the third scenario may lead to a duplicate message being published to NATS, as the message was already sent before the failure occurred. To mitigate this, the connector utilizes the current resume token as the Nats-Msg-Id header. This ensures that service consumers are not burdened with handling duplicate messages, as NATS JetStream automatically discards any duplicates based on the Nats-Msg-Id header. Additionally, existing change stream services have a limitation in immediately storing seen resume tokens. If not all tokens are successfully processed due to an error or service restart, events may be lost. The connector addresses this by always resuming from the last successfully processed resume token, thereby preventing event loss. Furthermore, it limits concurrent processing to prevent memory spikes, particularly during batch updates.

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
