// Copyright 2015-2017, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package dbr

import (
	"bytes"
	"context"
	"testing"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSelect_BasicToSQL(t *testing.T) {
	t.Parallel()

	t.Run("no table no args", func(t *testing.T) {
		sel := NewSelect().AddColumnsConditions(Expr("1").Alias("n")).AddColumnsAliases("abc", "str")
		compareToSQL(t, sel, nil,
			"SELECT 1 AS `n`, `abc` AS `str`",
			"",
		)
	})
	t.Run("no table with args", func(t *testing.T) {
		sel := NewSelect().
			AddColumnsConditions(
				Expr("?").Alias("n").Int64(1),
				Expr("CAST(? AS CHAR(20))").Alias("str").Str("a'bc"),
			)
		compareToSQL(t, sel, nil,
			"SELECT ? AS `n`, CAST(? AS CHAR(20)) AS `str`",
			"SELECT 1 AS `n`, CAST('a\\'bc' AS CHAR(20)) AS `str`",
			int64(1), "a'bc",
		)
	})

	t.Run("two cols, one table, one condition", func(t *testing.T) {
		sel := NewSelect("a", "b").From("c").Where(Column("id").Equal().Int(1))
		compareToSQL(t, sel, nil,
			"SELECT `a`, `b` FROM `c` WHERE (`id` = ?)",
			"SELECT `a`, `b` FROM `c` WHERE (`id` = 1)",
			int64(1),
		)
	})

	t.Run("column right expression without arguments", func(t *testing.T) {
		sel := NewSelect("sku", "name").From("products").Where(
			Column("id").NotBetween().Ints(4, 7),
			Column("name").NotEqual().Expr("CONCAT('Canon','E0S 5D Mark III')"),
		)
		compareToSQL(t, sel, nil,
			"SELECT `sku`, `name` FROM `products` WHERE (`id` NOT BETWEEN ? AND ?) AND (`name` != CONCAT('Canon','E0S 5D Mark III'))",
			"SELECT `sku`, `name` FROM `products` WHERE (`id` NOT BETWEEN 4 AND 7) AND (`name` != CONCAT('Canon','E0S 5D Mark III'))",
			int64(4), int64(7),
		)
	})

	t.Run("column right expression with one argument", func(t *testing.T) {
		sel := NewSelect("sku", "name").From("products").Where(
			Column("id").NotBetween().Ints(4, 7),
			Column("name").Like().Expr("CONCAT('Canon',?,'E0S 7D Mark VI')").Str("Camera"),
		)
		compareToSQL(t, sel, nil,
			"SELECT `sku`, `name` FROM `products` WHERE (`id` NOT BETWEEN ? AND ?) AND (`name` LIKE CONCAT('Canon',?,'E0S 7D Mark VI'))",
			"SELECT `sku`, `name` FROM `products` WHERE (`id` NOT BETWEEN 4 AND 7) AND (`name` LIKE CONCAT('Canon','Camera','E0S 7D Mark VI'))",
			int64(4), int64(7), "Camera",
		)
	})

	t.Run("column right expression with slice argument", func(t *testing.T) {
		sel := NewSelect("sku", "name").From("products").Where(
			Column("id").NotBetween().Ints(4, 7),
			Column("name").NotLike().Expr("CONCAT('Canon',?,'E0S 8D Mark VII')").Strs("Camera", "Photo", "Flash"),
		)
		compareToSQL(t, sel, nil,
			"SELECT `sku`, `name` FROM `products` WHERE (`id` NOT BETWEEN ? AND ?) AND (`name` NOT LIKE CONCAT('Canon',?,?,?,'E0S 8D Mark VII'))",
			"SELECT `sku`, `name` FROM `products` WHERE (`id` NOT BETWEEN 4 AND 7) AND (`name` NOT LIKE CONCAT('Canon','Camera','Photo','Flash','E0S 8D Mark VII'))",
			int64(4), int64(7), "Camera", "Photo", "Flash",
		)
	})

	//t.Run("column left and right expression without arguments", func(t *testing.T) {
	//	sel := NewSelect("sku", "name").From("products").Where(
	//		Column("id").NotBetween().Ints(4, 7),
	//		Column("name").NotEqual().Expr("CONCAT('Canon','E0S 5D Mark III')"),
	//	)
	//	compareToSQL(t, sel, nil,
	//		"SELECT `sku`, `name` FROM `products` WHERE (`id` NOT BETWEEN ? AND ?) AND (`name` != CONCAT('Canon','E0S 5D Mark III'))",
	//		"SELECT `sku`, `name` FROM `products` WHERE (`id` NOT BETWEEN 4 AND 7) AND (`name` != CONCAT('Canon','E0S 5D Mark III'))",
	//		int64(4), int64(7),
	//	)
	//})
}

func TestSelectFullToSQL(t *testing.T) {
	t.Parallel()

	sel := NewSelect("a", "b").
		Distinct().
		FromAlias("c", "cc").
		Where(
			ParenthesisOpen(),
			Column("d").Int(1),
			Column("e").Str("wat").Or(),
			ParenthesisClose(),
			Column("f").Int(2),
			Column("g").Int(3),
			Column("h").In().Int64s(4, 5, 6),
		).
		GroupBy("ab").
		Having(
			ParenthesisOpen(),
			Column("m").Int(33),
			Column("n").Str("wh3r3").Or(),
			ParenthesisClose(),
			Expr("j = k"),
		).
		OrderBy("l").
		Limit(7).
		Offset(8)

	compareToSQL(t, sel, nil,
		"SELECT DISTINCT `a`, `b` FROM `c` AS `cc` WHERE ((`d` = ?) OR (`e` = ?)) AND (`f` = ?) AND (`g` = ?) AND (`h` IN (?,?,?)) GROUP BY `ab` HAVING ((`m` = ?) OR (`n` = ?)) AND (j = k) ORDER BY `l` LIMIT 7 OFFSET 8",
		"SELECT DISTINCT `a`, `b` FROM `c` AS `cc` WHERE ((`d` = 1) OR (`e` = 'wat')) AND (`f` = 2) AND (`g` = 3) AND (`h` IN (4,5,6)) GROUP BY `ab` HAVING ((`m` = 33) OR (`n` = 'wh3r3')) AND (j = k) ORDER BY `l` LIMIT 7 OFFSET 8",
		int64(1), "wat", int64(2), int64(3), int64(4), int64(5), int64(6), int64(33), "wh3r3",
	)
}

func TestSelect_Interpolate(t *testing.T) {
	t.Parallel()

	t.Run("with paranthesis", func(t *testing.T) {
		sel := NewSelect("a", "b").
			Distinct().
			FromAlias("c", "cc").
			Where(
				ParenthesisOpen(),
				Column("d").Int(1),
				Column("e").Str("wat").Or(),
				ParenthesisClose(),
				Column("f").Int(2),
				Column("g").Int(3),
				Column("h").In().Int64s(4, 5, 6),
			).
			GroupBy("ab").
			Having(
				ParenthesisOpen(),
				Column("m").Int(33),
				Column("n").Str("wh3r3").Or(),
				ParenthesisClose(),
				Expr("j = k"),
			).
			OrderBy("l").
			Limit(7).
			Offset(8)
		compareToSQL(t, sel, nil,
			"SELECT DISTINCT `a`, `b` FROM `c` AS `cc` WHERE ((`d` = ?) OR (`e` = ?)) AND (`f` = ?) AND (`g` = ?) AND (`h` IN (?,?,?)) GROUP BY `ab` HAVING ((`m` = ?) OR (`n` = ?)) AND (j = k) ORDER BY `l` LIMIT 7 OFFSET 8",
			"SELECT DISTINCT `a`, `b` FROM `c` AS `cc` WHERE ((`d` = 1) OR (`e` = 'wat')) AND (`f` = 2) AND (`g` = 3) AND (`h` IN (4,5,6)) GROUP BY `ab` HAVING ((`m` = 33) OR (`n` = 'wh3r3')) AND (j = k) ORDER BY `l` LIMIT 7 OFFSET 8",
			int64(1), "wat", int64(2), int64(3), int64(4), int64(5), int64(6), int64(33), "wh3r3",
		)
	})

	t.Run("two args in one condition", func(t *testing.T) {
		sel := NewSelect("a", "b", "z", "y", "x").From("c").
			Distinct().
			Where(Expr("`d` = ? OR `e` = ?").Int64(1).Str("wat")).
			Where(
				Column("g").Int(3),
				Column("h").In().Int64s(1, 2, 3),
			).
			GroupBy("ab").GroupBy("ii").GroupBy("iii").
			Having(Expr("j = k"), Column("jj").Int64(1)).
			Having(Column("jjj").Int64(2)).
			OrderBy("l1").OrderBy("l2").OrderBy("l3").
			Limit(7).Offset(8)

		compareToSQL(t, sel, nil,
			"SELECT DISTINCT `a`, `b`, `z`, `y`, `x` FROM `c` WHERE (`d` = ? OR `e` = ?) AND (`g` = ?) AND (`h` IN (?,?,?)) GROUP BY `ab`, `ii`, `iii` HAVING (j = k) AND (`jj` = ?) AND (`jjj` = ?) ORDER BY `l1`, `l2`, `l3` LIMIT 7 OFFSET 8",
			"SELECT DISTINCT `a`, `b`, `z`, `y`, `x` FROM `c` WHERE (`d` = 1 OR `e` = 'wat') AND (`g` = 3) AND (`h` IN (1,2,3)) GROUP BY `ab`, `ii`, `iii` HAVING (j = k) AND (`jj` = 1) AND (`jjj` = 2) ORDER BY `l1`, `l2`, `l3` LIMIT 7 OFFSET 8",
			int64(1), "wat", int64(3), int64(1), int64(2), int64(3), int64(1), int64(2),
		)
	})
}

func TestSelect_Paginate(t *testing.T) {
	t.Parallel()

	t.Run("asc", func(t *testing.T) {
		compareToSQL(t,
			NewSelect("a", "b").
				From("c").
				Where(Column("d").Int(1)).
				Paginate(3, 30).
				OrderBy("id"),
			nil,
			"SELECT `a`, `b` FROM `c` WHERE (`d` = ?) ORDER BY `id` LIMIT 30 OFFSET 60",
			"SELECT `a`, `b` FROM `c` WHERE (`d` = 1) ORDER BY `id` LIMIT 30 OFFSET 60",
			int64(1),
		)
	})
	t.Run("desc", func(t *testing.T) {
		compareToSQL(t,
			NewSelect("a", "b").
				From("c").
				Where(Column("d").Int(1)).
				Paginate(1, 20).
				OrderByDesc("id"),
			nil,
			"SELECT `a`, `b` FROM `c` WHERE (`d` = ?) ORDER BY `id` DESC LIMIT 20 OFFSET 0",
			"SELECT `a`, `b` FROM `c` WHERE (`d` = 1) ORDER BY `id` DESC LIMIT 20 OFFSET 0",
			int64(1),
		)
	})
}

func TestSelectWithoutWhere(t *testing.T) {
	t.Parallel()

	compareToSQL(t,
		NewSelect("a", "b").From("c"),
		nil,
		"SELECT `a`, `b` FROM `c`",
		"SELECT `a`, `b` FROM `c`",
	)
}

func TestSelectMultiHavingSQL(t *testing.T) {
	t.Parallel()

	compareToSQL(t,
		NewSelect("a", "b").From("c").
			Where(Column("p").Int(1)).
			GroupBy("z").Having(Column("z`z").Int(2), Column("y").Int(3)),
		nil,
		"SELECT `a`, `b` FROM `c` WHERE (`p` = ?) GROUP BY `z` HAVING (`zz` = ?) AND (`y` = ?)",
		"SELECT `a`, `b` FROM `c` WHERE (`p` = 1) GROUP BY `z` HAVING (`zz` = 2) AND (`y` = 3)",
		int64(1), int64(2), int64(3),
	)
}

func TestSelectMultiOrderSQL(t *testing.T) {
	t.Parallel()
	compareToSQL(t,
		NewSelect("a", "b").From("c").OrderBy("name").OrderByDesc("id"),
		nil,
		"SELECT `a`, `b` FROM `c` ORDER BY `name`, `id` DESC",
		"SELECT `a`, `b` FROM `c` ORDER BY `name`, `id` DESC",
	)
}

func TestSelect_OrderByDeactivated(t *testing.T) {
	t.Parallel()
	compareToSQL(t,
		NewSelect("a", "b").From("c").OrderBy("name").OrderByDeactivated(),
		nil,
		"SELECT `a`, `b` FROM `c` ORDER BY NULL",
		"SELECT `a`, `b` FROM `c` ORDER BY NULL",
	)
}

func TestSelect_ConditionColumn(t *testing.T) {
	t.Parallel()
	// TODO rewrite test to use every type which implements interface Argument and every operator

	runner := func(wf *Condition, wantSQL string, wantVal []interface{}) func(*testing.T) {
		return func(t *testing.T) {
			compareToSQL(t,
				NewSelect("a", "b").From("c").Where(wf),
				nil,
				wantSQL,
				"",
				wantVal...,
			)
		}
	}
	t.Run("single int64", runner(
		Column("d").Int64(33),
		"SELECT `a`, `b` FROM `c` WHERE (`d` = ?)",
		[]interface{}{int64(33)},
	))
	t.Run("IN int64", runner(
		Column("d").In().Int64s(33, 44),
		"SELECT `a`, `b` FROM `c` WHERE (`d` IN (?,?))",
		[]interface{}{int64(33), int64(44)},
	))
	t.Run("single float64", runner(
		Column("d").Float64(33),
		"SELECT `a`, `b` FROM `c` WHERE (`d` = ?)",
		[]interface{}{float64(33)},
	))
	t.Run("IN float64", runner(
		Column("d").In().Float64s(33, 44),
		"SELECT `a`, `b` FROM `c` WHERE (`d` IN (?,?))",
		[]interface{}{float64(33), float64(44)},
	))
	t.Run("NOT IN float64", runner(
		Column("d").NotIn().Float64s(33, 44),
		"SELECT `a`, `b` FROM `c` WHERE (`d` NOT IN (?,?))",
		[]interface{}{float64(33), float64(44)},
	))
	t.Run("single int", runner(
		Column("d").Equal().Int(33),
		"SELECT `a`, `b` FROM `c` WHERE (`d` = ?)",
		[]interface{}{int64(33)},
	))
	t.Run("IN int", runner(
		Column("d").In().Ints(33, 44),
		"SELECT `a`, `b` FROM `c` WHERE (`d` IN (?,?))",
		[]interface{}{int64(33), int64(44)},
	))
	t.Run("single string", runner(
		Column("d").Str("w"),
		"SELECT `a`, `b` FROM `c` WHERE (`d` = ?)",
		[]interface{}{"w"},
	))
	t.Run("IN string", runner(
		Column("d").In().Strs("x", "y"),
		"SELECT `a`, `b` FROM `c` WHERE (`d` IN (?,?))",
		[]interface{}{"x", "y"},
	))

	t.Run("BETWEEN int64", runner(
		Column("d").Between().Int64s(5, 6),
		"SELECT `a`, `b` FROM `c` WHERE (`d` BETWEEN ? AND ?)",
		[]interface{}{int64(5), int64(6)},
	))
	t.Run("NOT BETWEEN int64", runner(
		Column("d").NotBetween().Int64s(5, 6),
		"SELECT `a`, `b` FROM `c` WHERE (`d` NOT BETWEEN ? AND ?)",
		[]interface{}{int64(5), int64(6)},
	))

	t.Run("LIKE string", runner(
		Column("d").Like().Str("x%"),
		"SELECT `a`, `b` FROM `c` WHERE (`d` LIKE ?)",
		[]interface{}{"x%"},
	))
	t.Run("NOT LIKE string", runner(
		Column("d").NotLike().Str("x%"),
		"SELECT `a`, `b` FROM `c` WHERE (`d` NOT LIKE ?)",
		[]interface{}{"x%"},
	))

	t.Run("Less float64", runner(
		Column("d").Less().Float64(5.1),
		"SELECT `a`, `b` FROM `c` WHERE (`d` < ?)",
		[]interface{}{float64(5.1)},
	))
	t.Run("Greater float64", runner(
		Column("d").Greater().Float64(5.1),
		"SELECT `a`, `b` FROM `c` WHERE (`d` > ?)",
		[]interface{}{float64(5.1)},
	))
	t.Run("LessOrEqual float64", runner(
		Column("d").LessOrEqual().Float64(5.1),
		"SELECT `a`, `b` FROM `c` WHERE (`d` <= ?)",
		[]interface{}{float64(5.1)},
	))
	t.Run("GreaterOrEqual float64", runner(
		Column("d").GreaterOrEqual().Float64(5.1),
		"SELECT `a`, `b` FROM `c` WHERE (`d` >= ?)",
		[]interface{}{float64(5.1)},
	))

}

func TestSelect_Null(t *testing.T) {
	t.Parallel()

	t.Run("col is null", func(t *testing.T) {
		compareToSQL(t,
			NewSelect("a", "b").From("c").Where(Column("r").Null()),
			nil,
			"SELECT `a`, `b` FROM `c` WHERE (`r` IS NULL)",
			"SELECT `a`, `b` FROM `c` WHERE (`r` IS NULL)",
		)
	})

	t.Run("col is not null", func(t *testing.T) {
		compareToSQL(t,
			NewSelect("a", "b").From("c").Where(Column("r").NotNull()),
			nil,
			"SELECT `a`, `b` FROM `c` WHERE (`r` IS NOT NULL)",
			"SELECT `a`, `b` FROM `c` WHERE (`r` IS NOT NULL)",
		)
	})

	t.Run("complex", func(t *testing.T) {
		compareToSQL(t,
			NewSelect("a", "b").From("c").
				Where(
					Column("r").Null(),
					Column("d").Int(3),
					Column("ab").Null(),
					Column("w").NotNull(),
				),
			nil,
			"SELECT `a`, `b` FROM `c` WHERE (`r` IS NULL) AND (`d` = ?) AND (`ab` IS NULL) AND (`w` IS NOT NULL)",
			"SELECT `a`, `b` FROM `c` WHERE (`r` IS NULL) AND (`d` = 3) AND (`ab` IS NULL) AND (`w` IS NOT NULL)",
			int64(3),
		)
	})
}

func TestSelectWhereNULL(t *testing.T) {
	t.Parallel()

	t.Run("one nil", func(t *testing.T) {
		compareToSQL(t,
			NewSelect("a").From("b").Where(Column("a")),
			nil,
			"SELECT `a` FROM `b` WHERE (`a` IS NULL)",
			"SELECT `a` FROM `b` WHERE (`a` IS NULL)",
		)
	})

	t.Run("no values", func(t *testing.T) {
		var args = []interface{}{}
		compareToSQL(t,
			NewSelect("a").From("b").Where(Column("a").PlaceHolder()),
			nil,
			"SELECT `a` FROM `b` WHERE (`a` = ?)",
			"",
			args...,
		)
	})

	t.Run("empty Ints trigger invalid SQL", func(t *testing.T) {
		var iVal []int
		compareToSQL(t,
			NewSelect("a").From("b").Where(Column("a").In().Ints(iVal...)),
			nil,
			"SELECT `a` FROM `b` WHERE (`a` IN )",
			"",
		)
	})

	t.Run("Map nil arg", func(t *testing.T) {
		s := NewSelect("a").From("b").
			Where(
				Column("a"),
				Column("b").Bool(false),
				Column("c").Null(),
				Column("d").NotNull(),
			)
		compareToSQL(t, s, nil,
			"SELECT `a` FROM `b` WHERE (`a` IS NULL) AND (`b` = ?) AND (`c` IS NULL) AND (`d` IS NOT NULL)",
			"SELECT `a` FROM `b` WHERE (`a` IS NULL) AND (`b` = 0) AND (`c` IS NULL) AND (`d` IS NOT NULL)",
			false,
		)
	})
}

func TestSelectBySQL(t *testing.T) {
	t.Parallel()

	s := createFakeSession()

	compareToSQL(t,
		s.SelectBySQL("SELECT * FROM users WHERE x = 1", nil),
		nil,
		"SELECT * FROM users WHERE x = 1",
		"SELECT * FROM users WHERE x = 1",
	)
	compareToSQL(t,
		s.SelectBySQL("SELECT * FROM users WHERE x = ? AND y IN (?)", MakeArgs(2).Int(9).Ints(5, 6, 7)),
		nil,
		"SELECT * FROM users WHERE x = ? AND y IN (?)",
		"SELECT * FROM users WHERE x = 9 AND y IN (5,6,7)",
		int64(9), int64(5), int64(6), int64(7),
	)
	compareToSQL(t,
		s.SelectBySQL("wat", MakeArgs(2).Int(9).Ints(5, 6, 7)),
		nil,
		"wat",
		"",
		int64(9), int64(5), int64(6), int64(7),
	)
}

func TestSelectVarieties(t *testing.T) {
	t.Parallel()

	// This would be incorrect SQL!
	compareToSQL(t, NewSelect("id, name, email").From("users"), nil,
		"SELECT `id, name, email` FROM `users`",
		"SELECT `id, name, email` FROM `users`",
	)
	// With unsafe it still gets quoted because unsafe has been applied after
	// the column names has been added.
	compareToSQL(t, NewSelect("id, name, email").Unsafe().From("users"), nil,
		"SELECT `id, name, email` FROM `users`",
		"SELECT `id, name, email` FROM `users`",
	)
	// correct way to handle it
	compareToSQL(t, NewSelect("id", "name", "email").From("users"), nil,
		"SELECT `id`, `name`, `email` FROM `users`",
		"SELECT `id`, `name`, `email` FROM `users`",
	)
}

func TestSelect_Load_Slice_Scanner(t *testing.T) {
	s := createRealSessionWithFixtures(t, nil)

	var people dbrPersons
	count, err := s.Select("id", "name", "email").From("dbr_people").OrderBy("id").Load(context.TODO(), &people)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)

	assert.Equal(t, len(people.Data), 2)
	if len(people.Data) == 2 {
		// Make sure that the Ids are isSet. It'ab possible (maybe?) that different DBs isSet ids differently so
		// don't assume they're 1 and 2.
		assert.True(t, people.Data[0].ID > 0)
		assert.True(t, people.Data[1].ID > people.Data[0].ID)

		assert.Equal(t, "Jonathan", people.Data[0].Name)
		assert.True(t, people.Data[0].Email.Valid)
		assert.Equal(t, "jonathan@uservoice.com", people.Data[0].Email.String)
		assert.Equal(t, "Dmitri", people.Data[1].Name)
		assert.True(t, people.Data[1].Email.Valid)
		assert.Equal(t, "zavorotni@jadius.com", people.Data[1].Email.String)
	}
}

