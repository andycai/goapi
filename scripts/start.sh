#!/bin/bash

SERVICE_NAME=$1
if [ -z "$1" ]; then
	SERVICE_NAME="unitool_serve_linux"
fi

nohup ./$SERVICE_NAME -port 8080 > output.log 2>&1 &
