package main

import "log"
import "fmt"

func main() {
	flags := initRequestFlags()
	if flags.printEnv {
		fmt.Println(configFile())
		return
	}

	// Quick shortcut if no task is requested from the tool.
	if !flags.printCurrent && !flags.changing() {
		flags.showUsage()
		return
	}

	context := newContext()
	if flags.printCurrent {
		fmt.Println("Current stream information:")
		context.streamInfo.printInfo()
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

func newContext() *Context {
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
