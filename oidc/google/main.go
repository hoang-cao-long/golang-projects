package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	clientID     = "1096864804548-6t3uujnbi7j7g8vn0o6rdamcf4ngjvp1.apps.googleusercontent.com"
	clientSecret = "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"
	redirectURL  = "http://localhost:8080/callback"
	state        = "random-state-string" // for CSRF protection
)

var oauth2Config = oauth2.Config{
	ClientID:     clientID,
	ClientSecret: clientSecret,
	RedirectURL:  redirectURL,
	Endpoint:     google.Endpoint,
	Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
}

func main() {
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleCallback)
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Step 1: Redirect the user to the Google login page
func handleLogin(w http.ResponseWriter, r *http.Request) {
	url := oauth2Config.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

// Step 2: Handle the callback from Google and exchange the authorization code for a token
func handleCallback(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("state") != state {
		http.Error(w, "State is invalid", http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch the ID token and verify it
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No ID token in the token response", http.StatusInternalServerError)
		return
	}

	provider, err := oidc.NewProvider(context.Background(), "https://accounts.google.com")
	if err != nil {
		http.Error(w, "Failed to create OIDC provider: "+err.Error(), http.StatusInternalServerError)
		return
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: clientID})
	idToken, err := verifier.Verify(context.Background(), rawIDToken)
	if err != nil {
		http.Error(w, "Failed to verify ID token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse the ID token claims to extract user information
	var claims struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
	}
	if err := idToken.Claims(&claims); err != nil {
		http.Error(w, "Failed to parse ID token claims: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Login successful!\n")
	fmt.Fprintf(w, "Email: %s\n", claims.Email)
	fmt.Fprintf(w, "Name: %s\n", claims.Name)
	fmt.Fprintf(w, "Picture: %s\n", claims.Picture)
}
