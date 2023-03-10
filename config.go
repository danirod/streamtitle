package main

import "github.com/adrg/xdg"
import "github.com/joho/godotenv"

const CONFIG_FILE = "StreamTitle/.StreamTitle.env"
const CLIENT_ID = "CLIENT_ID"
const CLIENT_SECRET = "CLIENT_SECRET"
const CHANNEL_LOGIN = "CHANNEL_LOGIN"
const REFRESH_TOKEN = "REFRESH_TOKEN"

type staticConfig struct {
	clientId     string // The client ID used in OAuth
	clientSecret string // the client secret used in OAuth
	channelLogin string // The login (username) to interact
	refreshToken string // The refresh token for the session
}

func configFile() string {
	file, err := xdg.ConfigFile(CONFIG_FILE)
	if err != nil {
		panic(err)
	}
	return file
}

func (cfg *staticConfig) read() error {
	file := configFile()
	conf, err := godotenv.Read(file)
	if err != nil {
		return err
	}

	cfg.clientId = conf[CLIENT_ID]
	cfg.clientSecret = conf[CLIENT_SECRET]
	cfg.channelLogin = conf[CHANNEL_LOGIN]
	cfg.refreshToken = conf[REFRESH_TOKEN]
	return nil
}

func (cfg *staticConfig) write() error {
	path := configFile()
	data := map[string]string{
		CLIENT_ID:     cfg.clientId,
		CLIENT_SECRET: cfg.clientSecret,
		CHANNEL_LOGIN: cfg.channelLogin,
		REFRESH_TOKEN: cfg.refreshToken,
	}
	return godotenv.Write(data, path)
}
