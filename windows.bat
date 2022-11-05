
@echo off
SET CGO_ENABLED=0
SET GOOS=windows
SET GOARCH=amd64
go build -ldflags "-s -w -H=windowsgui"  -o winsfs.exe 

