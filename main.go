// SPDX-License-Identifier: GPL-3.0-or-later

package main

import (
	"fmt"
	"io"
	"log"
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
	if !flags.printCurrent && !flags.changing() {
		flags.showUsage()
		return
	}

	context := newContext()
	if flags.printCurrent {
		context.streamInfo.printInfo()
		return
	}
	if flags.changing() {
		context.streamInfo.setTitle(flags.newTitle)
		context.streamInfo.setGame(flags.newCategory)
		context.streamInfo.setTagString(flags.newTags)
		context.streamInfo.setLanguage(flags.newLanguage)
		if flags.printCurrent {
			fmt.Println("New stream information:")
			context.streamInfo.printInfo()
		}
		context.SendStreamInfo()
	}
}

func newContext() *Client {
	context, err := NewContext()
	if err != nil {
		log.Fatal(err)
	}
	if err := context.Login(); err != nil {
		log.Fatal(err)
	}
	if err := context.FetchStreamInfo(); err != nil {
		log.Fatal(err)
	}
	return context
}
