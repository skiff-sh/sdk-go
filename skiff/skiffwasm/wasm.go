package skiffwasm

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os"

	"github.com/skiff-sh/api/go/skiff/plugin/v1alpha1"
	"github.com/skiff-sh/sdk-go/skiff"
	"google.golang.org/protobuf/proto"
)

const (
	// WASMFuncHandleRequestName the name of the exported WASM func to handle requests.
	WASMFuncHandleRequestName = "handleRequest"
)

const (
	// EnvVarProjectPath the absolute path to the root of the project
	EnvVarProjectPath = "__PROJECT_PATH"
	// EnvVarMessageDelimiter the delimiter for input/output messages read from stdin or stdout. Every message you read or write must end with this string
	EnvVarMessageDelimiter = "__MESSAGE_DELIM"
	// EnvVarMessageDelimiterDefaultValue is the default value for EnvVarMessageDelimiter.
	EnvVarMessageDelimiterDefaultValue = '\n'
)

// RequestType the input to distinguish which request is being targeted.
type RequestType uint64

const (
	RequestTypeNone RequestType = iota
	RequestTypeWriteFile
)

type ExitCode uint64

func (e ExitCode) String() string {
	switch e {
	case ExitCodeOK:
		return "ok"
	case ExitCodePluginNotRegistered:
		return "no plugin registered"
	case ExitCodeFailedToReadRequest:
		return "failed to read request"
	case ExitCodeFailedToUnmarshalRequest:
		return "failed to unmarshal request"
	case ExitCodePluginErr:
		return "plugin error"
	case ExitCodeFailedToMarshalResponse:
		return "failed to marshal response"
	case ExitCodeFailedToWriteResponse:
		return "failed to write response"
	case ExitCodeMessageDelimInvalid:
		return "message delimiter must be a single byte"
	case ExitCodeRootPathEnvVarMissing:
		return EnvVarProjectPath + " env var missing"
	}
	return ""
}

const (
	ExitCodeOK ExitCode = iota
	ExitCodeRootPathEnvVarMissing
	ExitCodePluginNotRegistered
	ExitCodeFailedToReadRequest
	ExitCodeFailedToUnmarshalRequest
	ExitCodePluginErr
	ExitCodeFailedToMarshalResponse
	ExitCodeFailedToWriteResponse
	ExitCodeMessageDelimInvalid
)

//go:wasmexport handleRequest
func handleRequest(typ uint64) uint64 {
	evs, code := parseEnvVars()
	if code != ExitCodeOK {
		return uint64(code)
	}

	root := os.DirFS(evs.RootPath)

	ctx := &skiff.Context{
		Ctx:  context.Background(),
		Root: root,
	}

	var resp proto.Message
	var err error
	switch RequestType(typ) {
	case RequestTypeWriteFile:
		req := new(v1alpha1.WriteFileRequest)
		code = parseRequest(os.Stdin, evs.MessageDelim, req)
		if code != ExitCodeOK {
			return uint64(code)
		}
		resp, err = skiff.GetPlugin().WriteFile(ctx, req)
	default:
		return uint64(ExitCodeOK)
	}
	if err != nil {
		return uint64(ExitCodePluginErr)
	}

	return uint64(writeResponse(os.Stdout, evs.MessageDelim, resp))
}

func writeResponse(writer io.Writer, delim byte, resp proto.Message) ExitCode {
	raw, err := proto.Marshal(resp)
	if err != nil {
		return ExitCodeFailedToMarshalResponse
	}

	_, err = io.Copy(writer, bytes.NewBuffer(append(raw, delim)))
	if err != nil {
		return ExitCodeFailedToWriteResponse
	}

	return ExitCodeOK
}

type envVars struct {
	MessageDelim byte
	RootPath     string
}

func parseEnvVars() (*envVars, ExitCode) {
	out := &envVars{}
	msgDelim, ok := os.LookupEnv(EnvVarMessageDelimiter)
	if !ok {
		out.MessageDelim = EnvVarMessageDelimiterDefaultValue
	} else {
		if len(msgDelim) != 1 {
			return nil, ExitCodeMessageDelimInvalid
		}
		out.MessageDelim = msgDelim[0]
	}

	out.RootPath, ok = os.LookupEnv(EnvVarProjectPath)
	if !ok {
		return nil, ExitCodeRootPathEnvVarMissing
	}

	return out, ExitCodeOK
}

func parseRequest(reader io.Reader, delim byte, msg proto.Message) ExitCode {
	if skiff.GetPlugin() == nil {
		return ExitCodePluginNotRegistered
	}

	b, err := bufio.NewReader(reader).ReadBytes(delim)
	if err != nil {
		return ExitCodeFailedToReadRequest
	}

	err = proto.Unmarshal(b, msg)
	if err != nil {
		return ExitCodeFailedToUnmarshalRequest
	}

	return ExitCodeOK
}
