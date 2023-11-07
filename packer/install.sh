#!/bin/bash

# Installing updates and required packages
sudo apt-get update
sudo apt-get upgrade -y
sudo apt-get install zip unzip -y

# Create user and groups
sudo groupadd $LINUX_GROUP
sudo useradd -s /bin/false -g $LINUX_GROUP -d $USER_HOME_DIR -m $LINUX_USER

# Install AWS Cloudwatch agent
sudo wget https://amazoncloudwatch-agent.s3.amazonaws.com/debian/amd64/latest/amazon-cloudwatch-agent.deb
sudo dpkg -i -E ./amazon-cloudwatch-agent.deb
sudo rm ./amazon-cloudwatch-agent.deb

# Unzip files from tmp directory and copy to necessary locations
cd /tmp/
sudo unzip app.zip
sudo cp webapp-service.service /etc/systemd/system/webapp-service.service
sudo cp app $USER_HOME_DIR
sudo cp users.csv $ACCOUNT_CSV_PATH
sudo cp aws-cloudwatch-config.json $CLOUDWATCH_CONFIG_FILE

# Create log files
sudo touch $STANDARD_LOG_FILE
sudo touch $ERROR_LOG_FILE
sudo chown $LINUX_USER:$LINUX_GROUP $STANDARD_LOG_FILE $ERROR_LOG_FILE
sudo chmod 600 $STANDARD_LOG_FILE $ERROR_LOG_FILE

# Navigate to webapp
cd $USER_HOME_DIR

# Set permissions for server
sudo chown $LINUX_USER:$LINUX_GROUP ./app
sudo chmod +x ./app

# Start server as service
sudo systemctl daemon-reload
sudo systemctl enable webapp-service.service

sudo apt-get clean


