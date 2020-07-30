package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type problem struct {
	q, a string
}

func timmer(timeLimit time.Duration, timeout chan<- bool) {
	<-time.After(timeLimit)
	timeout <- true
}

func parseQuestions(file *os.File, done chan<- bool, allQuestions, rightQuestions *int) {
	csvReader := csv.NewReader(file)

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			done <- true
			break
		}
		if err != nil {
			log.Fatalf("reading record: %v", err)
		}

		r2 := problem{q: record[0], a: strings.TrimSpace(record[1])}
		*allQuestions++
		fmt.Println("Question: ", r2.q)
		fmt.Print("Answer: ")

		var answer string
		__, err := fmt.Scanf("%s", &answer)
		_ = __

		if err != nil {
			log.Fatalf("scan stdin: %s", err)
		}

		if answer == r2.a {
			*rightQuestions++
			fmt.Println("Right!")
		} else {
			fmt.Println("Wrong!")
		}
	}
}

func main() {
	fileName := flag.String("file", "problems.csv", "file contains the problem set")
	timeLimit := flag.Duration("time", 5*time.Second, "quiz timeLimit")
	done, timeout := make(chan bool), make(chan bool)
	allQuestions, rightQuestions := 0, 0

	flag.Parse()
	file, err := os.Open(*fileName)
	if err != nil {
		log.Fatalf("open file %v: %v", fileName, err)
	}
	defer file.Close()

	go timmer(*timeLimit, timeout)
	go parseQuestions(file, done, &allQuestions, &rightQuestions)

	select {
	case <-done:
		fmt.Println("\n Answered all questions!!!")
	case <-timeout:
		fmt.Println("\nThe quiz is over!")
	}

	fmt.Println("All questions: ", allQuestions)
	fmt.Println("Right questions: ", rightQuestions)
}
