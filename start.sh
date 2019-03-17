#!/bin/bash

echo "Waiting for mongo to launch on 27017..."

while ! nc -z mongo 27017; do  
  sleep 0.1 # wait for 1/10 of a second before checking again
done

echo "Mongo launched"
./app

