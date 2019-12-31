#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

# This script is designed to be run in the cli container as the third
# step of the EYFN tutorial. It installs the chaincode as version 2.0
# on peer0.org1 and peer0.org2, and uprage the chaincode on the
# channel to version 2.0, thus completing the addition of org3 to the
# network previously setup in the BYFN tutorial.
#

# 必须<=9个参数
NEW_ORG_NAME="$1"
NEW_ORG_MSPID="$2"
NEW_ORG_DOMAIN="$3"
PEER_NAME="$4"
PORT="$5"
CC_SRC_PATH="$6"
CC_NAME="$7"
CC_VERSION="$8"

LANGUAGE="golang"

: ${CHANNEL_NAME:="mychannel"}
: ${DELAY:="3"}
: ${TIMEOUT:="10"}
: ${VERBOSE:="false"}
COUNTER=1
MAX_RETRY=5

echo
echo "========= Finish adding $NEW_ORG_NAME to your first network ========= "
echo

# import utils
. scripts/utils.sh

echo "===================== Installing chaincode $CC_VERSION on peer0.org1 ===================== "
installChaincodeWithArgs Org1MSP org1.example.com peer0 7051 $CC_NAME $CC_VERSION
echo "===================== Installing chaincode $CC_VERSION on peer0.org2 ===================== "
installChaincodeWithArgs Org2MSP org2.example.com peer0 9051 $CC_NAME $CC_VERSION

echo "===================== Upgrading chaincode on peer0.org1 ===================== "
# TODO: 确认合约具体有多少组织需要介入认可
upgradeChaincodeWithArgs Org1MSP org1.example.com peer0 7051 \
  '{"Args":["init","a","90","b","210"]}' "AND ('Org1MSP.peer','Org2MSP.peer','$NEW_ORG_MSPID.peer')" $CC_NAME $CC_VERSION

echo
echo "========= Finished adding $NEW_ORG_NAME to your first network! ========= "
echo

exit 0
