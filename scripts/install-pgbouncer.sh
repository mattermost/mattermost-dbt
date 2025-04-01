#!/bin/bash

# Basic pgbouncer setup script for standard Ubuntu EC2 instance.
# Usage: install-pgbouncer.sh HOST_IP_OR_DNS PGBOUNCER_INI PGBOUNCER_USERLIST

if [ -z "$DBT_PEM" ]; then
    echo "Error: DBT_PEM is empty; set with PEM file location"
    exit 1
fi

HOST=$1
INI=$2
USERLIST=$3

if [ -z "$HOST" ] || [ -z "$INI" ] || [ -z "$USERLIST"]; then
    echo "Usage: install-pgbouncer.sh HOST_IP_OR_DNS PGBOUNCER_INI PGBOUNCER_USERLIST"
    echo "Error: missing input value"
    exit 1
fi

echo "### INSTALLING PGBOUNCER ON HOST $HOST"
ssh -i "$DBT_PEM" ubuntu@$HOST 'sudo apt update && sudo apt install -y pgbouncer'
scp -i "$DBT_PEM" $INI ubuntu@$HOST:/home/ubuntu/pgbouncer.ini
ssh -i "$DBT_PEM" ubuntu@$HOST 'sudo mv /home/ubuntu/pgbouncer.ini /etc/pgbouncer/pgbouncer.ini'
scp -i "$DBT_PEM" $USERLIST ubuntu@$HOST:/home/ubuntu/userlist.txt
ssh -i "$DBT_PEM" ubuntu@$HOST 'sudo mv /home/ubuntu/userlist.txt /etc/pgbouncer/userlist.txt'
ssh -i "$DBT_PEM" ubuntu@$HOST 'sudo systemctl reload pgbouncer.service && sudo systemctl restart pgbouncer.service'
echo ""

echo "DONE"
