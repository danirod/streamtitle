// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

var (
	logger = log.Default()

	appVersion = "1.0.12"
)

func main() {
	flags := initRequestFlags()
	if !flags.verbose {
		logger.SetOutput(io.Discard)
	}

	if flags.help {
		fmt.Printf(`%s [...flags]

StreamTitle can be used to show and to change the current information of your
Twitch channel. On first run, it will ask for login. Once logged in, you can
use the set of flags to either print the current channel information, or to
modify the channel information, allowing you to change multiple parameters.

Note that it is currently not possible to change the notification text, since
the Twitch API does not expose that field. Remember to check on your dashboard
before going live to make sure that the notification is the one that you want
to use.

Flags:
`, os.Args[0])
		flags.showUsage()
		return
	}

	if flags.version {
		fmt.Println("StreamTitle", appVersion)
		return
	}

	// Quick shortcut if no task is requested from the tool.
	if !flags.printCurrent && !flags.listProfiles && !flags.changing() {
		flags.showUsage()
		return
	}

	context := newContext()

	if flags.listProfiles {
		for _, p := range context.config.profiles() {
			fmt.Println(p)
		}
		return
	}

	initContext(context)

	if flags.printCurrent {
		context.streamInfo.printInfo()
		return
	}
	if flags.changing() {
		if flags.useProfile != "" {
			profile, found := context.config.profile(flags.useProfile)
			if found {
				context.streamInfo.apply(profile)
			} else {
				fmt.Println("Error: Profile", flags.useProfile, "not found")
				os.Exit(1)
			}
		}
		context.streamInfo.setTitle(flags.newTitle)
		context.streamInfo.setGame(flags.newCategory)
		context.streamInfo.setTagString(flags.newTags)
		context.streamInfo.setLanguage(flags.newLanguage)
		if flags.printCurrent {
			fmt.Println("New stream information:")
			context.streamInfo.printInfo()
		}
		context.SendStreamInfo(flags.dryRun)
	}
}

func newContext() *Client {
	context, err := NewContext()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	return context
}

func initContext(context *Client) {
	if err := context.Login(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if err := context.FetchStreamInfo(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