func TestSelect_Load_Rows(t *testing.T) {
	s := createRealSessionWithFixtures(t, nil)

	t.Run("found", func(t *testing.T) {
		var person dbrPerson
		_, err := s.Select("id", "name", "email").From("dbr_people").
			Where(Column("email").Str("jonathan@uservoice.com")).Load(context.TODO(), &person)
		assert.NoError(t, err)
		assert.True(t, person.ID > 0)
		assert.Equal(t, "Jonathan", person.Name)
		assert.True(t, person.Email.Valid)
		assert.Equal(t, "jonathan@uservoice.com", person.Email.String)
	})

	t.Run("not found", func(t *testing.T) {
		var person2 dbrPerson
		count, err := s.Select("id", "name", "email").From("dbr_people").
			Where(Column("email").Str("dontexist@uservoice.com")).Load(context.TODO(), &person2)

		require.NoError(t, err, "%+v", err)
		assert.Exactly(t, dbrPerson{}, person2)
		assert.Empty(t, count, "Should have no rows loaded")
	})
}

func TestSelectBySQL_Load_Slice(t *testing.T) {
	s := createRealSessionWithFixtures(t, nil)

	t.Run("single slice item", func(t *testing.T) {
		var people dbrPersons
		count, err := s.SelectBySQL("SELECT name FROM dbr_people WHERE email = ?", MakeArgs(1).Str("jonathan@uservoice.com")).Load(context.TODO(), &people)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)
		if len(people.Data) == 1 {
			assert.Equal(t, "Jonathan", people.Data[0].Name)
			assert.Equal(t, uint64(0), people.Data[0].ID)      // not set
			assert.Equal(t, false, people.Data[0].Email.Valid) // not set
			assert.Equal(t, "", people.Data[0].Email.String)   // not set
		}
	})

	t.Run("IN Clause", func(t *testing.T) {
		ids, err := s.Select("id").From("dbr_people").
			Where(Column("id").In().Int64s(1, 2, 3)).LoadInt64s(context.TODO())
		assert.NoError(t, err)
		assert.Exactly(t, []int64{1, 2}, ids)
	})
	t.Run("NOT IN Clause", func(t *testing.T) {
		ids, err := s.Select("id").From("dbr_people").
			Where(Column("id").NotIn().Int64s(2, 3)).LoadInt64s(context.TODO())
		assert.NoError(t, err)
		assert.Exactly(t, []int64{1}, ids)
	})
}

