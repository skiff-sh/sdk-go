package wasm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime/debug"

	"github.com/skiff-sh/api/go/skiff/plugin/v1alpha1"
	"github.com/skiff-sh/sdk-go/skiff"
	"github.com/skiff-sh/sdk-go/skiff/issue"
	"github.com/skiff-sh/sdk-go/skiff/pluginapi"
)

//go:wasmexport handleRequest
func handleRequest() uint64 {
	stdout := os.Stdout
	// All log statements must go to stderr. Stdout is for returning the response of the call.
	os.Stdout = os.Stderr
	return uint64(runRequest(skiff.GetPlugin(), os.Stdin, stdout))
}

func runRequest(plug skiff.Plugin, in io.Reader, out io.Writer) pluginapi.ExitCode {
	if plug == nil {
		return pluginapi.ExitCodePluginNotRegistered
	}

	logger, _ := newLogger(logSpec{Outputs: []string{"stderr"}})
	evs, code := parseEnvVars()
	if code != pluginapi.ExitCodeOK {
		return code
	}

	logger.Info("Parsing request.")
	req, code := parseRequest(in, evs.MessageDelim)
	if code != pluginapi.ExitCodeOK {
		return code
	}
	if req == nil {
		return code
	}

	var cwd *skiff.VolumeMount
	if evs.CWDPath != "" {
		cwd = &skiff.VolumeMount{
			FS:       os.DirFS(evs.CWDPath),
			HostPath: evs.CWDHostPath,
		}
	}
	ctx := &skiff.Context{
		Ctx:      context.Background(),
		CWD:      cwd,
		Data:     req.Data,
		Metadata: req.Metadata,
	}

	logger.Info("Handling request.")
	resp, err := runPlugin(ctx, plug, req)
	if err != nil {
		logger.Error("Failed to handle request.", "err", err.Error())
		resp = &v1alpha1.Response{Issues: issues(err)}
	}

	logger.Info("Returning response.")
	return writeResponse(out, evs.MessageDelim, resp)
}

func runPlugin(ctx *skiff.Context, plug skiff.Plugin, req *v1alpha1.Request) (resp *v1alpha1.Response, err error) {
	resp = &v1alpha1.Response{}
	defer func() {
		if recovered := recover(); recovered != nil {
			slog.Error("Panic occurred.", "panic", recovered)
			fmt.Println("-- Stack trace --")
			fmt.Println(string(debug.Stack()))
			err = fmt.Errorf("runtime error: %v", recovered)
		}
	}()
	if req.WriteFile != nil {
		resp.WriteFile, err = plug.WriteFile(ctx, req.WriteFile)
	}
	return resp, err
}

func issues(err error) []*v1alpha1.Issue {
	switch typ := err.(type) {
	case issue.PluginIssue:
		if iss := typ.Issue(); iss != nil {
			return []*v1alpha1.Issue{iss}
		}
	case interface{ Unwrap() []error }:
		errs := typ.Unwrap()
		out := make([]*v1alpha1.Issue, 0, len(errs))
		for _, v := range errs {
			out = append(out, issues(v)...)
		}
	default:
		return []*v1alpha1.Issue{{
			Level:   v1alpha1.IssueLevel_LEVEL_ERROR,
			Message: err.Error(),
		}}
	}
	return []*v1alpha1.Issue{}
}

func writeResponse(writer io.Writer, delim byte, resp *v1alpha1.Response) pluginapi.ExitCode {
	raw, err := json.Marshal(resp)
	if err != nil {
		slog.Error("Failed to marshal response.", "err", err.Error())
		return pluginapi.ExitCodeFailedToMarshalResponse
	}

	_, err = io.Copy(writer, bytes.NewBuffer(append(raw, delim)))
	if err != nil {
		slog.Error("Failed to copy byte buffer for response.", "err", err.Error())
		return pluginapi.ExitCodeFailedToWriteResponse
	}

	return pluginapi.ExitCodeOK
}

type envVars struct {
	MessageDelim byte
	CWDPath      string
	CWDHostPath  string
	LogLevel     string
}

func parseEnvVars() (*envVars, pluginapi.ExitCode) {
	out := &envVars{
		LogLevel: os.Getenv(pluginapi.EnvVarLogLevel),
	}
	msgDelim, ok := os.LookupEnv(pluginapi.EnvVarMessageDelimiter)
	if !ok {
		out.MessageDelim = pluginapi.EnvVarMessageDelimiterDefaultValue
	} else {
		if len(msgDelim) != 1 {
			return nil, pluginapi.ExitCodeMessageDelimInvalid
		}
		out.MessageDelim = msgDelim[0]
	}

	out.CWDPath, ok = os.LookupEnv(pluginapi.EnvVarCWD)
	if ok {
		out.CWDHostPath, ok = os.LookupEnv(pluginapi.EnvVarCWDHost)
		if !ok {
			return nil, pluginapi.ExitCodeCWDHostPathMissing
		}
	}

	return out, pluginapi.ExitCodeOK
}

func parseRequest(reader io.Reader, delim byte) (*v1alpha1.Request, pluginapi.ExitCode) {
	b, err := bufio.NewReader(reader).ReadBytes(delim)
	if err != nil {
		slog.Error("Failed to read request.", "err", err.Error())
		return nil, pluginapi.ExitCodeFailedToReadRequest
	}
	if len(b) == 0 {
		slog.Info("Received an empty request. Returning.")
		return nil, pluginapi.ExitCodeOK
	}

	// Drop the delimiter
	b = b[:len(b)-1]

	req := new(v1alpha1.Request)
	err = json.Unmarshal(b, req)
	if err != nil {
		slog.Error("Failed to unmarshal request.", "err", err.Error())
		return nil, pluginapi.ExitCodeFailedToUnmarshalRequest
	}

	return req, pluginapi.ExitCodeOK
}

// logSpec represents logging config.
type logSpec struct {
	Level string
	// Valid values are:
	// * stdout
	// * stderr
	// * fullfile path
	Outputs []string
}

// Copied from the config library. Want to avoid external dependency to reduce overall binary size.
func newLogger(log logSpec) (*slog.Logger, error) {
	w := make([]io.Writer, 0, len(log.Outputs))
	for _, v := range log.Outputs {
		switch v {
		case "stdout":
			w = append(w, os.Stdout)
		case "stderr":
			w = append(w, os.Stderr)
		default:
			f, err := os.OpenFile(v, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
			if err != nil {
				return nil, err
			}
			w = append(w, f)
		}
	}

	logger := slog.New(slog.NewJSONHandler(io.MultiWriter(w...), &slog.HandlerOptions{
		AddSource: true,
		Level:     parseLevel(log.Level),
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				source, _ := a.Value.Any().(*slog.Source)
				if source != nil {
					source.Function = ""
					source.File = filepath.Base(source.File)
				}
			}
			return a
		},
	}))

	slog.SetDefault(logger)

	return logger, nil
}

var validLogLevels = map[string]slog.Level{
	"info":    slog.LevelInfo,
	"debug":   slog.LevelDebug,
	"warn":    slog.LevelWarn,
	"warning": slog.LevelWarn,
	"error":   slog.LevelError,
	"err":     slog.LevelError,
}

func parseLevel(lvl string) slog.Level {
	o, ok := validLogLevels[lvl]
	if !ok {
		return slog.LevelInfo
	}
	return o
}
