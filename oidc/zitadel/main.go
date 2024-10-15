package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

// Variables for your Zitadel client credentials and other settings
var (
	clientID    = "289375988414941237"                                                                // Replace with your Zitadel Client ID
	issuer      = "https://hoang-cao-long-instance-aki4x4.us1.zitadel.cloud"                          // Replace with the correct Zitadel issuer
	redirectURL = "https://hoang-cao-long-instance-aki4x4.us1.zitadel.cloud/ui/console/auth/callback" // Replace with your redirect URI
)

var (
	oauth2Config *oauth2.Config
	verifier     *oidc.IDTokenVerifier
)

func main() {
	// Set up the OIDC provider
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		log.Fatal(err)
	}

	// Set up OAuth2 configuration with PKCE
	oauth2Config = &oauth2.Config{
		ClientID:    clientID,
		RedirectURL: redirectURL,
		Endpoint:    provider.Endpoint(),
		Scopes:      []string{oidc.ScopeOpenID, "profile", "email"},
	}

	// Create an ID token verifier to verify tokens received from Zitadel
	verifier = provider.Verifier(&oidc.Config{ClientID: clientID})

	// Set up HTTP routes
	http.HandleFunc("/", handleHome)
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/callback", handleCallback)

	// Start the server
	log.Println("Listening on http://localhost:8080/")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Handle the home page with a login link
func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, `<html><body><a href="/login">Log in with Zitadel</a></body></html>`)
}

// Handle the login request, initiate the PKCE flow
func handleLogin(w http.ResponseWriter, r *http.Request) {
	codeVerifier, codeChallenge := generateCodeVerifierAndChallenge()
	state := "some-random-state"

	// Generate the authorization URL with PKCE code challenge
	authURL := oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline,
		oauth2.SetAuthURLParam("code_challenge", codeChallenge),
		oauth2.SetAuthURLParam("code_challenge_method", "S256"),
	)

	// Store the code verifier in a cookie or session
	http.SetCookie(w, &http.Cookie{Name: "code_verifier", Value: codeVerifier, Path: "/"})

	// Redirect the user to the Zitadel login page
	http.Redirect(w, r, authURL, http.StatusFound)
}

// Handle the callback from Zitadel
func handleCallback(w http.ResponseWriter, r *http.Request) {
	// Verify the state to prevent CSRF (in a real app, you should implement state verification)
	if r.URL.Query().Get("state") != "some-random-state" {
		http.Error(w, "State did not match", http.StatusBadRequest)
		return
	}

	// Retrieve the code from the query string
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	// Get the code verifier from the cookie
	codeVerifierCookie, err := r.Cookie("code_verifier")
	if err != nil {
		http.Error(w, "Code verifier not found", http.StatusBadRequest)
		return
	}
	codeVerifier := codeVerifierCookie.Value

	// Exchange the authorization code for tokens
	ctx := context.Background()
	token, err := oauth2Config.Exchange(ctx, code,
		oauth2.SetAuthURLParam("code_verifier", codeVerifier),
	)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract the ID token
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "No ID token found", http.StatusInternalServerError)
		return
	}

	// Verify the ID token
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		http.Error(w, "Failed to verify ID token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract the user claims from the ID token
	var claims struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
	}
	if err := idToken.Claims(&claims); err != nil {
		http.Error(w, "Failed to extract claims: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Display the user's email
	fmt.Fprintf(w, "User logged in! Email: %s, Verified: %t", claims.Email, claims.EmailVerified)
}

// Generate PKCE code verifier and challenge
func generateCodeVerifierAndChallenge() (string, string) {
	// Generate a random code verifier
	codeVerifier := make([]byte, 32)
	_, err := rand.Read(codeVerifier)
	if err != nil {
		log.Fatal(err)
	}

	// Encode the code verifier to base64
	encodedVerifier := base64.RawURLEncoding.EncodeToString(codeVerifier)

	// Create a SHA-256 hash of the verifier as the code challenge
	challenge := sha256Base64URLEncode(encodedVerifier)

	return encodedVerifier, challenge
}

// Helper function to encode a string using SHA-256 and Base64 URL encoding
func sha256Base64URLEncode(s string) string {
	hash := sha256.New()
	hash.Write([]byte(s))
	return base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	// Clear session cookies or other session data
	http.SetCookie(w, &http.Cookie{
		Name:   "id_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

	// Redirect to Zitadel's logout endpoint
	http.Redirect(w, r, getLogoutURL(), http.StatusFound)
}

func getLogoutURL() string {
	// Here you should retrieve the ID token from your session store or cookies
	idToken := "ID_TOKEN_STORED_IN_SESSION" // For example, this could be from a cookie or session storage

	// Define where to redirect the user after logout
	postLogoutRedirectURL := "http://localhost:8080" // Replace with your app’s URL

	// Zitadel logout endpoint
	logoutEndpoint := "https://issuer.zitadel.ch/oidc/v1/logout"

	// Create the logout URL
	logoutURL := fmt.Sprintf("%s?id_token_hint=%s&post_logout_redirect_uri=%s", logoutEndpoint, idToken, postLogoutRedirectURL)
	return logoutURL
}