func TestSelect_LoadType_Single(t *testing.T) {
	s := createRealSessionWithFixtures(t, nil)

	t.Run("LoadString", func(t *testing.T) {
		name, err := s.Select("name").From("dbr_people").Where(Column("email").Str("jonathan@uservoice.com")).LoadString(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, "Jonathan", name)
	})
	t.Run("LoadString too many columns", func(t *testing.T) {
		name, err := s.Select("name", "email").From("dbr_people").Where(Expr("email = 'jonathan@uservoice.com'")).LoadString(context.TODO())
		assert.Error(t, err, "%+v", err)
		assert.Empty(t, name)
	})
	t.Run("LoadString not found", func(t *testing.T) {
		name, err := s.Select("name").From("dbr_people").Where(Expr("email = 'notfound@example.com'")).LoadString(context.TODO())
		assert.True(t, errors.IsNotFound(err), "%+v", err)
		assert.Empty(t, name)
	})

	t.Run("LoadInt64", func(t *testing.T) {
		id, err := s.Select("id").From("dbr_people").Limit(1).LoadInt64(context.TODO())
		assert.NoError(t, err)
		assert.True(t, id > 0)
	})
	t.Run("LoadInt64 too many columns", func(t *testing.T) {
		id, err := s.Select("id", "email").From("dbr_people").Limit(1).LoadInt64(context.TODO())
		assert.Error(t, err, "%+v", err)
		assert.Empty(t, id)
	})
	t.Run("LoadInt64 not found", func(t *testing.T) {
		id, err := s.Select("id").From("dbr_people").Where(Expr("id=236478326")).LoadInt64(context.TODO())
		assert.True(t, errors.IsNotFound(err), "%+v", err)
		assert.Empty(t, id)
	})

	t.Run("LoadUint64", func(t *testing.T) {
		id, err := s.Select("id").From("dbr_people").Limit(1).LoadUint64(context.TODO())
		assert.NoError(t, err)
		assert.True(t, id > 0)
	})
	t.Run("LoadUint64 too many columns", func(t *testing.T) {
		id, err := s.Select("id", "email").From("dbr_people").Limit(1).LoadUint64(context.TODO())
		assert.Error(t, err, "%+v", err)
		assert.Empty(t, id)
	})
	t.Run("LoadUint64 not found", func(t *testing.T) {
		id, err := s.Select("id").From("dbr_people").Where(Expr("id=236478326")).LoadUint64(context.TODO())
		assert.True(t, errors.IsNotFound(err), "%+v", err)
		assert.Empty(t, id)
	})

	t.Run("LoadFloat64", func(t *testing.T) {
		id, err := s.Select("id").From("dbr_people").Limit(1).LoadFloat64(context.TODO())
		assert.NoError(t, err)
		assert.True(t, id > 0)
	})
	t.Run("LoadFloat64 too many columns", func(t *testing.T) {
		id, err := s.Select("id", "email").From("dbr_people").Limit(1).LoadFloat64(context.TODO())
		assert.Error(t, err, "%+v", err)
		assert.Empty(t, id)
	})
	t.Run("LoadFloat64 not found", func(t *testing.T) {
		id, err := s.Select("id").From("dbr_people").Where(Expr("id=236478326")).LoadFloat64(context.TODO())
		assert.True(t, errors.IsNotFound(err), "%+v", err)
		assert.Empty(t, id)
	})
}

