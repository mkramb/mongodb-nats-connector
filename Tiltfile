config.define_bool("only-infra")
cfg = config.parse()

docker_prune_settings(
    disable=False,
    num_builds=3,
    keep_recent=2
)

k8s_yaml([".k8s/nats.yml"])
k8s_resource("nats", port_forwards=["4222:4222"])
