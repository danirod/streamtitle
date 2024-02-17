// SPDX-License-Identifier: GPL-3.0-or-later

package main

import "flag"

type requestFlags struct {
	printCurrent bool
	printEnv     bool
	newTitle     string
	newCategory  string
	newTags      string
	newLanguage  string
}

func initRequestFlags() requestFlags {
	flags := requestFlags{}
	flag.BoolVar(&flags.printEnv, "env", false, "Whether to print the ENV file location")
	flag.BoolVar(&flags.printCurrent, "verbose", false, "Whether to print the current and new stream info")
	flag.StringVar(&flags.newTitle, "title", "", "The new title to use for the stream")
	flag.StringVar(&flags.newCategory, "game", "", "The game ID to use for the stream")
	flag.StringVar(&flags.newTags, "tags", "", "The new tags to use for the stream, comma separated")
	flag.StringVar(&flags.newLanguage, "language", "", "The new language code to use for the stream")
	flag.Parse()
	return flags
}

func (r *requestFlags) showUsage() {
	flag.PrintDefaults()
}

func (r *requestFlags) changing() bool {
	return r.newTitle != "" || r.newCategory != "" || r.newTags != "" || r.newLanguage != ""
}
