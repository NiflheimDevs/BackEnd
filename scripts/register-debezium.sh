#!/bin/sh

set -e

echo "Waiting for Kafka Connect to be ready..."
until curl -s http://debezium:8083/connectors; do
  echo "Kafka Connect not ready yet. Sleeping..."
  sleep 5
done

echo "Registering Debezium connector..."
curl -X POST http://debezium:8083/connectors \
  -H "Content-Type: application/json" \
  -d '{
    "name": "postgres-connector",
    "config": {
      "connector.class": "io.debezium.connector.postgresql.PostgresConnector",
      "database.hostname": "'"${DB_HOST}"'",
      "database.port": '"${DB_PORT}"',
      "database.user": "'"${DB_USER}"'",
      "database.password": "'"${DB_PASS}"'",
      "database.dbname": "'"${DB_NAME}"'",
      "topic.prefix": "'"${TOPIC_PREFIX_DB}"'",
      "plugin.name": "pgoutput",
      "slot.name": "debezium_slot",
      "publication.name": "debezium_pub",
      "table.include.list": "public.project,public.project_tag,public.users,public.users_career_tag,public.team",
      "key.converter": "org.apache.kafka.connect.json.JsonConverter",
      "value.converter": "org.apache.kafka.connect.json.JsonConverter",
      "key.converter.schemas.enable": "false",
      "value.converter.schemas.enable": "false",
      "producer.max.request.size": "1048576",      
      "producer.buffer.memory": "33554432"         
    }
  }'

sleep 30