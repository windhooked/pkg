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

package text_test

import (
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"testing"

	"errors"
	"github.com/corestoreio/csfw/storage/text"
	"github.com/stretchr/testify/assert"
)

// These checks if a type implements an interface belong into the test package
// and not into its "main" package. Otherwise you would also compile each time
// al the package with their interfaces.
var _ encoding.TextMarshaler = (*text.Chars)(nil)
var _ encoding.TextUnmarshaler = (*text.Chars)(nil)
var _ sql.Scanner = (*text.Chars)(nil)
var _ driver.Valuer = (*text.Chars)(nil)

func TestEqual(t *testing.T) {
	tests := []struct {
		a    text.Chars
		b    text.Chars
		want bool
	}{
		{nil, nil, true},
		{text.Chars("a"), text.Chars("a"), true},
		{text.Chars("a"), text.Chars("b"), false},
		{text.Chars("a\x80"), text.Chars("a"), false},
	}
	for i, test := range tests {
		assert.Exactly(t, test.want, test.a.Equal(test.b), "Index %d", i)
	}
}

func TestChars(t *testing.T) {
	t.Parallel()
	const have string = `Hello fellow Gpher's`
	l := text.Chars(have)
	var l1 text.Chars
	assert.True(t, l1.IsEmpty())
	assert.False(t, l.IsEmpty())
	assert.Exactly(t, have, l.String())

	l2 := l.Copy()
	assert.Exactly(t, l, l2)
	l2 = nil
	assert.True(t, l2.IsEmpty())
	assert.False(t, l.IsEmpty())
}

func TestTextMarshal(t *testing.T) {
	const have = `admin/security/passwrd_lifetime`
	t.Parallel()
	r := text.Chars(have)
	j, err := json.Marshal(r)
	assert.NoError(t, err)
	assert.Exactly(t, `"`+have+`"`, string(j))
}

func TestUnmarshalTextOk(t *testing.T) {
	t.Parallel()
	const have = `admin/security/passwörd_lif‹time`
	var r text.Chars
	err := json.Unmarshal([]byte(`"`+have+`"`), &r)
	assert.NoError(t, err)
	assert.Exactly(t, have, r.String())
}

func TestScan(t *testing.T) {
	tests := []struct {
		want    string
		val     interface{}
		wantErr error
	}{
		{"", nil, nil},
		{"hello", "hello", nil},
		{"h€llo", []byte("h€llo"), nil},
		{"", 8, errors.New("Cannot convert value 8 to []byte")},
	}
	for i, test := range tests {
		var l text.Chars
		haveErr := l.Scan(test.val)
		if test.wantErr != nil {
			assert.EqualError(t, haveErr, test.wantErr.Error(), "Index %d", i)
			assert.Nil(t, l, "Index %d", i)
			continue
		}
		assert.NoError(t, haveErr, "Index %d", i)
		assert.Exactly(t, test.want, l.String(), "Index %d", i)
	}
}

func TestValue(t *testing.T) {
	t.Parallel()
	l1 := text.Chars(`x`)
	v, err := l1.Value()
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Exactly(t, l1, v.(text.Chars))

	var l2 text.Chars
	v, err = l2.Value()
	assert.NoError(t, err)
	assert.Nil(t, v)
}

func TestRuneCount(t *testing.T) {
	t.Parallel()
	const have string = "Hello F€llow Gphers"
	l1 := text.Chars(have)
	assert.Exactly(t, 20, l1.RuneCount())
	assert.Exactly(t, 24, len(l1))
}