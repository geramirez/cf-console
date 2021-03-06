#!/bin/bash

# This script will install the autopilot plugin, login, pick the right manifest and deploy the app with 0 downtime.

# Command line arguments
# $1 = specified org
# $2 = manifest file path
# $3 = file path to project

set -e

CF_PATH="."
if [ "$TRAVIS_BRANCH" == "master" ]
then
	CF_MANIFEST="manifests/manifest-master.yml"
	CF_SPACE="staging"
elif [ "$TRAVIS_BRANCH" == "production" ]
then
	CF_MANIFEST="manifests/manifest-production.yml"
	CF_SPACE="prod"
fi

echo $CF_MANIFEST
echo $CF_SPACE

export CF_BIN=$CLIPATH/out/cf
# Log in
$CF_BIN api $CF_API
$CF_BIN login --u $CF_USERNAME --p $CF_PASSWORD --o $CF_ORGANIZATION --s $CF_SPACE && $CF_BIN install-plugin $GOAUTOPILOT_BIN
# Run autopilot plugin
$CF_BIN zero-downtime-push cf-deck -f $CF_MANIFEST -p $CF_PATH
