#!/bin/bash

echo "copy ./docker/dev/.env.example to docker/dev/.env"
cp -f ${PWD}/docker/dev/.env.example  ${PWD}/docker/dev/.env

echo "copy ./docker/dev/docker-compose.example.yaml to docker/dev/docker-compose.yaml"
cp -f ${PWD}/docker/dev/docker-compose.example.yaml  ${PWD}/docker/dev/docker-compose.yaml
