#!/bin/bash

echo "note: the directories are for MY directories so..."

if [ "$#" -eq 0 ]; then
    echo "Usage: $0 <entity_name>"
    exit 1
fi

ORIG_DIR="$PWD"

lower_entity_name="${1,,}"

upper_entity_name="${lower_entity_name^}"

echo "first letter upper: $upper_entity_name"

echo "first letter lower: $lower_entity_name"

# cd "${ORIG_DIR}"

cat <<EOF >> ./internal/domain/models/${lower_entity_name}_model.go
package models

type ${upper_entity_name}Model struct {}
EOF

# cd "${ORIG_DIR}"

cat <<EOF >> ./internal/application/services/impl/${lower_entity_name}_service_impl.go
package servicesimpl

import repositories "github.com/niflheimdevs/backend/internal/domain/repositories/postgres"

type ${upper_entity_name}Service struct {
    ${upper_entity_name}Repo repositories.${upper_entity_name}Repo
}

func New${upper_entity_name}Service (
    ${lower_entity_name}Repo repositories.${upper_entity_name}Repo,
    ) *${upper_entity_name}Service {
    return &${upper_entity_name}Service{
        ${upper_entity_name}Repo: ${lower_entity_name}Repo,
    }
} 
EOF

# cd "${ORIG_DIR}"

cat <<EOF >> ./internal/application/services/${lower_entity_name}_service.go
package services

type ${upper_entity_name}Service interface {
}

EOF

# cd "${ORIG_DIR}"

cat <<EOF >> ./internal/infrastructure/repositories/postgres/${lower_entity_name}_repo.go
package repositoriesimpl

import "github.com/jackc/pgx/v5/pgxpool"

type ${upper_entity_name}Repo struct {
    PG *pgxpool.Pool
}

func New${upper_entity_name}Repo (
    PG *pgxpool.Pool,
    ) *${upper_entity_name}Repo {
    return &${upper_entity_name}Repo{
        PG: PG,
    }
}
EOF

cat <<EOF >> ./internal/domain/repositories/postgres/${lower_entity_name}_repo.go
package repositories

type ${upper_entity_name}Repo interface {
}

EOF

cat <<EOF >> ./internal/delivery/handlers/${lower_entity_name}_handler.go
package handlers

import (
	"github.com/go-playground/validator/v10"
    "github.com/niflheimdevs/backend/bootstrap"
	"github.com/niflheimdevs/backend/internal/application/services"
)

type ${upper_entity_name}Handler struct {
	${upper_entity_name}Service services.${upper_entity_name}Service
	Constants      *bootstrap.Constants
	Validator      *validator.Validate
}

func New${upper_entity_name}Handler(
	${lower_entity_name}Service services.${upper_entity_name}Service,
	constants *bootstrap.Constants,
	validator *validator.Validate,
) *${upper_entity_name}Handler {
	return &${upper_entity_name}Handler{
		${upper_entity_name}Service: ${lower_entity_name}Service,
		Constants:      constants,
		Validator:      validator,
	}
}

EOF


echo "DON'T FORGET ABOUT WIRE! I CAN'T BE BOTHERED NO MORE YOU LAZY FUCK!"

echo "\nif you are unhappy with the changes..."
echo "know that i do not give one single fuck"
echo "good news though, i know who does!"
echo "its name is slay_files.sh"