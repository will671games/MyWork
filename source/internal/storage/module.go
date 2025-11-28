package storage

import (
	"TestProject/source/config"
	"TestProject/source/internal/application"
	"TestProject/source/internal/storage/wallet"
	"context"
	"fmt"
	"log/slog"
	"time"

	slogGorm "github.com/orandin/slog-gorm"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const moduleName = "storage"

var Module = fx.Module(
	moduleName,
	fx.Provide(
		NewDatabaseConnection,
		wallet.NewRepo,
		func(repo *wallet.Repo) application.WalletRepo {
			return repo
		},
	),
	fx.Decorate(
		func(log *slog.Logger) *slog.Logger {
			return log.With(slog.String("module", moduleName))
		},
	),
)

const (
	LogErr    = slog.Level(3)
	LogNotice = slog.Level(5)
)

func NewDatabaseConnection(lc fx.Lifecycle, conf config.DBConfig, logger *slog.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Europe/Moscow",
		conf.Host,
		conf.User,
		conf.Password,
		conf.Database,
		conf.Port,
		conf.SSLMode,
	)

	gormLogger := slogGorm.New(
		slogGorm.WithHandler(logger.Handler()),
		slogGorm.SetLogLevel(slogGorm.ErrorLogType, LogErr),
		slogGorm.SetLogLevel(slogGorm.SlowQueryLogType, LogNotice),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                 gormLogger,
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(50)
	sqlDB.SetMaxOpenConns(200)
	sqlDB.SetConnMaxLifetime(time.Minute * 15)

	err = db.AutoMigrate(&wallet.Wallet{})
	if err != nil {
		logger.Error("cannot auto migrate", slog.Any("error", err))
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			sqlDB, err := db.DB()
			if err != nil {
				return err
			}
			return sqlDB.Close()
		},
	})

	return db, nil
}
