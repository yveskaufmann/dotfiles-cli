package github

import (
	"sync"

	"yv35.com/dotfiles-cli/internal/config"
	"yv35.com/dotfiles-cli/internal/types"
)

const maxConcurrentDownloads = 4

// processWithWorkerPool processes GitHub release installations concurrently using a worker pool
func (p *Provider) processWithWorkerPool(specs []config.GithubSpec, onComplete types.OnTaskComplete) error {
	jobs := make(chan config.GithubSpec, len(specs))
	results := make(chan types.TaskResult, len(specs))
	var wg sync.WaitGroup

	// Track first error
	var firstErr error
	var errMutex sync.Mutex

	// Start worker goroutines
	for i := 0; i < maxConcurrentDownloads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for spec := range jobs {
				err := p.installGithubRelease(spec)
				result := types.TaskResult{
					Name: spec.Name,
				}
				if err != nil {
					result.Status = types.StatusFailed
					result.Error = err

					// Capture first error
					errMutex.Lock()
					if firstErr == nil {
						firstErr = err
					}
					errMutex.Unlock()
				} else {
					result.Status = types.StatusSuccess
				}
				results <- result
			}
		}()
	}

	// Send jobs
	for _, spec := range specs {
		jobs <- spec
	}
	close(jobs)

	// Wait for completion
	wg.Wait()
	close(results)

	// Collect results
	for result := range results {
		onComplete(result)
	}

	return firstErr
}
