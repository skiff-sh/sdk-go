package pluginapi

const (
	// WASMFuncHandleRequestName the name of the exported WASM func to handle requests.
	WASMFuncHandleRequestName = "handleRequest"
)

const (
	// EnvVarCWD the absolute path to the current working directory. Only set if the user provides the permission.
	EnvVarCWD = "__CWD"

	// EnvVarCWDHost the absolute path to the current working directory on the host machine. Only set if the user provides the permission. Useful for more informative logging or error messages.
	EnvVarCWDHost = "__CWD_HOST"

	// EnvVarMessageDelimiter the delimiter for input/output messages read from stdin or stdout. Every message you read or write must end with this string
	EnvVarMessageDelimiter = "__MESSAGE_DELIM"

	// EnvVarMessageDelimiterDefaultValue is the default value for EnvVarMessageDelimiter.
	EnvVarMessageDelimiterDefaultValue = '\n'
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
	case ExitCodeCWDHostPathMissing:
		return EnvVarCWD + " was set but not " + EnvVarCWDHost
	}
	return ""
}

const (
	ExitCodeOK ExitCode = iota
	ExitCodePluginNotRegistered
	ExitCodeFailedToReadRequest
	ExitCodeFailedToUnmarshalRequest
	ExitCodePluginErr
	ExitCodeFailedToMarshalResponse
	ExitCodeFailedToWriteResponse
	ExitCodeMessageDelimInvalid
	ExitCodeCWDHostPathMissing
)
