@ECHO OFF
go build
set iscc="C:\\Program Files (x86)\\Inno Setup 6\\ISCC.EXE"
set version=2.0
%ISCC% /dversion=%VERSION% install.iss