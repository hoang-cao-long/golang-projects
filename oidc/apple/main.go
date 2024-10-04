package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

// Replace these with your Apple credentials
const (
	clientID = "hclong2k@gmail.com" // The service identifier (App ID)
	// teamID      = "YOUR_TEAM_ID"       // Your Apple Developer Team ID
	redirectURL = "https://your-app.com/callback"
	keyID       = "YOUR_KEY_ID" // Your Apple Services Key ID
	privateKey  = `-----BEGIN PRIVATE KEY-----
YOUR_PRIVATE_KEY_CONTENT_HERE
-----END PRIVATE KEY-----`
)

// Set up the OAuth2 config
var oauth2Config = oauth2.Config{
	ClientID:    clientID,
	RedirectURL: redirectURL,
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://appleid.apple.com/auth/authorize",
		TokenURL: "https://appleid.apple.com/auth/token",
	},
	Scopes: []string{"openid", "email", "name"},
}

func main() {
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleCallback)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Step 1: Redirect user to Apple Sign-In page
func handleLogin(w http.ResponseWriter, r *http.Request) {
	state := "random-state-string"
	url := oauth2Config.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

// Step 2: Handle the callback and exchange code for tokens
func handleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// Verify state for CSRF protection
	if r.URL.Query().Get("state") != "random-state-string" {
		http.Error(w, "State mismatch", http.StatusBadRequest)
		return
	}

	// Exchange the authorization code for tokens
	code := r.URL.Query().Get("code")
	token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Verify ID token and extract claims
	idToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No id_token in the token response", http.StatusInternalServerError)
		return
	}

	provider, err := oidc.NewProvider(ctx, "https://appleid.apple.com")
	if err != nil {
		http.Error(w, "Failed to create provider: "+err.Error(), http.StatusInternalServerError)
		return
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: clientID})
	verifiedIDToken, err := verifier.Verify(ctx, idToken)
	if err != nil {
		http.Error(w, "Failed to verify ID token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var claims map[string]interface{}
	if err := verifiedIDToken.Claims(&claims); err != nil {
		http.Error(w, "Failed to parse claims: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Login successful! User info: %+v", claims)
}
