# Installers for Zwibbler Collaboration Server

This project packages the zwibserve collaboration server as a system service that you can install on Linux.

Please refer to the instructions at https://github.com/smhanov/zwibserve, which also includes the link to the installers.

## Building
To build, you need to have installed:

* make
* [go](https://golang.org/)
* [fpm](https://fpm.readthedocs.io/en/latest/installing.html)

## Structure
It is just a project that pulls in go modules from other places and connects them together. It pulls in the main zwibserve collaboration code from https://github.com/smhanov/zwibserve and go code to run a system service and provide logging. Then I use a Makefile to tell fpm to create an installer.

