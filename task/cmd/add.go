package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"goercises/task/db"
	"strings"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a task to your task list",
	Run: func(cmd *cobra.Command, args []string) {
		task := strings.Join(args, " ")
		_, err := db.CreateTask(task)
		if err != nil {
			fmt.Println("cannot add new task: ", err)
			return
		}
		fmt.Printf("Added \"%s\" to your task list.\n", task)
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
