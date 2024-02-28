package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	zbaction "github.com/zeabur/action"
	"github.com/zeabur/action/environment"
	"github.com/zeabur/action/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

type pair struct {
	Name string

	MetadataPath string
	Metadata     environment.SoftwareList
}

type metadata []pair

func (m *metadata) String() string {
	if m == nil {
		return ""
	}

	pairs := make([]string, 0, len(*m))

	for _, p := range *m {
		pairs = append(pairs, p.Name+":"+p.MetadataPath)
	}

	return strings.Join(pairs, ",")
}

func (m *metadata) Set(s string) error {
	before, after, found := strings.Cut(s, ":")
	if !found {
		return flag.ErrHelp
	}

	f, err := os.Open(after)
	if err != nil {
		return err
	}

	var md environment.SoftwareList
	err = json.NewDecoder(f).Decode(&md)
	if err != nil {
		return err
	}

	*m = append(*m, pair{
		Name:         before,
		Metadata:     md,
		MetadataPath: after,
	})
	return nil
}

func metadataFlag(name string, usage string) *metadata {
	f := &metadata{}
	flag.CommandLine.Var(f, name, usage)
	return f
}

var metadataFile = metadataFlag("metadata", "metadata name:metadata path")
var actionFile = flag.String("action", "", "action file")

func main() {
	flag.Parse()

	if *actionFile == "" {
		panic("action file is required")
	}

	actionPbRaw, err := os.ReadFile(*actionFile)
	if err != nil {
		panic(err)
	}

	var actionPb proto.Action
	err = protojson.Unmarshal(actionPbRaw, &actionPb)
	if err != nil {
		panic(err)
	}

	action, err := zbaction.ActionFromProto(&actionPb)
	if err != nil {
		panic(err)
	}

	compiledAction, err := zbaction.CompileActionRequirement(action)
	if err != nil {
		panic(err)
	}

	for _, metadata := range *metadataFile {
		fmt.Println("check requirement", metadata.Name)
		err := compiledAction.CheckRequirement(metadata.Metadata)
		if err != nil {
			fmt.Printf("%s's requirement not met: %s\n", metadata.Name, err.Error())
		}
		fmt.Printf("%s's requirement met\n", metadata.Name)
	}
}
