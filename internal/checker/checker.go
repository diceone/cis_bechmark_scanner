package checker

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

const (
	ResultPass  = "PASS"
	ResultFail  = "FAIL"
	ResultError = "ERROR"
)

type Check struct {
	ID          string `yaml:"id"`
	Description string `yaml:"description"`
	Command     string `yaml:"command"`
	Expected    string `yaml:"expected"`
	Remediation string `yaml:"remediation,omitempty"`
}

type Result struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Command     string `json:"command"`
	Expected    string `json:"expected"`
	Result      string `json:"result"`
	Output      string `json:"output,omitempty"`
	Remediation string `json:"remediation,omitempty"`
}

type Executor interface {
	Execute(ctx context.Context, cmd string) (string, error)
}

type Checker struct {
	executor  Executor
	checks    []Check
	results   []Result
	mu        sync.Mutex
	parallel  int
	matchMode MatchMode
}

type MatchMode int

const (
	MatchExact    MatchMode = iota
	MatchContains MatchMode = iota
)

type Option func(*Checker)

func WithParallel(workers int) Option {
	return func(c *Checker) {
		if workers > 0 {
			c.parallel = workers
		}
	}
}

func WithMatchMode(mode MatchMode) Option {
	return func(c *Checker) {
		c.matchMode = mode
	}
}

func NewChecker(executor Executor, opts ...Option) *Checker {
	c := &Checker{
		executor:  executor,
		parallel:  1,
		matchMode: MatchExact,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Checker) LoadChecks(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("opening checks file: %w", err)
	}
	defer file.Close()

	c.checks = nil
	c.results = nil

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&c.checks); err != nil {
		return fmt.Errorf("parsing yaml: %w", err)
	}

	return c.validateChecks()
}

func (c *Checker) validateChecks() error {
	seen := make(map[string]bool)
	for _, check := range c.checks {
		if check.ID == "" {
			return fmt.Errorf("check missing ID")
		}
		if seen[check.ID] {
			return fmt.Errorf("duplicate check ID: %s", check.ID)
		}
		seen[check.ID] = true
	}
	return nil
}

func (c *Checker) GetChecks() []Check {
	return c.checks
}

func (c *Checker) RunChecks(ctx context.Context, checkIDs []string) []Result {
	if len(c.checks) == 0 {
		return nil
	}

	checksToRun := c.getChecksToRun(checkIDs)
	c.results = make([]Result, 0, len(checksToRun))

	if c.parallel > 1 && len(checksToRun) > 1 {
		c.runChecksParallel(ctx, checksToRun)
	} else {
		c.runChecksSequential(ctx, checksToRun)
	}

	return c.results
}

func (c *Checker) getChecksToRun(checkIDs []string) []int {
	if len(checkIDs) == 0 {
		indices := make([]int, len(c.checks))
		for i := range c.checks {
			indices[i] = i
		}
		return indices
	}

	indices := make([]int, 0, len(checkIDs))
	seen := make(map[string]bool)
	for _, id := range checkIDs {
		if seen[id] {
			continue
		}
		seen[id] = true
		for i := range c.checks {
			if c.checks[i].ID == id {
				indices = append(indices, i)
				break
			}
		}
	}
	return indices
}

func (c *Checker) runChecksSequential(ctx context.Context, indices []int) {
	for _, i := range indices {
		c.runCheck(ctx, i)
	}
}

func (c *Checker) runChecksParallel(ctx context.Context, indices []int) {
	var wg sync.WaitGroup
	wg.Add(len(indices))

	workerCount := c.parallel
	if workerCount > len(indices) {
		workerCount = len(indices)
	}

	jobs := make(chan int, len(indices))
	for _, i := range indices {
		jobs <- i
	}
	close(jobs)

	for w := 0; w < workerCount; w++ {
		go func() {
			for idx := range jobs {
				c.runCheck(ctx, idx)
				wg.Done()
			}
		}()
	}

	wg.Wait()
}

func (c *Checker) runCheck(ctx context.Context, index int) {
	check := &c.checks[index]
	output, err := c.executor.Execute(ctx, check.Command)

	result := Result{
		ID:          check.ID,
		Description: check.Description,
		Command:     check.Command,
		Expected:    check.Expected,
		Output:      output,
		Remediation: check.Remediation,
	}

	if err != nil {
		if isExitCodeError(err) {
			if check.Expected == "" {
				result.Result = ResultPass
			} else {
				result.Result = ResultFail
			}
		} else {
			result.Result = fmt.Sprintf("%s: %v", ResultError, err)
		}
	} else if c.matchesExpected(output, check.Expected) {
		result.Result = ResultPass
	} else {
		result.Result = ResultFail
	}

	c.mu.Lock()
	c.results = append(c.results, result)
	c.mu.Unlock()
}

func (c *Checker) matchesExpected(output, expected string) bool {
	if c.matchMode == MatchContains {
		return strings.Contains(output, expected)
	}
	return output == expected
}

func isExitCodeError(err error) bool {
	type exitCoder interface {
		ExitStatus() int
	}
	_, ok := err.(exitCoder)
	return ok
}
