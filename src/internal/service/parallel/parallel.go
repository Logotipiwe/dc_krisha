package parallel

import (
	"sync"
)

func DoJobs[T any](jobs []func() T, workersCount int) []T {
	wg := sync.WaitGroup{}
	count := len(jobs)
	jobsCh := make(chan func() T, count)
	resultsCh := make(chan T, count)
	for _, job := range jobs {
		jobsCh <- job
	}
	close(jobsCh)
	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobsCh {
				resultsCh <- job()
			}
		}()
	}
	wg.Wait()
	close(resultsCh)

	results := make([]T, 0)

	for result := range resultsCh {
		results = append(results, result)
	}
	return results
}
