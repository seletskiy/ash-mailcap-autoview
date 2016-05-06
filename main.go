package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/docopt/docopt-go"
)

const usage = `Supplementary program for ash-mailcap.

Reads ash-formatted overview diff and outputs section which contains
specified comment.

Exits with 1 if comment not found.
Exits with 2 if there are problems reading file.

Usage:
  $0 -h | --help
  $0 <file> <comment-id>
`

var reComment = regexp.MustCompile(`^(#\s+)(\[(\d+)@\d+\])`)

func main() {
	args, err := docopt.Parse(usage, nil, true, "1.0", false)
	if err != nil {
		panic(err)
	}

	file, err := os.Open(args["<file>"].(string))
	if err != nil {
		log.Printf("can't open specified file: %s", err)
	}

	commentFound := false

	section := []string{}

	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Printf(
					"error: can't find comment with ID '%s'\n",
					args["<comment-id>"].(string),
				)

				os.Exit(1)
			}

			log.Printf("error while reading file: %s", err)
			os.Exit(2)
		}

		matches := reComment.FindStringSubmatch(line)
		if len(matches) > 0 && matches[3] == args["<comment-id>"].(string) {
			commentFound = true
			line = reComment.ReplaceAllString(line, "$1\033[7m$2\033[0m")
		}

		if line == "\n" {
			if commentFound {
				fmt.Print(strings.Join(section, ""))
				os.Exit(0)
			}

			section = make([]string, 0)
		} else {
			section = append(section, line)
		}
	}
}
