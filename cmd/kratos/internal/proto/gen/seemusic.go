package gen

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/SeeMusic/kratos/cmd/kratos/v2/internal/base"
	"github.com/spf13/cobra"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var long = `
Generate the proto Golang code. You can use this command anywhere 
inside of your project, kratos will find the proto file you want.

example:
	kratos proto gen --all

this command will find all proto files in your project, and let 
you choose which proto file you want to generate. then, kratos 
will generate go, grpc, http, openapi, kratos error code for you.

you can use options to generate specific code you want.

example:
	kratos proto gen -e tag* --http --grpc 
	
this will find all proto files that matches 'tag*.proto$', and let
you choose too. and then will generate go, http, grpc code for 
you
`

// CmdGen represents the source command.
var CmdGen = &cobra.Command{
	Use:   "gen",
	Short: "Generate the proto Golang code",
	Long:  long,
	Run:   seeMusic,
}

func init() {
	CmdGen.Flags().StringVarP(&expr, "exp", "e", expr, "use regexp")
	CmdGen.Flags().StringVarP(&dir, "dir", "d", dir, "dir to search protobuf files")

	CmdGen.Flags().BoolVar(&genAll, "all", genAll, "generate all code")
	CmdGen.Flags().BoolVar(&genOpenapi, "openapi", genOpenapi, "generate openapi code")
	CmdGen.Flags().BoolVar(&genError, "error", genError, "generate kratos error code")
	CmdGen.Flags().BoolVar(&genHttp, "http", genHttp, "generate http code")
	CmdGen.Flags().BoolVar(&genGrpc, "grpc", genGrpc, "generate grpc code")

	CmdGen.Flags().BoolVarP(&verbose, "verbose", "v", verbose, "show more information")
}

// commands
var (
	expr       string
	dir        string
	verbose    = false
	genGrpc    = true
	genHttp    = false
	genOpenapi = false
	genError   = false
	genAll     = false
)

func seeMusic(cmd *cobra.Command, args []string) {
	fmt.Printf("expr: %s", expr)
	fmt.Printf("dir: %s", dir)
	fmt.Printf("verbose: %v", verbose)
	fmt.Printf("genGrpc: %v", genGrpc)
	fmt.Printf("genHttp: %v", genHttp)
	fmt.Printf("genOpenapi: %v", genOpenapi)
	fmt.Printf("genError: %v", genError)
	fmt.Printf("genAll: %v", genAll)
	var (
		baseDir    string
		currentDir string
		err        error
	)

	if expr == "" {
		expr = ".*.proto$"
	} else {
		expr = expr + ".proto$"
	}

	baseDir, err = modDir()
	if err != nil {
		fmt.Println(err)
		return
	}
	currentDir, err = os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	if dir == "" {
		dir = baseDir
	} else {
		dir = filepath.Join(baseDir, dir)
	}

	protoFiles, err := findProto(currentDir, dir, expr)

	var target string

	op := survey.Select{
		Message: "which proto file you want to generate?",
		Options: protoFiles,
	}
	if err = survey.AskOne(&op, &target); err != nil {
		fmt.Println(err)
		return
	}

	if err = gen(baseDir, currentDir, target, args); err != nil {
		fmt.Println(err)
	}
}

func findProto(curDir, dir string, expr string) ([]string, error) {
	reg, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}

	var protoFiles []string
	err = filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if strings.Contains(path, "third_party") {
			return nil
		}
		if reg.MatchString(path) {
			p, _ := filepath.Rel(curDir, path)
			protoFiles = append(protoFiles, p)
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	return protoFiles, nil
}

func modDir() (string, error) {
	mod, err := base.FindModulePath(".")
	if err != nil {
		return "", err
	}
	return filepath.Dir(mod), nil
}

// generate is used to execute the generate command for the specified proto file
func gen(baseDir, curDir, proto string, args []string) error {
	thirdParty, err := filepath.Rel(curDir, filepath.Join(baseDir, "third_party"))
	if err != nil {
		return err
	}
	proto, err = filepath.Abs(proto)
	if err != nil {
		return err
	}
	proto, err = filepath.Rel(curDir, proto)
	if err != nil {
		return err
	}

	inputDir := filepath.Dir(proto)
	var input = []string{"--proto_path=" + inputDir}

	if genAll {
		inputExt := []string{
			"--proto_path=" + thirdParty,
			"--go_out=paths=source_relative:" + inputDir,
			"--go-grpc_out=paths=source_relative:" + inputDir,
			"--go-http_out=paths=source_relative:" + inputDir,
			"--go-errors_out=paths=source_relative:" + inputDir,
			"--openapi_out=paths=source_relative:" + inputDir,
		}
		input = append(input, inputExt...)
	} else {
		if genGrpc {
			input = append(input, "--go-grpc_out=paths=source_relative:"+inputDir)
		}
		if genHttp {
			input = append(input, "--go-http_out=paths=source_relative:"+inputDir)
		}
		if genError {
			input = append(input, "--go-errors_out=paths=source_relative:"+inputDir)
		}
		if genOpenapi {
			input = append(input, "--openapi_out=paths=source_relative:"+inputDir)
		}
	}

	protoBytes, err := os.ReadFile(proto)
	if err == nil && len(protoBytes) > 0 {
		if ok, _ := regexp.Match(`\n[^/]*(import)\s+"validate/validate.proto"`, protoBytes); ok {
			input = append(input, "--validate_out=lang=go,paths=source_relative:"+inputDir)
		}
	}
	input = append(input, proto)
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			input = append(input, a)
		}
	}
	fd := exec.Command("protoc", input...)
	fd.Stdout = os.Stdout
	fd.Stderr = os.Stderr
	fd.Dir = "."
	if verbose {
		fmt.Println(fd.String())
	}
	if err := fd.Run(); err != nil {
		return err
	}
	fmt.Printf("proto: %s\n", proto)
	return nil
}
