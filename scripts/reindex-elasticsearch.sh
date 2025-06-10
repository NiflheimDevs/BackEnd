#!/bin/bash

until curl -s http://elasticsearch:9200; do
  echo "Waiting for Elasticsearch..."
  sleep 5
done

curl -X PUT http://elasticsearch:9200/users_v2 -H 'Content-Type: application/json' -d @create_index_users.json
curl -X PUT http://elasticsearch:9200/teams_v2 -H 'Content-Type: application/json' -d @create_index_teams.json
curl -X PUT http://elasticsearch:9200/projects_v2 -H 'Content-Type: application/json' -d @create_index_projects.json

curl -X POST http://elasticsearch:9200/_reindex -H 'Content-Type: application/json' -d '{
  "source": {
    "index": "users"
  },
  "dest": {
    "index": "users_v2"
  }
}'
curl -X POST http://elasticsearch:9200/_reindex -H 'Content-Type: application/json' -d '{
  "source": {
    "index": "projects"
  },
  "dest": {
    "index": "projects_v2"
  }
}'
curl -X POST http://elasticsearch:9200/_reindex -H 'Content-Type: application/json' -d '{
  "source": {
    "index": "teams"
  },
  "dest": {
    "index": "teams_v2"
  }
}'

curl -X DELETE http://elasticsearch:9200/users
curl -X DELETE http://elasticsearch:9200/teams
curl -X DELETE http://elasticsearch:9200/projects

curl -X PUT http://elasticsearch:9200/users -H 'Content-Type: application/json' -d @create_index_users.json
curl -X PUT http://elasticsearch:9200/teams -H 'Content-Type: application/json' -d @create_index_teams.json
curl -X PUT http://elasticsearch:9200/projects -H 'Content-Type: application/json' -d @create_index_projects.json

curl -X POST http://elasticsearch:9200/_reindex -H 'Content-Type: application/json' -d '{
  "source": {
    "index": "users_v2"
  },
  "dest": {
    "index": "users"
  }
}'
curl -X POST http://elasticsearch:9200/_reindex -H 'Content-Type: application/json' -d '{
  "source": {
    "index": "projects_v2"
  },
  "dest": {
    "index": "projects"
  }
}'
curl -X POST http://elasticsearch:9200/_reindex -H 'Content-Type: application/json' -d '{
  "source": {
    "index": "teams_v2"
  },
  "dest": {
    "index": "teams"
  }
}'

curl -X DELETE http://elasticsearch:9200/users_v2
curl -X DELETE http://elasticsearch:9200/teams_v2
curl -X DELETE http://elasticsearch:9200/projects_v2


echo "done"