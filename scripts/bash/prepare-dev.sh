#!/bin/bash

echo "copy ./docker/dev/.env.example to docker/dev/.env"
cp -f ${PWD}/docker/dev/.env.example  ${PWD}/docker/dev/.env

echo "copy ./docker/dev/docker-compose.example.yaml to docker/dev/docker-compose.yaml"
cp -f ${PWD}/docker/dev/docker-compose.example.yaml  ${PWD}/docker/dev/docker-compose.yaml


echo "copy ./config.dev.example.yaml to config.dev.yaml"
cp -f ${PWD}/config.dev.example.yaml ${PWD}/config.dev.yaml

