package main

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"
	"strings"
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

	argIndexesToReplace := findIndexesInArr(args, replacer)

	if len(argIndexesToReplace) == 0 {
		stdinIterator(args, scanner)
	} else {
		replaceIterator(args, scanner, argIndexesToReplace)
	}

}

func findIndexesInArr(arr []string, strToMatch string) (r []int) {
	for i, s := range arr {
		if strings.Contains(s, strToMatch) {
			r = append(r, i)
		}
	}
	return
}

func stdinIterator(args []string, scanner *bufio.Scanner) {
	for scanner.Scan() {
		l := scanner.Text()
		cmd := exec.Command(args[0], args[1:]...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmdStdin, _ := cmd.StdinPipe()

		// simply call command once per line
		// e.g. printf does its own substitution
		// relies on data coming from stdin
		cmdStdin.Write([]byte(l))
		cmdStdin.Close()

		cmd.Run()
	}
}

func replaceIterator(args []string, scanner *bufio.Scanner, indexesToReplace []int) {
	for scanner.Scan() {
		inputLine := scanner.Text()
		replacedArgs := make([]string, len(args))
		copy(replacedArgs, args)
		for _, indexToReplace := range indexesToReplace {
			argWithReplacer := replacedArgs[indexToReplace]
			argReplaced := strings.Replace(argWithReplacer, replacer, inputLine, -1)
			replacedArgs[indexToReplace] = argReplaced
		}
		cmd := exec.Command(replacedArgs[0], replacedArgs[1:]...)
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
