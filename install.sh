#!/bin/bash

# Installing updates and required packages
sudo apt update
sudo apt-get update
sudo apt-get install zip unzip -y

# Check if the required environment variables are set
if [[ -z "$POSTGRES_USER" || -z "$POSTGRES_PASSWORD" ]]; then
    echo "Error: Please set POSTGRES_USER and POSTGRES_PASSWORD environment variables."
    exit 1
fi

# Install PostgreSQL
sudo apt install postgresql postgresql-contrib -y


# Start the PostgreSQL server
sudo service postgresql start
sudo pg_isready

# Check if the user already exists
if sudo -u postgres psql -t -c "SELECT 1 FROM pg_roles WHERE rolname='$POSTGRES_USER'" | grep -q 1; then
    # User already exists, so alter the user's password
    sudo -u postgres psql -c "ALTER ROLE $POSTGRES_USER WITH PASSWORD '$POSTGRES_PASSWORD';"
    echo "User '$POSTGRES_USER' password updated."
else
    # User does not exist, create the PostgreSQL user with the specified username and password
    sudo -u postgres psql -c "CREATE USER $POSTGRES_USER WITH PASSWORD '$POSTGRES_PASSWORD';"
    echo "User '$POSTGRES_USER' created."
fi

# Restart the PostgreSQL server
sudo service postgresql restart
sudo pg_isready

echo "PostgreSQL installed, user handling completed, and PostgreSQL server started."

# Define the desired Go version
GO_VERSION="1.21.1"

# Set the download URL for the Go binary (adjust the version if needed)
GO_URL="https://golang.org/dl/go$GO_VERSION.linux-amd64.tar.gz"

# Set the installation directory
INSTALL_DIR="/usr/local"

# Download and extract the Go binary
sudo curl -O -L $GO_URL
sudo tar -C $INSTALL_DIR -xzf go$GO_VERSION.linux-amd64.tar.gz

# Add Go binary directory to PATH (if not already added)
if [[ ":$PATH:" != *":$INSTALL_DIR/go/bin:"* ]]; then
    echo "export PATH=$PATH:$INSTALL_DIR/go/bin" >> ~/.bashrc
    source ~/.bashrc
fi

# Clean up downloaded archive
rm go$GO_VERSION.linux-amd64.tar.gz

# Verify the installation
go version

# Unzip files from tmp directory and copy to necessary locations
cd /tmp/
sudo unzip app.zip
sudo cp webapp-service.service /etc/systemd/system/webapp-service.service
sudo mkdir /usr/webapp
sudo cp app /usr/webapp/
sudo cp users.csv $ACCOUNT_CSV_PATH

# Navigate to webapp
cd /usr/webapp

# Set permissions for server
sudo chmod +x ./app

# Start server as service
sudo systemctl daemon-reload
sudo systemctl enable webapp-service.service
sudo systemctl start webapp-service.service
sudo systemctl status webapp-service.service


