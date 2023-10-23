#!/bin/bash

# Installing updates and required packages
sudo apt-get update
sudo apt-get upgrade -y
sudo apt-get install zip unzip -y

# Create user and groups
sudo groupadd $LINUX_GROUP
sudo useradd -s /bin/false -g $LINUX_GROUP -d $USER_HOME_DIR -m $LINUX_USER

# Unzip files from tmp directory and copy to necessary locations
cd /tmp/
sudo unzip app.zip
sudo cp webapp-service.service /etc/systemd/system/webapp-service.service
sudo cp app $USER_HOME_DIR
sudo cp users.csv $ACCOUNT_CSV_PATH

# Navigate to webapp
cd /opt/webapp

# Set permissions for server
sudo chown $LINUX_USER:$LINUX_GROUP ./app
sudo chmod +x ./app

# Start server as service
sudo systemctl daemon-reload
sudo systemctl enable webapp-service.service

sudo apt-get clean


