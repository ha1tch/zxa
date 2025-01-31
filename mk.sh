// file: /mk.sh

#!/bin/bash
mkdir -p bin

echo "Building all binaries for zxa..."
pushd cmd/zxa
GOOS=windows GOARCH=amd64 go build -x -o ../../bin/zxa.win64.exe   main.go
GOOS=windows GOARCH=386   go build -x -o ../../bin/zxa.win32.exe   main.go
GOOS=linux   GOARCH=amd64 go build -x -o ../../bin/zxa.linux64     main.go
GOOS=linux   GOARCH=386   go build -x -o ../../bin/zxa.linux32     main.go
GOOS=darwin  GOARCH=arm64 go build -x -o ../../bin/zxa.mac64.m1    main.go
GOOS=darwin  GOARCH=amd64 go build -x -o ../../bin/zxa.mac64.intel main.go
popd
