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

package dml

import (
	"bytes"
	"database/sql"
	"math"
	"testing"
	"time"

	"github.com/corestoreio/pkg/util/bufferpool"
)

var preprocessSink string

// BenchmarkInterpolate-4   	  500000	      4013 ns/o	     174 B/o	      11 allocs/o with reflection
// BenchmarkInterpolate-4   	  500000	      3591 ns/o	     174 B/o	      11 allocs/o string
// BenchmarkInterpolate-4   	  500000	      3599 ns/o	     174 B/o	      11 allocs/o []byte
// BenchmarkInterpolate-4   	  500000	      2684 ns/op	     160 B/op	       1 allocs/op
func BenchmarkInterpolate(b *testing.B) {
	ipBuf := bufferpool.Get()
	defer bufferpool.Put(ipBuf)

	const want = `SELECT * FROM x WHERE a = 1 AND b = -2 AND c = 3 AND d = 4 AND e = 5 AND f = 6 AND g = 7 AND h = 8 AND i = 9 AND j = 10 AND k = 'Hello' AND l = 1`
	const sqlBytes = `SELECT * FROM x WHERE a = ? AND b = ? AND c = ? AND d = ? AND e = ? AND f = ? AND g = ? AND h = ? AND i = ? AND j = ? AND k = ? AND l = ?`
	args := MakeArgs(3).
		Int(1).Int(-2).Int(3).Int(4).Int(5).Int(6).Int(7).Int(8).Int(9).Int(10).
		String("Hello").
		Bool(true)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := writeInterpolate(ipBuf, sqlBytes, args.arguments); err != nil {
			b.Fatal(err)
		}
		preprocessSink = ipBuf.String()
		ipBuf.Reset()
	}
	if preprocessSink != want {
		b.Fatalf("Have: %v Want: %v", ipBuf.String(), want)
	}
}

var benchmarkIsValidIdentifier int8

// BenchmarkIsValidIdentifier-4   	20000000	        92.0 ns/o	       0 B/o	       0 allocs/o
// BenchmarkIsValidIdentifier-4   	 5000000	       280 ns/o	       0 B/o	       0 allocs/o
func BenchmarkIsValidIdentifier(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkIsValidIdentifier = isValidIdentifier(`store_owner.catalog_product_entity_varchar`)
	}
	if benchmarkIsValidIdentifier != 0 {
		b.Fatalf("Should be zero but got %d", benchmarkIsValidIdentifier)
	}
}

func BenchmarkQuoteAlias(b *testing.B) {
	const want = "`e`.`price` AS `final_price`"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if have := Quoter.NameAlias("e.price", "final_price"); have != want {
			b.Fatalf("Have %s\nWant %s\n", have, want)
		}
	}
}

// BenchmarkConditions_writeOnDuplicateKey-4   	 5000000	       337 ns/o	       0 B/o	       0 allocs/o
func BenchmarkConditions_writeOnDuplicateKey(b *testing.B) {
	buf := new(bytes.Buffer)
	dk := Conditions{
		Column("name").Str("E0S 5D Mark III"),
		Column("sku").Values(),
		Column("stock").Int64(14),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := dk.writeOnDuplicateKey(buf, nil); err != nil {
			b.Fatalf("%+v", err)
		}
		buf.Reset()
	}
}

var benchmarkDialectEscapeTimeBuf = new(bytes.Buffer)

func BenchmarkDialectEscapeTime(b *testing.B) {
	date := now()
	for i := 0; i < b.N; i++ {
		dialect.EscapeTime(benchmarkDialectEscapeTimeBuf, date)
		benchmarkDialectEscapeTimeBuf.Reset()
	}
}

var benchmarkArgEnc argEncoded

