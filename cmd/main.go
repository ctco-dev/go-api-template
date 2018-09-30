package main

import (
	"context"
	"ctco-dev/go-api-template/internal/app"
	"ctco-dev/go-api-template/internal/log"
	"encoding/json"
	"github.com/kelseyhightower/envconfig"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func main() {

	rootCtx := context.Background()

	var env app.Specification
	err := envconfig.Process("", &env)
	if err != nil {
		log.WithCtx(rootCtx).Panicf("env vars error: '%v'", err)
	}

	someApp := app.NewApp(env)

	http.HandleFunc(
		"/",
		func(w http.ResponseWriter, r *http.Request) {

			reqID := uuid.NewV4().String()[0:8]
			reqCtx := log.NewContext(rootCtx, logrus.Fields{"reqID": reqID})
			reqCtx, cancel := context.WithTimeout(reqCtx, time.Second*10)
			defer cancel()

			resp, err := someApp.DoSomething(reqCtx)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.WithCtx(reqCtx).Error(err)
				w.Write([]byte("Cant't get new joke"))
				return
			}

			bytes, err := json.Marshal(resp)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Can't encode response"))
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write(bytes)

		})

	log.WithCtx(rootCtx).Info("Server is running at: http://localhost:3006")
	log.WithCtx(rootCtx).Fatal(http.ListenAndServe(":3006", nil))

}
