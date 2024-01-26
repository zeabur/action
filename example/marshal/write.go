package marshal

import (
	"fmt"
	"os"

	"github.com/zeabur/builder/zbaction"
	"google.golang.org/protobuf/encoding/protojson"
)

func Write(action zbaction.Action, file string) error {
	// turn it to protobuf
	actionPb, err := zbaction.ActionToProto(action)
	if err != nil {
		return fmt.Errorf("convert action: %w", err)
	}

	// write to file
	marshaled, err := protojson.Marshal(actionPb)
	if err != nil {
		return fmt.Errorf("marshal action: %w", err)
	}

	if err := os.WriteFile(file, marshaled, 0644); err != nil {
		return fmt.Errorf("write action to file: %w", err)
	}

	return nil
}
