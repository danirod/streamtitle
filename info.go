// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
	"strings"

	"github.com/nicklaw5/helix/v2"
)

type streamInfo struct {
	title    string
	game     string
	language string
	tags     []string
}

func (info *streamInfo) printInfo() {
	fmt.Println("       Title:", info.title)
	fmt.Println("        Game:", info.game)
	fmt.Println("    Language:", info.language)
	fmt.Println("        Tags:", strings.Join(info.tags, ","))
}

func (info *streamInfo) load(resp *helix.ChannelInformation) {
	info.title = resp.Title
	info.game = resp.GameID
	info.language = resp.BroadcasterLanguage
	info.tags = resp.Tags
}

func (info *streamInfo) apply(profile *appProfile) {
	info.setTitle(profile.Title)
	info.setGame(fmt.Sprintf("%d", profile.Game))
	info.setLanguage(profile.Language)
	info.setTagArray(profile.Tags)
}

func (info *streamInfo) setTitle(title string) {
	if title != "" {
		info.title = title
	}
}

func (info *streamInfo) setGame(game string) {
	if game != "" {
		info.game = game
	}
}

func (info *streamInfo) setLanguage(language string) {
	if language != "" {
		info.language = language
	}
}

func (info *streamInfo) setTagArray(tags []string) {
	if len(tags) > 0 {
		info.tags = tags
	}
}

func (info *streamInfo) setTagString(tagString string) {
	if tagString != "" {
		tags := strings.Split(tagString, ",")
		for i, tag := range tags {
			tags[i] = strings.TrimSpace(tag)
		}
		info.setTagArray(tags)
	}
}
