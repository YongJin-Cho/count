@echo off
set TAG=%1
if "%TAG%"=="" set TAG=latest

echo Building count-processing-service:%TAG%...
pushd src\count-processing-service
go mod tidy
popd
docker build -t count-processing-service:%TAG% src/count-processing-service
