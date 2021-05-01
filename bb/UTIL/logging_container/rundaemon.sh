#!/usr/bin/env bash

docker rm -f BB-DAEMON-A 
docker run --name BB-DAEMON-A -d -v /bitlog:/bitlog --restart=always -t bblogger /bin/logger -log_dir /bitlog/BB -flag_file  /bitlog/BBFLAG


