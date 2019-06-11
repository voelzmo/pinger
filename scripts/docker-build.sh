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
# Due to a bug in Docker we need to analyse the log to find out if build passed (see https://github.com/dotcloud/docker/issues/1875)
docker build --no-cache=true -t $REPOSITORY/$PROJECT:${VERSION} . | tee /tmp/docker_build_result.log
RESULT=$(cat /tmp/docker_build_result.log | tail -n 1)
if [[ "$RESULT" != *Successfully* ]];
then
  exit -1
fi


echo '>>> Push new image'
docker push $REPOSITORY/$PROJECT:${VERSION}