func BenchmarkArgumentEncoding(b *testing.B) {

	b.Run("all types without warm up", func(b *testing.B) {
		t1 := now()
		t2 := now().Add(time.Minute * 2)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			benchmarkArgEnc = makeArgBytes().
				appendInt(3).
				appendInts(4, 5, 6).
				appendInt64(30).
				appendInt64s(40, 50, 60).
				appendUint64(math.MaxUint32).
				appendUint64s(800, 900).
				appendFloat64(math.MaxFloat32).
				appendFloat64s(80.5490, math.Pi).
				appendString("Finally, how will we ship and deliver Go 2?").
				appendStrings("Finally, how will we fly and deliver Go 1?", "Finally, how will we run and deliver Go 3?", "Finally, how will we walk and deliver Go 3?").
				appendBool(true).
				appendBool(false).
				appendBools(false, true, true, false, true).
				appendTime(t1).
				appendTimes(t1, t2, t1).
				appendNullString(sql.NullString{}, sql.NullString{Valid: true, String: "Hello"}).
				appendNullFloat64(sql.NullFloat64{Valid: true, Float64: math.E}, sql.NullFloat64{}).
				appendNullInt64(sql.NullInt64{Valid: true, Int64: 987654321}, sql.NullInt64{}).
				appendNullBool(sql.NullBool{}, sql.NullBool{Valid: true, Bool: true}, sql.NullBool{Valid: false, Bool: true}).
				appendNullTime(NullTime{Valid: true, Time: t1}, NullTime{})
		}
	})

	b.Run("all types with warm up", func(b *testing.B) {

		t1 := now()
		t2 := now().Add(time.Minute * 2)

		benchmarkArgEnc = makeArgBytes().
			appendInt(3).
			appendInts(4, 5, 6).
			appendInt64(30).
			appendInt64s(40, 50, 60).
			appendUint64(math.MaxUint32).
			appendUint64s(800, 900).
			appendFloat64(math.MaxFloat32).
			appendFloat64s(80.5490, math.Pi).
			appendString("Finally, how will we ship and deliver Go 2?").
			appendStrings("Finally, how will we fly and deliver Go 1?", "Finally, how will we run and deliver Go 3?", "Finally, how will we walk and deliver Go 3?").
			appendBool(true).
			appendBool(false).
			appendBools(false, true, true, false, true).
			appendTime(t1).
			appendTimes(t1, t2, t1).
			appendNullString(sql.NullString{}, sql.NullString{Valid: true, String: "Hello"}).
			appendNullFloat64(sql.NullFloat64{Valid: true, Float64: math.E}, sql.NullFloat64{}).
			appendNullInt64(sql.NullInt64{Valid: true, Int64: 987654321}, sql.NullInt64{}).
			appendNullBool(sql.NullBool{}, sql.NullBool{Valid: true, Bool: true}, sql.NullBool{Valid: false, Bool: true}).
			appendNullTime(NullTime{Valid: true, Time: t1}, NullTime{})

		ns := []sql.NullString{{}, {Valid: true, String: "Hello"}}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			benchmarkArgEnc.
				reset().
				appendInt(13).
				appendInts(14, 15, 16).
				appendInt64(130).
				appendInt64s(140, 150, 160).
				appendUint64(math.MaxUint16).
				appendUint64s(1800, 1900).
				appendFloat64(math.MaxFloat32).
				appendFloat64s(84.5490, math.Pi).
				appendString("F1nally, how will we ship and deliver Go 1?").
				appendStrings("F1nally, how will we fly and deliver Go 2?", "Finally, how will we run and deliver Go 3?", "Finally, how will we walk and deliver Go 4?").
				appendBool(false).
				appendBool(true).
				appendBools(false, true, true, false, true).
				appendTime(t1).
				appendTimes(t1, t2, t1).
				appendNullString(ns...).
				appendNullFloat64(sql.NullFloat64{Valid: true, Float64: math.E}, sql.NullFloat64{}).
				appendNullInt64(sql.NullInt64{Valid: true, Int64: 987654321}, sql.NullInt64{}).
				appendNullBool(sql.NullBool{}, sql.NullBool{Valid: true, Bool: true}, sql.NullBool{Valid: false, Bool: true}).
				appendNullTime(NullTime{Valid: true, Time: t1}, NullTime{})
			// b.Fatal(benchmarkArgEnc.DebugBytes())
		}
	})
	b.Run("number slices without warm up", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			benchmarkArgEnc = makeArgBytes().
				appendInt(3).
				appendInts(4, 5, 6).
				appendInt64(30).
				appendInt64s(40, 50, 60).
				appendUint64(math.MaxUint32).
				appendUint64s(800, 900).
				appendFloat64(math.MaxFloat32).
				appendFloat64s(80.5490, math.Pi)
		}
	})
	b.Run("number slices with warm up", func(b *testing.B) {
		benchmarkArgEnc = makeArgBytes().
			appendInt(3).
			appendInts(4, 5, 6).
			appendInt64(30).
			appendInt64s(40, 50, 60).
			appendUint64(math.MaxUint32).
			appendUint64s(800, 900).
			appendFloat64(math.MaxFloat32).
			appendFloat64s(80.5490, math.Pi)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			benchmarkArgEnc = benchmarkArgEnc.
				reset().
				appendInt(13).
				appendInts(14, 15, 16).
				appendInt64(130).
				appendInt64s(140, 150, 160).
				appendUint64(math.MaxUint32).
				appendUint64s(1800, 1900).
				appendFloat64(math.MaxFloat32).
				appendFloat64s(180.5490, math.Pi)
		}
	})

	b.Run("numbers without warm up", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			benchmarkArgEnc = makeArgBytes().
				appendInt(3).
				appendInt64(30).
				appendUint64(math.MaxUint32).
				appendFloat64(math.MaxFloat32)
		}
	})
	b.Run("numbers with warm up", func(b *testing.B) {
		benchmarkArgEnc = makeArgBytes().
			appendInt(3).
			appendInt64(30).
			appendUint64(math.MaxUint32).
			appendFloat64(math.MaxFloat32)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			benchmarkArgEnc = benchmarkArgEnc.
				reset().
				appendInt(9).
				appendInt64(130).
				appendUint64(math.MaxUint64).
				appendFloat64(math.MaxFloat64)
		}
	})
}

