#!/bin/bash
#from: http://stackoverflow.com/questions/28244869/creating-a-table-in-single-user-mode-in-postgres

echo "******CREATING DOCKER DATABASE******"

echo "starting postgres"
gosu postgres pg_ctl -w start

echo "initializing tables"
gosu postgres psql -h localhost -p 5432 -U postgres -a -f /tmp/setup.sql

echo "stopping postgres"
gosu postgres pg_ctl stop

echo "stopped postgres"

echo ""
echo "******DOCKER DATABASE CREATED******"
