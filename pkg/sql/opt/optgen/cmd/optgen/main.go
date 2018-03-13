// Copyright 2018 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/cockroachdb/cockroach/pkg/sql/opt/optgen/lang"
)

type globResolver func(pattern string) (matches []string, err error)

type genFunc func(compiled *lang.CompiledExpr, w io.Writer)

var (
	errInvalidArgCount     = errors.New("invalid number of arguments")
	errUnrecognizedCommand = errors.New("unrecognized command")
)

func main() {
	gen := optgen{useGoFmt: true, maxErrors: 10, stdErr: os.Stderr}
	if !gen.run(os.Args[1:]...) {
		os.Exit(2)
	}
}

type optgen struct {
	// useGoFmt runs the go fmt tool on code generated by optgen, if this
	// setting is true.
	useGoFmt bool

	// maxErrors is the max number of errors that will be printed by optgen
	// before showing the "too many errors" message.
	maxErrors int

	// stdErr is the writer to which all standard error output will be
	// redirected.
	stdErr io.Writer

	// globResolver is called to map from source arguments to a set of file
	// names, using filepath.Glob syntax. The files will then be resolved by
	// fileResolver. Tests can hook this in order to avoid actually listing
	// files/directories on disk.
	globResolver globResolver

	// fileResolver is called to open an input file of the specified name. Tests
	// can hook this in order to avoid actually opening files on disk.
	fileResolver lang.FileResolver

	// cmdLine stores the set of flags used to invoke the Optgen tool.
	cmdLine *flag.FlagSet
}

func (g *optgen) run(args ...string) bool {
	// Parse command line.
	g.cmdLine = flag.NewFlagSet("optgen", flag.ContinueOnError)
	g.cmdLine.SetOutput(g.stdErr)
	g.cmdLine.Usage = g.usage
	g.cmdLine.String("out", "", "output file name of generated code")
	err := g.cmdLine.Parse(args)
	if err != nil {
		return false
	}

	// Get remaining args after any flags have been parsed.
	args = g.cmdLine.Args()
	if len(args) < 2 {
		g.cmdLine.Usage()
		g.reportError(errInvalidArgCount)
		return false
	}

	cmd := args[0]
	sources := g.cmdLine.Args()[1:]

	switch cmd {
	case "compile":
	case "exprs":
	case "factory":
	case "ifactory":
	case "ops":

	default:
		g.cmdLine.Usage()
		g.reportError(errUnrecognizedCommand)
		return false
	}

	// Set glob resolver if it hasn't yet been set.
	if g.globResolver == nil {
		g.globResolver = filepath.Glob
	}

	// Map sources to a set of files using the glob resolver.
	files := make([]string, 0, len(sources))
	for _, source := range sources {
		matches, err := g.globResolver(source)
		if err != nil {
			g.reportError(err)
			return false
		}
		files = append(files, matches...)
	}

	// Sort the files so that output is stable.
	sort.Strings(files)

	compiler := lang.NewCompiler(files...)

	if g.fileResolver != nil {
		// Use caller-provided custom file resolver.
		compiler.SetFileResolver(g.fileResolver)
	}

	var errors []error
	compiled := compiler.Compile()
	if compiled == nil {
		errors = compiler.Errors()
	} else {
		// Do additional validation checks.
		errors = g.validate(compiled)
	}

	if errors != nil {
		for i, err := range errors {
			if i >= g.maxErrors-1 {
				count := len(errors) - g.maxErrors + 1
				if count > 1 {
					fmt.Fprintf(g.stdErr, "... too many errors (%d more)\n", count)
					break
				}
			}

			fmt.Fprintf(g.stdErr, "%v\n", err)
		}
		return false
	}

	switch cmd {
	case "compile":
		err = g.writeOutputFile([]byte(compiled.String()))

	case "exprs":
		var gen exprsGen
		err = g.generate(compiled, gen.generate)

	case "factory":
		var gen factoryGen
		err = g.generate(compiled, gen.generate)

	case "ifactory":
		var gen ifactoryGen
		err = g.generate(compiled, gen.generate)

	case "ops":
		var gen opsGen
		err = g.generate(compiled, gen.generate)
	}

	if err != nil {
		g.reportError(err)
		return false
	}
	return true
}