func BenchmarkSQLScanner(b *testing.B) {
	var valInt64 int64
	var valFloat64 float64
	var valUint64 uint64
	var valString string
	var valTime time.Time
	b.Run("NullInt64", func(b *testing.B) {
		val := []byte(`12345678`)
		for i := 0; i < b.N; i++ {
			var nv NullInt64
			if err := nv.Scan(val); err != nil {
				b.Fatal(err)
			}
			if nv.Int64 != 12345678 {
				b.Fatalf("Have %d Want %d", nv.Int64, 12345678)
			}
			valInt64 = nv.Int64
		}
	})
	b.Run("NullFloat64", func(b *testing.B) {
		val := []byte(`-1234.5678`)
		for i := 0; i < b.N; i++ {
			var nv NullFloat64
			if err := nv.Scan(val); err != nil {
				b.Fatal(err)
			}
			if nv.Float64 != -1234.5678 {
				b.Fatalf("Have %d Want %d", nv.Float64, -1234.5678)
			}
			valFloat64 = nv.Float64
		}
	})
	b.Run("NullUint64", func(b *testing.B) {
		val := []byte(`12345678910`)
		for i := 0; i < b.N; i++ {
			var nv NullUint64
			if err := nv.Scan(val); err != nil {
				b.Fatal(err)
			}
			if nv.Uint64 != 12345678910 {
				b.Fatalf("Have %d Want %d", nv.Uint64, 12345678910)
			}
			valUint64 = nv.Uint64
		}
	})
	b.Run("NullString", func(b *testing.B) {
		const want = `12345678910`
		val := []byte(want)
		for i := 0; i < b.N; i++ {
			var nv NullString
			if err := nv.Scan(val); err != nil {
				b.Fatal(err)
			}
			if nv.String != want {
				b.Fatalf("Have %q Want %q", nv.String, want)
			}
			valString = nv.String
		}
	})
	b.Run("NullTime", func(b *testing.B) {
		const want = `2006-01-02 19:04:05`
		val := []byte(want)
		for i := 0; i < b.N; i++ {
			var nv NullTime
			if err := nv.Scan(val); err != nil {
				b.Fatal(err)
			}
			if nv.Time.IsZero() {
				b.Fatalf("Time cannot be zero %s", nv.String())
			}
			valTime = nv.Time
		}
	})
}
