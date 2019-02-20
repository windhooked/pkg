// Auto generated via github.com/corestoreio/pkg/sql/dmlgen

package testdata

import (
	"context"
	"fmt"
	"github.com/corestoreio/pkg/sql/ddl"
	"github.com/corestoreio/pkg/sql/dml"
	"github.com/corestoreio/pkg/sql/dmltest"
	"github.com/corestoreio/pkg/util/assert"
	"github.com/corestoreio/pkg/util/pseudo"
	"sort"
	"testing"
	"time"
)

func TestNewTables(t *testing.T) {
	db := dmltest.MustConnectDB(t)
	defer dmltest.Close(t, db)

	defer dmltest.SQLDumpLoad(t, "test_*_tables.sql", &dmltest.SQLDumpOptions{
		SkipDBCleanup: true,
	}).Deferred()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()
	tbls, err := NewTables(ctx, ddl.WithConnPool(db))
	assert.NoError(t, err)

	tblNames := tbls.Tables()
	sort.Strings(tblNames)
	assert.Exactly(t, []string{"catalog_product_index_eav_decimal_idx", "core_config_data", "customer_address_entity", "customer_entity", "dmlgen_types", "sales_order_status_state", "view_customer_auto_increment", "view_customer_no_auto_increment"}, tblNames)

	err = tbls.Validate(ctx)
	assert.NoError(t, err)
	var ps *pseudo.Service
	ps = pseudo.MustNewService(0, &pseudo.Options{Lang: "de", FloatMaxDecimals: 6},
		pseudo.WithTagFakeFunc("website_id", func(maxLen int) (interface{}, error) {
			return 1, nil
		}),
		pseudo.WithTagFakeFunc("store_id", func(maxLen int) (interface{}, error) {
			return 1, nil
		}),
		pseudo.WithTagFakeFunc("testdata.CustomerAddressEntity.ParentID", func(maxLen int) (interface{}, error) {
			return nil, nil
		}),
		pseudo.WithTagFakeFunc("col_date1", func(maxLen int) (interface{}, error) {
			if ps.Intn(1000)%3 == 0 {
				return nil, nil
			}
			return ps.Dob18(), nil
		}),
		pseudo.WithTagFakeFunc("col_date2", func(maxLen int) (interface{}, error) {
			return ps.Dob18().MarshalText()
		}),
		pseudo.WithTagFakeFunc("col_decimal101", func(maxLen int) (interface{}, error) {
			return fmt.Sprintf("%.1f", ps.Price()), nil
		}),
		pseudo.WithTagFakeFunc("price124b", func(maxLen int) (interface{}, error) {
			return fmt.Sprintf("%.4f", ps.Price()), nil
		}),
		pseudo.WithTagFakeFunc("col_decimal123", func(maxLen int) (interface{}, error) {
			return fmt.Sprintf("%.3f", ps.Float64()), nil
		}),
		pseudo.WithTagFakeFunc("col_decimal206", func(maxLen int) (interface{}, error) {
			return fmt.Sprintf("%.6f", ps.Float64()), nil
		}),
		pseudo.WithTagFakeFunc("col_decimal2412", func(maxLen int) (interface{}, error) {
			return fmt.Sprintf("%.12f", ps.Float64()), nil
		}),
		pseudo.WithTagFakeFuncAlias(
			"col_decimal124", "price124b",
			"price124a", "price124b",
			"col_float", "col_decimal206",
		),
	)
	t.Run("CatalogProductIndexEAVDecimalIDX_Entity", func(t *testing.T) {
		tbl := tbls.MustTable(TableNameCatalogProductIndexEAVDecimalIDX)

		entSELECT := tbl.SelectByPK("*")
		entSELECTStmtA := entSELECT.WithArgs().ExpandPlaceHolders() // WithArgs generates the cached SQL string with key ""

		entSELECT.WithCacheKey("select_10").Wheres.Reset()
		_, _, err := entSELECT.Where().ToSQL() // ToSQL generates the new cached SQL string with key select_10
		assert.NoError(t, err)
		entCol := NewCatalogProductIndexEAVDecimalIDXCollection()

		// this table/view does not support auto_increment
		rowCount, err := entSELECTStmtA.WithCacheKey("select_10").Load(ctx, entCol)
		assert.NoError(t, err)
		t.Logf("SELECT queries: %#v", entSELECT.CachedQueries())
		t.Logf("Collection load rowCount: %d", rowCount)
	})
	t.Run("CoreConfigData_Entity", func(t *testing.T) {
		tbl := tbls.MustTable(TableNameCoreConfigData)

		entSELECT := tbl.SelectByPK("*")
		entSELECTStmtA := entSELECT.WithArgs().ExpandPlaceHolders() // WithArgs generates the cached SQL string with key ""

		entSELECT.WithCacheKey("select_10").Wheres.Reset()
		_, _, err := entSELECT.Where(
			dml.Column("config_id").LessOrEqual().Int(10),
		).ToSQL() // ToSQL generates the new cached SQL string with key select_10
		assert.NoError(t, err)
		entCol := NewCoreConfigDataCollection()

		entINSERT := tbl.Insert().BuildValues()
		entINSERTStmtA := entINSERT.PrepareWithArgs(ctx)

		for i := 0; i < 9; i++ {
			entIn := new(CoreConfigData)
			if err := ps.FakeData(entIn); err != nil {
				t.Errorf("IDX[%d]: %+v", i, err)
				return
			}

			lID := dmltest.CheckLastInsertID(t, "Error: TestNewTables.CoreConfigData_Entity")(entINSERTStmtA.Record("", entIn).ExecContext(ctx))
			entINSERTStmtA.Reset()

			entOut := new(CoreConfigData)
			rowCount, err := entSELECTStmtA.Int64s(lID).Load(ctx, entOut)
			assert.NoError(t, err)
			assert.Exactly(t, uint64(1), rowCount, "IDX%d: RowCount did not match", i)
			assert.Exactly(t, entIn.ConfigID, entOut.ConfigID, "IDX%d: ConfigID should match", lID)
			assert.ExactlyLength(t, 8, &entIn.Scope, &entOut.Scope, "IDX%d: Scope should match", lID)
			assert.Exactly(t, entIn.ScopeID, entOut.ScopeID, "IDX%d: ScopeID should match", lID)
			assert.Exactly(t, entIn.Expires, entOut.Expires, "IDX%d: Expires should match", lID)
			assert.ExactlyLength(t, 255, &entIn.Path, &entOut.Path, "IDX%d: Path should match", lID)
			assert.ExactlyLength(t, 65535, &entIn.Value, &entOut.Value, "IDX%d: Value should match", lID)
			// ignoring: version_ts
			// ignoring: version_te
		}
		dmltest.Close(t, entINSERTStmtA)

		rowCount, err := entSELECTStmtA.WithCacheKey("select_10").Load(ctx, entCol)
		assert.NoError(t, err)
		t.Logf("Collection load rowCount: %d", rowCount)

		entINSERTStmtA = entINSERT.WithCacheKey("row_count_%d", len(entCol.Data)).Replace().SetRowCount(len(entCol.Data)).PrepareWithArgs(ctx)
		lID := dmltest.CheckLastInsertID(t, "Error: CoreConfigDataCollection")(entINSERTStmtA.Record("", entCol).ExecContext(ctx))
		dmltest.Close(t, entINSERTStmtA)
		t.Logf("Last insert ID into: %d", lID)
		t.Logf("INSERT queries: %#v", entINSERT.CachedQueries())
		t.Logf("SELECT queries: %#v", entSELECT.CachedQueries())
	})
	t.Run("CustomerAddressEntity_Entity", func(t *testing.T) {
		tbl := tbls.MustTable(TableNameCustomerAddressEntity)

		entSELECT := tbl.SelectByPK("*")
		entSELECTStmtA := entSELECT.WithArgs().ExpandPlaceHolders() // WithArgs generates the cached SQL string with key ""

		entSELECT.WithCacheKey("select_10").Wheres.Reset()
		_, _, err := entSELECT.Where(
			dml.Column("entity_id").LessOrEqual().Int(10),
		).ToSQL() // ToSQL generates the new cached SQL string with key select_10
		assert.NoError(t, err)
		entCol := NewCustomerAddressEntityCollection()

		entINSERT := tbl.Insert().BuildValues()
		entINSERTStmtA := entINSERT.PrepareWithArgs(ctx)

		for i := 0; i < 9; i++ {
			entIn := new(CustomerAddressEntity)
			if err := ps.FakeData(entIn); err != nil {
				t.Errorf("IDX[%d]: %+v", i, err)
				return
			}

			lID := dmltest.CheckLastInsertID(t, "Error: TestNewTables.CustomerAddressEntity_Entity")(entINSERTStmtA.Record("", entIn).ExecContext(ctx))
			entINSERTStmtA.Reset()

			entOut := new(CustomerAddressEntity)
			rowCount, err := entSELECTStmtA.Int64s(lID).Load(ctx, entOut)
			assert.NoError(t, err)
			assert.Exactly(t, uint64(1), rowCount, "IDX%d: RowCount did not match", i)
			assert.Exactly(t, entIn.EntityID, entOut.EntityID, "IDX%d: EntityID should match", lID)
			assert.ExactlyLength(t, 50, &entIn.IncrementID, &entOut.IncrementID, "IDX%d: IncrementID should match", lID)
			assert.Exactly(t, entIn.ParentID, entOut.ParentID, "IDX%d: ParentID should match", lID)
			assert.Exactly(t, entIn.CreatedAt, entOut.CreatedAt, "IDX%d: CreatedAt should match", lID)
			assert.Exactly(t, entIn.UpdatedAt, entOut.UpdatedAt, "IDX%d: UpdatedAt should match", lID)
			assert.Exactly(t, entIn.IsActive, entOut.IsActive, "IDX%d: IsActive should match", lID)
			assert.ExactlyLength(t, 255, &entIn.City, &entOut.City, "IDX%d: City should match", lID)
			assert.ExactlyLength(t, 255, &entIn.Company, &entOut.Company, "IDX%d: Company should match", lID)
			assert.ExactlyLength(t, 255, &entIn.CountryID, &entOut.CountryID, "IDX%d: CountryID should match", lID)
			assert.ExactlyLength(t, 255, &entIn.Fax, &entOut.Fax, "IDX%d: Fax should match", lID)
			assert.ExactlyLength(t, 255, &entIn.Firstname, &entOut.Firstname, "IDX%d: Firstname should match", lID)
			assert.ExactlyLength(t, 255, &entIn.Lastname, &entOut.Lastname, "IDX%d: Lastname should match", lID)
			assert.ExactlyLength(t, 255, &entIn.Middlename, &entOut.Middlename, "IDX%d: Middlename should match", lID)
			assert.ExactlyLength(t, 255, &entIn.Postcode, &entOut.Postcode, "IDX%d: Postcode should match", lID)
			assert.ExactlyLength(t, 40, &entIn.Prefix, &entOut.Prefix, "IDX%d: Prefix should match", lID)
			assert.ExactlyLength(t, 255, &entIn.Region, &entOut.Region, "IDX%d: Region should match", lID)
			assert.Exactly(t, entIn.RegionID, entOut.RegionID, "IDX%d: RegionID should match", lID)
			assert.ExactlyLength(t, 65535, &entIn.Street, &entOut.Street, "IDX%d: Street should match", lID)
			assert.ExactlyLength(t, 40, &entIn.Suffix, &entOut.Suffix, "IDX%d: Suffix should match", lID)
			assert.ExactlyLength(t, 255, &entIn.Telephone, &entOut.Telephone, "IDX%d: Telephone should match", lID)
			assert.ExactlyLength(t, 255, &entIn.VatID, &entOut.VatID, "IDX%d: VatID should match", lID)
			assert.Exactly(t, entIn.VatIsValid, entOut.VatIsValid, "IDX%d: VatIsValid should match", lID)
			assert.ExactlyLength(t, 255, &entIn.VatRequestDate, &entOut.VatRequestDate, "IDX%d: VatRequestDate should match", lID)
			assert.ExactlyLength(t, 255, &entIn.VatRequestID, &entOut.VatRequestID, "IDX%d: VatRequestID should match", lID)
			assert.Exactly(t, entIn.VatRequestSuccess, entOut.VatRequestSuccess, "IDX%d: VatRequestSuccess should match", lID)
		}
		dmltest.Close(t, entINSERTStmtA)

		rowCount, err := entSELECTStmtA.WithCacheKey("select_10").Load(ctx, entCol)
		assert.NoError(t, err)
		t.Logf("Collection load rowCount: %d", rowCount)

		entINSERTStmtA = entINSERT.WithCacheKey("row_count_%d", len(entCol.Data)).Replace().SetRowCount(len(entCol.Data)).PrepareWithArgs(ctx)
		lID := dmltest.CheckLastInsertID(t, "Error: CustomerAddressEntityCollection")(entINSERTStmtA.Record("", entCol).ExecContext(ctx))
		dmltest.Close(t, entINSERTStmtA)
		t.Logf("Last insert ID into: %d", lID)
		t.Logf("INSERT queries: %#v", entINSERT.CachedQueries())
		t.Logf("SELECT queries: %#v", entSELECT.CachedQueries())
	})
	t.Run("CustomerEntity_Entity", func(t *testing.T) {
		tbl := tbls.MustTable(TableNameCustomerEntity)

		entSELECT := tbl.SelectByPK("*")
		entSELECTStmtA := entSELECT.WithArgs().ExpandPlaceHolders() // WithArgs generates the cached SQL string with key ""

		entSELECT.WithCacheKey("select_10").Wheres.Reset()
		_, _, err := entSELECT.Where(
			dml.Column("entity_id").LessOrEqual().Int(10),
		).ToSQL() // ToSQL generates the new cached SQL string with key select_10
		assert.NoError(t, err)
		entCol := NewCustomerEntityCollection()

		entINSERT := tbl.Insert().BuildValues()
		entINSERTStmtA := entINSERT.PrepareWithArgs(ctx)

		for i := 0; i < 9; i++ {
			entIn := new(CustomerEntity)
			if err := ps.FakeData(entIn); err != nil {
				t.Errorf("IDX[%d]: %+v", i, err)
				return
			}

			lID := dmltest.CheckLastInsertID(t, "Error: TestNewTables.CustomerEntity_Entity")(entINSERTStmtA.Record("", entIn).ExecContext(ctx))
			entINSERTStmtA.Reset()

			entOut := new(CustomerEntity)
			rowCount, err := entSELECTStmtA.Int64s(lID).Load(ctx, entOut)
			assert.NoError(t, err)
			assert.Exactly(t, uint64(1), rowCount, "IDX%d: RowCount did not match", i)
			assert.Exactly(t, entIn.EntityID, entOut.EntityID, "IDX%d: EntityID should match", lID)
			assert.Exactly(t, entIn.WebsiteID, entOut.WebsiteID, "IDX%d: WebsiteID should match", lID)
			assert.ExactlyLength(t, 255, &entIn.Email, &entOut.Email, "IDX%d: Email should match", lID)
			assert.Exactly(t, entIn.GroupID, entOut.GroupID, "IDX%d: GroupID should match", lID)
			assert.ExactlyLength(t, 50, &entIn.IncrementID, &entOut.IncrementID, "IDX%d: IncrementID should match", lID)
			assert.Exactly(t, entIn.StoreID, entOut.StoreID, "IDX%d: StoreID should match", lID)
			assert.Exactly(t, entIn.CreatedAt, entOut.CreatedAt, "IDX%d: CreatedAt should match", lID)
			assert.Exactly(t, entIn.UpdatedAt, entOut.UpdatedAt, "IDX%d: UpdatedAt should match", lID)
			assert.Exactly(t, entIn.IsActive, entOut.IsActive, "IDX%d: IsActive should match", lID)
			assert.Exactly(t, entIn.DisableAutoGroupChange, entOut.DisableAutoGroupChange, "IDX%d: DisableAutoGroupChange should match", lID)
			assert.ExactlyLength(t, 255, &entIn.CreatedIn, &entOut.CreatedIn, "IDX%d: CreatedIn should match", lID)
			assert.ExactlyLength(t, 40, &entIn.Prefix, &entOut.Prefix, "IDX%d: Prefix should match", lID)
			assert.ExactlyLength(t, 255, &entIn.Firstname, &entOut.Firstname, "IDX%d: Firstname should match", lID)
			assert.ExactlyLength(t, 255, &entIn.Middlename, &entOut.Middlename, "IDX%d: Middlename should match", lID)
			assert.ExactlyLength(t, 255, &entIn.Lastname, &entOut.Lastname, "IDX%d: Lastname should match", lID)
			assert.ExactlyLength(t, 40, &entIn.Suffix, &entOut.Suffix, "IDX%d: Suffix should match", lID)
			assert.Exactly(t, entIn.Dob, entOut.Dob, "IDX%d: Dob should match", lID)
			assert.ExactlyLength(t, 128, &entIn.passwordHash, &entOut.passwordHash, "IDX%d: passwordHash should match", lID)
			assert.ExactlyLength(t, 128, &entIn.RpToken, &entOut.RpToken, "IDX%d: RpToken should match", lID)
			assert.Exactly(t, entIn.RpTokenCreatedAt, entOut.RpTokenCreatedAt, "IDX%d: RpTokenCreatedAt should match", lID)
			assert.Exactly(t, entIn.DefaultBilling, entOut.DefaultBilling, "IDX%d: DefaultBilling should match", lID)
			assert.Exactly(t, entIn.DefaultShipping, entOut.DefaultShipping, "IDX%d: DefaultShipping should match", lID)
			assert.ExactlyLength(t, 50, &entIn.Taxvat, &entOut.Taxvat, "IDX%d: Taxvat should match", lID)
			assert.ExactlyLength(t, 64, &entIn.Confirmation, &entOut.Confirmation, "IDX%d: Confirmation should match", lID)
			assert.Exactly(t, entIn.Gender, entOut.Gender, "IDX%d: Gender should match", lID)
			assert.Exactly(t, entIn.FailuresNum, entOut.FailuresNum, "IDX%d: FailuresNum should match", lID)
			assert.Exactly(t, entIn.FirstFailure, entOut.FirstFailure, "IDX%d: FirstFailure should match", lID)
			assert.Exactly(t, entIn.LockExpires, entOut.LockExpires, "IDX%d: LockExpires should match", lID)
		}
		dmltest.Close(t, entINSERTStmtA)

		rowCount, err := entSELECTStmtA.WithCacheKey("select_10").Load(ctx, entCol)
		assert.NoError(t, err)
		t.Logf("Collection load rowCount: %d", rowCount)

		entINSERTStmtA = entINSERT.WithCacheKey("row_count_%d", len(entCol.Data)).Replace().SetRowCount(len(entCol.Data)).PrepareWithArgs(ctx)
		lID := dmltest.CheckLastInsertID(t, "Error: CustomerEntityCollection")(entINSERTStmtA.Record("", entCol).ExecContext(ctx))
		dmltest.Close(t, entINSERTStmtA)
		t.Logf("Last insert ID into: %d", lID)
		t.Logf("INSERT queries: %#v", entINSERT.CachedQueries())
		t.Logf("SELECT queries: %#v", entSELECT.CachedQueries())
	})
	t.Run("DmlgenTypes_Entity", func(t *testing.T) {
		tbl := tbls.MustTable(TableNameDmlgenTypes)

		entSELECT := tbl.SelectByPK("*")
		entSELECTStmtA := entSELECT.WithArgs().ExpandPlaceHolders() // WithArgs generates the cached SQL string with key ""

		entSELECT.WithCacheKey("select_10").Wheres.Reset()
		_, _, err := entSELECT.Where(
			dml.Column("id").LessOrEqual().Int(10),
		).ToSQL() // ToSQL generates the new cached SQL string with key select_10
		assert.NoError(t, err)
		entCol := NewDmlgenTypesCollection()

		entINSERT := tbl.Insert().BuildValues()
		entINSERTStmtA := entINSERT.PrepareWithArgs(ctx)

		for i := 0; i < 9; i++ {
			entIn := new(DmlgenTypes)
			if err := ps.FakeData(entIn); err != nil {
				t.Errorf("IDX[%d]: %+v", i, err)
				return
			}

			lID := dmltest.CheckLastInsertID(t, "Error: TestNewTables.DmlgenTypes_Entity")(entINSERTStmtA.Record("", entIn).ExecContext(ctx))
			entINSERTStmtA.Reset()

			entOut := new(DmlgenTypes)
			rowCount, err := entSELECTStmtA.Int64s(lID).Load(ctx, entOut)
			assert.NoError(t, err)
			assert.Exactly(t, uint64(1), rowCount, "IDX%d: RowCount did not match", i)
			assert.Exactly(t, entIn.ID, entOut.ID, "IDX%d: ID should match", lID)
			assert.Exactly(t, entIn.ColBigint1, entOut.ColBigint1, "IDX%d: ColBigint1 should match", lID)
			assert.Exactly(t, entIn.ColBigint2, entOut.ColBigint2, "IDX%d: ColBigint2 should match", lID)
			assert.Exactly(t, entIn.ColBigint3, entOut.ColBigint3, "IDX%d: ColBigint3 should match", lID)
			assert.Exactly(t, entIn.ColBigint4, entOut.ColBigint4, "IDX%d: ColBigint4 should match", lID)
			assert.ExactlyLength(t, 65535, &entIn.ColBlob, &entOut.ColBlob, "IDX%d: ColBlob should match", lID)
			assert.Exactly(t, entIn.ColDate1, entOut.ColDate1, "IDX%d: ColDate1 should match", lID)
			assert.Exactly(t, entIn.ColDate2, entOut.ColDate2, "IDX%d: ColDate2 should match", lID)
			assert.Exactly(t, entIn.ColDatetime1, entOut.ColDatetime1, "IDX%d: ColDatetime1 should match", lID)
			assert.Exactly(t, entIn.ColDatetime2, entOut.ColDatetime2, "IDX%d: ColDatetime2 should match", lID)
			assert.Exactly(t, entIn.ColDecimal101, entOut.ColDecimal101, "IDX%d: ColDecimal101 should match", lID)
			assert.Exactly(t, entIn.ColDecimal124, entOut.ColDecimal124, "IDX%d: ColDecimal124 should match", lID)
			assert.Exactly(t, entIn.Price124a, entOut.Price124a, "IDX%d: Price124a should match", lID)
			assert.Exactly(t, entIn.Price124b, entOut.Price124b, "IDX%d: Price124b should match", lID)
			assert.Exactly(t, entIn.ColDecimal123, entOut.ColDecimal123, "IDX%d: ColDecimal123 should match", lID)
			assert.Exactly(t, entIn.ColDecimal206, entOut.ColDecimal206, "IDX%d: ColDecimal206 should match", lID)
			assert.Exactly(t, entIn.ColDecimal2412, entOut.ColDecimal2412, "IDX%d: ColDecimal2412 should match", lID)
			assert.Exactly(t, entIn.ColInt1, entOut.ColInt1, "IDX%d: ColInt1 should match", lID)
			assert.Exactly(t, entIn.ColInt2, entOut.ColInt2, "IDX%d: ColInt2 should match", lID)
			assert.Exactly(t, entIn.ColInt3, entOut.ColInt3, "IDX%d: ColInt3 should match", lID)
			assert.Exactly(t, entIn.ColInt4, entOut.ColInt4, "IDX%d: ColInt4 should match", lID)
			assert.ExactlyLength(t, 4294967295, &entIn.ColLongtext1, &entOut.ColLongtext1, "IDX%d: ColLongtext1 should match", lID)
			assert.ExactlyLength(t, 4294967295, &entIn.ColLongtext2, &entOut.ColLongtext2, "IDX%d: ColLongtext2 should match", lID)
			assert.ExactlyLength(t, 16777215, &entIn.ColMediumblob, &entOut.ColMediumblob, "IDX%d: ColMediumblob should match", lID)
			assert.ExactlyLength(t, 16777215, &entIn.ColMediumtext1, &entOut.ColMediumtext1, "IDX%d: ColMediumtext1 should match", lID)
			assert.ExactlyLength(t, 16777215, &entIn.ColMediumtext2, &entOut.ColMediumtext2, "IDX%d: ColMediumtext2 should match", lID)
			assert.Exactly(t, entIn.ColSmallint1, entOut.ColSmallint1, "IDX%d: ColSmallint1 should match", lID)
			assert.Exactly(t, entIn.ColSmallint2, entOut.ColSmallint2, "IDX%d: ColSmallint2 should match", lID)
			assert.Exactly(t, entIn.ColSmallint3, entOut.ColSmallint3, "IDX%d: ColSmallint3 should match", lID)
			assert.Exactly(t, entIn.ColSmallint4, entOut.ColSmallint4, "IDX%d: ColSmallint4 should match", lID)
			assert.Exactly(t, entIn.HasSmallint5, entOut.HasSmallint5, "IDX%d: HasSmallint5 should match", lID)
			assert.Exactly(t, entIn.IsSmallint5, entOut.IsSmallint5, "IDX%d: IsSmallint5 should match", lID)
			assert.ExactlyLength(t, 65535, &entIn.ColText, &entOut.ColText, "IDX%d: ColText should match", lID)
			assert.Exactly(t, entIn.ColTimestamp1, entOut.ColTimestamp1, "IDX%d: ColTimestamp1 should match", lID)
			assert.Exactly(t, entIn.ColTimestamp2, entOut.ColTimestamp2, "IDX%d: ColTimestamp2 should match", lID)
			assert.Exactly(t, entIn.ColTinyint1, entOut.ColTinyint1, "IDX%d: ColTinyint1 should match", lID)
			assert.ExactlyLength(t, 1, &entIn.ColVarchar1, &entOut.ColVarchar1, "IDX%d: ColVarchar1 should match", lID)
			assert.ExactlyLength(t, 100, &entIn.ColVarchar100, &entOut.ColVarchar100, "IDX%d: ColVarchar100 should match", lID)
			assert.ExactlyLength(t, 16, &entIn.ColVarchar16, &entOut.ColVarchar16, "IDX%d: ColVarchar16 should match", lID)
			assert.ExactlyLength(t, 21, &entIn.ColChar1, &entOut.ColChar1, "IDX%d: ColChar1 should match", lID)
			assert.ExactlyLength(t, 17, &entIn.ColChar2, &entOut.ColChar2, "IDX%d: ColChar2 should match", lID)
		}
		dmltest.Close(t, entINSERTStmtA)

		rowCount, err := entSELECTStmtA.WithCacheKey("select_10").Load(ctx, entCol)
		assert.NoError(t, err)
		t.Logf("Collection load rowCount: %d", rowCount)

		entINSERTStmtA = entINSERT.WithCacheKey("row_count_%d", len(entCol.Data)).Replace().SetRowCount(len(entCol.Data)).PrepareWithArgs(ctx)
		lID := dmltest.CheckLastInsertID(t, "Error: DmlgenTypesCollection")(entINSERTStmtA.Record("", entCol).ExecContext(ctx))
		dmltest.Close(t, entINSERTStmtA)
		t.Logf("Last insert ID into: %d", lID)
		t.Logf("INSERT queries: %#v", entINSERT.CachedQueries())
		t.Logf("SELECT queries: %#v", entSELECT.CachedQueries())
	})
	t.Run("SalesOrderStatusState_Entity", func(t *testing.T) {
		tbl := tbls.MustTable(TableNameSalesOrderStatusState)

		entSELECT := tbl.SelectByPK("*")
		entSELECTStmtA := entSELECT.WithArgs().ExpandPlaceHolders() // WithArgs generates the cached SQL string with key ""

		entSELECT.WithCacheKey("select_10").Wheres.Reset()
		_, _, err := entSELECT.Where().ToSQL() // ToSQL generates the new cached SQL string with key select_10
		assert.NoError(t, err)
		entCol := NewSalesOrderStatusStateCollection()

		// this table/view does not support auto_increment
		rowCount, err := entSELECTStmtA.WithCacheKey("select_10").Load(ctx, entCol)
		assert.NoError(t, err)
		t.Logf("SELECT queries: %#v", entSELECT.CachedQueries())
		t.Logf("Collection load rowCount: %d", rowCount)
	})
	t.Run("ViewCustomerAutoIncrement_Entity", func(t *testing.T) {
		tbl := tbls.MustTable(TableNameViewCustomerAutoIncrement)

		entSELECT := tbl.SelectByPK("*")
		entSELECTStmtA := entSELECT.WithArgs().ExpandPlaceHolders() // WithArgs generates the cached SQL string with key ""

		entSELECT.WithCacheKey("select_10").Wheres.Reset()
		_, _, err := entSELECT.Where().ToSQL() // ToSQL generates the new cached SQL string with key select_10
		assert.NoError(t, err)
		entCol := NewViewCustomerAutoIncrementCollection()

		// this table/view does not support auto_increment
		rowCount, err := entSELECTStmtA.WithCacheKey("select_10").Load(ctx, entCol)
		assert.NoError(t, err)
		t.Logf("SELECT queries: %#v", entSELECT.CachedQueries())
		t.Logf("Collection load rowCount: %d", rowCount)
	})
	t.Run("ViewCustomerNoAutoIncrement_Entity", func(t *testing.T) {
		tbl := tbls.MustTable(TableNameViewCustomerNoAutoIncrement)

		entSELECT := tbl.SelectByPK("*")
		entSELECTStmtA := entSELECT.WithArgs().ExpandPlaceHolders() // WithArgs generates the cached SQL string with key ""

		entSELECT.WithCacheKey("select_10").Wheres.Reset()
		_, _, err := entSELECT.Where().ToSQL() // ToSQL generates the new cached SQL string with key select_10
		assert.NoError(t, err)
		entCol := NewViewCustomerNoAutoIncrementCollection()

		// this table/view does not support auto_increment
		rowCount, err := entSELECTStmtA.WithCacheKey("select_10").Load(ctx, entCol)
		assert.NoError(t, err)
		t.Logf("SELECT queries: %#v", entSELECT.CachedQueries())
		t.Logf("Collection load rowCount: %d", rowCount)
	})
}
