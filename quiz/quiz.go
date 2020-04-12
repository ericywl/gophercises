package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Quiz struct to contain the parsed CSV
type Quiz struct {
	question, answer string
}

// readQuizCSV reads a CSV file containing quiz questions and answers
func readQuizCSV(filename string) []Quiz {
	csvFile, _ := os.Open(filename)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	var quizzes []Quiz
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		question := strings.Join(line[:len(line)-1], ",")
		answer := line[len(line)-1]
		quizzes = append(quizzes, Quiz{question, answer})
	}

	return quizzes
}

// shuffleQuizzes shuffles a quiz array randomly
func shuffleQuizzes(quizzes []Quiz) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(quizzes), func(i, j int) { quizzes[i], quizzes[j] = quizzes[j], quizzes[i] })
}

// quizGame starts a quiz game with the given quizzes array
func quizGame(quizzes []Quiz, scorePtr *int, ch1 chan string) {
	reader := bufio.NewReader(os.Stdin)
	for _, quiz := range quizzes {
		fmt.Println("Question:", quiz.question)
		fmt.Printf("Enter answer: ")
		ans, _ := reader.ReadString('\n')
		ans = strings.TrimSpace(ans)
		if strings.ToLower(ans) == strings.ToLower(quiz.answer) {
			*scorePtr++
		}
	}
	// Send result through channel
	ch1 <- fmt.Sprintf("\nYour score: %d / %d", *scorePtr, len(quizzes))
	close(ch1)
}

func main() {
	// Set up flags
	shufflePtr := flag.Bool("shuffle", true, "toggle random shuffling of the questions")
	timePtr := flag.Int("time", 30, "time limit for quiz")
	flag.Usage = func() {
		fmt.Println("Usage of ./quiz [csvFileName]:")
		flag.PrintDefaults()
	}
	// Parse command line arguments
	flag.Parse()
	var csvFileName string
	switch len(flag.Args()) {
	case 0:
		csvFileName = "problems.csv"
	case 1:
		csvFileName = flag.Args()[0]
	default:
		flag.Usage()
		return
	}
	// Read and shuffle quizzes
	quizzes := readQuizCSV(csvFileName)
	if *shufflePtr {
		shuffleQuizzes(quizzes)
	}
	// Welcome message
	fmt.Println("Welcome to the Quiz Game!")
	fmt.Println("Press ENTER to start the quiz.")
	text, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	if text != "\n" {
		fmt.Println()
	}
	// Run quiz game
	score := 0
	ch1 := make(chan string, 1)
	go quizGame(quizzes, &score, ch1)
	// Wait for result or timeout
	select {
	case result := <-ch1:
		fmt.Println(result)
	case <-time.After(time.Duration(*timePtr) * time.Second):
		fmt.Println("\n\nYou ran out of time!")
		fmt.Printf("Your score: %d / %d\n", score, len(quizzes))
	}
}
