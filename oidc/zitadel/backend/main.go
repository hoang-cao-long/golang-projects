package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	"github.com/zitadel/zitadel-go/v3/pkg/http/middleware"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
	"golang.org/x/exp/slog"
)

var (
	// flags to be provided for running the example server
	domain = flag.String("domain", "localhost", "your ZITADEL instance domain (in the form: <instance>.zitadel.cloud or <yourdomain>)")
	key    = flag.String("key", "292578997605236738.json", "path to your key.json")
	port   = flag.String("port", "8091", "port to run the server on (default is 8089)")

	// tasks are used to store an in-memory list used in the protected endpoint
	tasks []string
)

func main() {
	flag.Parse()

	ctx := context.Background()

	// init authorize
	authZ, err := authorization.New(ctx, zitadel.New(*domain), oauth.DefaultAuthorization(*key))
	if err != nil {
		slog.Error("zitadel sdk could not initialize", "error", err)
		os.Exit(1)
	}

	// init middleware
	mw := middleware.New(authZ)

	router := http.NewServeMux()

	router.Handle("/api/healthz", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			err = jsonResponse(w, "OK", http.StatusOK)
			if err != nil {
				slog.Error("error writing response", "error", err)
			}
		}))

	router.Handle("GET /api/tasks", mw.RequireAuthorization()(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			authCtx := mw.Context(r.Context())
			slog.Info("user accessed task list", "id", authCtx.UserID(), "username", authCtx.Username)

			list := tasks
			if authCtx.IsGrantedRole("admin") {
				list = append(list, "create a new task on /api/add-task")
			}

			err = jsonResponse(w, &taskList{Tasks: list}, http.StatusOK)
			if err != nil {
				slog.Error("error writing response", "error", err)
			}
		})))

	// router.Handle("POST /api/add-task", mw.RequireAuthorization(authorization.WithRole(`admin`))(http.HandlerFunc(
	router.Handle("POST /api/add-task", mw.RequireAuthorization()(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			task := strings.TrimSpace(r.FormValue("task"))
			if task == "" {
				err = jsonResponse(w, "task must not be empty", http.StatusBadRequest)
				if err != nil {
					slog.Error("error writing invalid task response", "error", err)
					return
				}
				return
			}

			tasks = append(tasks, task)

			slog.Info("admin added task", "id", authorization.UserID(r.Context()), "task", task)

			err = jsonResponse(w, fmt.Sprintf("task `%s` added", task), http.StatusOK)
			if err != nil {
				slog.Error("error writing task added response", "error", err)
				return
			}
		})))

	lis := fmt.Sprintf(":%s", *port)
	slog.Info("server listening, press ctrl+c to stop", "addr", "http://localhost"+lis)
	err = http.ListenAndServe(lis, router)
	if !errors.Is(err, http.ErrServerClosed) {
		slog.Error("server terminated", "error", err)
		os.Exit(1)
	}
}

func jsonResponse(w http.ResponseWriter, resp any, status int) error {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	data, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

type taskList struct {
	Tasks []string `json:"tasks,omitempty"`
}
