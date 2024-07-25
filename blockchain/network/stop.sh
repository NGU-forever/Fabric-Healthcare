echo "-------Stopping-------"
# 关闭区块链浏览器
docker compose -f explorer/docker-compose.yaml down -v > /dev/null 2>&1
# 关闭区块链网络
./network.sh down  > /dev/null 2>&1
# 删除organizations
rm -rf explorer/organizations 
# 删除mysql容器
docker rm -f fabrictrace-mysql > /dev/null 2>&1
echo "-------Closing-------"