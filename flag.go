// SPDX-License-Identifier: GPL-3.0-or-later

package main

import "flag"

type requestFlags struct {
	// No-op auxiliary flags
	help    bool
	version bool

	// Flags that affect the program operation
	dryRun  bool
	verbose bool

	// Flags that only present information
	printCurrent bool
	listProfiles bool

	// Flags that change information
	useProfile  string
	newTitle    string
	newCategory string
	newTags     string
	newLanguage string
}

func initRequestFlags() requestFlags {
	flags := requestFlags{}

	flag.BoolVar(&flags.help, "help", false, "Print usage information and exit.")
	flag.BoolVar(&flags.version, "version", false, "Print version information and exit.")

	flag.BoolVar(&flags.dryRun, "dry-run", false, "Do not actually commit the new stream information when it changes.")
	flag.BoolVar(&flags.verbose, "verbose", false, "Print additional information to output showing what is StreamTitle doing.")

	flag.BoolVar(&flags.printCurrent, "current", false, "Print the current stream information reported by the Twitch API and exit.")
	flag.BoolVar(&flags.listProfiles, "list-profiles", false, "List the known profiles available to use and exit.")

	flag.StringVar(&flags.useProfile, "profile", "", "Set the stream information from an existing profile.")
	flag.StringVar(&flags.newTitle, "title", "", "The new title to use for the channel information, string.")
	flag.StringVar(&flags.newCategory, "game", "", "The game ID to use for the channel information, numeric.")
	flag.StringVar(&flags.newTags, "tags", "", "The new tags to use for the channel information, comma separated.")
	flag.StringVar(&flags.newLanguage, "language", "", "The new language code to use for the channel information, ISO code.")

	flag.Parse()
	return flags
}

func (r *requestFlags) showUsage() {
	flag.PrintDefaults()
}

func (r *requestFlags) changing() bool {
	return r.useProfile != "" || r.newTitle != "" || r.newCategory != "" || r.newTags != "" || r.newLanguage != ""
}
