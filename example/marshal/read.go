package marshal

import (
	"fmt"
	"os"

	"github.com/zeabur/action"
	actionProto "github.com/zeabur/action/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

func Read(file string) (zbaction.Action, error) {
	// read file
	marshaled, err := os.ReadFile(file)
	if err != nil {
		return zbaction.Action{}, fmt.Errorf("read file: %w", err)
	}

	// turn it to protobuf
	actionPb := &actionProto.Action{}
	if err := protojson.Unmarshal(marshaled, actionPb); err != nil {
		return zbaction.Action{}, fmt.Errorf("unmarshal action: %w", err)
	}

	// turn it to zbaction.Action
	action, err := zbaction.ActionFromProto(actionPb)
	if err != nil {
		return zbaction.Action{}, fmt.Errorf("convert action: %w", err)
	}

	return action, nil
}
