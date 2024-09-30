package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

var (
	clientID     = os.Getenv("CLIENT_ID")
	clientSecret = os.Getenv("CLIENT_SECRET")
	redirectURL  = "http://localhost:8080/callback"
	// Google OIDC endpoints
	providerURL = "https://accounts.google.com"
	// OAuth2 config
	oauth2Config *oauth2.Config
	oidcProvider *oidc.Provider
	verifier     *oidc.IDTokenVerifier
)

func main() {
	// Set up the OIDC provider and OAuth2 config
	ctx := context.Background()
	var err error

	oidcProvider, err = oidc.NewProvider(ctx, providerURL)
	if err != nil {
		log.Fatalf("Failed to get provider: %v", err)
	}

	oauth2Config = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     oidcProvider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "offline_access"}, // Add offline_access for refresh tokens
	}

	verifier = oidcProvider.Verifier(&oidc.Config{ClientID: clientID})

	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleCallback)

	log.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	html := `<html><body><a href="/login">Login with OIDC</a></body></html>`
	fmt.Fprint(w, html)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, oauth2Config.AuthCodeURL("state"), http.StatusFound)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("state") != "state" {
		http.Error(w, "state did not match", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No code in callback", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Save refresh token
	refreshToken := token.RefreshToken
	fmt.Printf("Refresh Token: %s\n", refreshToken)

	// Print token expiry
	fmt.Printf("Access Token Expiry: %s\n", token.Expiry)

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token field in OAuth2 token", http.StatusInternalServerError)
		return
	}

	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var claims struct {
		Email string `json:"email"`
	}

	if err := idToken.Claims(&claims); err != nil {
		http.Error(w, "Failed to parse ID token claims: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Display logged-in user's email and refresh the access token if needed
	fmt.Fprintf(w, "Hello, %s!\n\n", claims.Email)
	fmt.Fprintf(w, "Access Token Expiry: %v\n\n", token.Expiry)
	fmt.Fprintf(w, "Refresh Token: %v\n", refreshToken)

	// Refresh the token if necessary
	if time.Now().After(token.Expiry) {
		fmt.Fprint(w, "\nRefreshing token...\n")
		refreshedToken, err := refreshAccessToken(ctx, token)
		if err != nil {
			http.Error(w, "Failed to refresh token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "New Access Token Expiry: %v\n", refreshedToken.Expiry)
		fmt.Fprintf(w, "New Refresh Token: %v\n", refreshedToken.RefreshToken)
	}
}

// refreshAccessToken uses the refresh token to get a new access token
func refreshAccessToken(ctx context.Context, token *oauth2.Token) (*oauth2.Token, error) {
	src := oauth2Config.TokenSource(ctx, token)
	newToken, err := src.Token() // Automatically uses refresh token if expired
	if err != nil {
		return nil, fmt.Errorf("could not refresh access token: %w", err)
	}

	return newToken, nil
}
