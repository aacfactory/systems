package commons

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func ReadLines(filename string) (lines []string, err error) {
	file, openErr := os.Open(filename)
	if openErr != nil {
		err = openErr
		return
	}
	lines = make([]string, 0, 1)
	reader := bufio.NewReader(file)
	for {
		line, readErr := reader.ReadString('\n')
		if len(line) > 0 {
			lines = append(lines, strings.Trim(line, "\n"))
		}
		if readErr != nil {
			if readErr == io.EOF {
				break
			}
			err = readErr
			return
		}
	}
	return
}
