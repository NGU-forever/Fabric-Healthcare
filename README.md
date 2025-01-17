# Fabric-Healthcare

## If the configuration environment is unsuccessful

In the job zip file we uploaded there is a private key to connect to the cloud server, you can connect to the server to see that. The server address is ```47.237.161.142``` and the User is ```root```.
setting
```
Host Fabric
  HostName 47.237.161.142
  User root
  IdentityFile ~/.ssh/NGU.pem
```

## Project github address
https://github.com/NGU-forever/Fabric-Healthcare

Explorer: http://47.237.161.142:8080 (User:exploreradmin, Passsword:exploreradminpw)

<img width="1440" alt="截屏2024-08-01 11 46 58" src="https://github.com/user-attachments/assets/6c90b615-e5fd-4ef3-90fb-04b1ff9057f3">


Backend:http://47.237.161.142:8888

## Project Strcture and Member contribution
```
Fabric-Healthcare/
├── app/
│   ├── backend/             (backend)
│   ├── frontend/            (Web, GUI)
├── blockchain/
│   ├── chaincode/           (Chaincode)
│   ├── network/             (Blockchain)
├── README.md
```

Frank Xiong(z5503242): Blockchain architecture and framwork part

Zenglin Zhong(z5360071):  Chaincode part

Danniel Deng(z5275904): Backend part

Zeyu Wang(z5453260): Frontend

Yifan Hong(z5414347): Off-chain componet

## The part of Blockchain

The blockchain platform used in this project is Hyperledger Fabric, the version of which is V2.5, which has better performance and stability, and the Fabric-gateway mode is used to invoke and use chaincode. The technology stack used in this blockchain is as follows: CouchDB is used to view the blockchain data status and world status, and Hyperledger explorer is used to view the blockchain node status, on-chain situation, chaincode definition, transactions, and so on. The blockchain is deployed on cloud servers throughout to ensure its stability and discharge the numerous problems that occur in virtual machines.

## Installation Instructions (Blockchain)

### 1. Install Docker and Docker Compose

#### Update package lists and install dependencies.

```sh
sudo apt update
sudo apt install apt-transport-https ca-certificates curl software-properties-common
```

#### Add Docker's GPG key and official repository

```sh
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
```

#### Install Docker

```sh
sudo apt update
sudo apt install docker-ce
```

#### Start and enable Docker service

```sh
sudo systemctl start docker
sudo systemctl enable docker
```

#### Download Docker Compose binary

```sh
sudo curl -L "https://github.com/docker/compose/releases/download/1.29.2/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
```

#### Set execute permissions

```sh
sudo chmod +x /usr/local/bin/docker-compose
```

#### Verify installation

```sh
docker --version
docker-compose --version
```

### 2. Install Go

#### Download and Extract Go package

```sh
wget https://golang.org/dl/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
```

#### Set environment variables

```sh
echo "export GOPATH=$HOME/go
export GOROOT=/usr/local/go
export PATH=$GOROOT/bin:$PATH
export PATH=$GOPATH/bin:$PATH" >> ~/.profile
source ~/.profile
```

#### Verify installation

```sh
go version
```

### 3. Install Node and jq

#### Download NVM

```sh
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.3/install.sh | bash
```

#### Download and Use node

```sh
nvm install 16
nvm use 16
```

#### Verify installation

```sh
node -v
```

#### Download jq

```sh
sudo apt install jq
```

### 4. Clone or Unzip Project

#### Clone

```sh
git clone https://github.com/NGU-forever/Fabric-Healthcare.git
```

#### Or unzip
```sh
unzip Fabric-Healthcare-main.zip
```

#### Download the binary files

```sh
cd Fabric-Healthcare/blockchain
curl -sSL https://github.com/hyperledger/fabric/releases/download/v2.5.0/hyperledger-fabric-linux-amd64-2.5.0.tar.gz -o fabric-bin.tar.gz
```

#### Extract the binary files

```sh
tar -xvf fabric-bin.tar.gz
```

```sh
mv bin network/
```

### 5. Blockchain

#### Start Blockchain

```sh
cd blockchain/network
./start.sh
```

For ease of deployment, we wrote a special startup script to launch the blockchain.

#### start.sh

```sh
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
```

#### stop.sh

```sh
echo "-------Stopping-------"
# stop explorer
docker compose -f explorer/docker-compose.yaml down -v > /dev/null 2>&1
# stop network
./network.sh down  > /dev/null 2>&1
# delete organizations
rm -rf explorer/organizations 
# delete mysql image
docker rm -f fabrictrace-mysql > /dev/null 2>&1
echo "-------Closing-------"
```

#### If an error is reported, you can try the following: 

```sh
./network.sh down
docker rm -f $(docker ps -aq)
```

#### Then start again


## Chaincode

smart_contract_addresses can check Hyperledger explorer: http://47.237.161.142:8080 (User:exploreradmin, Passsword:exploreradminpw)

#### Installing go mod
```sh
cd ../chaincode
go mod vendor
```

## Backend

#### Starting the backend requires a Java environment

Because the java file is too large, you need to manually download the jar file and add it in ```Fabric-Healthcare/app/backend``` directory. The download link is as follows: https://drive.google.com/file/d/1-19k2ru4iJzruFTwnOXhVDGPT3rOAKCN/view?usp=sharing

#### Start Backend

```sh
cd ../../app/backend
./app.start
```

#### Stop Backend

```sh
cd java-fabric
./app.stop
```

## Frontend

#### Starting the frontend requires Node.js, npm and Vue-CLI

#### if environment need to been configure

```sh
sudo apt update
sudo apt install nodejs npm
sudo npm install -g @vue/cli
```

### Start fronted

```sh
cd frontend
npm run serve
```


