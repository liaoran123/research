
@echo off
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=386
go build  -o winresearch32.exe

SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build  -o winresearch.exe 

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=386
go build  -o linuxresearch32 

SET CGO_ENABLED=0
SET GOOS=linux
SET GOARCH=amd64
go build  -o linuxresearch 

SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=amd64
go build  -o macresearch

SET CGO_ENABLED=0
SET GOOS=darwin
SET GOARCH=386
go build  -o macresearch32