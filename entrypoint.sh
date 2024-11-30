#!/bin/bash

if [ "$1" == "app" ]; then
  echo "Starting app..."
  exec ./app $2 $3
elif [ "$1" == "crawler" ]; then
  echo "Starting crawler..."
  exec ./crawler $2 $3
else
  echo "Usage: $0 {app|crawler}"
  exit 1
fi