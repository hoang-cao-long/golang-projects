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

	"github.com/zitadel/oidc/v3/pkg/oidc"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization"
	"github.com/zitadel/zitadel-go/v3/pkg/authorization/oauth"
	"github.com/zitadel/zitadel-go/v3/pkg/client"
	objectV2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/object/v2"
	userV2 "github.com/zitadel/zitadel-go/v3/pkg/client/zitadel/user/v2"
	"github.com/zitadel/zitadel-go/v3/pkg/http/middleware"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
	"golang.org/x/exp/slog"
)

var (
	// flags to be provided for running the example server
	domain = flag.String("domain", "localhost", "your ZITADEL instance domain (in the form: <instance>.zitadel.cloud or <yourdomain>)")
	key    = flag.String("key", "292578997605236738.json", "path to your key.json")
	port   = flag.String("port", "8091", "port to run the server on (default is 8089)")

	keyUserService = flag.String("keyUserService", "295610152923430914.json", "path to your key.json")

	// tasks are used to store an in-memory list used in the protected endpoint
	tasks  []string
	userID = "user123"
)

func main() {
	flag.Parse()

	ctx := context.Background()

	// authz
	authZ, err := authorization.New(ctx, zitadel.New(*domain), oauth.DefaultAuthorization(*key))
	if err != nil {
		slog.Error("zitadel sdk could not initialize", "error", err)
		os.Exit(1)
	}

	// client
	api, err := client.New(ctx, zitadel.New(*domain),
		client.WithAuth(client.DefaultServiceUserAuthentication(*keyUserService, oidc.ScopeOpenID, client.ScopeZitadelAPI())),
	)
	if err != nil {
		slog.Error("Failed to create Zitadel client: %v", err)
	}

	// init middleware
	mw := middleware.New(authZ)

	router := http.NewServeMux()

	router.Handle("GET /api/healthz", http.HandlerFunc(
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

	router.Handle("POST /api/user/create", mw.RequireAuthorization()(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userName := "testUserName1"

			// create user
			userRespCreated, err := api.UserServiceV2().AddHumanUser(ctx, &userV2.AddHumanUserRequest{
				UserId:       &userID,
				Username:     &userName,
				Organization: &objectV2.Organization{},
				Profile: &userV2.SetHumanProfile{
					GivenName:         "testGivenName1",
					FamilyName:        "testFamilyName1",
					NickName:          new(string),
					DisplayName:       new(string),
					PreferredLanguage: new(string),
					Gender:            userV2.Gender_GENDER_MALE.Enum(),
				},
				Email: &userV2.SetHumanEmail{
					Email:        "hclong2k@gmail.com",
					Verification: nil,
				},
				// PasswordType: &userV2.AddHumanUserRequest_Password{
				// 	Password: &userV2.Password{
				// 		Password:       "Zitadel@123",
				// 		ChangeRequired: true,
				// 	},
				// },
			})
			if err != nil {
				slog.Error("Failed to create Zitadel client: %v", err)
			}
			slog.Info("User created successfully: %v", userRespCreated)

			// set password
			// _, err = api.UserServiceV2().SetPassword(ctx, &userV2.SetPasswordRequest{
			// 	UserId: userResp.UserId,
			// 	NewPassword: &userV2.Password{
			// 		Password:       "Zitadel@123",
			// 		ChangeRequired: false,
			// 	},
			// 	Verification: nil,
			// })
			// if err != nil {
			// 	slog.Error("Failed to create Zitadel client: %v", err)
			// }
		})))

	router.Handle("POST /api/user/update", mw.RequireAuthorization()(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userName := "testUserName2"

			// update user
			userRespUpdated, err := api.UserServiceV2().UpdateHumanUser(ctx, &userV2.UpdateHumanUserRequest{
				UserId:   userID,
				Username: &userName,
				Profile: &userV2.SetHumanProfile{
					GivenName:  "testGivenName2",
					FamilyName: "testFamilyName2",
					NickName:   new(string),
					// DisplayName:       new(string),
					PreferredLanguage: new(string),
					Gender:            userV2.Gender_GENDER_DIVERSE.Enum(),
				},
				Email: &userV2.SetHumanEmail{
					Email:        "hclong2k@gmail.com",
					Verification: nil,
				},
				// Phone: &userV2.SetHumanPhone{
				// 	Phone:        "",
				// 	Verification: nil,
				// },
				Password: &userV2.SetPassword{
					PasswordType: &userV2.SetPassword_Password{
						Password: &userV2.Password{
							Password:       "Zitadel@321",
							ChangeRequired: false,
						},
					},
					Verification: nil,
				},
			})
			if err != nil {
				slog.Error("Failed to update Zitadel client: %v", err)
			}
			slog.Info("User updated successfully: %v", userRespUpdated)
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
