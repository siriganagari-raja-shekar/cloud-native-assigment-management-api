#!/bin/bash

# Check if the required environment variables are set
if [[ -z "$POSTGRES_USER" || -z "$POSTGRES_PASSWORD" ]]; then
    echo "Error: Please set POSTGRES_USER and POSTGRES_PASSWORD environment variables."
    exit 1
fi

# Install PostgreSQL
sudo apt update
sudo apt install postgresql postgresql-contrib


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
