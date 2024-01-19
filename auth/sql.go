package auth

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"strings"

	"github.com/autoapev1/indexer/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func NewSqlDB(uri string) *bun.DB {
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(uri)))
	return bun.NewDB(pgdb, pgdialect.New())
}

type SqlAuthProvider struct {
	db      *bun.DB
	keyType KeyType
}

type sqlKey struct {
	bun.BaseModel `bun:"keys,alias:ku"` // Table name: keys
	Key           string                `bun:"key,pk"`
	Iat           int64                 `bun:"iat"`
	Exp           int64                 `bun:"exp"`
	LastIP        string                `bun:"last_ip"`
	LastAccess    int64                 `bun:"last_access"`
	CallCount     int64                 `bun:"call_count"`
	MethodUsages  []*sqlMethodUsage     `bun:"rel:has-many,join:key=key"`
}

type sqlMethodUsage struct {
	bun.BaseModel `bun:"method_usages,alias:mu"` // Table name: method_usages
	MethodUsageID int64                          `bun:",pk,autoincrement"`
	Key           string                         `bun:"key"`
	MethodName    string                         `bun:"method_name"`
	UsageCount    int64                          `bun:"usage_count"`
}

func NewSqlAuthProvider(db *bun.DB) *SqlAuthProvider {
	provider := &SqlAuthProvider{
		db:      db,
		keyType: KeyTypeHex64,
	}

	provider.initTables()
	return provider
}

func (a *SqlAuthProvider) WithKeyType(keyType KeyType) *SqlAuthProvider {
	a.keyType = keyType
	return a
}

// create table and index
func (a *SqlAuthProvider) initTables() {
	_, err := a.db.NewCreateTable().
		Model((*sqlKey)(nil)).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		slog.Error("error creating table keys", "err", err)
	}

	_, err = a.db.NewCreateIndex().
		Model((*sqlKey)(nil)).
		Index("key_inx").
		Column("key").
		Unique().
		Exec(context.Background())
	if err != nil {
		slog.Error("error creating index key_inx", "err", err)
	}

}

func (a *SqlAuthProvider) Authenticate(r *http.Request) (AuthLevel, error) {
	var sqlKey sqlKey

	// get key from request
	key := r.Header.Get("Authentication")
	key = strings.TrimPrefix(key, "Bearer ")

	// check master
	master := config.Get().API.AuthMasterKey
	if master != "" && key == master {
		return AuthLevelMaster, nil
	}

	// search db for key
	err := a.db.NewSelect().
		Model(&sqlKey).
		Where("key = ?", key).
		Scan(context.Background())

	if err != nil {
		return AuthLevelUnauthorized, ErrCheckingAuth
	}

	return AuthLevelBasic, nil
}

func (a *SqlAuthProvider) Register() (string, error) {
	key, err := GenerateKey(a.keyType)
	if err != nil {
		return "", err
	}

	sqlKey := &sqlKey{
		Key:          key,
		MethodUsages: make([]*sqlMethodUsage, 0),
	}

	_, err = a.db.NewInsert().
		Model(sqlKey).
		Exec(context.Background())
	if err != nil {
		return "", err
	}

	return key, nil
}
func (a *SqlAuthProvider) UpdateUsage(key string, usageDelta KeyUsage) error {
	// Start a transaction
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}

	// Rollback in case of error
	defer tx.Rollback()

	// Update the key usage
	_, err = tx.NewUpdate().
		Model((*sqlKey)(nil)).
		Where("key = ?", key).
		Set("call_count = call_count + ?", usageDelta.Requests).
		Set("last_access = ?", usageDelta.AccessedAt).
		Set("last_ip = ?", usageDelta.IP).
		Exec(context.Background())
	if err != nil {
		return err
	}

	// Update method usages
	for method, count := range usageDelta.MethodUsage {
		var methodUsage sqlMethodUsage
		_, err := tx.NewSelect().
			Model(&methodUsage).
			Where("key = ? AND method_name = ?", key, method).
			Exec(context.Background())

		if err != nil && err != sql.ErrNoRows {
			// Handle errors other than 'no rows found'
			return err
		}

		if methodUsage.MethodUsageID == 0 {
			// Insert new method usage
			_, err = tx.NewInsert().
				Model(&sqlMethodUsage{
					Key:        key,
					MethodName: method,
					UsageCount: count,
				}).
				Exec(context.Background())
		} else {
			// Update existing method usage
			_, err = tx.NewUpdate().
				Model(&sqlMethodUsage{}).
				Set("usage_count = usage_count + ?", count).
				Where("key = ? AND method_name = ?", key, method).
				Exec(context.Background())
		}

		if err != nil {
			return err
		}
	}

	// Commit the transaction
	return tx.Commit()
}

var _ Provider = (*SqlAuthProvider)(nil)
