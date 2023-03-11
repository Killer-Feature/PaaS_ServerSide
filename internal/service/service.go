package service

import (
	"KillerFeature/ServerSide/internal"
)

type Service struct {
}

func NewService() internal.Usecase {
	return &Service{}
}
