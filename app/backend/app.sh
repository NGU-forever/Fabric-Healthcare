#!/bin/bash
#当前路径
CURRENT_PATH=`pwd`
# 定义Spring Boot应用的名称
APP_NAME="fabric-0.0.1-SNAPSHOT"

# 定义Spring Boot应用的JAR文件
APP_JAR="$APP_NAME.jar"

# 定义Spring Boot应用的PID文件
PID_FILE="application.pid"

# Spring Boot应用的内存设置
MEM_OPTS="-Xms512m -Xmx1024m"

# GC日志的输出路径和文件名
GC_LOG_FILE="$CURRENT_PATH/gc.log"
GC_LOG_OPTS="-XX:+UseG1GC -XX:+PrintGCDetails -Xloggc:gc.log -Xlog:gc+heap=trace"

sudo rm -rf ./crypto-config/ordererOrganizations
sudo rm -rf ./crypto-config/peerOrganizations
sudo cp ../Fabric-Healthcare/blockchain/network/organizations/ordererOrganizations
sudo cp ../Fabric-Healthcare/blockchain/network/organizations/peerOrganizations
# 定义启动Spring Boot应用的命令
start() {
    stop
    echo "Starting $APP_NAME..."
    nohup java $MEM_OPTS $GC_LOG_OPTS -jar -Dotel.traces.exporter=none -Dotel.metrics.exporter=none $APP_JAR  > /dev/null 2>&1 &
    # nohup java $MEM_OPTS $GC_LOG_OPTS -jar $APP_JAR -Dotel.traces.exporter=none -Dotel.metrics.exporter=none  &
    echo "$APP_NAME started successfully."
}

# 定义停止Spring Boot应用的命令
stop() {
    PID=`jps |grep $APP_NAME |awk  '{print $1}'`
    if [ -z "$PID" ]; then
        echo "$APP_NAME not start"
    else
        echo "Stopping $APP_NAME..."
        kill -9 $PID
        echo "$APP_NAME stopped successfully."
    fi
}

# 根据传入的参数调用相应的函数
case "$1" in
    start)
        start
        ;;
    stop)
        stop
        ;;
    *)
        echo "Usage: $0 {start|stop}"
        exit 1
esac

exit 0