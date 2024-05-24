#!/bin/bash

nats context add nats --server localhost:4222 --description "locahost" --select
nats str add cs --subjects "cs.*.*.*" --defaults

echo "Nats cs stream created successfully"