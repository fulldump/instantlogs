package api

import (
	"context"

	"github.com/fulldump/instantlogs/service"
)

func getService(ctx context.Context) *service.Service {

	value := ctx.Value("instantlogs-service")

	return value.(*service.Service) // this will raise a panic on error!!
}
