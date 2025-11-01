package sqlstorage

import (
	"context"
	"path/filepath"

	"github.com/google/uuid"                                                                //nolint:depguard
	_ "github.com/jackc/pgx/stdlib"                                                         //nolint:depguard
	"github.com/jmoiron/sqlx"                                                               //nolint:depguard
	"github.com/pressly/goose/v3"                                                           //nolint:depguard
	"github.com/pressly/goose/v3/database"                                                  //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/config"                 //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/storage/common"         //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/storage/sql/migrations" //nolint:depguard
)

type Storage struct {
	db  *sqlx.DB
	c   config.StorageConfig
	err error
	ctx context.Context
}

func New(ctx context.Context, c config.StorageConfig) common.StorageDriverInterface {
	s := &Storage{c: c, ctx: ctx}
	s.err = s.Connect(ctx)
	return s
}

func (s *Storage) Connect(ctx context.Context) error {
	s.db, s.err = sqlx.ConnectContext(ctx, "pgx", s.c.Dsn)
	if s.err != nil {
		return s.err
	}
	s.err = s.db.PingContext(ctx)
	return s.err
}

func (s *Storage) Close() error {
	s.err = s.db.Close()
	return s.err
}

func (s *Storage) Add(event common.Event) (common.Event, error) {
	if event.ID.(string) == "" {
		event.ID = uuid.New().String()
	}
	sql := `INSERT INTO events("id","title","date_time","duration","description","user","notify_time") 
VALUES(:id, :title, :date_time, :duration, :description, :user, :notify_time)`
	_, err := s.db.NamedExecContext(s.ctx, sql, event)
	if err != nil {
		return event, err
	}
	return event, err
}

func (s *Storage) Update(event common.Event) error {
	sql := `UPDATE events SET "title" = :title,"date_time" = :date_time,"duration" = :duration,
                  "description" = :description,"user" = :user,
                  "notify_time" = :notify_time WHERE id = :id`
	_, err := s.db.NamedExecContext(s.ctx, sql, event)
	if err != nil {
		return err
	}
	return err
}

func (s *Storage) Delete(id interface{}) error {
	sql := `DELETE FROM events WHERE id = $1`
	_, err := s.db.ExecContext(s.ctx, sql, id)
	return err
}

func (s *Storage) GetByID(id interface{}) (common.Event, error) {
	event := common.Event{}
	sql := `SELECT "id","title","date_time","duration","description","user","notify_time" FROM events WHERE id = $1`
	err := s.db.GetContext(s.ctx, &event, sql, id)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return event, common.ErrEventNotFound
	}
	return event, err
}

func (s *Storage) List() ([]common.Event, error) {
	event := make([]common.Event, 0)
	sql := `SELECT "id","title","date_time","duration","description","user","notify_time" FROM events`
	err := s.db.SelectContext(s.ctx, &event, sql)
	if err != nil && err.Error() == "sql: no rows in result set" {
		return event, common.ErrEventNotFound
	}
	return event, err
}

func (s *Storage) PrepareStorage(log common.LoggerInterface) error {
	provider, err := goose.NewProvider(database.DialectPostgres, s.db.DB, migrations.Embed)
	if err != nil {
		log.Error("init goose", "error", err)
	}
	sources := provider.ListSources()
	for _, s := range sources {
		log.Info("Migration item", "type", s.Type, "version", s.Version, "path", filepath.Base(s.Path))
	}

	stats, err := provider.Status(s.ctx)
	if err != nil {
		log.Error("status", "error", err)
	}
	for _, s := range stats {
		log.Info("Migrate status", "type", s.Source.Type, "version", s.Source.Version, "duration", s.State)
	}
	results, err := provider.Up(s.ctx)
	if err != nil {
		log.Error("up", "error", err)
	}
	for _, r := range results {
		log.Info("Migrate done", "type", r.Source.Type, "version", r.Source.Version, "duration", (r.Duration).String())
	}

	return nil
}
