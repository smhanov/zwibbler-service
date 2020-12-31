; -- Example1.iss --
; SEE THE DOCUMENTATION FOR DETAILS ON CREATING .ISS SCRIPT FILES!

[Setup]
AppName=Zwibbler Collaboration Service
AppVerName=Zwibbler Collaboration Service {#version}
DefaultDirName={pf}\Zwibbler Collaboration Service
DefaultGroupName=Zwibbler
OutputDir=.
Compression=lzma
SolidCompression=yes
OutputBaseFilename=ZwibblerCollaborationService{#version}

[Dirs]

[Files]
Source: "zwibbler.exe"; DestDir: "{app}"; 
Source: "zwibbler.conf"; DestDir: "\zwibbler"; Flags: onlyifdoesntexist uninsneveruninstall;

[Icons]
Name: "{group}\Test your server"; Filename: "https://zwibbler.com/collaboration/testing.html"
Name: "{group}\Configuration file"; Filename: "\zwibbler\zwibbler.conf"
Name: "{group}\Log file"; Filename: "\zwibbler\zwibbler.log"
Name: "{group}\Uninstall Zwibbler Collaboration Service"; Filename: "{app}\unins000.exe"

[Registry]

[Run]
Filename: "{app}\ZWIBBLER.EXE"; Parameters: "--install"

[UninstallRun]
Filename: "{app}\ZWIBBLER.EXE"; Parameters: "--uninstall"
