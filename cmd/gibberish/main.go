package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/prophittcorey/gibberish"
)

var (
	errNoFiles = fmt.Errorf("error: one or more files are required for training")
)

func main() {
	var classifierfile string
	var trainingfile string
	var check string
	var repl bool

	flag.StringVar(&classifierfile, "classifier", "", "a path to a classifier")
	flag.StringVar(&trainingfile, "train", "", "one or more text files to train a classifier with (plain text)")
	flag.StringVar(&check, "check", "", "a string to check using the specified classifier")
	flag.BoolVar(&repl, "repl", false, "if specified, a repl will be started")

	flag.Parse()

	if len(classifierfile) > 0 {
		if repl {
			classifier := gibberish.New()

			if err := classifier.LoadFile(classifierfile); err == nil {
				reader := bufio.NewScanner(os.Stdin)
				fmt.Printf("> Write some text to check...\n\n")

				for reader.Scan() {
					text := reader.Text()

					if text == "quit" || text == "exit" {
						break
					}

					result := classifier.Analyze(text)

					if result.IsGibberish {
						fmt.Printf("\n => Gibberish (%.2f%% / %.2f%%)\n", result.Probability, result.Threshold)
					} else {
						fmt.Printf("\n => Good (%.2f%% / %.2f%%)\n", result.Probability, result.Threshold)
					}

					fmt.Printf("\n> Write some text to check...\n\n")
				}
			}

			return
		}

		/* training */
		if len(trainingfile) > 0 {
			if err := train(classifierfile, trainingfile); err != nil {
				log.Fatal(err)
			}

			return
		}

		/* checking */
		if len(check) > 0 {
			classifier := gibberish.New()

			if err := classifier.LoadFile(classifierfile); err == nil {
				result := classifier.Analyze(check)

				fmt.Printf("Gibberish? %v (%.2f%%)\n", result.IsGibberish, result.Probability)
			}

			return
		}
	}

	flag.Usage()
}

func train(classifierfile, glob string) error {
	files, err := filepath.Glob(glob)

	if err != nil {
		return err
	}

	if len(files) == 0 {
		return errNoFiles
	}

	classifier := gibberish.New()

	for _, f := range files {
		(func() {
			f, err := os.Open(f)

			if err != nil {
				log.Fatal(err)
			}

			defer f.Close()

			classifier.Train(f)
		})()
	}

	if err = classifier.SaveFile(classifierfile); err != nil {
		return err
	}

	return nil
}