func TestSelect_LoadUint64(t *testing.T) {
	s := createRealSessionWithFixtures(t, &installFixturesConfig{
		AddPeopleWithMaxUint64: true,
	})
	// Despite it seems that Go can support large uint64 values ... the down side is that
	// the byte encoded uint64 gets transferred as a string and MySQL/MariaDB must convert that
	// string into a bigint.
	const bigID uint64 = 18446744073700551613 // see also file dbr_test.go MaxUint64

	sel := s.Select("id").From("dbr_people").Where(Column("id").Uint64(bigID))

	t.Run("MaxUint64 prepared stmt o:equal", func(t *testing.T) {
		id, err := sel.LoadUint64(context.TODO())
		require.NoError(t, err)
		assert.Exactly(t, bigID, id)
	})
	t.Run("MaxUint64 interpolated o:equal", func(t *testing.T) {
		id, err := sel.Interpolate().LoadUint64(context.TODO())
		require.NoError(t, err)
		assert.Exactly(t, bigID, id)
	})

}

func TestSelect_LoadType_Slices(t *testing.T) {
	s := createRealSessionWithFixtures(t, nil)

	t.Run("LoadStrings", func(t *testing.T) {
		names, err := s.Select("name").From("dbr_people").LoadStrings(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, []string{"Jonathan", "Dmitri"}, names)
	})
	t.Run("LoadStrings too many columns", func(t *testing.T) {
		vals, err := s.Select("name", "email").From("dbr_people").LoadStrings(context.TODO())
		assert.Error(t, err, "%+v", err)
		assert.Exactly(t, []string(nil), vals)
	})
	t.Run("LoadStrings not found", func(t *testing.T) {
		names, err := s.Select("name").From("dbr_people").Where(Expr("name ='jdhsjdf'")).LoadStrings(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, []string{}, names)
	})

	t.Run("LoadInt64s", func(t *testing.T) {
		names, err := s.Select("id").From("dbr_people").LoadInt64s(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, []int64{1, 2}, names)
	})
	t.Run("LoadInt64s too many columns", func(t *testing.T) {
		vals, err := s.Select("id", "email").From("dbr_people").LoadInt64s(context.TODO())
		assert.Error(t, err, "%+v", err)
		assert.Exactly(t, []int64(nil), vals)
	})
	t.Run("LoadInt64s not found", func(t *testing.T) {
		names, err := s.Select("id").From("dbr_people").Where(Expr("name ='jdhsjdf'")).LoadInt64s(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, []int64{}, names)
	})

	t.Run("LoadUint64s", func(t *testing.T) {
		names, err := s.Select("id").From("dbr_people").LoadUint64s(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, []uint64{1, 2}, names)
	})
	t.Run("LoadUint64s too many columns", func(t *testing.T) {
		vals, err := s.Select("id", "email").From("dbr_people").LoadUint64s(context.TODO())
		assert.Error(t, err, "%+v", err)
		assert.Exactly(t, []uint64(nil), vals)
	})
	t.Run("LoadUint64s not found", func(t *testing.T) {
		names, err := s.Select("id").From("dbr_people").Where(Expr("name ='jdhsjdf'")).LoadUint64s(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, []uint64{}, names)
	})

	t.Run("LoadFloat64s", func(t *testing.T) {
		names, err := s.Select("id").From("dbr_people").LoadFloat64s(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, []float64{1, 2}, names)
	})
	t.Run("LoadFloat64s too many columns", func(t *testing.T) {
		vals, err := s.Select("id", "email").From("dbr_people").LoadFloat64s(context.TODO())
		assert.Error(t, err, "%+v", err)
		assert.Exactly(t, []float64(nil), vals)
	})
	t.Run("LoadFloat64s not found", func(t *testing.T) {
		names, err := s.Select("id").From("dbr_people").Where(Expr("name ='jdhsjdf'")).LoadFloat64s(context.TODO())
		assert.NoError(t, err)
		assert.Equal(t, []float64{}, names)
	})

}

func TestSelectJoin(t *testing.T) {
	t.Parallel()
	s := createRealSessionWithFixtures(t, nil)

	t.Run("inner, distinct, no cache, high prio", func(t *testing.T) {
		sqlObj := s.
			Select("p1.*", "p2.*").
			Distinct().StraightJoin().SQLNoCache().
			FromAlias("dbr_people", "p1").
			Join(
				MakeIdentifier("dbr_people").Alias("p2"),
				Expr("`p2`.`id` = `p1`.`id`"),
				Column("p1.id").Int(42),
			)

		compareToSQL(t, sqlObj, nil,
			"SELECT DISTINCT STRAIGHT_JOIN SQL_NO_CACHE `p1`.*, `p2`.* FROM `dbr_people` AS `p1` INNER JOIN `dbr_people` AS `p2` ON (`p2`.`id` = `p1`.`id`) AND (`p1`.`id` = ?)",
			"SELECT DISTINCT STRAIGHT_JOIN SQL_NO_CACHE `p1`.*, `p2`.* FROM `dbr_people` AS `p1` INNER JOIN `dbr_people` AS `p2` ON (`p2`.`id` = `p1`.`id`) AND (`p1`.`id` = 42)",
			int64(42),
		)

	})

	t.Run("inner", func(t *testing.T) {
		sqlObj := s.
			Select("p1.*", "p2.*").
			FromAlias("dbr_people", "p1").
			Join(
				MakeIdentifier("dbr_people").Alias("p2"),
				Expr("`p2`.`id` = `p1`.`id`"),
				Column("p1.id").Int(42),
			)

		compareToSQL(t, sqlObj, nil,
			"SELECT `p1`.*, `p2`.* FROM `dbr_people` AS `p1` INNER JOIN `dbr_people` AS `p2` ON (`p2`.`id` = `p1`.`id`) AND (`p1`.`id` = ?)",
			"SELECT `p1`.*, `p2`.* FROM `dbr_people` AS `p1` INNER JOIN `dbr_people` AS `p2` ON (`p2`.`id` = `p1`.`id`) AND (`p1`.`id` = 42)",
			int64(42),
		)
	})

	t.Run("left", func(t *testing.T) {
		sqlObj := s.
			Select("p1.*", "p2.name").
			FromAlias("dbr_people", "p1").
			LeftJoin(
				MakeIdentifier("dbr_people").Alias("p2"),
				Expr("`p2`.`id` = `p1`.`id`"),
				Column("p1.id").Int(42),
			)

		compareToSQL(t, sqlObj, nil,
			"SELECT `p1`.*, `p2`.`name` FROM `dbr_people` AS `p1` LEFT JOIN `dbr_people` AS `p2` ON (`p2`.`id` = `p1`.`id`) AND (`p1`.`id` = ?)",
			"SELECT `p1`.*, `p2`.`name` FROM `dbr_people` AS `p1` LEFT JOIN `dbr_people` AS `p2` ON (`p2`.`id` = `p1`.`id`) AND (`p1`.`id` = 42)",
			int64(42),
		)
	})

	t.Run("right", func(t *testing.T) {
		sqlObj := s.
			Select("p1.*").
			AddColumnsAliases("p2.name", "p2Name", "p2.email", "p2Email", "id", "internalID").
			FromAlias("dbr_people", "p1").
			RightJoin(
				MakeIdentifier("dbr_people").Alias("p2"),
				Expr("`p2`.`id` = `p1`.`id`"),
			)
		compareToSQL(t, sqlObj, nil,
			"SELECT `p1`.*, `p2`.`name` AS `p2Name`, `p2`.`email` AS `p2Email`, `id` AS `internalID` FROM `dbr_people` AS `p1` RIGHT JOIN `dbr_people` AS `p2` ON (`p2`.`id` = `p1`.`id`)",
			"SELECT `p1`.*, `p2`.`name` AS `p2Name`, `p2`.`email` AS `p2Email`, `id` AS `internalID` FROM `dbr_people` AS `p1` RIGHT JOIN `dbr_people` AS `p2` ON (`p2`.`id` = `p1`.`id`)",
		)
	})

	t.Run("using", func(t *testing.T) {
		sqlObj := s.
			Select("p1.*").
			AddColumnsAliases("p2.name", "p2Name", "p2.email", "p2Email").
			FromAlias("dbr_people", "p1").
			RightJoin(
				MakeIdentifier("dbr_people").Alias("p2"),
				Columns("id", "email"),
			)
		compareToSQL(t, sqlObj, nil,
			"SELECT `p1`.*, `p2`.`name` AS `p2Name`, `p2`.`email` AS `p2Email` FROM `dbr_people` AS `p1` RIGHT JOIN `dbr_people` AS `p2` USING (`id`,`email`)",
			"SELECT `p1`.*, `p2`.`name` AS `p2Name`, `p2`.`email` AS `p2Email` FROM `dbr_people` AS `p1` RIGHT JOIN `dbr_people` AS `p2` USING (`id`,`email`)",
		)
	})
}

