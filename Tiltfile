config.define_bool("only-infra")

cfg = config.parse()
only_infra = cfg.get('only-infra', False)

docker_prune_settings(
    disable=False,
    num_builds=3,
    keep_recent=2
)

k8s_yaml([
    ".k8s/mongodb.yml",
    ".k8s/nats.yml"
])

k8s_resource("nats", port_forwards=["4222:4222"])
k8s_resource("mongodb", port_forwards=["27017:27017"])

local_resource(
    'init-nats',
    cmd='.scripts/init_nats.sh',
    resource_deps=['nats'],
    labels=["scripts"],
    auto_init=True
)

local_resource(
    'init-mongo',
    cmd='.scripts/init_mongo.sh',
    resource_deps=['mongodb'],
    labels=["scripts"],
    auto_init=True
)

if not only_infra:
    k8s_yaml(".k8s/mongodb-nats-connector.yml")
    k8s_resource(
        'mongodb-nats-connector',
        resource_deps=[
            'nats',
            'mongodb',
            'init-nats',
            'init-mongo'
        ]
    )

    docker_build(
        'mongodb-nats-connector',
        dockerfile="Dockerfile",
        context=".",
        network="host"
    )

