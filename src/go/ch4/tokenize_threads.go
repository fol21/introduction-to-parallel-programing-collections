package ch4

import (
	"errors"
	"runtime"
	"strings"
	"sync"
)

// TokenizeResult captures the thread and tokens for one input line.
type TokenizeResult struct {
	LineIndex int
	Thread    int
	Line      string
	Tokens    []string
}

type unsafeTokenizerState struct {
	parts []string
	idx   int
}

func (s *unsafeTokenizerState) start(line string) {
	s.parts = strings.Fields(line)
	s.idx = 0
}

func (s *unsafeTokenizerState) next() (string, bool) {
	if s.idx >= len(s.parts) {
		return "", false
	}
	tok := s.parts[s.idx]
	s.idx++
	return tok, true
}

// TokenizeLinesUnsafe simulates pth_tokenize.c by using a shared non-reentrant tokenizer state.
func TokenizeLinesUnsafe(lines []string, threadCount int) ([]TokenizeResult, error) {
	return tokenizeLines(lines, threadCount, false)
}

// TokenizeLinesSafe simulates pth_tokenize_r.c by using a per-goroutine tokenizer state.
func TokenizeLinesSafe(lines []string, threadCount int) ([]TokenizeResult, error) {
	return tokenizeLines(lines, threadCount, true)
}

func tokenizeLines(lines []string, threadCount int, safe bool) ([]TokenizeResult, error) {
	if threadCount <= 0 {
		return nil, errors.New("threadCount must be positive")
	}

	sems := make([]chan struct{}, threadCount)
	for i := range sems {
		sems[i] = make(chan struct{}, 1)
	}
	sems[0] <- struct{}{}

	results := make([]TokenizeResult, 0, len(lines))
	var resultsMu sync.Mutex
	var indexMu sync.Mutex
	nextIndex := 0
	var wg sync.WaitGroup

	shared := &unsafeTokenizerState{}

	for rank := 0; rank < threadCount; rank++ {
		wg.Add(1)
		go func(r int) {
			defer wg.Done()
			nextRank := (r + 1) % threadCount
			for {
				<-sems[r]

				indexMu.Lock()
				if nextIndex >= len(lines) {
					indexMu.Unlock()
					sems[nextRank] <- struct{}{}
					return
				}
				idx := nextIndex
				line := lines[nextIndex]
				nextIndex++
				indexMu.Unlock()

				sems[nextRank] <- struct{}{}

				tokens := make([]string, 0)
				if safe {
					tokens = append(tokens, strings.Fields(line)...)
				} else {
					shared.start(line)
					for {
						tok, ok := shared.next()
						if !ok {
							break
						}
						tokens = append(tokens, tok)
						// Encourage overlap between workers to expose shared-state interference.
						runtime.Gosched()
					}
				}

				resultsMu.Lock()
				results = append(results, TokenizeResult{
					LineIndex: idx,
					Thread:    r,
					Line:      line,
					Tokens:    tokens,
				})
				resultsMu.Unlock()
			}
		}(rank)
	}

	wg.Wait()
	if len(results) > len(lines) {
		results = results[:len(lines)]
	}
	return results, nil
}
