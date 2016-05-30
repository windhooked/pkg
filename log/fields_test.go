// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package log

import (
	"bytes"
	"encoding/json"
	"errors"
	"math"
	"net/http"
	"testing"
	"time"

	"github.com/corestoreio/csfw/storage/text"
	"github.com/stretchr/testify/assert"
)

const testKey = "MyTestKey"

func TestField_Bool(t *testing.T) {
	f := Bool(testKey, true)
	assert.Exactly(t, typeBool, f.fieldType)
	assert.Exactly(t, int64(1), f.int64)
	assert.Exactly(t, testKey, f.key)
}

func TestField_Float64(t *testing.T) {
	f := Float64(testKey, math.Pi)
	assert.Exactly(t, typeFloat64, f.fieldType)
	assert.Exactly(t, math.Pi, f.float64)
	assert.Exactly(t, testKey, f.key)
}

func TestField_Int(t *testing.T) {
	f := Int(testKey, math.MaxInt32)
	assert.Exactly(t, typeInt, f.fieldType)
	assert.Exactly(t, int64(math.MaxInt32), f.int64)
	assert.Exactly(t, testKey, f.key)
}

func TestField_Int64(t *testing.T) {
	f := Int64(testKey, math.MaxInt64)
	assert.Exactly(t, typeInt64, f.fieldType)
	assert.Exactly(t, int64(math.MaxInt64), f.int64)
	assert.Exactly(t, testKey, f.key)
}

func TestField_Uint(t *testing.T) {
	f := Uint(testKey, math.MaxUint32)
	assert.Exactly(t, typeInt, f.fieldType)
	assert.Exactly(t, int64(math.MaxUint32), f.int64)
	assert.Exactly(t, testKey, f.key)
}

func TestField_Uint64(t *testing.T) {
	f := Uint64(testKey, math.MaxUint64)
	assert.Exactly(t, typeInt64, f.fieldType)
	assert.Exactly(t, int64(math.MaxInt64), f.int64)
	assert.Exactly(t, testKey, f.key)
}

func TestField_String(t *testing.T) {
	const data = `16. “One is never alone with a rubber duck.” Douglas Adams`
	f := String(testKey, data)
	assert.Exactly(t, typeString, f.fieldType)
	assert.Exactly(t, data, f.string)
	assert.Exactly(t, testKey, f.key)
}

func TestField_Stringer(t *testing.T) {
	const data = `27. “Anything invented after you're thirty-five is against the natural order of things.” Douglas Adams`
	f := Stringer(testKey, bytes.NewBufferString(data))
	assert.Exactly(t, typeStringer, f.fieldType)
	assert.Empty(t, f.string)
	assert.Exactly(t, testKey, f.key)

	buf := &bytes.Buffer{}
	wt := WriteTypes{W: buf}
	if err := f.AddTo(wt); err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, " MyTestKey: \"27. “Anything invented after you're thirty-five is against the natural order of things.” Douglas Adams\"", buf.String())
}

type gs struct {
	data interface{}
	err  error
}

func (g gs) MarshalText() ([]byte, error) {
	if g.err != nil {
		return nil, g.err
	}
	return g.data.([]byte), nil
}
func (gs) GoString() string { return "gs struct {}" }
func (g gs) MarshalJSON() ([]byte, error) {
	d, err := json.Marshal(g.data)
	if err != nil {
		g.err = err
	}
	return d, g.err
}
func (g gs) MarshalLog() (Field, error) {
	if g.err != nil {
		return Field{}, g.err
	}
	return String("ignored", "Val1x"), nil
}

func TestField_GoStringer(t *testing.T) {
	f := GoStringer(testKey, gs{})
	assert.Exactly(t, typeGoStringer, f.fieldType)
	assert.Empty(t, f.string)
	assert.Exactly(t, testKey, f.key)
	buf := &bytes.Buffer{}
	wt := WriteTypes{W: buf}
	if err := f.AddTo(wt); err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, " MyTestKey: \"gs struct {}\"", buf.String())
}

