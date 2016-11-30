// Package implements tail -n which returns last n lines of file
package tail_n

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

// Tail returns slice of last n strings from file in path
func Tail(path string, n int) ([]string, error) {
	tail, _, err := tail(path, n, nil, true)
	return tail, err
}

// TailReverse returns reversed slice of last n strings from file in path
func TailReverse(path string, n int) ([]string, error) {
	tail, _, err := tail(path, n, nil, false)
	return tail, err
}

// TailBytes returns slice bytes divided by n new lines from file in path
func TailBytes(path string, n int) ([]byte, error) {
	_, tail, err := tail(path, n, nil, true)
	return tail, err
}

// TailBytesReverse returns reversed slice of bytes divided by n new lines from file in path
func TailBytesReverse(path string, n int) ([]byte, error) {
	_, tail, err := tail(path, n, nil, false)
	return tail, err
}

func tail(path string, n int, logger *log.Logger, keepOrder bool) (tail []string, tailBytes []byte, err error) {
	if n <= 0 {
		return
	}
	if logger == nil {
		logger = log.New(ioutil.Discard, "", log.Ldate)
	}
	file, err := os.Open(path)
	if err != nil {
		logger.Printf("Failed to open file %s: %s\n", path, err)
		return
	}
	defer file.Close()

	nl := []byte("\n")
	offsetEnd, err := file.Seek(0, io.SeekEnd)
	newStringEnd := offsetEnd
	cursor := make([]byte, 1)
	var tmpBytes [][]byte
	for i := offsetEnd - 1; i >= 0; i-- {
		_, err = file.ReadAt(cursor, i)
		if err != nil {
			logger.Printf("Failed to read at %d: %s\n", i, err)
			break

		}

		if cursor[0] == nl[0] {
			_, err = file.Seek(i+1, io.SeekStart)
			if err != nil {
				logger.Printf("Failed to seek at %d: %s\n", i, err)
				break
			}
			newString := make([]byte, newStringEnd-i)
			_, err = file.Read(newString)
			if err != nil {
				logger.Printf("Failed to read new line at %d: %s\n", i, err)
				break
			}
			tail = append(tail, string(newString))
			tmpBytes = append(tmpBytes, newString)
			if len(tail) >= n {
				break
			}
			newStringEnd = i
		}
	}

	if keepOrder {
		reverse(tail)
		reverseBytes(tmpBytes)
	}

	tailBytes = mergeBytes(tmpBytes)

	return
}

func reverse(list []string) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func reverseBytes(list [][]byte) {
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
}

func mergeBytes(list [][]byte) (merged []byte) {
	for _, item := range list {
		merged = append(merged, item...)
	}

	return
}
