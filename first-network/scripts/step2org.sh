#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

# TODO: description
# This script is designed to be run in the org3cli container as the
# second step of the EYFN tutorial. It joins the org3 peers to the
# channel previously setup in the BYFN tutorial and install the
# chaincode as version 2.0 on peer0.org3.
#

NEW_ORG_NAME="$1"
NEW_ORG_MSPID="$2"
NEW_ORG_DOMAIN="$3"
PEER_NAME="$4"
PORT="$5"
CHANNEL_NAME="$6"
CC_SRC_PATH="$7"
CC_VERSION="$8"
LANGUAGE="$9"
VERBOSE="$10"

echo "$VERBOSE"

# CHANNEL_NAME="$1"
# DELAY="$2"
# LANGUAGE="$3"
# TIMEOUT="$4"
# VERBOSE="$5"
: ${CHANNEL_NAME:="mychannel"}
: ${CC_SRC_PATH:="github.com/chaincode/crp"}
: ${DELAY:="3"}
: ${LANGUAGE:="golang"}
: ${TIMEOUT:="10"}
: ${VERBOSE:="false"}
LANGUAGE=`echo "$LANGUAGE" | tr [:upper:] [:lower:]`
COUNTER=1
MAX_RETRY=5

echo
echo "========= Getting ${NEW_ORG_NAME} on to your first network ========= "
echo

# CC_SRC_PATH="github.com/chaincode/chaincode_example02/go/"
# if [ "$LANGUAGE" = "node" ]; then
# 	CC_SRC_PATH="/opt/gopath/src/github.com/chaincode/chaincode_example02/node/"
# fi

# import utils
. scripts/utils.sh

echo "Fetching channel config block from orderer..."
set -x
peer channel fetch 0 $CHANNEL_NAME.block -o orderer.example.com:7050 -c $CHANNEL_NAME --tls --cafile $ORDERER_CA >&log.txt
res=$?
set +x
cat log.txt
verifyResult $res "Fetching config block from orderer has Failed"

joinChannelWithArgsAndRetry $NEW_ORG_MSPID $NEW_ORG_DOMAIN $PEER_NAME $PORT
echo "===================== ${PEER_NAME}.${NEW_ORG_NAME} joined channel '$CHANNEL_NAME' ===================== "
echo "Installing chaincode 2.0 on ${PEER_NAME}.${NEW_ORG_NAME}..."
installChaincode 0 3 2.0

echo
echo "========= ${NEW_ORG_NAME} is now halfway onto your first network ========= "
echo

exit 0
