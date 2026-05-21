package jetbrains

import (
	"sync"

	"yv35.com/dotfiles-cli/internal/types"
)

// maxConcurrentDownloads limits parallel IDE downloads.  IDEs are large
// (~1 GB each) so a conservative limit avoids saturating the network.
const maxConcurrentDownloads = 2

// processWithWorkerPool installs the given candidates concurrently using a
// bounded worker pool, then calls onComplete for every result.
func (p *Provider) processWithWorkerPool(candidates []installCandidate, onComplete types.OnTaskComplete) error {
	jobs := make(chan installCandidate, len(candidates))
	results := make(chan types.TaskResult, len(candidates))

	var wg sync.WaitGroup
	var firstErr error
	var errMu sync.Mutex

	// Start workers.
	for i := 0; i < maxConcurrentDownloads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for c := range jobs {
				displayName := c.Spec.IDE
				if c.Spec.Name != "" {
					displayName = c.Spec.Name
				}

				err := p.installIDE(c)
				result := types.TaskResult{Name: displayName}
				if err != nil {
					result.Status = types.StatusFailed
					result.Error = err
					errMu.Lock()
					if firstErr == nil {
						firstErr = err
					}
					errMu.Unlock()
				} else {
					result.Status = types.StatusSuccess
				}
				results <- result
			}
		}()
	}

	// Enqueue all jobs.
	for _, c := range candidates {
		jobs <- c
	}
	close(jobs)

	// Wait for all workers, then close results.
	wg.Wait()
	close(results)

	// Deliver results to the engine callback.
	for r := range results {
		onComplete(r)
	}

	return firstErr
}
