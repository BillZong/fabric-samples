#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

# This script is designed to be run in the org3cli container as the
# final step of the EYFN tutorial. It simply issues a couple of
# chaincode requests through the org3 peers to check that org3 was
# properly added to the network previously setup in the BYFN tutorial.
#

echo
echo " ____    _____      _      ____    _____ "
echo "/ ___|  |_   _|    / \    |  _ \  |_   _|"
echo "\___ \    | |     / _ \   | |_) |   | |  "
echo " ___) |   | |    / ___ \  |  _ <    | |  "
echo "|____/    |_|   /_/   \_\ |_| \_\   |_|  "
echo
echo "Extend your first network (EYFN) test"
echo

# 必须<=9个参数
NEW_ORG_NAME="$1"
NEW_ORG_MSPID="$2"
NEW_ORG_DOMAIN="$3"
PEER_NAME="$4"
PORT="$5"
CC_NAME="$6"

DELAY=3
: ${CHANNEL_NAME:="mychannel"}
: ${TIMEOUT:="10"}
: ${VERBOSE:="false"}

echo "Channel name : "$CHANNEL_NAME

# import functions
. scripts/utils.sh

# Query chaincode on peer0.org3, check if the result is 90
echo "Querying chaincode on  ${PEER_NAME}.${NEW_ORG_NAME}..."
chaincodeQueryWithArgs $NEW_ORG_MSPID $NEW_ORG_DOMAIN $PEER_NAME $PORT \
  $CC_NAME '{"Args":["query","a"]}' 90

exit 0

# Invoke chaincode on peer0.org1, peer0.org2, and peer0.org3
echo "Sending invoke transaction on peer0.org1 peer0.org2 peer0.org3..."
chaincodeInvoke 0 1 0 2 0 3

# Query on chaincode on peer0.org3, peer0.org2, peer0.org1 check if the result is 80
# We query a peer in each organization, to ensure peers from all organizations are in sync
# and there is no state fork between organizations.
echo "Querying chaincode on peer0.org3..."
chaincodeQuery 0 3 80

echo "Querying chaincode on peer0.org2..."
chaincodeQuery 0 2 80

echo "Querying chaincode on peer0.org1..."
chaincodeQuery 0 1 80


echo
echo "========= All GOOD, EYFN test execution completed =========== "
echo

echo
echo " _____   _   _   ____   "
echo "| ____| | \ | | |  _ \  "
echo "|  _|   |  \| | | | | | "
echo "| |___  | |\  | | |_| | "
echo "|_____| |_| \_| |____/  "
echo

exit 0
