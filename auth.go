package main

import "fmt"

import "net/http"
import "github.com/matoous/go-nanoid/v2"

// An internal data structure just to pass two variables between channels.
// There are other solutions in Go to pass multiple variables through the
// same channel, but I don't like anonymous structs.
type loginTokens struct {
	accessToken, refreshToken string
}

// Use this goroutine to asyncally fetch an access token. The caller thread
// will be blocked until the channel yields the proper access token. In
// case of emergency, this goroutine will handle the panic.
func createTokenProvider(client *Context, done chan string) {
	if client.config.refreshToken == "" {
		// Find a fresh pair of tokens since we don't have one.
		receiver := make(chan loginTokens)
		go spawnAuthorizationServer(client, receiver)
		loginToken := <-receiver

		// Store the received refresh token and yield the access token.
		client.config.refreshToken = loginToken.refreshToken
		client.config.write()
		done <- loginToken.accessToken
	} else {
		// Refresh the token to fetch a fresh access token.
		tokens, err := useRefreshToken(client)
		if err != nil {
			// Something happened with the refresh token, start over.
			client.config.refreshToken = ""
			client.config.write()
			createTokenProvider(client, done)
		} else {
			// The refresh token was valid and we have an access token.
			client.config.refreshToken = tokens.refreshToken
			client.config.write()
			done <- tokens.accessToken
		}
	}
}

// Use the refresh token stored in the client context to get a new pair of
// access and refresh token that we can use for this session.
func useRefreshToken(client *Context) (*loginTokens, error) {
	resp, err := client.client.RefreshUserAccessToken(client.config.refreshToken)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		// Treat this as invalid token.
		return nil, fmt.Errorf(resp.Error)
	}
	return &loginTokens{
		accessToken:  resp.Data.AccessToken,
		refreshToken: resp.Data.RefreshToken,
	}, nil
}

// A tricky function that has to spawn a full HTTP server so that an OAuth
// flow can be made from the command line. This goroutine will block the
// caller, print a Twitch authorization URL to the stdout, and wait until
// the internal web server receives the redirection from the Twitch OAuth.
func spawnAuthorizationServer(client *Context, done chan loginTokens) {
	// State parameter to be used in the OAuth flow.
	state, _ := gonanoid.New()

	http.HandleFunc("/st-callback", func(w http.ResponseWriter, h *http.Request) {
		// Validate the state parameter (why do I even bother?)
		if h.URL.Query().Get("state") != state {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Add("Content-Type", "text/plain")
			fmt.Fprintln(w, "Invalid state parameter, is everything OK?")
			panic("Invalid state parameter, is everything OK?")
		}

		// Send affirmative message to the browser.
		w.Header().Add("Content-Type", "text/plain")
		fmt.Fprintln(w, "You should be able to close this window now!")

		// Handle the token and send the credentials to the caller.
		code := h.URL.Query().Get("code")
		resp, err := client.client.RequestUserAccessToken(code)
		if err != nil {
			panic(err)
		}
		done <- loginTokens{
			accessToken:  resp.Data.AccessToken,
			refreshToken: resp.Data.RefreshToken,
		}
	})

	// Present the URL via stdout.
	url := client.AuthorizationURL(state)
	fmt.Println("Please visit", url, "to continue.")

	// Start the HTTP server to handle the OAuth callback.
	http.ListenAndServe(":9300", nil)
}
