#!/usr/bin/env bash

while inotifywait -e modify,move,create,delete -r ./; do
  docker-compose up --build -d
done
