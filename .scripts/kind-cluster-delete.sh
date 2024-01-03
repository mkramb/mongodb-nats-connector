#!/bin/bash

ctlptl delete cluster kind-kind
kind delete cluster
docker container rm kind-registry --force