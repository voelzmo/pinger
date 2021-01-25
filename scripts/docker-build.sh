#!/usr/bin/env bash
# build the docker image
VERSION=$1
REPOSITORY=$2
PROJECT=$3


# causes the shell to exit if any subcommand or pipeline returns a non-zero status.
set -e

# set debug mode
#set -x

# build the new docker image
#
echo '>>> Building new image'
docker build --no-cache=true -t $REPOSITORY/$PROJECT:${VERSION} . 

echo '>>> Push new image'
docker push $REPOSITORY/$PROJECT:${VERSION}


