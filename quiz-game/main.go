package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

type problem struct {
	Question string
	Answer   string
}

func parseProblemsFromData(data [][]string) []problem {
	problems := make([]problem, len(data))

	for i, line := range data {
		problems[i] = problem{Question: line[0], Answer: line[1]}
	}

	return problems
}

func main() {

	csvFileName := flag.String("csv", "problems.csv", "Path to a csv with format 'question,answer'")
	timeLimit := flag.Int("timeLimit", 30, "The time limit for the quiz in seconds")
	flag.Parse()

	// read file
	f, err := os.Open(*csvFileName)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	problems := parseProblemsFromData(data)

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	countOfCorrect := 0
	for i, question := range problems {
		fmt.Printf("Question #%d: %s\n", i+1, question.Question)
		answerCh := make(chan string)
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Printf("You scored %d out of %d.", countOfCorrect, len(problems))
			return
		case answer := <-answerCh:
			if answer == question.Answer {
				countOfCorrect++
			}
		}
	}

	fmt.Printf("You scored %d out of %d.", countOfCorrect, len(problems))
}
