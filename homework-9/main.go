package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	k := flag.Int("k", 5, "максимальное число одновременных запросов")
	q := flag.String("q", "Go", "искомая строка")
	flag.Parse()

	if *k <= 0 {
		fmt.Fprintln(os.Stderr, "Error: количество одновременных запросов (-k) должно быть больше нуля")
		os.Exit(1)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var totalCount atomic.Int64

	sem := make(chan struct{}, *k)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if url == "" {
			continue
		}

		sem <- struct{}{}

		wg.Add(1)

		go func(targetURL string) {
			defer func() {
				<-sem
				wg.Done()
			}()

			count, err := fetchAndCount(targetURL, *q)

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error for %s: %v\n", targetURL, err)
				return
			}

			fmt.Printf("Count for %s: %d\n", targetURL, count)
			totalCount.Add(int64(count))
		}(url)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}

	wg.Wait()
	fmt.Printf("Total: %d\n", totalCount.Load())
}

func fetchAndCount(url string, searchStr string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	return strings.Count(string(body), searchStr), nil
}
