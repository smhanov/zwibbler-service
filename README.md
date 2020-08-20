# Installers for Zwibbler Collaboration Server

This project packages the zwibserve collaboration server as a system service that you can install on Linux.

After installation, it will be running on port 3000 as non-https. You can check this by going to http://yourserver.com:3000 in a web browser. You should receive a 404 error if it is working.

The next step is to enable HTTPS using your certificate and private key file.

Edit /etc/zwibbler.conf and change it:

    ServerBindAddress=0.0.0.0
    ServerPort=443
    CertFile=
    KeyFile=

Change CertFile and KeyFile to be the path to your SSL certificate information on the system. CertFile is your certificate, and KeyFile is your private key.

Next, restart the service using

    systemctl restart zwibbler

You can view the logs using

    sudo tail -f /var/log/zwibbler/zwibbler.log

You should now be able to test using https://zwibbler.com/collaboration and entering wss://yourserver/socket in the URL with no port.

## Building
To build, you need to have installed:

* make
* go
* fpm https://fpm.readthedocs.io/en/latest/installing.html

## Structure
It is just a project that pulls in go modules from other places and connects them together. It pulls in the main zwibserve collaboration code from https://github.com/smhanov/zwibserve and go code to run a system service and provide logging. Then I use a Makefile to tell fpm to create an installer.

