package main

import (
	"ch2/todo"
	"flag"
	"fmt"
	"os"
)

const todoFileName = ".todo.json"

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%s tool. Developed for The Pragmatic Bookshelf\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Copyright 2024\n")
		fmt.Fprintln(flag.CommandLine.Output(), "Usage information:")
		flag.PrintDefaults()
	}
	// Parsing command line flags
	task := flag.String("task", "", "Task to be included in the todo list")
	list := flag.Bool("list", false, "List all the tasks")
	complete := flag.Int("complete", 0, "Item to be completed")
	flag.Parse()

	// Initialize the todo list from file
	l := &todo.List{}
	if err := l.Get(todoFileName); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	saveList := func() {
		if err := l.Save(todoFileName); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	// Decide what to do based on the flags
	switch {
	// Show list
	case *list:
		for _, item := range *l {
			if !item.Done {
				fmt.Println(item)
			}
		}
	// Complete item
	case *complete > 0:
		if err := l.Complete(*complete); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		saveList()
	// Add task
	case *task != "":
		l.Add(*task)
		saveList()
	// Invalid flag
	default:
		fmt.Fprintln(os.Stderr, "Invalid option")
		os.Exit(1)
	}
}
