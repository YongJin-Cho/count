@echo off
set TAG=%1
if "%TAG%"=="" set TAG=latest
set IMAGE_NAME=count-api-service

echo Building Docker image %IMAGE_NAME%:%TAG%...

cd src\count-api-service
call go mod tidy
docker build -t %IMAGE_NAME%:%TAG% .

echo Build complete: %IMAGE_NAME%:%TAG%
cd ..\..
