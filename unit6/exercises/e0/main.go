package main

import (
	"context"
	"fmt"
	"time"
)

func process(iteration int) string {
	time.Sleep(2 * time.Second)
	return fmt.Sprintf("iteration %d is processed", iteration)
}

func cycle(ctx context.Context, size int) []string {
	lines := []string{}
	for i := 0; i < size; i++ {
		select {
		case <-ctx.Done():
			// if context was canceled or expired the channel will be closed,
			// so code flow will be here
			return lines
		default:
			// if context is still valid, other cases are not valid, so proceed with default one
			lines = append(lines, process(i))
		}
	}
	return lines
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// note that cancel function is deferred to avoid context leakage:
	// Failing to call the CancelFunc leaks the child and its children until the parent is canceled or the timer fires.

	lines := cycle(ctx, 5)
	// call cycle with 5 iterations, but due to timeout only 2 will be returned

	for _, line := range lines {
		fmt.Println(line)
	}
}
