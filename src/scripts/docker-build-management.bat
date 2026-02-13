@echo off
set TAG=%1
if "%TAG%"=="" set TAG=latest

echo Building count-management-service:%TAG%...
pushd src\count-management-service
go mod tidy
popd
docker build -t count-management-service:%TAG% src/count-management-service
