package artifact

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"path/filepath"

	"github.com/moby/buildkit/client"
	dockerfile "github.com/moby/buildkit/frontend/dockerfile/builder"
	"github.com/moby/buildkit/util/progress/progressui"
	"github.com/nwtgck/go-fakelish"
	"github.com/zeabur/builder/zbaction"
	"golang.org/x/sync/errgroup"
)

func init() {
	zbaction.RegisterProcedure("action/artifact/docker", func(args zbaction.ProcStepArgs) (zbaction.ProcedureStep, error) {
		tag, ok := args["tag"]
		if !ok {
			key := fakelish.GenerateFakeWord(12, 36)
			tag = "zeabur/built-resource-" + key + ":latest"
		}

		contextInput, ok := args["context"]
		if !ok {
			return nil, zbaction.NewErrRequiredArgument("context")
		}

		dockerfileInput, ok := args["dockerfile"]
		if !ok {
			dockerfileInput = `FROM docker.io/library/alpine:latest
COPY . .`
		}

		cache, ok := args["cache"]
		if !ok {
			cache = "true"
		}

		return &DockerArtifactAction{
			Tag:        zbaction.NewArgumentStr(tag),
			Context:    zbaction.NewArgumentStr(contextInput),
			Dockerfile: zbaction.NewArgumentStr(dockerfileInput),
			Cache:      zbaction.NewArgumentBool(cache),
		}, nil
	})
}

type DockerArtifactAction struct {
	// Tag is the tag of this artifact.
	Tag zbaction.Argument[string]
	// Context is the directory to run the build in.
	Context zbaction.Argument[string]
	// Dockerfile is the content of the Dockerfile for runtime.
	Dockerfile zbaction.Argument[string]
	// Cache indicates whether to use cache when building the image.
	// By default, it is true.
	Cache zbaction.Argument[bool]
}

var buildKitAddress = os.Getenv("BUILDKIT_HOST")

func (d DockerArtifactAction) Run(ctx context.Context, sc *zbaction.StepContext) (zbaction.CleanupFn, error) {
	contextDirectory := d.Context.Value(sc.ExpandString)
	tag := d.Tag.Value(sc.ExpandString)
	dockerFileContent := d.Dockerfile.Value(sc.ExpandString)
	cache := d.Cache.Value(sc.ExpandString)

	cleanupStack := zbaction.CleanupStack{}
	cleanupFn := cleanupStack.WrapRun()

	if buildKitAddress == "" {
		return cleanupFn, errors.New("BUILDKIT_HOST is not set")
	}

	builderTmpDir, err := os.MkdirTemp("", "zbaction-docker-builder-*")
	if err != nil {
		return cleanupFn, fmt.Errorf("make temp dir: %w", err)
	}
	cleanupStack.Push(func() {
		_ = os.RemoveAll(builderTmpDir)
	})

	// Write the users' dockerfile to the dockerfile path.
	dockerfilePath := path.Join(builderTmpDir, "Dockerfile")
	err = os.WriteFile(dockerfilePath, []byte(dockerFileContent), 0644)
	if err != nil {
		return cleanupFn, fmt.Errorf("write dockerfile: %w", err)
	}

	// Open a docker.tar which is the artifact of this action.
	artifactTarPath := path.Join(builderTmpDir, "docker.tar")
	artifactTar, err := os.OpenFile(artifactTarPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return cleanupFn, fmt.Errorf("open artifact tar: %w", err)
	}

	c, err := client.New(ctx, buildKitAddress)
	if err != nil {
		return cleanupFn, fmt.Errorf("connect to buildkit: %w", err)
	}

	// FIXME: build-args
	frontendAttrs := map[string]string{
		"filename": filepath.Base(dockerfilePath),
	}

	if !cache {
		frontendAttrs["no-cache"] = ""
	}

	solveOpt := client.SolveOpt{
		Exports: []client.ExportEntry{
			{
				Type: client.ExporterDocker,
				Attrs: map[string]string{
					"name": tag,
				},
				Output: func(_ map[string]string) (io.WriteCloser, error) {
					return artifactTar, nil
				},
			},
		},
		LocalDirs: map[string]string{
			"context":    contextDirectory,
			"dockerfile": filepath.Dir(dockerfilePath),
		},
		Frontend:      "dockerfile.v0",
		FrontendAttrs: frontendAttrs,
	}

	ch := make(chan *client.SolveStatus)
	eg, ectx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		_, err := c.Build(ectx, solveOpt, "zbaction-docker-builder", dockerfile.Build, ch)
		return err
	})
	eg.Go(func() error {
		d, err := progressui.NewDisplay(os.Stderr, progressui.AutoMode)
		_, err = d.UpdateFrom(context.TODO(), ch)
		return err
	})

	if err := eg.Wait(); err != nil {
		return cleanupFn, err
	}

	sc.SetThisOutput("tag", tag)
	sc.SetThisOutput("context", contextDirectory)
	sc.SetThisOutput("dockerfile", dockerfilePath)
	sc.SetThisOutput("artifact", artifactTarPath)

	slog.Info("Docker image built",
		slog.String("tag", tag),
		slog.String("context", contextDirectory),
		slog.String("dockerfile", dockerfilePath),
		slog.String("artifact", artifactTarPath))

	return nil, nil // FIXME: cleanupFn, nil
}

var _ zbaction.ProcedureStep = (*DockerArtifactAction)(nil)
