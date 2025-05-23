#!/bin/sh

set -e

echo "Waiting for Kafka Connect to be ready..."
sleep 10

echo "Registering Debezium connector..."
curl -X POST http://debezium:8083/connectors \
  -H "Content-Type: application/json" \
  -d '{
    "name": "postgres-connector",
    "config": {
      "connector.class": "io.debezium.connector.postgresql.PostgresConnector",
      "database.hostname": "db",
      "database.port": "5432",
      "database.user": "'"${DB_USER}"'",
      "database.password": "'"${DB_PASS}"'",
      "database.dbname": "'"${DB_NAME}"'",
      "database.server.name": "postgres",
      "plugin.name": "pgoutput",
      "slot.name": "debezium_slot",
      "publication.name": "debezium_pub",
      "table.include.list": "public.project,public.project_tag,public.users,public.users_career_tag,public.team",
      "key.converter": "org.apache.kafka.connect.json.JsonConverter",
      "value.converter": "org.apache.kafka.connect.json.JsonConverter",
      "key.converter.schemas.enable": "false",
      "value.converter.schemas.enable": "false"
    }
  }'