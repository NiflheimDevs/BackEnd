package servicesimpl

import repositories "github.com/niflheimdevs/backend/internal/domain/repositories/postgres"

type RoleService struct {
    RoleRepo repositories.RoleRepo
}

func NewRoleService (
    roleRepo repositories.RoleRepo,
    ) *RoleService {
    return &RoleService{
        RoleRepo: roleRepo,
    }
} 