func TestSelect_Locks(t *testing.T) {
	t.Parallel()
	t.Run("LOCK IN SHARE MODE", func(t *testing.T) {
		s := NewSelect("p1.*").
			AddColumnsAliases("p2.name", "p2Name", "p2.email", "p2Email").
			FromAlias("dbr_people", "p1").LockInShareMode()
		compareToSQL(t, s, nil,
			"SELECT `p1`.*, `p2`.`name` AS `p2Name`, `p2`.`email` AS `p2Email` FROM `dbr_people` AS `p1` LOCK IN SHARE MODE",
			"SELECT `p1`.*, `p2`.`name` AS `p2Name`, `p2`.`email` AS `p2Email` FROM `dbr_people` AS `p1` LOCK IN SHARE MODE",
		)
	})
	t.Run("FOR UPDATE", func(t *testing.T) {
		s := NewSelect("p1.*").
			AddColumnsAliases("p2.name", "p2Name", "p2.email", "p2Email").
			FromAlias("dbr_people", "p1").ForUpdate()
		compareToSQL(t, s, nil,
			"SELECT `p1`.*, `p2`.`name` AS `p2Name`, `p2`.`email` AS `p2Email` FROM `dbr_people` AS `p1` FOR UPDATE",
			"SELECT `p1`.*, `p2`.`name` AS `p2Name`, `p2`.`email` AS `p2Email` FROM `dbr_people` AS `p1` FOR UPDATE",
		)
	})
}

func TestSelect_Events(t *testing.T) {
	t.Parallel()

	t.Run("Stop Propagation", func(t *testing.T) {
		d := NewSelect("a", "b").FromAlias("tableA", "tA")
		d.OrderBy("col3")

		d.Log = log.BlackHole{EnableInfo: true, EnableDebug: true}
		d.Listeners.Add(
			Listen{
				Name:      "listener1",
				EventType: OnBeforeToSQL,
				SelectFunc: func(b *Select) {
					b.OrderByDesc("col1")
				},
			},
			Listen{
				Name:      "listener2",
				EventType: OnBeforeToSQL,
				SelectFunc: func(b *Select) {
					b.OrderByDesc("col2")
					b.PropagationStopped = true
				},
			},
			Listen{
				Name:      "listener3",
				EventType: OnBeforeToSQL,
				SelectFunc: func(b *Select) {
					panic("Should not get called")
				},
			},
		)
		compareToSQL(t, d, nil,
			"SELECT `a`, `b` FROM `tableA` AS `tA` ORDER BY `col3`, `col1` DESC, `col2` DESC",
			"SELECT `a`, `b` FROM `tableA` AS `tA` ORDER BY `col3`, `col1` DESC, `col2` DESC, `col1` DESC, `col2` DESC",
		)
	})

	t.Run("Missing EventType", func(t *testing.T) {
		s := NewSelect("a", "b").FromAlias("tableA", "tA")
		s.OrderBy("col3")
		s.Listeners.Add(Listen{
			Name: "a col1",
			SelectFunc: func(s2 *Select) {
				s2.Where(Column("a").Float64(3.14159))
				s2.OrderByDesc("col1")
			},
		})
		compareToSQL(t, s, errors.IsEmpty,
			"",
			"",
		)
	})

	t.Run("Should Dispatch", func(t *testing.T) {
		s := NewSelect("a", "b").FromAlias("tableA", "tA")
		s.OrderBy("col3")
		s.Listeners.Add(Listen{
			Name:      "a col1",
			Once:      true,
			EventType: OnBeforeToSQL,
			SelectFunc: func(s2 *Select) {
				s2.Where(Column("a").Float64(3.14159))
				s2.OrderByDesc("col1")
			},
		})
		s.Listeners.Add(Listen{
			Name:      "b col2",
			EventType: OnBeforeToSQL,
			SelectFunc: func(s2 *Select) {
				s2.OrderByDesc("col2").
					Where(Column("b").Str("a"))
			},
		})

		compareToSQL(t, s, nil,
			"SELECT `a`, `b` FROM `tableA` AS `tA` WHERE (`a` = ?) AND (`b` = ?) ORDER BY `col3`, `col1` DESC, `col2` DESC",
			"SELECT `a`, `b` FROM `tableA` AS `tA` WHERE (`a` = 3.14159) AND (`b` = 'a') AND (`b` = 'a') ORDER BY `col3`, `col1` DESC, `col2` DESC, `col2` DESC",
			3.14159, "a",
		)

		assert.Exactly(t, `a col1; b col2`, s.Listeners.String())
	})
}

func TestSelect_Columns(t *testing.T) {
	t.Parallel()

	t.Run("AddColumns, multiple args", func(t *testing.T) {
		s := NewSelect("a", "b")
		s.FromAlias("tableA", "tA")
		s.AddColumns("d,e, f", "g", "h", "i,j ,k")
		compareToSQL(t, s, nil,
			"SELECT `a`, `b`, `d,e, f`, `g`, `h`, `i,j ,k` FROM `tableA` AS `tA`",
			"SELECT `a`, `b`, `d,e, f`, `g`, `h`, `i,j ,k` FROM `tableA` AS `tA`",
		)
	})
	t.Run("AddColumns, each column itself", func(t *testing.T) {
		s := NewSelect("a", "b")
		s.FromAlias("tableA", "tA")
		s.AddColumns("d", "e", "f")
		compareToSQL(t, s, nil,
			"SELECT `a`, `b`, `d`, `e`, `f` FROM `tableA` AS `tA`",
			"SELECT `a`, `b`, `d`, `e`, `f` FROM `tableA` AS `tA`",
		)
	})
	t.Run("AddColumnsAliases Expression Quoted", func(t *testing.T) {
		s := NewSelect().From("t3").
			AddColumnsAliases("x", "u", "y", "v").
			AddColumnsAliases("SUM(price)", "total_price")
		compareToSQL(t, s, nil,
			"SELECT `x` AS `u`, `y` AS `v`, `SUM(price)` AS `total_price` FROM `t3`",
			"SELECT `x` AS `u`, `y` AS `v`, `SUM(price)` AS `total_price` FROM `t3`",
		)
	})
	t.Run("AddColumns+AddColumnsExprAlias", func(t *testing.T) {
		s := NewSelect().From("t3").
			AddColumns("t3.name", "sku").
			AddColumnsConditions(Expr("SUM(price)").Alias("total_price"))
		compareToSQL(t, s, nil,
			"SELECT `t3`.`name`, `sku`, SUM(price) AS `total_price` FROM `t3`",
			"SELECT `t3`.`name`, `sku`, SUM(price) AS `total_price` FROM `t3`",
		)
	})

	t.Run("AddColumnsAliases multi", func(t *testing.T) {
		s := NewSelect().From("t3").
			AddColumnsAliases("t3.name", "t3Name", "t3.sku", "t3SKU")
		compareToSQL(t, s, nil,
			"SELECT `t3`.`name` AS `t3Name`, `t3`.`sku` AS `t3SKU` FROM `t3`",
			"SELECT `t3`.`name` AS `t3Name`, `t3`.`sku` AS `t3SKU` FROM `t3`",
		)
	})
	t.Run("AddColumnsAliases imbalanced", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				if err, ok := r.(error); ok {
					assert.True(t, errors.IsMismatch(err), "%+v", err)
				} else {
					t.Errorf("Panic should contain an error but got:\n%+v", r)
				}
			} else {
				t.Error("Expecting a panic but got nothing")
			}
		}()
		NewSelect().From("t3").
			AddColumnsAliases("t3.name", "t3Name", "t3.sku")

	})

	t.Run("AddColumnsExprAlias", func(t *testing.T) {
		s := NewSelect().FromAlias("sales_bestsellers_aggregated_daily", "t3").
			AddColumnsConditions(Expr("DATE_FORMAT(t3.period, '%Y-%m-01')").Alias("period"))
		compareToSQL(t, s, nil,
			"SELECT DATE_FORMAT(t3.period, '%Y-%m-01') AS `period` FROM `sales_bestsellers_aggregated_daily` AS `t3`",
			"SELECT DATE_FORMAT(t3.period, '%Y-%m-01') AS `period` FROM `sales_bestsellers_aggregated_daily` AS `t3`",
		)
	})
	t.Run("AddColumns with expression incorrect", func(t *testing.T) {
		s := NewSelect().AddColumns(" `t.value`", "`t`.`attribute_id`", "t.{column} AS `col_type`").FromAlias("catalog_product_entity_{type}", "t")
		compareToSQL(t, s, nil,
			"SELECT ` t`.`value`, `t`.`attribute_id`, `t`.`{column} AS col_type` FROM `catalog_product_entity_{type}` AS `t`",
			"SELECT ` t`.`value`, `t`.`attribute_id`, `t`.`{column} AS col_type` FROM `catalog_product_entity_{type}` AS `t`",
		)
	})
}

