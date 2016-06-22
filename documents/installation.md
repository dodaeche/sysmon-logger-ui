# Installation

The following outlines the instructions for installing the User Interface (UI) server on Ubuntu linux.

The UI server uses the items installed with the analysis server and provides it's own HTTP server and therefore no other software is required.

## Logging
The server logs to it's own specific lfile, which requires the following steps:

- Make a directory for the log
```
sudo mkdir /var/log/sml-ui
```
- Change the directory owner to the user that will run the server
```
sudo chown YOURUSER /var/log/sml-ui
```
- Change the directory group to the user that will run the server
```
sudo chgrp YOURUSER /var/log/sml-ui
```

## Systemd Service

To enable the AutoRun UI server to run on boot up, a Systemd service must be created. The following instructions detail the steps required. The instructions show the installation into the **/opt** directory.

- Create a directory in **/opt**
```
sudo mkdir /opt/sml-ui
```
- Change the directory owner to the user we want to run the application as
```
sudo chown USERNAME /opt/sml-ui
```
- Change the directory group to the user we want to run the application as
```
sudo chgrp USERNAME /opt/sml-ui
```
- Copy the required files to the directory as listed below:
```
sml-ui.config
sml-ui-setbind.sh
sml-ui
static
templates
```
- Set the appropriate configuration options in the config file (arl.config). The configuration values are detailed in the **configuration.pdf** document.
- If the server is to be run on the lower TCP ports such as 80 or 81, then edit the **sml-ui-setbind.sh** file  with the correct paths (if they differ). The file is used to set execute options on the **sml-ui** binary and to allow it to bind to lower TCP ports. The file should be run each time the binary is changed. The files default contents are:
```
#!/bin/sh
chmod +x /opt/sml/sml-ui
sudo setcap cap_net_bind_service+ep /opt/sml/sml-ui
```
- Next define the Systemd service config by creating a new service file:
```
sudo nano /etc/systemd/system/sml-ui.service
```
- Copy the following to the service file or use the file located in the **configuration** directory:
```
[Unit]
Description=SML-UI

[Service]
ExecStart=/opt/sml/sml -c /opt/sml-ui/sml-ui.config

[Install]
WantedBy=multi-user.target
```
- To test the service, run the following command:
```
systemctl start sml-ui.service
```
- To check the service output, run the following command:
```
systemctl status sml-ui.service
```
- If there are problems, then the service can be stopped by running the command:
```
systemctl stop sml-ui.service
```
- Once the configuration is correct and the server is running, the service can be made to run at boot, by running the following command:
```
systemctl enable sml-ui.service
```
