#!/bin/bash

echo "******INITIALIZING DATABASE******"

echo "host all \"dodec\" 0.0.0.0/0 trust" >> /var/lib/postgresql/data/pg_hba.conf

echo "starting postgres"
gosu postgres pg_ctl -w start

echo "initializing tables"
gosu postgres psql -h localhost -p 5432 -d dodec -U postgres -a -f /tmp/setup.sql

echo "stopping postgres"
gosu postgres pg_ctl stop

echo "stopped postgres"

echo ""
echo "******DATABASE INITIALIZED******"
