package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fulldump/box"

	"github.com/fulldump/instantlogs/service"
	"github.com/fulldump/instantlogs/statics"
)

func NewApi(service *service.Service, staticsDir string) *box.B {

	b := box.NewBox()

	b.WithInterceptors(func(next box.H) box.H {
		return func(ctx context.Context) {
			ctx = context.WithValue(ctx, "instantlogs-service", service)
			next(ctx)
		}
	})

	b.Resource("/ingest").WithActions(
		box.Post(ingest),
	)

	b.Resource("/filter").WithActions(
		box.Get(filter),
	)

	b.Resource("/stats").WithActions(
		box.Get(stats),
	)

	// Mount statics
	b.Resource("/*").WithActions(
		box.Get(statics.ServeStatics(staticsDir)),
	)

	return b
}

func ingest(w http.ResponseWriter, r *http.Request, ctx context.Context) interface{} {

	n, err := getService(ctx).Ingest(r.Body)

	response := map[string]interface{}{
		"n": n,
	}
	if err != nil {
		response["error"] = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}

	return response
}

func filter(w http.ResponseWriter, r *http.Request, ctx context.Context) {

	follow := len(r.URL.Query()["follow"]) > 0

	go func() {
		<-r.Context().Done()
		follow = false
	}()

	err := getService(ctx).Filter(
		w,
		r.URL.Query()["regex"],
		&follow,
	)
	if err != nil {
		fmt.Println("ERROR:", err.Error())
	}
}

func stats(ctx context.Context) interface{} {
	return getService(ctx).Stats()
}
