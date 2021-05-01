#!/usr/bin/env bash

docker rm -f BB-DAEMON-A 
docker rm -f BB-DAEMON-B

docker run --name BB-DAEMON-A -d -v /bitlog:/bitlog --restart=always -t bblogger /bin/logger -log_dir /bitlog/BB -flag_file  /bitlog/BBFLAG 120 -exit_wait 120

sleep 120

docker run --name BB-DAEMON-B -d -v /bitlog:/bitlog --restart=always -t bblogger /bin/logger -log_dir /bitlog/BB -flag_file  /bitlog/BBFLAG 120 -exit_wait 120





