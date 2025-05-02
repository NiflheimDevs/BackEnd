#!/bin/bash

if [ "$#" -eq 0 ]; then
    echo "Usage: $0 <entity_name>"
    exit 1
fi

lower_entity_name="${1,,}"

upper_entity_name="${lower_entity_name^}"

echo "first letter upper: $upper_entity_name"

echo "first letter lower: $lower_entity_name"

rm ./internal/domain/models/${lower_entity_name}_model.go

rm ./internal/application/services/impl/${lower_entity_name}_service_impl.go

rm ./internal/application/services/${lower_entity_name}_service.go

rm ./internal/infrastructure/repositories/postgres/${lower_entity_name}_repo.go

rm ./internal/domain/repositories/postgres/${lower_entity_name}_repo.go

rm ./internal/delivery/handlers/${lower_entity_name}_handler.go
