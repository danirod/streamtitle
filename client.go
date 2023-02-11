package main

import "github.com/nicklaw5/helix/v2"

// Context is the data structure that groups all the program state.
type Context struct {
	// A facade for the following collaborators.
	client     *helix.Client // The HTTP client used to interact with Twitch.
	config     staticConfig  // The static config information in the environment.
	streamInfo streamInfo    // The stream information object.

	// Some transient information that is only valid for this session.
	accessToken string // The access token received during bot startup.
	broadcastId string // The broadcast ID, which depends on the login.
}

// NewContext builds a new context and returns the outcome.
func NewContext() (*Context, error) {
	var ctx Context

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
func (ctx *Context) Login() error {
	callback := make(chan string)
	go createTokenProvider(ctx, callback)
	ctx.accessToken = <-callback
	ctx.client.SetUserAccessToken(ctx.accessToken)

	// Since any API method will require the broadcast ID, do it now.
	return ctx.fetchBroadcastId()
}

func (ctx *Context) fetchBroadcastId() error {
	resp, err := ctx.client.GetUsers(&helix.UsersParams{
		Logins: []string{ctx.config.channelLogin},
	})
	if err == nil {
		ctx.broadcastId = resp.Data.Users[0].ID
	}
	return err
}

// AuthorizationURL will build the URL that the user has to visit in order
// to authorize the bot when the refresh token does not exist or it is
// either invalid or expired.
func (ctx *Context) AuthorizationURL(state string) string {
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
func (ctx *Context) FetchStreamInfo() error {
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
func (ctx *Context) SendStreamInfo() error {
	_, err := ctx.client.EditChannelInformation(&helix.EditChannelInformationParams{
		BroadcasterID:       ctx.broadcastId,
		Title:               ctx.streamInfo.title,
		GameID:              ctx.streamInfo.game,
		BroadcasterLanguage: ctx.streamInfo.language,
		Tags:                ctx.streamInfo.tags,
	})
	return err
}
