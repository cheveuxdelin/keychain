#!/bin/bash
action=$1

if [[ $action = 'build' ]];
then
    # windows
    env GOOS=windows GOARCH=amd64 go build -o ./build/keychain-windows-amd64.exe
    # macos
    env GOOS=darwin GOARCH=arm64 go build -o ./build/keychain-macos-sillicon
    env GOOS=darwin GOARCH=amd64 go build -o ./build/keychain-macos-intel
    # linux
    env GOOS=linux GOARCH=amd64 go build -o ./build/keychain-linux-amd64
    env GOOS=linux GOARCH=arm go build -o ./build/keychain-linux-arm
    env GOOS=linux GOARCH=arm64 go build -o ./build/keychain-linux-arm64
elif [[ $action = 'test' ]];
then
    go test ./... | grep github.com/cheveuxdelin/keychain/
fi