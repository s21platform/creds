package service

import (
	"context"

	"github.com/s21platform/creds/internal/model"
)

type DbRepo interface {
	GetToken(ctx context.Context, tokenUUID string) (string, error)
	GetCreds(ctx context.Context, service string) (model.CredsData, error)
}
