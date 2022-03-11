package base

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

// FindModulePath returns go.mod file path
func FindModulePath(current string) (string, error) {
	target, err := filepath.Abs(current)
	if err != nil {
		return "", err
	}
	name := filepath.Join(target, "go.mod")
	_, err = os.Stat(name)
	if err != nil {
		return FindModulePath(filepath.Join(target, ".."))
	}
	return name, nil
}

// ModulePath returns go module path.
func ModulePath(filename string) (string, error) {
	modBytes, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return modfile.ModulePath(modBytes), nil
}

// ModuleVersion returns module version.
func ModuleVersion(path string) (string, error) {
	stdout := &bytes.Buffer{}
	fd := exec.Command("go", "mod", "graph")
	fd.Stdout = stdout
	fd.Stderr = stdout
	if err := fd.Run(); err != nil {
		return "", err
	}
	rd := bufio.NewReader(stdout)
	for {
		line, _, err := rd.ReadLine()
		if err != nil {
			return "", err
		}
		str := string(line)
		i := strings.Index(str, "@")
		if strings.Contains(str, path+"@") && i != -1 {
			return path + str[i:], nil
		}
	}
}

// KratosMod returns kratos mod.
func KratosMod() string {
	// go 1.15+ read from env GOMODCACHE
	cacheOut, _ := exec.Command("go", "env", "GOMODCACHE").Output()
	cachePath := strings.Trim(string(cacheOut), "\n")
	pathOut, _ := exec.Command("go", "env", "GOPATH").Output()
	gopath := strings.Trim(string(pathOut), "\n")
	if cachePath == "" {
		cachePath = filepath.Join(gopath, "pkg", "mod")
	}
	if path, err := ModuleVersion("github.com/SeeMusic/kratos/v2"); err == nil {
		// $GOPATH/pkg/mod/github.com/SeeMusic/kratos@v2
		return filepath.Join(cachePath, path)
	}
	// $GOPATH/src/github.com/SeeMusic/kratos
	return filepath.Join(gopath, "src", "github.com", "go-kratos", "kratos")
}
