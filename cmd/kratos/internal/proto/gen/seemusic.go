package gen

import (
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/SeeMusic/kratos/cmd/kratos/v2/internal/base"
	"github.com/spf13/cobra"
)

var long = `
生成 protobuf 的 Golang 代码

默认为查找<项目目录>下所有的 proto 文件，由调用者选择需要生成的 proto 文件。
这个选项为多选，每个选中的 proto 文件都会单独执行一次 protoc 命令。

默认只会生成 <filename>.pb.go 文件，如果需要生成更多类型的 Golang 代码，
可以使用 --all / --grpc / --http / --error / --openapi

	example: kratos proto gen -e 'admin.*'

	会查找所有文件路径上带有 admin 的 proto 文件，然后你可以自由选择需要的
	proto 文件来生成对应的 Golang 代码

可以使用 -d / --dir 指定要生成 go 代码的 proto 文件所在的目录，用一条命令
生成 dir 下所有的 proto 文件的 go 代码。当同一个目录下有多个 proto 文件但
想把定义的 openapi 输出到一个 openapi.yaml 时很有用(这大概是唯一的方式)

	example: kratos proto gen -d api/core/admin/v1 --openapi

	这会在 api/core/admin/v1 目录下生成对应 proto 文件的 pb.go 文件和
	一个 openapi.yaml 文件
`

var CmdGen = &cobra.Command{
	Use:   "gen",
	Short: "Generate the proto Golang code",
	Long:  long,
	Run:   run,
}

func init() {
	log.SetFlags(log.LstdFlags)

	CmdGen.Flags().StringVarP(&expr, "expr", "e", expr, "使用正则表达式来匹配文件")
	CmdGen.Flags().StringVarP(&dir, "dir", "d", dir, "生成该文件夹下所有 proto 文件的 go 代码")
	CmdGen.Flags().BoolVar(&genAll, "all", genAll, "生成所有类型的 go 代码，包括 grpc, http, openapi, kratos error code")
	CmdGen.Flags().BoolVar(&genOpenapi, "openapi", genOpenapi, "生成 openapi.yaml 文件")
	CmdGen.Flags().BoolVar(&genError, "error", genError, "生成 kratos error code")
	CmdGen.Flags().BoolVar(&genHTTP, "http", genHTTP, "生成 http code")
	CmdGen.Flags().BoolVar(&genGrpc, "grpc", genGrpc, "生成 grpc code")
	CmdGen.Flags().BoolVarP(&verbose, "verbose", "v", verbose, "显示详细更多日志信息")
}

var (
	currentDir string
	rootDir    string
	expr       string
	dir        string
	verbose    bool
	genGrpc    bool
	genHTTP    bool
	genOpenapi bool
	genError   bool
	genAll     bool
)

func initDirs() {
	var (
		err    error
		modDir string
	)
	currentDir, err = os.Getwd()
	if err != nil {
		log.Fatalf("get current dir failed: %s\n", err)
	}

	modDir, err = base.FindModulePath(currentDir)
	if err != nil {
		log.Fatalf("find module path failed: %s\n", err)
	}
	rootDir = filepath.Dir(modDir)
}

