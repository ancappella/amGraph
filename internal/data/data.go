package data

import (
	"database/sql"
	"fmt"

	"amGraph/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewGreeterRepo)

// Data .
type Data struct {
	DB *sql.DB
}

// NewData .
func NewData(c *conf.Data) (*Data, func(), error) {
	cleanup := func() {
		log.Info("closing the data resources")
	}
	if c == nil || c.Database == nil {
		return nil, cleanup, fmt.Errorf("data.database is not configured")
	}
	db, err := sql.Open(c.Database.Driver, c.Database.Source)
	if err != nil {
		return nil, cleanup, fmt.Errorf("open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, cleanup, fmt.Errorf("ping database: %w", err)
	}
	log.Infof("connected to database (%s)", c.Database.Driver)
	d := &Data{DB: db}
	cleanup = func() {
		log.Info("closing the data resources")
		_ = db.Close()
	}
	return d, cleanup, nil
}
