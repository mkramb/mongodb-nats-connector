# Development

## Setting up local cluster

```
.scripts/kind_cluster_delete.sh
.scripts/kind_cluster_create.sh
```

## Running locally

```
tilt up
tilt down
```

## To only start Nats / Mongo server

```
tilt up -- --only-infra
```
