package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type problem struct {
	q, a string
}

func main() {
	fileName := flag.String("file", "problems.csv", "file contains the problem set")
	flag.Parse()

	file, err := os.Open(*fileName)
	if err != nil {
		log.Fatalf("open file %v: %v", fileName, err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	allQuestions, rightQuestion := 0, 0

	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("reading record: %v", err)
		}

		r2 := problem{q: record[0], a: strings.TrimSpace(record[1])}
		allQuestions++
		fmt.Println("Question: ", r2.q)
		fmt.Print("Answer: ")


		var answer string
		__, err := fmt.Scanf("%s", &answer)
		_ = __

		if err!= nil {
			log.Fatalf("scan stdin: %s", err)
		}

		if answer == r2.a {
			rightQuestion++
			fmt.Println("Right!")
		} else {
			fmt.Println("Wrong!")
		}
	}

	fmt.Println("All questions: ", allQuestions)
	fmt.Println("Right questions: ", rightQuestion)
}
