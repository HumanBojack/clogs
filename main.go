package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/araddon/dateparse"
)

type logEntry struct {
	timestamp time.Time
	color     string
	line      string
}

func main() {
	noColor := flag.Bool("no-color", false, "Disable color output")
	flag.Parse()

	args := flag.Args()
	if len(args)%3 != 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-no-color] <file> <start> <end> [<file> <start> <end> ...]\n", os.Args[0])
		os.Exit(1)
	}

	var entries []logEntry
	colors := []string{"\033[30m", "\033[31m", "\033[32m", "\033[33m", "\033[34m", "\033[35m", "\033[36m", "\033[37m"}
	colorIndex := 0

	for i := 0; i < len(args); i += 3 {
		arg := args[i]
		start, _ := strconv.Atoi(args[i+1])
		end, _ := strconv.Atoi(args[i+2])

		file, err := os.Open(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()

		entries = append(entries, readFrom(file, start, end, colors[colorIndex])...)
		colorIndex = (colorIndex + 1) % len(colors)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].timestamp.Before(entries[j].timestamp)
	})

	for _, entry := range entries {
		if *noColor {
			fmt.Println(entry.line)
		} else {
			fmt.Println(entry.color + entry.line + "\033[0m")
		}
	}
}

func readFrom(reader io.Reader, start int, end int, color string) []logEntry {
	var entries []logEntry
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		rawDT := line[start:end]

		timestamp, err := dateparse.ParseAny(rawDT)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing timestamp: %v\n", err)
			os.Exit(1)
		}
		entries = append(entries, logEntry{timestamp, color, line})
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
	return entries
}
