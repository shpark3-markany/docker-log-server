#!/bin/bash
docker build -t docker-log-server .
docker save -o docker-log-server.tar docker-log-server