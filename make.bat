@ECHO OFF
del zwibbler.exe
go build
set iscc="C:\\Program Files (x86)\\Inno Setup 6\\ISCC.EXE"
set version=11.0
%ISCC% /dversion=%VERSION% install.iss