package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/alecthomas/template"
)

type testCase struct {
	Plugin         string
	ExpectedOutput string
	ExecutedCode   string
}

func TestWillGenerateComplexPluginWithoutErrors(t *testing.T) {

	testCases := []testCase{
		{
			Plugin:         "complex_plugin",
			ExecutedCode:   "fmt.Println(pl)",
			ExpectedOutput: "",
		},
		{
			Plugin:         "basic_plugin",
			ExecutedCode:   "fmt.Println(pl.ReturningIntArray())",
			ExpectedOutput: "[1 0 1]",
		},
	}

	for i, testCase := range testCases {
		t.Logf("[Test %d] Generating %s plugin...", i, testCase.Plugin)

		config := Config{
			PluginPackage:      fmt.Sprintf("../internal/test_fixtures/%s", testCase.Plugin),
			OutputPath:         fmt.Sprintf("../internal/test_fixtures/generated/%s/plugin.go", testCase.Plugin),
			PluginPath:         fmt.Sprintf("../internal/test_fixtures/generated/%s/plugin.so", testCase.Plugin),
			FormatCode:         true,
			CheckSha256:        true,
			ForcePluginRebuild: true,
			OutputPackage:      "main",
			OutputName:         "TestWrapper",
		}

		client, err := New(config, log.New(ioutil.Discard, "", 0))
		if err != nil {
			t.Fatalf("[Test %d] Expected err to be nil, actual: %s", i, err)
		}

		if generateErr := client.GenerateFile(); err != nil {
			t.Fatalf("[Test %d] Expected err to be nil, actual: %s", i, generateErr)
		}

		runFile := fmt.Sprintf("../internal/test_fixtures/generated/%s/plugin.go", testCase.Plugin)

		t.Logf("[Test %d] Running plugin via %s", i, runFile)
		output, err := runPlugin(testCase.ExecutedCode, runFile, config)

		if err != nil {
			t.Fatalf("[Test %d] Expected err to be nil, actual: %s", i, err)
		}

		if !strings.Contains(output, testCase.ExpectedOutput) {
			t.Fatalf("[Test %d] Expected output to contain %s, actual output:\n=======\n%s\n=======\n", i, testCase.ExpectedOutput, output)
		}
	}
}

func runPlugin(code string, path string, config Config) (string, error) {
	file, err := os.OpenFile(config.OutputPath, os.O_APPEND|os.O_WRONLY, 0700)

	if err != nil {
		return "", err
	}

	tmp, err := template.New("test_case").Parse(runTemplate)
	if err != nil {
		return "", err
	}

	if err := tmp.Execute(file, struct {
		Config Config
		Code   string
	}{
		Config: config,
		Code:   code,
	}); err != nil {
		return "", err
	}

	var outBuffer bytes.Buffer

	cmd := exec.Command("go", "run", config.OutputPath)
	cmd.Stdout = bufio.NewWriter(&outBuffer)
	cmd.Stderr = os.Stderr

	runErr := cmd.Run()
	return string(outBuffer.Bytes()), runErr
}

var runTemplate = `

func main() {
  pl, err := Bind{{.Config.OutputName}}("{{.Config.PluginPath}}")

  if err != nil {
    fmt.Println(err)
    os.Exit(1)
  }

  {{.Code}}
}
`
