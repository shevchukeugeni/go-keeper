package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"

	"keeper-project/internal/server"
	"keeper-project/internal/store/file"
	"keeper-project/internal/store/file/storage/minio"
	"keeper-project/internal/store/postgres"
	"keeper-project/internal/store/postgres/secrets/cards"
	"keeper-project/internal/store/postgres/secrets/creds"
	"keeper-project/internal/store/postgres/secrets/notes"
	"keeper-project/internal/store/postgres/users"
)

type config struct {
	Address        string `env:"ADDRESS" `
	DatabaseDSN    string `env:"DATABASE_DSN"`
	MinioURL       string `env:"MINIO_URL"`
	MinioAccessKey string `env:"MINIO_ACCESS_KEY"`
	MinioSecretKey string `env:"MINIO_SECRET_KEY"`
}

var cfg config

func init() {
	flag.StringVar(&cfg.Address, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&cfg.DatabaseDSN, "d", "postgresql://postgres:123456@localhost:5432/keeper", "database connection url")
	flag.StringVar(&cfg.MinioURL, "m-url", "localhost:9000", "minio URL")
	flag.StringVar(&cfg.MinioAccessKey, "m-access", "minio", "minio access key")
	flag.StringVar(&cfg.MinioSecretKey, "m-secret", "minio123", "minio secret key")
}

func main() {
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	logger.Info("Starting...")

	var (
		router http.Handler
		wg     sync.WaitGroup
	)

	ctx, cancelCtx := context.WithCancel(context.Background())

	db, err := postgres.NewPostgresDB(postgres.Config{URL: cfg.DatabaseDSN})
	if err != nil {
		logger.Fatal("failed to initialize db: " + err.Error())
	}
	defer func() {
		if db != nil {
			db.Close()
		}
	}()

	userStore := users.NewRepository(db)
	notesStore := notes.NewRepository(db)
	credsStore := creds.NewRepository(db)
	cardsStore := cards.NewRepository(db)

	fileStore, err := minio.NewStorage(logger, cfg.MinioURL, cfg.MinioAccessKey, cfg.MinioSecretKey)
	if err != nil {
		logger.Fatal("unable to create minio storage", zap.Error(err))
		return
	}

	fileService, err := file.NewService(fileStore, logger)
	if err != nil {
		logger.Fatal("unable to create file service", zap.Error(err))
		return
	}

	router = server.SetupRouter(logger, userStore, notesStore, credsStore, cardsStore, fileService)

	logger.Info("Running HTTP server on", zap.String("address", cfg.Address))
	srv := http.Server{Addr: cfg.Address, Handler: router}
	// через этот канал сообщим основному потоку, что соединения закрыты
	idleConnsClosed := make(chan struct{})
	// канал для перенаправления прерываний
	// поскольку нужно отловить всего одно прерывание,
	// ёмкости 1 для канала будет достаточно
	sigint := make(chan os.Signal, 1)
	// регистрируем перенаправление прерываний
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	// запускаем горутину обработки пойманных прерываний
	go func() {
		// читаем из канала прерываний
		// поскольку нужно прочитать только одно прерывание,
		// можно обойтись без цикла
		<-sigint
		// получили сигнал os.Interrupt, запускаем процедуру graceful shutdown
		if err := srv.Shutdown(ctx); err != nil {
			// ошибки закрытия Listener
			logger.Error("HTTP server Shutdown", zap.Error(err))
		}
		// сообщаем основному потоку,
		// что все сетевые соединения обработаны и закрыты

		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		// ошибки старта или остановки Listener
		logger.Fatal("HTTP server ListenAndServe", zap.Error(err))
	}
	// ждём завершения процедуры graceful shutdown
	<-idleConnsClosed
	// получили оповещение о завершении
	// здесь можно освобождать ресурсы перед выходом,
	// например закрыть соединение с базой данных,
	// закрыть открытые файлы

	cancelCtx()
	wg.Wait()
	logger.Info("Server Shutdown gracefully")

}
