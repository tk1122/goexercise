package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"goercises/task/db"
	"os"
	"strconv"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Mark a task as complete",
	Run: func(cmd *cobra.Command, args []string) {
		var ids []int
		for _, a := range args {
			if id, err := strconv.Atoi(a); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "failed to parse arg %v: %v\n", a, err)
			} else {
				ids = append(ids, id)
			}
		}
		tasks, err := db.AllTasks()
		if err != nil {
			fmt.Println("cannot mark the task as complete: ", err)
			return
		}
		for _, id := range ids {
			if id < 0 || id > len(tasks) {
				fmt.Println("Invalid task number: ", id)
				continue
			}
			task := tasks[id - 1]
			err := db.DeleteTask(task.Key)
			if err != nil {
				fmt.Printf("cannot delte task id %v:%v\n", id, err)
				continue
			}
			fmt.Printf("Marked %d as completed\n", id)
		}
	},
}

func init() {
	RootCmd.AddCommand(doCmd)
}
