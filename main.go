package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"wasm/generator/gowasm"
	"wasm/generator/transform"
	"wasm/generator/types"

	"github.com/dennwc/webidl/ast"
	"github.com/dennwc/webidl/parser"
)

var args struct {
	outputPath  string
	packageBase string
	warnings    bool
	singlePkg   string
}

var stopErr = errors.New("stopping for previous error")
var currentFilename string

func main() {
	if msg := parseArgs(); msg != "" {
		fmt.Fprintln(os.Stderr, "command line error:", msg)
		os.Exit(1)
	}
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	if fi, err := os.Stat(args.outputPath); err != nil {
		return fmt.Errorf("trouble evaluate %s: %s", args.outputPath, err)
	} else if !fi.IsDir() {
		return fmt.Errorf("output path '%s' doesn't point to a directory", args.outputPath)
	}

	conv := types.NewConvert()
	setup := &types.Setup{
		Package: args.singlePkg,
		Error:   failing,
		Warning: warning,
	}

	for _, name := range flag.Args() {
		if err := processFile(name, conv, setup); err != nil {
			return err
		}
	}
	if err := conv.EvaluateInput(); err != nil {
		return err
	}
	if err := conv.EvaluateOutput(); err != nil {
		return err
	}
	transform.RenameOverrideMethods(conv)

	files, err := gowasm.WriteSource(conv)
	if err != nil {
		return err
	}

	for k, v := range files {
		path := filepath.Join(args.outputPath, k)
		dir := filepath.Dir(path)
		if !pathExist(dir) {
			fmt.Println("creating folder", dir)
			if err := os.MkdirAll(dir, 0775); err != nil {
				return err
			}
		}
		fmt.Println("saving ", path)
		if err := ioutil.WriteFile(path, v, 0666); err != nil {
			return err
		}
	}

	return nil
}

func processFile(filename string, conv *types.Convert, setup *types.Setup) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	file := parser.Parse(string(content))
	trouble := ast.GetAllErrorNodes(file)
	if len(trouble) > 0 {
		sort.SliceStable(trouble, func(i, j int) bool { return trouble[i].Line < trouble[j].Line })
		for _, e := range trouble {
			failing(e.NodeBase(), e.Message)
		}
		return stopErr
	}

	currentFilename = filename
	if args.singlePkg == "" {
		setup.Package = gowasm.FormatPkg(filename)
	}
	if err := conv.Process(file, setup); err != nil {
		return err
	}
	return nil
}

func failing(base *ast.Base, format string, args ...interface{}) {
	dst := os.Stderr
	fmt.Fprint(dst, "error:", currentFilename, ":")
	if base != nil {
		fmt.Fprint(dst, base.Line, ":")
	}
	fmt.Fprintf(dst, format, args...)
	fmt.Fprint(dst, "\n")
}

func warning(base *ast.Base, format string, values ...interface{}) {
	if !args.warnings {
		return
	}
	dst := os.Stderr
	fmt.Fprint(dst, "warning:", currentFilename, ":")
	if base != nil {
		fmt.Fprint(dst, base.Line, ":")
	}
	fmt.Fprintf(dst, format, values...)
	fmt.Fprint(dst, "\n")
}

func parseArgs() string {
	flag.BoolVar(&args.warnings, "log-warning", true, "log warnings")
	flag.StringVar(&args.outputPath, "output", "", "output path")
	flag.StringVar(&args.packageBase, "package-base", "", "package base name (e.g. github.com/foo/bar)")
	flag.StringVar(&args.singlePkg, "single-package", "", "all types to same package")
	flag.Parse()
	if len(flag.Args()) == 0 {
		return "no input files on command line"
	}
	if args.outputPath == "" {
		return "missing output path for file(s)"
	}
	return ""
}

func pathExist(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}
