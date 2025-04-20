package servicesimpl

import repositories "github.com/niflheimdevs/backend/internal/domain/repositories/postgres"

type TeamService struct {
    TeamRepo repositories.TeamRepo
}

func NewteamService (
    teamRepo repositories.TeamRepo,
    ) *TeamService {
    return &TeamService{
        TeamRepo: teamRepo,
    }
} 
