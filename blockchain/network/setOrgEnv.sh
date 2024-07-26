#!/bin/bash
#
# SPDX-License-Identifier: Apache-2.0

# default to using Org1
ORG=${1:-Org1}

# Exit on first error, print all commands.
set -e
set -o pipefail

# Log file
LOG_FILE="/tmp/setOrgEnv.log"

# Redirect stdout and stderr to the log file
exec > >(tee -i ${LOG_FILE})
exec 2>&1

# Where am I?
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )"

# Updated paths to match the new directory structure
ORDERER_CA=${DIR}/network/organizations/ordererOrganizations/example.com/tlsca/tlsca.example.com-cert.pem
PEER0_ORG1_CA=${DIR}/network/organizations/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem
PEER0_ORG2_CA=${DIR}/network/organizations/peerOrganizations/org2.example.com/tlsca/tlsca.org2.example.com-cert.pem
PEER0_ORG3_CA=${DIR}/network/organizations/peerOrganizations/org3.example.com/tlsca/tlsca.org3.example.com-cert.pem

echo "DIR: $DIR"
echo "ORDERER_CA: $ORDERER_CA"
echo "PEER0_ORG1_CA: $PEER0_ORG1_CA"
echo "PEER0_ORG2_CA: $PEER0_ORG2_CA"
echo "PEER0_ORG3_CA: $PEER0_ORG3_CA"

if [[ ${ORG,,} == "org1" || ${ORG,,} == "digibank" ]]; then
   CORE_PEER_LOCALMSPID=Org1MSP
   CORE_PEER_MSPCONFIGPATH=${DIR}/network/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
   CORE_PEER_ADDRESS=localhost:7051
   CORE_PEER_TLS_ROOTCERT_FILE=${DIR}/network/organizations/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem

elif [[ ${ORG,,} == "org2" || ${ORG,,} == "magnetocorp" ]]; then
   CORE_PEER_LOCALMSPID=Org2MSP
   CORE_PEER_MSPCONFIGPATH=${DIR}/network/organizations/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp
   CORE_PEER_ADDRESS=localhost:9051
   CORE_PEER_TLS_ROOTCERT_FILE=${DIR}/network/organizations/peerOrganizations/org2.example.com/tlsca/tlsca.org2.example.com-cert.pem

else
   echo "Unknown \"$ORG\", please choose Org1/Digibank or Org2/Magnetocorp"
   echo "For example to get the environment variables to set up an Org2 shell environment run:  ./setOrgEnv.sh Org2"
   echo
   echo "This can be automated to set them as well with:"
   echo
   echo 'export $(./setOrgEnv.sh Org2 | xargs)'
   exit 1
fi

echo "CORE_PEER_LOCALMSPID: $CORE_PEER_LOCALMSPID"
echo "CORE_PEER_MSPCONFIGPATH: $CORE_PEER_MSPCONFIGPATH"
echo "CORE_PEER_ADDRESS: $CORE_PEER_ADDRESS"
echo "CORE_PEER_TLS_ROOTCERT_FILE: $CORE_PEER_TLS_ROOTCERT_FILE"

# output the variables that need to be set
echo "CORE_PEER_TLS_ENABLED=true"
echo "ORDERER_CA=${ORDERER_CA}"
echo "PEER0_ORG1_CA=${PEER0_ORG1_CA}"
echo "PEER0_ORG2_CA=${PEER0_ORG2_CA}"
echo "PEER0_ORG3_CA=${PEER0_ORG3_CA}"

echo "CORE_PEER_MSPCONFIGPATH=${CORE_PEER_MSPCONFIGPATH}"
echo "CORE_PEER_ADDRESS=${CORE_PEER_ADDRESS}"
echo "CORE_PEER_TLS_ROOTCERT_FILE=${CORE_PEER_TLS_ROOTCERT_FILE}"

echo "CORE_PEER_LOCALMSPID=${CORE_PEER_LOCALMSPID}"
