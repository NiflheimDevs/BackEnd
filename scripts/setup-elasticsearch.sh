#!/bin/bash

until curl -s http://elasticsearch:9200; do
  echo "Waiting for Elasticsearch..."
  sleep 5
done

curl -X PUT http://elasticsearch:9200/users -H 'Content-Type: application/json' -d @/create_index_users.json
curl -X PUT http://elasticsearch:9200/teams -H 'Content-Type: application/json' -d @/create_index_teams.json
curl -X PUT http://elasticsearch:9200/projects -H 'Content-Type: application/json' -d @/create_index_projects.json

echo "done"
