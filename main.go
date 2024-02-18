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
)

func main() {
	flags := initRequestFlags()
	if !flags.verbose {
		logger.SetOutput(io.Discard)
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
