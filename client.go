// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"errors"

	"github.com/nicklaw5/helix/v2"
)

// Client is the data structure that groups all the program state.
type Client struct {
	client *helix.Client // The HTTP client used to interact with Twitch.

	config      appConfig // The static config information in the environment.
	state       appState  // The state (credentials and such)
	broadcastId string    // The broadcast ID for the current token

	streamInfo streamInfo // The stream information object.
}

// NewContext builds a new context and returns the outcome.
func NewContext() (*Client, error) {
	var ctx Client

	if err := ctx.config.read(); err != nil {
		return nil, err
	}
	if err := ctx.state.read(); err != nil {
		return nil, err
	}

	client, err := helix.NewClient(&helix.Options{
		ClientID:     ctx.config.ClientId,
		ClientSecret: ctx.config.ClientSecret,
		RedirectURI:  "http://localhost:9300/st-callback",
	})
	if err != nil {
		return nil, err
	}

	ctx.client = client
	return &ctx, nil
}

// refreshToken tries to use the given refresh token to get a new pair of
// access and refresh token. If the token cannot be refreshed because the
// token is not valid, a separate flag will be raised and the strings will
// be empty. If there is an issue during the token refresh, it will return
// the corresponding error.
func (ctx *Client) refreshToken(token string) (bool, string, string, error) {
	logger.Print("Using the refresh token...")

	if resp, err := ctx.client.RefreshUserAccessToken(token); err != nil {
		return false, "", "", err
	} else {
		if resp.Error == "" {
			return true, resp.Data.AccessToken, resp.Data.RefreshToken, nil
		} else {
			return false, "", "", nil
		}
	}
}

// Tries to finalize the login by checking the validity of the token. If the token
// is valid, it will also set the owner of the token as the broadcaster ID field
// in the local client structure.
func (ctx *Client) finishLogin(token string) (bool, error) {
	logger.Print("Validating access token...")
	if valid, response, err := ctx.client.ValidateToken(token); err != nil {
		return false, err
	} else if valid {
		logger.Print("Token is valid and belongs to ", response.Data.Login)
		ctx.broadcastId = response.Data.UserID
		ctx.client.SetUserAccessToken(token)
		return true, nil
	} else {
		logger.Print("Access token is not valid (maybe expired?")
		return false, nil
	}
}

// Login will initialise the client state so that there exists a valid access
// token in the client context which can be used by further API calls made
// within the client.
func (ctx *Client) Login() error {
	logger.Print("Prepare for login")

	// Check if the token in the app state is still valid
	if ctx.state.AccessToken != "" {
		if valid, err := ctx.finishLogin(ctx.state.AccessToken); err != nil {
			return err
		} else if valid {
			return nil
		}

		// So the access token is not valid. Is the refresh token valid?
		if ctx.state.RefreshToken != "" {
			if valid, access, refresh, err := ctx.refreshToken(ctx.state.RefreshToken); err != nil {
				return err
			} else if valid {
				logger.Print("A new token has been generated")
				ctx.state.AccessToken = access
				ctx.state.RefreshToken = refresh
				if err := ctx.state.write(); err != nil {
					return err
				}
				if valid, err := ctx.finishLogin(access); err != nil {
					return err
				} else if valid {
					return nil
				}
			} else {
				logger.Print("Refresh token is not valid")
			}
		}
	} else {
		logger.Print("No access token found in the application state")
	}

	// The user has to log in
	callback := make(chan loginTokens)
	go spawnAuthorizationServer(ctx, callback)
	tokens := <-callback
	ctx.state.AccessToken = tokens.accessToken
	ctx.state.RefreshToken = tokens.refreshToken
	if err := ctx.state.write(); err != nil {
		return err
	}
	if valid, err := ctx.finishLogin(ctx.state.AccessToken); err != nil {
		return err
	} else if valid {
		return nil
	} else {
		// I am out of ideas, man!
		return errors.New("cannot issue a valid access token")
	}
}

// AuthorizationURL will build the URL that the user has to visit in order
// to authorize the bot when the refresh token does not exist or it is
// either invalid or expired.
func (ctx *Client) AuthorizationURL(state string) string {
	return ctx.client.GetAuthorizationURL(&helix.AuthorizationURLParams{
		ResponseType: "code",
		Scopes:       []string{"channel:manage:broadcast"},
		State:        state,
		ForceVerify:  false,
	})
}

// FetchStreamInfo will download the current stream information via the API
// and it will store the information in the stream information part of the
// client context.
func (ctx *Client) FetchStreamInfo() error {
	logger.Print("Fetching stream information...")
	resp, err := ctx.client.GetChannelInformation(&helix.GetChannelInformationParams{
		BroadcasterIDs: []string{ctx.broadcastId},
	})
	if err == nil {
		ctx.streamInfo.load(&resp.Data.Channels[0])
	}
	return err
}

// SendStreamInfo generates an API call that will submit the currently stored
// stream information in the context. If this was previously set by calls
// to the setters, it will update the stream information.
func (ctx *Client) SendStreamInfo() error {
	logger.Print("Updating stream information...")
	logger.Print("New title: ", ctx.streamInfo.title)
	logger.Print("New game ID: ", ctx.streamInfo.game)
	logger.Print("New language: ", ctx.streamInfo.language)
	logger.Print("Tag list: ", ctx.streamInfo.tags)

	_, err := ctx.client.EditChannelInformation(&helix.EditChannelInformationParams{
		BroadcasterID:       ctx.broadcastId,
		Title:               ctx.streamInfo.title,
		GameID:              ctx.streamInfo.game,
		BroadcasterLanguage: ctx.streamInfo.language,
		Tags:                ctx.streamInfo.tags,
	})
	return err
}
