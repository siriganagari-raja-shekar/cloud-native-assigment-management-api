#!/bin/bash

# Check if the required environment variables are set
if [[ -z "$PG_USERNAME" || -z "$PG_PASSWORD" ]]; then
    echo "Error: Please set PG_USERNAME and PG_PASSWORD environment variables."
    exit 1
fi

# Install PostgreSQL
sudo apt update
sudo apt install postgresql postgresql-contrib

# Check if the user already exists
if sudo -u postgres psql -t -c "SELECT 1 FROM pg_roles WHERE rolname='$PG_USERNAME'" | grep -q 1; then
    # User already exists, so alter the user's password
    sudo -u postgres psql -c "ALTER ROLE $PG_USERNAME WITH PASSWORD '$PG_PASSWORD';"
    echo "User '$PG_USERNAME' password updated."
else
    # User does not exist, create the PostgreSQL user with the specified username and password
    sudo -u postgres psql -c "CREATE USER $PG_USERNAME WITH PASSWORD '$PG_PASSWORD';"
    echo "User '$PG_USERNAME' created."
fi


# Start the PostgreSQL server
sudo service postgresql start

echo "PostgreSQL installed, user handling completed, and PostgreSQL server started."