func TestField_Text(t *testing.T) {
	const data = `35. “My universe is my eyes and my ears. Anything else is hearsay.” Douglas Adams`
	f := Text(testKey, text.Chars(data))
	assert.Exactly(t, typeTextMarshaler, f.fieldType)
	assert.Empty(t, f.string)
	assert.Exactly(t, testKey, f.key)
	buf := &bytes.Buffer{}
	wt := WriteTypes{W: buf}
	if err := f.AddTo(wt); err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, " MyTestKey: \"35. “My universe is my eyes and my ears. Anything else is hearsay.” Douglas Adams\"", buf.String())
}
func TestField_TextError(t *testing.T) {
	var data = gs{data: nil, err: errors.New("Errr")}
	f := Text(testKey, data)
	assert.Exactly(t, typeTextMarshaler, f.fieldType)
	assert.Empty(t, f.string)
	assert.Exactly(t, testKey, f.key)
	buf := &bytes.Buffer{}
	wt := WriteTypes{W: buf}
	err := f.AddTo(wt)
	assert.Empty(t, buf.String())
	assert.EqualError(t, err, "[log] AddTo.TextMarshaler: Errr")

}

func TestField_JSON(t *testing.T) {
	const data = `12. “Reality is frequently inaccurate.” Douglas Adams`
	f := JSON(testKey, gs{data: data})
	assert.Exactly(t, typeString, f.fieldType)
	assert.Exactly(t, `"`+data+`"`, f.string)
	assert.Exactly(t, testKey, f.key)
}

func TestField_JSONError(t *testing.T) {
	f := JSON(testKey, gs{data: make(chan struct{})})
	assert.Exactly(t, typeString, f.fieldType)
	assert.Exactly(t, "[log] MarshalJSON: json: unsupported type: chan struct {}", f.string)
	assert.Exactly(t, ErrorKeyName, f.key)
}

func TestField_Time(t *testing.T) {
	now := time.Now()
	f := Time(testKey, now)
	assert.Exactly(t, typeInt64, f.fieldType)
	assert.Exactly(t, now.UnixNano(), f.int64)
	assert.Exactly(t, testKey, f.key)
}

func TestField_Duration(t *testing.T) {
	now := time.Hour * 2
	f := Duration(testKey, now)
	assert.Exactly(t, typeInt64, f.fieldType)
	assert.Exactly(t, now.Nanoseconds(), f.int64)
	assert.Exactly(t, testKey, f.key)
}

func TestField_Error(t *testing.T) {
	const data = `15. “There is no point in using the word 'impossible' to describe something that has clearly happened.” Douglas Adams`
	err := errors.New(data)
	f := Err(err)
	assert.Exactly(t, typeString, f.fieldType)
	assert.Exactly(t, data, f.string)
	assert.Exactly(t, ErrorKeyName, f.key)
}

func TestField_Error_Nil(t *testing.T) {
	f := Err(nil)
	assert.Exactly(t, typeString, f.fieldType)
	assert.Exactly(t, `<nil>`, f.string)
	assert.Exactly(t, ErrorKeyName, f.key)
}

func TestField_ErrorWithKey(t *testing.T) {
	const data = `15. “There is no point in using the word 'impossible' to describe something that has clearly happened.” Douglas Adams`
	err := errors.New(data)
	f := ErrWithKey("e1", err)
	assert.Exactly(t, typeString, f.fieldType)
	assert.Exactly(t, data, f.string)
	assert.Exactly(t, `e1`, f.key)
}

func TestField_ErrorWithKey_Nil(t *testing.T) {
	f := ErrWithKey(`e2`, nil)
	assert.Exactly(t, typeString, f.fieldType)
	assert.Exactly(t, `<nil>`, f.string)
	assert.Exactly(t, `e2`, f.key)
}

func TestField_Object(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://corestore.io", nil)
	req.RemoteAddr = "192.168.0.42"
	f := Object(testKey, req)
	assert.Exactly(t, typeObject, f.fieldType)
	assert.Exactly(t, req, f.obj)
	assert.Exactly(t, testKey, f.key)
}

func TestField_Nest(t *testing.T) {
	f := Nest("nest0", String("nest1", "1"), Float64("nest2", math.Log2E))
	assert.Exactly(t, typeMarshaler, f.fieldType)
	assert.Exactly(t, `nest0`, f.key)
	buf := &bytes.Buffer{}
	wt := WriteTypes{W: buf}
	if err := f.AddTo(wt); err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, ` nest1: "1" nest2: 1.4426950408889634`, buf.String())
}
