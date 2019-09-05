#!/usr/bin/env bash

#
#   This builds and pushes an image to dockerhub based on a tag
#   Requires docker credentials from circleci environment variables
#

# Get particular image being built from the tag name
# Expected tag format is: full-vendor-path_application-name_version
# e.g. php-7.1-fpm_drupal-8_12

# Exit if any of these commands error
# https://stackoverflow.com/questions/1378274/in-a-bash-script-how-can-i-exit-the-entire-script-if-a-certain-condition-occurs
set -e

DOCKER_LATEST=$(echo "dynamic-vmap:latest")
FULL_VERSION="10.0.0"

# Inject the version dynamically into the dockerfile based on the git TAG
# The empty echo adds a newline. Don't ask... -mjf
echo "Adding LABEL version=\"$FULL_VERSION\" to Dockerfile"
echo "" >> .docker/Dockerfile
echo "LABEL version=\"$FULL_VERSION\"" >> .docker/Dockerfile

echo "Logging in to dockerhub as $DOCKERHUB_USER"
docker login -u$DOCKERHUB_USER -p$DOCKERHUB_PASS

echo "Pulling image..."
# Pull image to 'cache' and potentially save time rebuilding...
# If error because image doesn't exist, let the script continue.
docker pull favish/$DOCKER_LATEST || true

echo "Building image..."
# Putting many tags on the initial build of a container is not supported in the version of Docker circleci runs
docker build -t favish/$DOCKER_LATEST ../.docker

echo "Tagging and pushing image..."
docker push favish/$DOCKER_LATEST
