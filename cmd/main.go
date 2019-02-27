package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ctco-dev/go-api-template/internal/app"
	"github.com/ctco-dev/go-api-template/internal/log"

	"github.com/kelseyhightower/envconfig"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

func main() {

	rootCtx := context.Background()

	var env app.Specification
	err := envconfig.Process("", &env)
	if err != nil {
		log.WithCtx(rootCtx).Panicf("env vars error: '%v'", err)
	}

	someApp := app.New(env)
	addr := fmt.Sprintf(":%d", env.Port)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := uuid.NewV4().String()[0:8]
		reqCtx := log.NewContext(rootCtx, logrus.Fields{"reqID": reqID})
		reqCtx, cancel := context.WithTimeout(reqCtx, time.Second*10)
		defer cancel()
		someApp.ServeHTTP(w, r.WithContext(reqCtx))
	})

	log.WithCtx(rootCtx).Infof("Server is running at: http://localhost:%d", env.Port)
	log.WithCtx(rootCtx).Fatal(http.ListenAndServe(addr, handler))
}