func TestSubSelect(t *testing.T) {
	t.Parallel()
	sub := NewSelect().From("catalog_category_product").
		AddColumns("entity_id").Where(Column("category_id").Int64(234))

	runner := func(op Op, wantSQL string) func(*testing.T) {
		c := Column("entity_id").Sub(sub)
		c.Operator = op
		return func(t *testing.T) {
			s := NewSelect("sku", "type_id").
				From("catalog_product_entity").
				Where(c)
			compareToSQL(t, s, nil, wantSQL, "", int64(234))
		}
	}
	t.Run("IN", runner(In,
		"SELECT `sku`, `type_id` FROM `catalog_product_entity` WHERE (`entity_id` IN (SELECT `entity_id` FROM `catalog_category_product` WHERE (`category_id` = ?)))",
	))
	t.Run("EXISTS", runner(Exists,
		"SELECT `sku`, `type_id` FROM `catalog_product_entity` WHERE (`entity_id` EXISTS (SELECT `entity_id` FROM `catalog_category_product` WHERE (`category_id` = ?)))",
	))
	t.Run("NOT EXISTS", runner(NotExists,
		"SELECT `sku`, `type_id` FROM `catalog_product_entity` WHERE (`entity_id` NOT EXISTS (SELECT `entity_id` FROM `catalog_category_product` WHERE (`category_id` = ?)))",
	))
	t.Run("NOT EQUAL", runner(NotEqual,
		"SELECT `sku`, `type_id` FROM `catalog_product_entity` WHERE (`entity_id` != (SELECT `entity_id` FROM `catalog_category_product` WHERE (`category_id` = ?)))",
	))
	t.Run("NOT EQUAL", runner(Equal,
		"SELECT `sku`, `type_id` FROM `catalog_product_entity` WHERE (`entity_id` = (SELECT `entity_id` FROM `catalog_category_product` WHERE (`category_id` = ?)))",
	))
}

func TestSelect_Subselect_Complex(t *testing.T) {
	t.Parallel()
	/* Something like:
	   SELECT
	     `t1`.`store_id`,
	     `t1`.`product_id`,
	     `t1`.`product_name`,
	     `t1`.`product_price`,
	     `t1`.`qty_ordered`
	   FROM (
	          SELECT
	            `t2`.`store_id`,
	            `t2`.`product_id`,
	            `t2`.`product_name`,
	            `t2`.`product_price`,
	            `t2`.`total_qty` AS `qty_ordered`
	          FROM (
	                 SELECT
	                   `t3`.`store_id`,
	                   `t3`.`product_id`,
	                   `t3`.`product_name`,
	                   AVG(`t3`.`product_price`) as `avg_price`,
	                   SUM(t3.qty_ordered) AS `total_qty`
	                 FROM `sales_bestsellers_aggregated_daily` AS `t3`
	                 GROUP BY `t3`.`store_id`,
	                   Date_format(t3.period, '%Y-%m-01'),
	                   `t3`.`product_id`
	                 ORDER BY `t3`.`store_id` ASC,
	                   Date_format(t3.period, '%Y-%m-01'),
	                   `total_qty` DESC
	               ) AS `t2`
	        ) AS `t1`
	*/

	t.Run("without args", func(t *testing.T) {
		sel3 := NewSelect().FromAlias("sales_bestsellers_aggregated_daily", "t3").
			Unsafe().
			AddColumnsConditions(Expr("DATE_FORMAT(t3.period, '%Y-%m-01')").Alias("period")).
			AddColumns("`t3`.`store_id`,`t3`.`product_id`,`t3`.`product_name`").
			AddColumnsConditions(
				Expr("AVG(`t3`.`product_price`)").Alias("avg_price"),
				Expr("SUM(t3.qty_ordered)").Alias("total_qty"),
			).
			GroupBy("t3.store_id").
			GroupBy("DATE_FORMAT(t3.period, '%Y-%m-01')").
			GroupBy("t3.product_id", "t3.product_name").
			OrderBy("t3.store_id").
			OrderBy("DATE_FORMAT(t3.period, '%Y-%m-01')").
			OrderByDesc("total_qty")

		sel2 := NewSelectWithDerivedTable(sel3, "t2").
			AddColumns("t2.period", "t2.store_id", "t2.product_id", "t2.product_name", "t2.avg_price").
			AddColumnsAliases("`t2`.`total_qty`", "`qty_ordered`")

		sel1 := NewSelectWithDerivedTable(sel2, "t1").
			AddColumns("t1.period", "t1.store_id", "t1.product_id", "t1.product_name", "t1.avg_price", "t1.qty_ordered").
			OrderBy("`t1`.period", "`t1`.product_id")

		compareToSQL(t, sel1, nil,
			"SELECT `t1`.`period`, `t1`.`store_id`, `t1`.`product_id`, `t1`.`product_name`, `t1`.`avg_price`, `t1`.`qty_ordered` FROM (SELECT `t2`.`period`, `t2`.`store_id`, `t2`.`product_id`, `t2`.`product_name`, `t2`.`avg_price`, `t2`.`total_qty` AS `qty_ordered` FROM (SELECT DATE_FORMAT(t3.period, '%Y-%m-01') AS `period`, `t3`.`store_id`,`t3`.`product_id`,`t3`.`product_name`, AVG(`t3`.`product_price`) AS `avg_price`, SUM(t3.qty_ordered) AS `total_qty` FROM `sales_bestsellers_aggregated_daily` AS `t3` GROUP BY `t3`.`store_id`, DATE_FORMAT(t3.period, '%Y-%m-01'), `t3`.`product_id`, `t3`.`product_name` ORDER BY `t3`.`store_id`, DATE_FORMAT(t3.period, '%Y-%m-01'), `total_qty` DESC) AS `t2`) AS `t1` ORDER BY `t1`.`period`, `t1`.`product_id`",
			"SELECT `t1`.`period`, `t1`.`store_id`, `t1`.`product_id`, `t1`.`product_name`, `t1`.`avg_price`, `t1`.`qty_ordered` FROM (SELECT `t2`.`period`, `t2`.`store_id`, `t2`.`product_id`, `t2`.`product_name`, `t2`.`avg_price`, `t2`.`total_qty` AS `qty_ordered` FROM (SELECT DATE_FORMAT(t3.period, '%Y-%m-01') AS `period`, `t3`.`store_id`,`t3`.`product_id`,`t3`.`product_name`, AVG(`t3`.`product_price`) AS `avg_price`, SUM(t3.qty_ordered) AS `total_qty` FROM `sales_bestsellers_aggregated_daily` AS `t3` GROUP BY `t3`.`store_id`, DATE_FORMAT(t3.period, '%Y-%m-01'), `t3`.`product_id`, `t3`.`product_name` ORDER BY `t3`.`store_id`, DATE_FORMAT(t3.period, '%Y-%m-01'), `total_qty` DESC) AS `t2`) AS `t1` ORDER BY `t1`.`period`, `t1`.`product_id`",
		)
	})

	t.Run("with args", func(t *testing.T) {
		// Full valid query which works in a M1 and M2 database.
		sel3 := NewSelect().FromAlias("sales_bestsellers_aggregated_daily", "t3").
			Unsafe().
			AddColumnsConditions(Expr("DATE_FORMAT(t3.period, '%Y-%m-01')").Alias("period")).
			AddColumns("`t3`.`store_id`,`t3`.`product_id`,`t3`.`product_name`").
			AddColumnsConditions(
				Expr("AVG(`t3`.`product_price`)").Alias("avg_price"),
				Expr("SUM(t3.qty_ordered)+?").Alias("total_qty").Float64(3.141),
			).
			GroupBy("t3.store_id").
			GroupBy("DATE_FORMAT(t3.period, '%Y-%m-01')").
			GroupBy("t3.product_id", "t3.product_name").
			Having(Expr("COUNT(*)>?").Int(3)).
			OrderBy("t3.store_id").
			OrderBy("DATE_FORMAT(t3.period, '%Y-%m-01')").
			OrderByDesc("total_qty").
			Where(Column("t3.store_id").In().Int64s(2, 3, 4))

		sel2 := NewSelectWithDerivedTable(sel3, "t2").
			AddColumns("t2.period", "t2.store_id", "t2.product_id", "t2.product_name", "t2.avg_price").
			AddColumnsAliases("t2.total_qty", "qty_ordered")

		sel1 := NewSelectWithDerivedTable(sel2, "t1").
			AddColumns("t1.period", "t1.store_id", "t1.product_id", "t1.product_name", "t1.avg_price", "t1.qty_ordered").
			OrderBy("`t1`.period", "`t1`.product_id")

		compareToSQL(t, sel1, nil,
			"SELECT `t1`.`period`, `t1`.`store_id`, `t1`.`product_id`, `t1`.`product_name`, `t1`.`avg_price`, `t1`.`qty_ordered` FROM (SELECT `t2`.`period`, `t2`.`store_id`, `t2`.`product_id`, `t2`.`product_name`, `t2`.`avg_price`, `t2`.`total_qty` AS `qty_ordered` FROM (SELECT DATE_FORMAT(t3.period, '%Y-%m-01') AS `period`, `t3`.`store_id`,`t3`.`product_id`,`t3`.`product_name`, AVG(`t3`.`product_price`) AS `avg_price`, SUM(t3.qty_ordered)+? AS `total_qty` FROM `sales_bestsellers_aggregated_daily` AS `t3` WHERE (`t3`.`store_id` IN (?,?,?)) GROUP BY `t3`.`store_id`, DATE_FORMAT(t3.period, '%Y-%m-01'), `t3`.`product_id`, `t3`.`product_name` HAVING (COUNT(*)>?) ORDER BY `t3`.`store_id`, DATE_FORMAT(t3.period, '%Y-%m-01'), `total_qty` DESC) AS `t2`) AS `t1` ORDER BY `t1`.`period`, `t1`.`product_id`",
			"SELECT `t1`.`period`, `t1`.`store_id`, `t1`.`product_id`, `t1`.`product_name`, `t1`.`avg_price`, `t1`.`qty_ordered` FROM (SELECT `t2`.`period`, `t2`.`store_id`, `t2`.`product_id`, `t2`.`product_name`, `t2`.`avg_price`, `t2`.`total_qty` AS `qty_ordered` FROM (SELECT DATE_FORMAT(t3.period, '%Y-%m-01') AS `period`, `t3`.`store_id`,`t3`.`product_id`,`t3`.`product_name`, AVG(`t3`.`product_price`) AS `avg_price`, SUM(t3.qty_ordered)+3.141 AS `total_qty` FROM `sales_bestsellers_aggregated_daily` AS `t3` WHERE (`t3`.`store_id` IN (2,3,4)) GROUP BY `t3`.`store_id`, DATE_FORMAT(t3.period, '%Y-%m-01'), `t3`.`product_id`, `t3`.`product_name` HAVING (COUNT(*)>3) ORDER BY `t3`.`store_id`, DATE_FORMAT(t3.period, '%Y-%m-01'), `total_qty` DESC) AS `t2`) AS `t1` ORDER BY `t1`.`period`, `t1`.`product_id`",
			3.141, int64(2), int64(3), int64(4), int64(3),
		)
	})
}

