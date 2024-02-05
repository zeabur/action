package main

import (
	"context"
	"flag"
	"io"
	"log/slog"
	"os"

	zbaction "github.com/zeabur/action"
	zbactionpb "github.com/zeabur/action/proto"
	"google.golang.org/protobuf/proto"

	// modules
	_ "github.com/zeabur/action/procedures"
	_ "github.com/zeabur/action/procedures/artifact"
	_ "github.com/zeabur/action/procedures/golang"
)

var file = flag.String("file", "-", "The filename of the action")

func main() {
	flag.Parse()

	// read from file
	slog.Info("read file", slog.String("file", *file))
	f, err := readFileOrStdin(*file)
	if err != nil {
		slog.Error("failed to read file", slog.String("file", *file), slog.String("error", err.Error()))
		os.Exit(1)
	}

	var actionPb zbactionpb.Action
	if err := proto.Unmarshal(f, &actionPb); err != nil {
		slog.Error("failed to unmarshal action", slog.String("file", *file), slog.String("error", err.Error()))
		os.Exit(1)
	}

	action, err := zbaction.ActionFromProto(&actionPb)
	if err != nil {
		slog.Error("failed to convert action from proto", slog.String("file", *file), slog.String("error", err.Error()))
		os.Exit(1)
	}

	err = zbaction.RunAction(context.Background(), action)
	if err != nil {
		slog.Error("failed to run action", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.Info("action completed successfully")
}

func readFileOrStdin(filename string) (data []byte, err error) {
	if filename == "-" {
		data, err = io.ReadAll(os.Stdin)
		return
	}

	data, err = os.ReadFile(filename)
	return
}
