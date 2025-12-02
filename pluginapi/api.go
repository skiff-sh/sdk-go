package pluginapi

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
type RequestType int64

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
