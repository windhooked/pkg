// Copyright 2015-present, Cyrill @ Schumacher.fm and the CoreStore contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package codegen

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"sort"

	"github.com/corestoreio/errors"
)

type Go struct {
	common
	BuildTags       string
	init            []string // Lines to emit in the init function.
	constantsString []string
}

// NewGo creates a new Go source code generator for a specific new package.
func NewGo(packageName string) *Go {
	return &Go{
		common: common{
			Buffer:       new(bytes.Buffer),
			packageName:  packageName,
			packageNames: map[string]string{},
		},
	}
}

func (g *Go) AddConstString(name, value string) {
	g.constantsString = append(g.constantsString, fmt.Sprintf("\t%s = %q\n", name, value))
}

// AddInitf stores the given statement to be printed inside the file's init
// function. The statement is given as a format specifier and arguments.
func (g *Go) AddInitf(stmt string, a ...interface{}) {
	g.init = append(g.init, fmt.Sprintf(stmt, a...))
}

func (g *Go) generateInitFunction() {
	if len(g.init) == 0 {
		return
	}
	g.P("func init() {")
	g.In()
	for _, l := range g.init {
		g.P(l)
	}
	g.Out()
	g.P("}")
	g.init = nil
}

func (g *Go) generateImports(w io.Writer) {
	fmt.Fprintln(w, "import (")
	pkgSorted := make([]string, 0, len(g.packageNames))
	for key := range g.packageNames {
		pkgSorted = append(pkgSorted, key)
	}
	sort.Strings(pkgSorted)
	for _, p := range pkgSorted {
		fmt.Fprintf(w, "\t%s %q\n", g.packageNames[p], p)
	}
	fmt.Fprintln(w, ")")
}

func (g *Go) generateConstants(w io.Writer) {
	fmt.Fprintln(w, "const (")
	sort.Strings(g.constantsString)
	for _, cs := range g.constantsString {
		fmt.Fprint(w, cs)
	}
	fmt.Fprintln(w, ")")
}

func (g *Go) GenerateFile(w io.Writer) error {

	var buf bytes.Buffer
	if g.BuildTags != "" {
		fmt.Fprintln(&buf, "// +build ", g.BuildTags)
		fmt.Fprint(&buf, "\n") // the extra line as required from the Go spec
	}
	fmt.Fprintf(&buf, "package %s\n", g.packageName)
	fmt.Fprintln(&buf, "// Auto generated source code")
	g.generateImports(&buf)
	g.generateConstants(&buf)
	g.generateInitFunction()

	g.Buffer.WriteTo(&buf)

	fmted, err := format.Source(buf.Bytes())
	if err != nil {
		return errors.NotAcceptable.New(err, "\nSource Code:\n%s\n", buf.String())
	}
	_, err = w.Write(fmted)
	return err
}