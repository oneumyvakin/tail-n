// Package implements tail -n which returns last n lines of file
package tail_n

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// Tail returns slice of last n strings from file in path
func Tail(path string, n int) ([]string, error) {
	tail, _, err := tail(path, n, true)
	return tail, err
}

// TailReverse returns reversed slice of last n strings from file
func TailReverse(path string, n int) ([]string, error) {
	tail, _, err := tail(path, n, false)
	return tail, err
}

// TailBytes returns bytes of last n strings from file
func TailBytes(path string, n int) ([]byte, error) {
	_, tail, err := tail(path, n, true)
	return tail, err
}

// TailBytesReverse returns bytes of last n string from file in reversed order
func TailBytesReverse(path string, n int) ([]byte, error) {
	_, tail, err := tail(path, n, false)
	return tail, err
}

// Ftail writes bytes of last n strings from file path and returns number of written bytes
func Ftail(w io.Writer, path string, n int) (int, error) {
	_, tail, err := tail(path, n, true)
	if err != nil {
		return 0, err
	}
	return w.Write(tail)
}

func tail(path string, n int, keepOrder bool) (tail []string, tailBytes []byte, err error) {
	if n <= 0 {
		return
	}
	file, err := os.Open(path)
	if err != nil {
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
			err = errors.New(fmt.Sprintf("Failed to read at %d: %s\n", i, err))
			break
		}

		if cursor[0] == nl[0] {
			if newStringEnd == i+1 {
				tail = append(tail, "\n")
				tmpBytes = append(tmpBytes, nl)
				continue
			}
			_, err = file.Seek(i+1, io.SeekStart)
			if err != nil {
				err = errors.New(fmt.Sprintf("Failed to seek at %d: %s\n", i, err))
				break
			}
			newString := make([]byte, newStringEnd-i)
			_, err = file.Read(newString)
			if err != nil {
				err = errors.New(fmt.Sprintf("Failed to read new line at %d: %s\n", i, err))
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
