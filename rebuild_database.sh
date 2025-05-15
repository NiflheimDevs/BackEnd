#!/bin/bash

read -sp "Enter DB password: " PGPASSWORD
export PGPASSWORD

pg_dump --data-only -U niflheim bidlancerdb > data_dump.sql

dropdb -U niflheim bidlancerdb

createdb -U niflheim bidlancerdb

# psql -U niflheim -d bidlancerdb -c "SET app.skip_seeds = 'on';" -f migrations/init.sql
psql -U niflheim -d bidlancerdb -f <(echo "SET app.skip_seeds = 'on'; \i 'migrations/init.sql'")

psql -U niflheim -d bidlancerdb -f data_dump.sql

psql -U niflheim -d bidlancerdb -f migrations/triggers.sql

unset PGPASSWORD


rm data_dump.sql