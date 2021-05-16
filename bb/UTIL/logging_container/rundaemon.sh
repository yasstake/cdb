#!/usr/bin/env bash

docker rm -f BB-DAEMON-A 
docker run --name BB-DAEMON-A -d -v /bitlog:/bitlog --restart=always -t bblogger /bin/logger -log_dir /bitlog/BB -flag_file  /bitlog/BBFLAG -exit_wait 240
sleep 240


docker rm -f BB-DAEMON-B
docker run --name BB-DAEMON-B -d -v /bitlog:/bitlog --restart=always -t bblogger /bin/logger -log_dir /bitlog/BB -flag_file  /bitlog/BBFLAG -exit_wait 240






