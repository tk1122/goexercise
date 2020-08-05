package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"goercises/task/db"
	"log"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all of your task",
	Run: func(cmd *cobra.Command, args []string) {
		tasks, err := db.AllTasks()
		if err != nil {
			log.Fatalf("cannot run list command: %v\n", err)
		}
		if len(tasks) == 0 {
			fmt.Println("You have no tasks to complete.")
			return
		}
		fmt.Println("You have the following tasks: ")
		for i, t := range tasks {
			fmt.Printf("%d. %s\n", i + 1, t.Value)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
