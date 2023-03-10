#!/bin/bash

PROGRAM=OCRServer

pid() {
	if [ -f "./pid" ];then
    	PID=$(cat "./pid")
    	EXIST=$(ps aux | awk '{print $2}' | grep -w $PID)
		if [ $EXIST ];then
			echo $PID
			return
		fi
	fi

    echo 0
}

build() {
	go build $PROGRAM.go
	if [ $? -ne 0 ];then
		echo "\033[31m build失败，启动终止 \033[0m"
		return 0
	else
		echo "build...成功"
		return 1
	fi
}

if [ $# -eq 0 ];then
    echo "Usage: $0 {start|stop|status|restart}"
    exit 0
fi


case "$1" in
	status)
		PID=$(pid)
		if [ $PID -gt 0 ];then
			echo "\033[32m 运行中...pid:$PID \033[0m" 
		else
			echo "\033[33m 未运行 \033[0m"
		fi
	;;
	start)
		PID=$(pid)
		if [ $PID -gt 0 ];then
			echo "\033[32m 运行中...pid:$PID \033[0m"
		else
			echo "启动中..."
			build
			if [ $? -ne 0 ];then
				nohup ./$PROGRAM >> ocr.log 2>&1 &
				sleep 2
				PID=$(pid)
				if [ $PID -gt 0 ];then
					echo "\033[32m 启动成功...pid:$PID \033[0m" 
				else
					echo "\033[33m 启动失败 \033[0m"
				fi
			fi
		fi
	;;
	stop)
		PID=$(pid)
		if [ $PID -gt 0 ];then
			echo "停止中..."
			kill -s TERM $PID
			sleep 3
			PID=$(pid)
			if [ $PID -gt 0 ];then
				echo "\033[33m 停止失败...pid:$PID \033[0m"
			else
				echo "\033[32m 已停止 \033[0m" 
			fi
		else
			echo "\033[33m 未运行 \033[0m"
		fi
	;;
	restart)
		PID=$(pid)
		if [ $PID -gt 0 ];then
			echo "重启中...原pid:$PID"
			build
			if [ $? -ne 0 ];then
				kill -s HUP $PID
				sleep 2
				PID=$(pid)
				if [ $PID -gt 0 ];then
					echo "\033[32m 重启成功...新pid:$PID \033[0m" 
				else
					echo "\033[33m 重启失败 \033[0m"
				fi
			fi
		else
			echo "\033[33m 未运行 \033[0m"
		fi
	;;
	*)
    	echo "Usage: $0 {start|stop|status|restart}"
        exit 1
    ;;
esac
