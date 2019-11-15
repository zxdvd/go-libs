package fs

import (
	"os"
)

// TODO add options to support mode and flags
// like echo hello >> test.txt
func AppendFile(path, data string) error {
	f, err := os.OpenFile(path, os.O_RDWR | os.O_APPEND, os.ModeAppend)
	if err != nil {
		return err
	}
	_, err = f.WriteString(data)
	f.Close()
	return err
}

// like echo hello > test.txt
func TruncFileWithString(path, data string) error {
	f, err := os.OpenFile(path, os.O_RDWR | os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = f.WriteString(data)
	f.Close()
	return err
}
