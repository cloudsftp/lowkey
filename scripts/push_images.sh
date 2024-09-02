#!/usr/bin/env bash

version="$1"
user="$2"
password="$3"

echo $password | docker login ghcr.io --username "$user" --password-stdin

tag_and_push() {
    service="$1"
    local_version="$2"

    docker tag "$service" "ghcr.io/cloudsftp/$service:$local_version"
    docker push "ghcr.io/cloudsftp/$service:$local_version"
}

tag_and_push "lowkey" "$version"
tag_and_push "lowkey" "latest"
