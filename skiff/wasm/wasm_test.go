package wasm

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/skiff-sh/api/go/skiff/plugin/v1alpha1"
	"github.com/skiff-sh/sdk-go/skiff"
	"github.com/skiff-sh/sdk-go/skiff/mocks/skiffmocks"
	"github.com/skiff-sh/sdk-go/skiff/pluginapi"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type WasmTestSuite struct {
	suite.Suite
}

func (w *WasmTestSuite) TestRunRequest() {
	type test struct {
		ExpectedExitCode pluginapi.ExitCode
		ExpectedResponse *v1alpha1.Response
		Given            *v1alpha1.Request
		Plugin           *skiffmocks.Plugin
		EnvVars          map[string]string
		Constructor      func() test
	}

	tests := map[string]test{
		"basic": {
			Constructor: func() test {
				plug := new(skiffmocks.Plugin)
				plug.EXPECT().WriteFile(mock.Anything, &v1alpha1.WriteFileRequest{}).Return(&v1alpha1.WriteFileResponse{}, nil)
				return test{
					ExpectedResponse: &v1alpha1.Response{WriteFile: &v1alpha1.WriteFileResponse{}},
					Given:            &v1alpha1.Request{WriteFile: &v1alpha1.WriteFileRequest{}},
					Plugin:           plug,
				}
			},
		},
		"handles panic": {
			Constructor: func() test {
				plug := new(skiffmocks.Plugin)
				plug.EXPECT().WriteFile(mock.Anything, &v1alpha1.WriteFileRequest{}).RunAndReturn(func(ctx *skiff.Context, req *v1alpha1.WriteFileRequest) (*v1alpha1.WriteFileResponse, error) {
					panic("panic!")
				})
				return test{
					ExpectedResponse: &v1alpha1.Response{Issues: []*v1alpha1.Issue{{Message: "runtime error: panic!", Level: v1alpha1.IssueLevel_LEVEL_ERROR}}},
					Given:            &v1alpha1.Request{WriteFile: &v1alpha1.WriteFileRequest{}},
					Plugin:           plug,
				}
			},
		},
	}

	for desc, v := range tests {
		w.Run(desc, func() {
			if v.Constructor != nil {
				v = v.Constructor()
			}

			for k, val := range v.EnvVars {
				_ = os.Setenv(k, val)
			}
			defer func() {
				for k := range v.EnvVars {
					_ = os.Unsetenv(k)
				}
			}()

			b, err := json.Marshal(v.Given)
			if !w.NoError(err) {
				return
			}

			out := bytes.NewBuffer(nil)
			code := runRequest(v.Plugin, bytes.NewBuffer(append(b, pluginapi.EnvVarMessageDelimiterDefaultValue)), out)
			if !w.Equal(v.ExpectedExitCode, code) {
				return
			}

			actual := new(v1alpha1.Response)
			err = json.Unmarshal(out.Bytes()[:len(out.Bytes())-1], actual)
			if !w.NoError(err) {
				return
			}

			w.Equal(v.ExpectedResponse, actual)
		})
	}
}

func TestWasmTestSuite(t *testing.T) {
	suite.Run(t, new(WasmTestSuite))
}
