#!/bin/bash

if [ "$1" == "app" ]; then
  echo "Starting app..."
  ./app
elif [ "$1" == "crawler" ]; then
  echo "Starting crawler..."
  ./crawler
else
  echo "Usage: $0 {app|crawler}"
  exit 1
fi