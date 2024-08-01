#!/bin/bash
./stop.sh
# mysql
docker run --name fabrihealth-mysql -p 3337:3306 -e MYSQL_ROOT_PASSWORD=fabrihealth -d mysql:8

#check images and pulling
image_versions=("2.5.9")

images=("hyperledger/fabric-tools" "hyperledger/fabric-peer" "hyperledger/fabric-orderer" "hyperledger/fabric-ccenv" "hyperledger/fabric-baseos")

for image in "${images[@]}"
do
    for version in "${image_versions[@]}"
    do
        if ! docker images -a | grep "$image" | grep "$version" &> /dev/null
        then
            echo "images $image:$version is none, pulling..."
            docker pull "$image:$version"
        fi
    done
done


# blockchain up and create channels with couchdb
 ./network.sh up createChannel -s couchdb

 # start explorer
cd explorer
export EXPLORER_CONFIG_FILE_PATH=./config.json
export EXPLORER_PROFILE_DIR_PATH=./connection-profile
export FABRIC_CRYPTO_PATH=./organizations
docker-compose down -v
cp -r ../organizations/ .
docker-compose up -d

# deploy chaincode
cd ~/Fabric-Healthcare/blockchain/network/
./network.sh deployCC -ccn mycc3 -ccp ../chaincode -ccl go