func TestSelect_Subselect_Compact(t *testing.T) {
	t.Parallel()

	sel2 := NewSelect().FromAlias("sales_bestsellers_aggregated_daily", "t3").
		AddColumns("`t3`.`product_name`").
		Where(Column("t3.store_id").In().Int64s(2, 3, 4)).
		GroupBy("t3.store_id").
		Having(Expr("COUNT(*)>?").Int(5))

	sel := NewSelectWithDerivedTable(sel2, "t2").
		AddColumns("t2.product_name").
		Where(Column("t2.floatcol").Equal().Float64(3.14159))

	compareToSQL(t, sel, nil,
		"SELECT `t2`.`product_name` FROM (SELECT `t3`.`product_name` FROM `sales_bestsellers_aggregated_daily` AS `t3` WHERE (`t3`.`store_id` IN (?,?,?)) GROUP BY `t3`.`store_id` HAVING (COUNT(*)>?)) AS `t2` WHERE (`t2`.`floatcol` = ?)",
		"SELECT `t2`.`product_name` FROM (SELECT `t3`.`product_name` FROM `sales_bestsellers_aggregated_daily` AS `t3` WHERE (`t3`.`store_id` IN (2,3,4)) GROUP BY `t3`.`store_id` HAVING (COUNT(*)>5)) AS `t2` WHERE (`t2`.`floatcol` = 3.14159)",
		int64(2), int64(3), int64(4), int64(5), 3.14159,
	)
}

func TestSelect_ParenthesisOpen_Close(t *testing.T) {
	t.Parallel()
	t.Run("beginning of WHERE", func(t *testing.T) {

		sel := NewSelect("a", "b").
			FromAlias("c", "cc").
			Where(
				ParenthesisOpen(),
				Column("d").Int(1),
				Column("e").Str("wat").Or(),
				ParenthesisClose(),
				Column("f").Float64(2.7182),
			).
			GroupBy("ab").
			Having(
				ParenthesisOpen(),
				Column("m").Int(33),
				Column("n").Str("wh3r3").Or(),
				ParenthesisClose(),
				Expr("j = k"),
			)
		compareToSQL(t, sel, nil,
			"SELECT `a`, `b` FROM `c` AS `cc` WHERE ((`d` = ?) OR (`e` = ?)) AND (`f` = ?) GROUP BY `ab` HAVING ((`m` = ?) OR (`n` = ?)) AND (j = k)",
			"SELECT `a`, `b` FROM `c` AS `cc` WHERE ((`d` = 1) OR (`e` = 'wat')) AND (`f` = 2.7182) GROUP BY `ab` HAVING ((`m` = 33) OR (`n` = 'wh3r3')) AND (j = k)",
			int64(1), "wat", 2.7182, int64(33), "wh3r3")

	})

	t.Run("end of WHERE", func(t *testing.T) {
		sel := NewSelect("a", "b").
			FromAlias("c", "cc").
			Where(
				Column("f").Float64(2.7182),
				ParenthesisOpen(),
				Column("d").Int(1),
				Column("e").Str("wat").Or(),
				ParenthesisClose(),
			).
			GroupBy("ab").
			Having(
				Expr("j = k"),
				ParenthesisOpen(),
				Column("m").Int(33),
				Column("n").Str("wh3r3").Or(),
				ParenthesisClose(),
			)
		compareToSQL(t, sel, nil,
			"SELECT `a`, `b` FROM `c` AS `cc` WHERE (`f` = ?) AND ((`d` = ?) OR (`e` = ?)) GROUP BY `ab` HAVING (j = k) AND ((`m` = ?) OR (`n` = ?))",
			"SELECT `a`, `b` FROM `c` AS `cc` WHERE (`f` = 2.7182) AND ((`d` = 1) OR (`e` = 'wat')) GROUP BY `ab` HAVING (j = k) AND ((`m` = 33) OR (`n` = 'wh3r3'))",
			2.7182, int64(1), "wat", int64(33), "wh3r3")
	})

	t.Run("middle of WHERE", func(t *testing.T) {
		sel := NewSelect("a", "b").
			FromAlias("c", "cc").
			Where(
				Column("f").Float64(2.7182),
				ParenthesisOpen(),
				Column("d").Int(1),
				Column("e").Str("wat").Or(),
				ParenthesisClose(),
				Column("p").Float64(3.141592),
			).
			GroupBy("ab").
			Having(
				Expr("j = k"),
				ParenthesisOpen(),
				Column("m").Int(33),
				Column("n").Str("wh3r3").Or(),
				ParenthesisClose(),
				Column("q").NotNull(),
			)
		compareToSQL(t, sel, nil,
			"SELECT `a`, `b` FROM `c` AS `cc` WHERE (`f` = ?) AND ((`d` = ?) OR (`e` = ?)) AND (`p` = ?) GROUP BY `ab` HAVING (j = k) AND ((`m` = ?) OR (`n` = ?)) AND (`q` IS NOT NULL)",
			"SELECT `a`, `b` FROM `c` AS `cc` WHERE (`f` = 2.7182) AND ((`d` = 1) OR (`e` = 'wat')) AND (`p` = 3.141592) GROUP BY `ab` HAVING (j = k) AND ((`m` = 33) OR (`n` = 'wh3r3')) AND (`q` IS NOT NULL)",
			2.7182, int64(1), "wat", 3.141592, int64(33), "wh3r3")
	})
}

