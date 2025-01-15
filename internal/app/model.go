package app

import (
	"gw-authorizer/internal/storages"
	"gw-authorizer/pkg/logs"
)

type App struct {
	log       *logs.Log
	storage   storages.Storage
	cost      int
	secretKey []byte
}
