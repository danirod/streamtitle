// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"errors"

	"github.com/nicklaw5/helix/v2"
)

// Client is the data structure that groups all the program state.
type Client struct {
	// A facade for the following collaborators.
	client     *helix.Client // The HTTP client used to interact with Twitch.
	config     staticConfig  // The static config information in the environment.
	streamInfo streamInfo    // The stream information object.

	// Some transient information that is only valid for this session.
	accessToken string // The access token received during bot startup.
	broadcastId string // The broadcast ID, which depends on the login.
}

// NewContext builds a new context and returns the outcome.
func NewContext() (*Client, error) {
	var ctx Client

	if err := ctx.config.read(); err != nil {
		return nil, err
	}

	client, err := helix.NewClient(&helix.Options{
		ClientID:     ctx.config.clientId,
		ClientSecret: ctx.config.clientSecret,
		RedirectURI:  "http://localhost:9300/st-callback",
	})
	if err != nil {
		return nil, err
	}

	ctx.client = client
	return &ctx, nil
}

// Login will initialise the client state so that there exists a valid access
// token in the client context which can be used by further API calls made
// within the client.
func (ctx *Client) Login() error {
	logger.Print("Prepare for login")

	callback := make(chan string)
	go createTokenProvider(ctx, callback)
	ctx.accessToken = <-callback
	ctx.client.SetUserAccessToken(ctx.accessToken)
	logger.Print("Access token has been set")

	bid, err := ctx.getTokenOwner()
	if err != nil {
		return err
	}
	ctx.broadcastId = bid
	logger.Print("Broadcast ID set to ", bid)
	return nil
}

// Uses the validation endpoint to check the token information, and as
// a side effect it serves as a validation that the token is OK.
func (ctx *Client) getTokenOwner() (string, error) {
	logger.Print("Validating access token...")
	valid, resp, err := ctx.client.ValidateToken(ctx.accessToken)
	if err != nil {
		return "", err
	}
	if !valid {
		return "", errors.New("The token is not valid")
	}
	logger.Print("Access token belongs to ", resp.Data.Login)
	return resp.Data.UserID, nil
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
