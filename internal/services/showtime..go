package services

import "context"

type IShowTime interface {
	AddShowTime(ctx context.Context) (int, error)
	DeleteShowTime(ctx context.Context) (int, error)
	GetShowTime(ctx context.Context) (interface{}, int, error)
}
