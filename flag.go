// SPDX-License-Identifier: GPL-3.0-or-later

package main

import "flag"

type requestFlags struct {
	dryRun       bool
	verbose      bool
	printCurrent bool
	useProfile   string
	listProfiles bool
	newTitle     string
	newCategory  string
	newTags      string
	newLanguage  string
}

func initRequestFlags() requestFlags {
	flags := requestFlags{}
	flag.BoolVar(&flags.dryRun, "dry-run", false, "Do not actually commit the new stream information")
	flag.BoolVar(&flags.verbose, "verbose", false, "Print log messages")
	flag.BoolVar(&flags.printCurrent, "current", false, "Whether to print the current stream information")
	flag.BoolVar(&flags.listProfiles, "list-profiles", false, "List the known profiles and quit")
	flag.StringVar(&flags.useProfile, "profile", "", "The profile to use")
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
	return r.useProfile != "" || r.newTitle != "" || r.newCategory != "" || r.newTags != "" || r.newLanguage != ""
}
