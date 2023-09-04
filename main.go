package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/makarski/progress-bot/cmd"
)

func main() {
	maxAttempts, attempt := 5, 0
	initialDelay := 1 * time.Second

	slog.Info("running report")
	for {
		err := cmd.Run()
		if err == nil {
			break
		}

		attempt++
		if attempt >= maxAttempts {
			slog.Error(fmt.Sprintf("failed to run command: %s", err))
			os.Exit(1)
		}

		// exponential backoff
		backoffTime := initialDelay * (1 << uint(attempt))
		slog.Error(fmt.Sprintf("failed to run command: %s. Retrying in %s", err, backoffTime))
		time.Sleep(backoffTime)
	}

	slog.Info("report completed")
}
