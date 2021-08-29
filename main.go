package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	fileName := flag.String("quiz", "problems.csv", "problems in csv format.")
	timer := flag.Duration("limit", 10*time.Second, "time limit to each question.")
	flag.Parse()

	f, err := os.Open(*fileName)
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(f)
	qa, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	var correct int
	var incorrect int
	problems := parseProblems(qa)

	for i, p := range problems {
		fmt.Printf("Problem #%d: %v\n", i+1, p.Question)
		answered := make(chan struct{})
		go func() {
			select {
			case <-time.After(*timer):
				log.Fatal("time is up.")
			case <-answered:
			}
		}()
		var answer string
		fmt.Scanf("%s\n", &answer)
		answered <- struct{}{}
		if answer != p.Answer {
			incorrect++
			continue
		}
		correct++
	}

	fmt.Printf("You got %d question(s) right and %d wrong", correct, incorrect)
}

type problem struct {
	Question string
	Answer   string
}

func parseProblems(input [][]string) []problem {
	problems := make([]problem, len(input))
	for i, line := range input {
		problems[i] = problem{
			Question: line[0],
			Answer:   line[1],
		}
	}
	return problems
}