func run(cmd *cobra.Command, args []string) {
	initDirs()

	if verbose {
		log.Printf("current dir: %s\n", currentDir)
		log.Printf("root dir: %s\n", rootDir)
	}

	protos, err := findProtos()
	if err != nil {
		log.Fatalf("find protos failed: %s\n", err)
	}

	if verbose {
		log.Printf("find %d proto files\n", len(protos))
	}

	// dir 单独处理
	if dir != "" {
		inputDir, err := filepath.Rel(currentDir, filepath.Join(rootDir, dir))
		if err != nil {
			log.Fatalf("get relative path failed: %s\n", err)
		}

		inputs, err := commandArgs(inputDir, protos...)
		if err != nil {
			log.Fatalf("get default proto path failed: %s\n", err)
		}
		gen(inputs)
		return
	}

	var targets []string

	// 只有一个 proto 文件时直接生成，不让用户选择了
	if len(protos) > 1 {
		q := &survey.MultiSelect{
			Message:  "📌 which protos do you want to generate?",
			Options:  protos,
			PageSize: 10,
		}

		err = survey.AskOne(q, &targets)
		if err != nil {
			log.Fatalf("ask proto failed: %s\n", err)
		}
	} else {
		targets = protos
	}

	for _, t := range targets {
		inputDir := filepath.Dir(t)
		inputs, err := commandArgs(inputDir, t)
		if err != nil {
			log.Fatalf("get default proto path failed: %s\n", err)
		}
		if err := gen(inputs); err != nil {
			log.Printf("gen proto failed: %s\n", err)
		}
	}
}

func getExpr() *regexp.Regexp {
	if expr == "" {
		return regexp.MustCompile(".*.proto$")
	}
	return regexp.MustCompile(expr + ".*.proto$")
}

func findProtos() ([]string, error) {
	lookPath := rootDir

	if dir != "" {
		lookPath = dir
		if !filepath.IsAbs(lookPath) {
			lookPath = filepath.Join(rootDir, lookPath)
		}
	}

	var err error
	reg := getExpr()

	var protoFiles []string
	err = filepath.Walk(lookPath, func(path string, info fs.FileInfo, err error) error {
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
			p, err := filepath.Rel(currentDir, path)
			if err != nil {
				return err
			}
			protoFiles = append(protoFiles, p)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
	}
	return protoFiles, nil
}

func commandArgs(inputDir string, target ...string) ([]string, error) {
	inputPath, err := filepath.Rel(currentDir, rootDir)
	if err != nil {
		return nil, err
	}

	thirdPartyInputPath, err := filepath.Rel(currentDir, filepath.Join(rootDir, "third_party"))
	if err != nil {
		return nil, err
	}

	args := []string{
		"--proto_path=" + inputPath,
		"--proto_path=" + thirdPartyInputPath,
		"--go_out=paths=source_relative:" + inputPath,
	}

	if genAll {
		args = append(args,
			"--go-grpc_out=paths=source_relative:"+inputPath,
			"--go-http_out=paths=source_relative:"+inputPath,
			"--go-errors_out=paths=source_relative:"+inputPath,
			"--openapi_out=paths=source_relative:"+inputDir,
			"--openapi_opt=naming=proto",
			"--openapi_opt=default_response=false",
		)
	} else {
		if genGrpc {
			args = append(args, "--go-grpc_out=paths=source_relative:"+inputPath)
		}
		if genHTTP {
			args = append(args, "--go-http_out=paths=source_relative:"+inputPath)
		}
		if genError {
			args = append(args, "--go-errors_out=paths=source_relative:"+inputPath)
		}
		if genOpenapi {
			args = append(
				args,
				"--openapi_out=paths=source_relative:"+inputDir,
				"--openapi_opt=naming=proto",
				"--openapi_opt=default_response=false",
			)
		}
	}

	// add validate generator
	reg := regexp.MustCompile(`\n[^/]*(import)\s+"validate/validate.proto"`)
	for _, t := range target {
		protoBytes, err := os.ReadFile(t)
		if err == nil && len(protoBytes) > 0 && reg.Match(protoBytes) {
			args = append(args, "--validate_out=lang=go,paths=source_relative:"+inputDir)
			break
		}
	}

	args = append(args, target...)
	return args, nil
}

func gen(args []string) error {
	fd := exec.Command("protoc", args...)
	fd.Stdout = os.Stdout
	fd.Stderr = os.Stderr
	fd.Dir = "."
	if verbose {
		log.Printf("command: %s\n", fd.String())
	}
	if err := fd.Run(); err != nil {
		return err
	}
	return nil
}
