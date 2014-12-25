#!/bin/bash

cd "${0%/*}"

echo "Building..."
go build "./handlers"
go build "../config"
go build

echo "Installing..."
go install

echo "Running..."
#sudo env "PATH=$PATH" dodecci
dodecci -port 8000
