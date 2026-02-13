@echo off
set TAG=%1
if "%TAG%"=="" set TAG=latest

echo Building all services...
call src\scripts\docker-build-management.bat %TAG%
call src\scripts\docker-build-processing.bat %TAG%