func TestSelect_Count(t *testing.T) {
	t.Parallel()
	t.Run("written count star gets quoted", func(t *testing.T) {
		compareToSQL(t,
			NewSelect("count(*)").From("dbr_people"),
			nil,
			"SELECT `count(*)` FROM `dbr_people`",
			"SELECT `count(*)` FROM `dbr_people`",
		)
	})
	t.Run("written count star gets not quoted Unsafe", func(t *testing.T) {
		compareToSQL(t,
			NewSelect().Unsafe().AddColumns("count(*)").From("dbr_people"),
			nil,
			"SELECT count(*) FROM `dbr_people`",
			"SELECT count(*) FROM `dbr_people`",
		)
	})
	t.Run("func count star", func(t *testing.T) {
		s := NewSelect("a", "b").Count().From("dbr_people")
		compareToSQL(t,
			s,
			nil,
			"SELECT COUNT(*) AS `counted` FROM `dbr_people`",
			"SELECT COUNT(*) AS `counted` FROM `dbr_people`",
		)
		var buf bytes.Buffer
		assert.NoError(t, s.Columns.WriteQuoted(&buf))
		assert.Exactly(t, "`a`, `b`", buf.String(), "Columns should be removed or changed when calling Count() function")
	})
}

func TestSelect_UseBuildCache(t *testing.T) {
	t.Parallel()

	sel := NewSelect("a", "b").
		Distinct().
		FromAlias("c", "cc").
		Where(
			ParenthesisOpen(),
			Column("d").Int(1),
			Column("e").Str("wat").Or(),
			ParenthesisClose(),
			Column("f").Int(2),
			Column("g").Int(3),
			Column("h").In().Int64s(4, 5, 6),
		).
		GroupBy("ab").
		Having(
			ParenthesisOpen(),
			Column("m").Int(33),
			Column("n").Str("wh3r3").Or(),
			ParenthesisClose(),
			Expr("j = k"),
		).
		OrderBy("l").
		Limit(7).
		Offset(8).BuildCache()

	const cachedSQLPlaceHolder = "SELECT DISTINCT `a`, `b` FROM `c` AS `cc` WHERE ((`d` = ?) OR (`e` = ?)) AND (`f` = ?) AND (`g` = ?) AND (`h` IN (?,?,?)) GROUP BY `ab` HAVING ((`m` = ?) OR (`n` = ?)) AND (j = k) ORDER BY `l` LIMIT 7 OFFSET 8"
	t.Run("without interpolate", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			compareToSQL(t, sel, nil,
				cachedSQLPlaceHolder,
				"",
				int64(1), "wat", int64(2), int64(3), int64(4), int64(5), int64(6), int64(33), "wh3r3",
			)
			assert.Equal(t, cachedSQLPlaceHolder, string(sel.cacheSQL))
		}
	})

	t.Run("with interpolate", func(t *testing.T) {
		sel.cacheSQL = nil
		const cachedSQLInterpolated = "SELECT DISTINCT `a`, `b` FROM `c` AS `cc` WHERE ((`d` = 1) OR (`e` = 'wat')) AND (`f` = 2) AND (`g` = 3) AND (`h` IN (4,5,6)) GROUP BY `ab` HAVING ((`m` = 33) OR (`n` = 'wh3r3')) AND (j = k) ORDER BY `l` LIMIT 7 OFFSET 8"
		for i := 0; i < 3; i++ {
			compareToSQL(t, sel, nil,
				cachedSQLPlaceHolder,
				cachedSQLInterpolated,
				int64(1), "wat", int64(2), int64(3), int64(4), int64(5), int64(6), int64(33), "wh3r3",
			)
			assert.Equal(t, cachedSQLPlaceHolder, string(sel.cacheSQL))
		}
	})
}

func TestSelect_SetRecord(t *testing.T) {
	t.Parallel()
	p := &dbrPerson{
		ID:    6666,
		Name:  "Hans Wurst",
		Email: MakeNullString("hans@wurst.com"),
	}

	t.Run("multiple args from record", func(t *testing.T) {
		sel := NewSelect("a", "b").
			FromAlias("dbr_person", "dp").
			Join(MakeIdentifier("dbr_group").Alias("dg"), Column("dp.id").PlaceHolder()).
			Where(
				ParenthesisOpen(),
				Column("name").PlaceHolder(),
				Column("e").Str("wat").Or(),
				ParenthesisClose(),
				Column("f").LessOrEqual().Int(2),
				Column("g").Greater().Int(3),
				Column("h").In().Int64s(4, 5, 6),
			).
			GroupBy("ab").
			Having(
				Column("email").PlaceHolder(),
				Column("n").Str("wh3r3"),
			).
			OrderBy("l").
			BindByQualifier("dp", p)

		compareToSQL(t, sel, nil,
			"SELECT `a`, `b` FROM `dbr_person` AS `dp` INNER JOIN `dbr_group` AS `dg` ON (`dp`.`id` = ?) WHERE ((`name` = ?) OR (`e` = ?)) AND (`f` <= ?) AND (`g` > ?) AND (`h` IN (?,?,?)) GROUP BY `ab` HAVING (`email` = ?) AND (`n` = ?) ORDER BY `l`",
			"",
			int64(6666), "Hans Wurst", "wat", int64(2), int64(3), int64(4), int64(5), int64(6), "hans@wurst.com", "wh3r3",
		)
	})
	t.Run("single arg JOIN", func(t *testing.T) {
		sel := NewSelect("a").FromAlias("dbr_people", "dp").
			Join(MakeIdentifier("dbr_group").Alias("dg"), Column("dp.id").PlaceHolder(), Column("dg.name").Strs("XY%")).
			Bind(p).OrderBy("id")

		compareToSQL(t, sel, nil,
			"SELECT `a` FROM `dbr_people` AS `dp` INNER JOIN `dbr_group` AS `dg` ON (`dp`.`id` = ?) AND (`dg`.`name` = ?) ORDER BY `id`",
			"SELECT `a` FROM `dbr_people` AS `dp` INNER JOIN `dbr_group` AS `dg` ON (`dp`.`id` = 6666) AND (`dg`.`name` = 'XY%') ORDER BY `id`",
			int64(6666), "XY%",
		)
	})
	t.Run("single arg WHERE", func(t *testing.T) {
		sel := NewSelect("a").From("dbr_people").
			Where(
				Column("id").PlaceHolder(),
			).
			Bind(p).OrderBy("id")

		compareToSQL(t, sel, nil,
			"SELECT `a` FROM `dbr_people` WHERE (`id` = ?) ORDER BY `id`",
			"SELECT `a` FROM `dbr_people` WHERE (`id` = 6666) ORDER BY `id`",
			int64(6666),
		)
	})
	t.Run("HAVING", func(t *testing.T) {
		sel := NewSelect("a").From("dbr_people").
			Having(
				Column("id").PlaceHolder(),
				Column("name").Like().PlaceHolder(),
			).
			Bind(p).OrderBy("id")

		compareToSQL(t, sel, nil,
			"SELECT `a` FROM `dbr_people` HAVING (`id` = ?) AND (`name` LIKE ?) ORDER BY `id`",
			"SELECT `a` FROM `dbr_people` HAVING (`id` = 6666) AND (`name` LIKE 'Hans Wurst') ORDER BY `id`",
			int64(6666), "Hans Wurst",
		)
	})

	t.Run("slice as record - nice feature", func(t *testing.T) {
		persons := &dbrPersons{
			Data: []*dbrPerson{
				{Name: "Muffin Hat", Email: MakeNullString("Muffin@Hat.head")},
				{Name: "Marianne Phyllis Finch", Email: MakeNullString("marianne@phyllis.finch")},
				{Name: "Daphne Augusta Perry", Email: MakeNullString("daphne@augusta.perry")},
			},
		}

		compareToSQL(t,
			NewSelect("name", "email").From("dbr_person").
				Where(
					Column("name").In().PlaceHolder(),
					Column("email").In().PlaceHolder(),
				).
				Bind(persons),
			nil,
			"SELECT `name`, `email` FROM `dbr_person` WHERE (`name` IN (?)) AND (`email` IN (?))",
			"SELECT `name`, `email` FROM `dbr_person` WHERE (`name` IN ('Muffin Hat','Marianne Phyllis Finch','Daphne Augusta Perry')) AND (`email` IN ('Muffin@Hat.head','marianne@phyllis.finch','daphne@augusta.perry'))",
			"Muffin Hat", "Marianne Phyllis Finch", "Daphne Augusta Perry",
			"Muffin@Hat.head", "marianne@phyllis.finch", "daphne@augusta.perry",
		)
	})

}
