#!/bin/sh
chmod +x /opt/sml-ui/sml-ui
sudo setcap cap_net_bind_service+ep /opt/sml-ui/sml-ui
