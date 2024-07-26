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