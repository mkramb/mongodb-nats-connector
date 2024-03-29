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

Start 3 separate connector instances (in separate terminals):

```
./mongodb-nats-connector -cluster=connector -size=3 -nats=nats://127.0.0.1:4222
./mongodb-nats-connector -cluster=connector -size=3 -nats=nats://127.0.0.1:4222
./mongodb-nats-connector -cluster=connector -size=3 -nats=nats://127.0.0.1:4222
```

## Resources

- [graft library](https://github.com/nats-io/graft)
- [example from argo-events](https://github.com/argoproj/argo-events/blob/bec39df2b40ed761d2f8786b944ce70b8d22e8fc/common/leaderelection/leaderelection.go#L160)

