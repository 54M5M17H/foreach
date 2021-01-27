package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"
)

const (
	replacer = "__"
)

func main() {

	if len(os.Args) <= 1 {
		return
	}
	args := os.Args[1:]

	separators := os.Getenv("IFS")
	if len(separators) == 0 {
		separators = "\n \t"
	}

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(createSplitter(separators))

	if indexWithoutFirst := indexOf(args[1:], replacer); indexWithoutFirst > -1 {
		replaceIterator(args, scanner, indexWithoutFirst+1)
	} else {
		stdinIterator(args, scanner)
	}
}

func stdinIterator(args []string, scanner *bufio.Scanner) {
	for scanner.Scan() {
		l := scanner.Text()
		cmd := exec.Command(args[0], args[1:]...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmdStdin, _ := cmd.StdinPipe()
		cmdStdin.Write([]byte(l))
		cmdStdin.Close()

		cmd.Run()
	}
}

func replaceIterator(args []string, scanner *bufio.Scanner, index int) {
	for scanner.Scan() {
		l := scanner.Text()
		replacedArgs := args
		replacedArgs[index] = l
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}

func createSplitter(separators string) bufio.SplitFunc {
	buffer := []byte{}
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		advance = len(data)

		if i := bytes.IndexAny(data, separators); i == 0 {
			advance = 1
			return
		} else if i >= 0 {
			included := data[:i]
			token = append(buffer, included...)
			buffer = []byte{}

			// +1 skips delimiter
			advance = i + 1
			return
		}

		buffer = append(buffer, data...)

		if atEOF {
			token = buffer
			err = io.EOF
		}

		return
	}
}

func indexOf(arr []string, toFind string) int {
	for i, s := range arr {
		if s == toFind {
			return i
		}
	}
	return -1
}
