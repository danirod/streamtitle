// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
	"net/http"

	gonanoid "github.com/matoous/go-nanoid/v2"
)

// An internal data structure just to pass two variables between channels.
// There are other solutions in Go to pass multiple variables through the
// same channel, but I don't like anonymous structs.
type loginTokens struct {
	accessToken, refreshToken string
}

// A tricky function that has to spawn a full HTTP server so that an OAuth
// flow can be made from the command line. This goroutine will block the
// caller, print a Twitch authorization URL to the stdout, and wait until
// the internal web server receives the redirection from the Twitch OAuth.
func spawnAuthorizationServer(client *Client, done chan loginTokens) {
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
