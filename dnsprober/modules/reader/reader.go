package reader

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/RevoltSecurities/Dnsprober/dnsprober/modules/utils"
)

func Reader(filename string) ([]string, error) {
	var urls []string
	filecontent, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer filecontent.Close()
	scanner := bufio.NewScanner(filecontent)
	for scanner.Scan() {
		line := scanner.Text()
		tline := strings.TrimSpace(line)
		urls = append(urls, tline)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return utils.Set(urls), nil
}

func StreamWithBuffered(filename string, limit int) (<-chan string, error) {
	if !utils.FileExists(filename) {
		return nil, fmt.Errorf("no such file or directory exist")
	}
	streamed := make(chan string, limit)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	go func() {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			s := strings.TrimSpace(scanner.Text())
			if s != "" {
				streamed <- s
			}
		}
		close(streamed)
	}()
	return streamed, nil
}

func Streamer(filename string) (chan string, error) {
	if !utils.FileExists(filename) {
		return nil, fmt.Errorf("no such file or directory exist")
	}
	streamed := make(chan string)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	go func() {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			s := strings.TrimSpace(scanner.Text())
			if s != "" {
				streamed <- s
			}
		}
		close(streamed)
	}()
	return streamed, nil
}

func IOStreamer(filename io.Reader) (chan string, error) {

	streamed := make(chan string)
	go func() {
		scanner := bufio.NewScanner(filename)
		for scanner.Scan() {
			s := strings.TrimSpace(scanner.Text())
			if s != "" {
				streamed <- s
			}
		}
		close(streamed)
	}()
	return streamed, nil
}
