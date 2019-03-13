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

package dmlgen_test

import (
	"bytes"
	"testing"

	"github.com/corestoreio/pkg/sql/dmlgen"
	"github.com/corestoreio/pkg/util/assert"
)

func TestNewGenerator(t *testing.T) {
	g := dmlgen.NewGenerator("config")
	g.BuildTags = "ignoring"
	g.AddImport("fmt", "")
	g.AddImport("github.com/corestoreio/pkg/storage/null", "null")
	g.C("These constants", "are used for testing.", "Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an AS IS BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.")
	g.AddConstString("TableA", "table_a")

	g.P("type", "CatalogProductEntity", "struct {")
	g.In()
	g.P("EntityID", "int64")
	g.P("StoreID", "uint32", `// store_id smallint(5) unsigned NOT NULL PRI   "Store ID"`)
	g.P("Value", "null.Decimal", `// value decimal(12,4) NOT NULL PRI   "Value"`)
	g.Out()
	g.P("}")

	var buf bytes.Buffer
	err := g.GenerateFile(&buf)
	assert.NoError(t, err)

	assert.Exactly(t, `// +build  ignoring

package config

// Auto generated via github.com/corestoreio/pkg/sql/dmlgen
import (
	"fmt"
	null "github.com/corestoreio/pkg/storage/null"
)

const (
	TableA = "table_a"
)

// These constants are used for testing. Unless required by applicable law or
// agreed to in writing, software distributed under the License is distributed on
// an AS IS BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
// or implied. See the License for the specific language governing permissions
// and limitations under the License.
type CatalogProductEntity struct {
	EntityID int64
	StoreID  uint32       // store_id smallint(5) unsigned NOT NULL PRI   "Store ID"
	Value    null.Decimal // value decimal(12,4) NOT NULL PRI   "Value"
}
`, buf.String())

}
