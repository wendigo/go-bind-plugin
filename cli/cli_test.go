package cli_test

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/alecthomas/template"
	"github.com/wendigo/go-bind-plugin/cli"
)

type testCase struct {
	Plugin         string
	ExpectedOutput string
	ExecutedCode   string
	AsInterface    bool
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
		{
			Plugin:         "plugin_as_interface",
			ExecutedCode:   "fmt.Println(pl.ReturningStringSlice())",
			ExpectedOutput: "hello world",
			AsInterface:    true,
		},
	}

	for i, testCase := range testCases {
		t.Logf("[Test %d] Generating %s plugin...", i, testCase.Plugin)

		config := cli.Config{
			PluginPackage:        fmt.Sprintf("./internal/test_fixtures/%s", testCase.Plugin),
			OutputPath:           fmt.Sprintf("./internal/test_fixtures/generated/%s/plugin.go", testCase.Plugin),
			PluginPath:           fmt.Sprintf("./internal/test_fixtures/generated/%s/plugin.so", testCase.Plugin),
			FormatCode:           true,
			CheckSha256:          true,
			ForcePluginRebuild:   true,
			OutputPackage:        "main",
			OutputName:           "TestWrapper",
			AsInterface:          testCase.AsInterface,
			DereferenceVariables: true,
		}

		t.Logf("[Test %d] Generator config: %+v", i, config)

		if err := generatePluginWithCli(config, t); err != nil {
			t.Fatalf("[Test %d] Expected error to be nil, actual: %s", i, err)
		}

		runFile := fmt.Sprintf("./internal/test_fixtures/generated/%s/plugin.go", testCase.Plugin)

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

// Switch to generatePluginWithCli when https://github.com/golang/go/issues/17928 is solved
func generatePluginWithCli(config cli.Config, t *testing.T) error {
	client, err := cli.New(config, log.New(os.Stdout, "", 0))

	if err != nil {
		return err
	}

	if generateErr := client.GenerateFile(); generateErr != nil {
		return generateErr
	}

	return nil
}

func generatePluginViaCommandLine(config cli.Config, t *testing.T) error {

	args := []string{"run", "../main.go"}
	args = append(args, strings.Split(config.String(), " ")...)
	cmd := exec.Command("go", args...)

	t.Logf("Generating plugin with config: %+v", cmd)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func runPlugin(code string, path string, config cli.Config) (string, error) {
	file, err := os.OpenFile(config.OutputPath, os.O_APPEND|os.O_WRONLY, 0700)

	if err != nil {
		return "", err
	}

	tmp, err := template.New("test_case").Parse(runTemplate)
	if err != nil {
		return "", err
	}

	if err := tmp.Execute(file, struct {
		Config cli.Config
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
