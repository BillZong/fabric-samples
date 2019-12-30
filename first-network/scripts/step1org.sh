#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

# This script is designed to be run in the ${NEW_ORG_NAME}cli container as the
# first step of the EYFN tutorial.  It creates and submits a
# configuration transaction to add ${NEW_ORG_NAME} to the network previously
# setup in the BYFN tutorial.
#

NEW_ORG_NAME="$1"
NEW_ORG_MSPID="$2"
CHANNEL_NAME="$3"
CC_SRC_PATH="$4"
DELAY="$5"
LANGUAGE="$6"
TIMEOUT="$7"
VERBOSE="$8"
: ${CHANNEL_NAME:="mychannel"}
: ${CC_SRC_PATH:="github.com/chaincode/crp"}
: ${DELAY:="3"}
: ${LANGUAGE:="golang"}
: ${TIMEOUT:="10"}
: ${VERBOSE:="false"}
LANGUAGE=`echo "$LANGUAGE" | tr [:upper:] [:lower:]`
COUNTER=1
MAX_RETRY=5

# CC_SRC_PATH="github.com/chaincode/chaincode_example02/go/"
# if [ "$LANGUAGE" = "node" ]; then
# 	CC_SRC_PATH="/opt/gopath/src/github.com/chaincode/chaincode_example02/node/"
# fi

# import utils
. scripts/utils.sh

echo
echo "========= Creating config transaction to add ${NEW_ORG_NAME} to network =========== "
echo

echo "Installing jq"
apt-get -y update && apt-get -y install jq

# Fetch the config for the channel, writing it to config.json
fetchChannelConfig ${CHANNEL_NAME} config.json

# Modify the configuration to append the new org
set -x
jq -s ".[0] * {\"channel_group\":{\"groups\":{\"Application\":{\"groups\": {\"$NEW_ORG_MSPID\":.[1]}}}}}" config.json ./channel-artifacts/${NEW_ORG_NAME}.json > modified_config.json
set +x

# Compute a config update, based on the differences between config.json and modified_config.json, write it as a transaction to ${NEW_ORG_NAME}_update_in_envelope.pb
createConfigUpdate ${CHANNEL_NAME} config.json modified_config.json ${NEW_ORG_NAME}_update_in_envelope.pb

echo
echo "========= Config transaction to add ${NEW_ORG_NAME} to network created ===== "
echo

echo "Signing config transaction"
echo
signConfigtxAsPeerOrg 1 ${NEW_ORG_NAME}_update_in_envelope.pb

echo
echo "========= Submitting transaction from a different peer (peer0.org2) which also signs it ========= "
echo
setGlobals 0 2
set -x
peer channel update -f ${NEW_ORG_NAME}_update_in_envelope.pb -c ${CHANNEL_NAME} -o orderer.example.com:7050 --tls --cafile ${ORDERER_CA}
set +x

echo
echo "========= Config transaction to add ${NEW_ORG_NAME} to network submitted! =========== "
echo

exit 0