// validate performs additional checks on the compiled Optgen expression. In
// particular, it checks the order and types of the fields in define
// expressions. The Optgen language itself allows any field order and types, so
// the compiler does not do these checks.
func (g *optgen) validate(compiled *lang.CompiledExpr) []error {
	var errors []error

	for _, define := range compiled.Defines {
		// Ensure that fields are defined in the following order:
		//   Expr*
		//   ExprList?
		//   Private?
		//
		// That is, there can be zero or more expression-typed fields, followed
		// by zero or one list-typed field, followed by zero or one private field.
		for i, field := range define.Fields {
			if isPrivateType(string(field.Type)) {
				if i != len(define.Fields)-1 {
					format := "private field '%s' is not the last field in '%s'"
					err := fmt.Errorf(format, field.Name, define.Name)
					err = addErrorSource(err, field.Source())
					errors = append(errors, err)
					break
				}
			}

			if isListType(string(field.Type)) {
				index := len(define.Fields) - 1
				if privateField(define) != nil {
					index--
				}

				if i != index {
					format := "list field '%s' is not the last non-private field in '%s'"
					err := fmt.Errorf(format, field.Name, define.Name)
					err = addErrorSource(err, field.Source())
					errors = append(errors, err)
					break
				}
			}
		}
	}

	return errors
}

func (g *optgen) generate(compiled *lang.CompiledExpr, genFunc genFunc) error {
	var buf bytes.Buffer

	buf.WriteString("// Code generated by optgen; DO NOT EDIT.\n\n")

	genFunc(compiled, &buf)

	var b []byte
	var err error

	if g.useGoFmt {
		b, err = format.Source(buf.Bytes())
		if err != nil {
			// Write out incorrect source for easier debugging.
			b = buf.Bytes()
			out := g.cmdLine.Lookup("out").Value.String()
			err = fmt.Errorf("Code formatting failed with Go parse error\n%s:%s", out, err)
		}
	} else {
		b = buf.Bytes()
	}

	if err != nil {
		// Ignore any write error if another error already occurred.
		_ = g.writeOutputFile(b)
	} else {
		err = g.writeOutputFile(b)
	}

	return err
}

func (g *optgen) writeOutputFile(b []byte) error {
	out := g.cmdLine.Lookup("out").Value.String()
	if out != "" {
		file, err := os.Create(out)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = file.Write(b)
		return err
	}

	// Send output to stderr.
	_, err := g.stdErr.Write(b)
	return err
}

// usage is a replacement usage function for the flags package.
func (g *optgen) usage() {
	fmt.Fprintf(g.stdErr, "Optgen is a tool for generating cost-based optimizers.\n\n")

	fmt.Fprintf(g.stdErr, "It compiles source files that use a custom syntax to define expressions,\n")
	fmt.Fprintf(g.stdErr, "match expression patterns, and generate replacement expressions.\n\n")

	fmt.Fprintf(g.stdErr, "Usage:\n")

	fmt.Fprintf(g.stdErr, "\toptgen [flags] command sources...\n\n")

	fmt.Fprintf(g.stdErr, "The commands are:\n\n")
	fmt.Fprintf(g.stdErr, "\tcompile    generate the optgen compiled format\n")
	fmt.Fprintf(g.stdErr, "\texprs      generate expression definitions and functions\n")
	fmt.Fprintf(g.stdErr, "\tfactory    generate expression tree creation and normalization functions\n")
	fmt.Fprintf(g.stdErr, "\tifactory   generate interface for factory construct methods\n")
	fmt.Fprintf(g.stdErr, "\tops        generate operator definitions and functions\n")
	fmt.Fprintf(g.stdErr, "\n")

	fmt.Fprintf(g.stdErr, "The sources can be file names and/or filepath.Glob patterns.\n")
	fmt.Fprintf(g.stdErr, "\n")

	fmt.Fprintf(g.stdErr, "Flags:\n")

	g.cmdLine.PrintDefaults()

	fmt.Fprintf(g.stdErr, "\n")
}

func (g *optgen) reportError(err error) {
	fmt.Fprintf(g.stdErr, "ERROR: %v\n", err)
}

func addErrorSource(err error, src *lang.SourceLoc) error {
	return fmt.Errorf("%s: %s", src, err)
}
