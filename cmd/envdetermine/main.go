package main

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/zeabur/action/environment"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	sl := environment.DetermineSoftwareList()
	result, err := json.Marshal(sl)
	if err != nil {
		panic(err)
	}

	_, _ = os.Stdout.Write(result)
}
