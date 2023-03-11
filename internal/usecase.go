package internal

import (
	"KillerFeature/ServerSide/internal/models"
)

type Usecase interface {
	DeployApp(creds *models.SshCreds) error
}
