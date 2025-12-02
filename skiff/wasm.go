package skiff

import (
	"os"

	"github.com/skiff-sh/sdk-go/pluginapi"
)

//go:wasmexport handleRequest
func handleRequest(typ uint64) uint64 {
	//evs, code := parseEnvVars()
	//if code != pluginapi.ExitCodeOK {
	//	return uint64(code)
	//}
	//
	//root := os.DirFS(evs.RootPath)
	//
	////ctx := &Context{
	////	Ctx:  context.Background(),
	////	Root: root,
	////}
	//
	//var resp proto.Message
	//var err error
	//switch pluginapi.RequestType(typ) {
	//case pluginapi.RequestTypeWriteFile:
	//	//req := new(v1alpha1.WriteFileRequest)
	//	//code = parseRequest(os.Stdin, evs.MessageDelim, req)
	//	if code != pluginapi.ExitCodeOK {
	//		return uint64(code)
	//	}
	//	//resp, err = plugin.WriteFile(ctx, req)
	//default:
	//	return uint64(pluginapi.ExitCodeOK)
	//}
	//if err != nil {
	//	return uint64(pluginapi.ExitCodePluginErr)
	//}

	//return uint64(writeResponse(os.Stdout, evs.MessageDelim, resp))
	return 0
}

//func writeResponse(writer io.Writer, delim byte, resp proto.Message) pluginapi.ExitCode {
//	raw, err := proto.Marshal(resp)
//	if err != nil {
//		return pluginapi.ExitCodeFailedToMarshalResponse
//	}
//
//	_, err = io.Copy(writer, bytes.NewBuffer(append(raw, delim)))
//	if err != nil {
//		return pluginapi.ExitCodeFailedToWriteResponse
//	}
//
//	return pluginapi.ExitCodeOK
//}

type envVars struct {
	MessageDelim byte
	RootPath     string
}

func parseEnvVars() (*envVars, pluginapi.ExitCode) {
	out := &envVars{}
	msgDelim, ok := os.LookupEnv(pluginapi.EnvVarMessageDelimiter)
	if !ok {
		out.MessageDelim = pluginapi.EnvVarMessageDelimiterDefaultValue
	} else {
		if len(msgDelim) != 1 {
			return nil, pluginapi.ExitCodeMessageDelimInvalid
		}
		out.MessageDelim = msgDelim[0]
	}

	out.RootPath, ok = os.LookupEnv(pluginapi.EnvVarProjectPath)
	if !ok {
		return nil, pluginapi.ExitCodeRootPathEnvVarMissing
	}

	return out, pluginapi.ExitCodeOK
}

//func parseRequest(reader io.Reader, delim byte, msg proto.Message) pluginapi.ExitCode {
//	if plugin == nil {
//		return pluginapi.ExitCodePluginNotRegistered
//	}
//
//	b, err := bufio.NewReader(reader).ReadBytes(delim)
//	if err != nil {
//		return pluginapi.ExitCodeFailedToReadRequest
//	}
//	if len(b) == 0 {
//		return 0
//	}
//
//	// Drop the delimiter
//	b = b[:len(b)-1]
//
//	err = proto.Unmarshal(b, msg)
//	if err != nil {
//		fmt.Println(err.Error())
//		return pluginapi.ExitCodeFailedToUnmarshalRequest
//	}
//
//	return pluginapi.ExitCodeOK
//}
