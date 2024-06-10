package main

import (
	"bufio"
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
	line      string
}

func main() {
	args := os.Args[1:]
	var entries []logEntry

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

		entries = append(entries, readFrom(file, start, end)...)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].timestamp.Before(entries[j].timestamp)
	})

	for _, entry := range entries {
		fmt.Println(entry.line)
	}
}

func readFrom(reader io.Reader, start int, end int) []logEntry {
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
		entries = append(entries, logEntry{timestamp, line})
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
	return entries
}
