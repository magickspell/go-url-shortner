package main

import (
  "fmt"
  "github.com/go-chi/chi/v5"
  "github.com/go-chi/chi/v5/middleware"
  _ "github.com/mattn/go-sqlite3" // init sqlite3 driver
  "go-url-shortner/internal/config"
  "go-url-shortner/internal/http-server/middleware/handlers/url/save"
  mwLogger "go-url-shortner/internal/http-server/middleware/logger"
  "go-url-shortner/internal/lib/logger/handlers/slogpretty"
  "go-url-shortner/internal/lib/logger/sl"
  "go-url-shortner/internal/storage/sqlite"
  "log/slog"
  "net/http"
  "os"
)

func main() {
  // TODO: init config: env
  cfg := config.MustLoad()
  fmt.Printf("CONIFG LOADED:\n%+v\n", cfg)

  // TODO: init logger: slog
  log := setupLogger(cfg.Env)
  log.Info("starting server", slog.String("env", cfg.Env), slog.Int64("version", 1))
  log.Debug("debug messages are enabled") // на проде дебага не будет

  // TODO: init storage: sqlite
  storage, err := sqlite.New(cfg.StoragePath)
  if err != nil {
    log.Error("error creating storage", sl.Err(err))
    os.Exit(1) // или просто return
  }

  //id, err := storage.SaveURL("https://google.com", "google")
  //if err != nil {
  //  log.Error("error saving url", sl.Err(err))
  //  os.Exit(1)
  //}
  //log.Info("sqved url", slog.Int64("id", id))
  //
  //id, err = storage.SaveURL("https://google.com", "google")
  //if err != nil {
  //  log.Error("error saving url", sl.Err(err))
  //  os.Exit(1)
  //}

  _ = storage

  // TODO: init router: chi, "chi render"
  router := chi.NewRouter()

  // middleware
  router.Use(middleware.RequestID) // добавляем ИДшники для дебага
  router.Use(middleware.RealIP)    // ip пользователя
  //router.Use(middleware.Logger)    // log chi
  router.Use(mwLogger.New(log))    // самописный логгер
  router.Use(middleware.Recoverer) // автовосстановление после падения
  router.Use(middleware.URLFormat) // форматирование урлов

  router.Post("/url", save.New(log, storage)) // методы везхде должны называться одинаково
  //router.Get("/{alias}", redirect.New(log, storage))

  // TODO: init server: run server
  log.Info("starting server", slog.String("env", cfg.Env), slog.Int64("version", 1), slog.String("address", cfg.Address))

  srv := &http.Server{
    Addr:         cfg.Address,
    Handler:      router,
    ReadTimeout:  cfg.HTTPServer.Timeout,
    WriteTimeout: cfg.HTTPServer.Timeout,
    IdleTimeout:  cfg.HTTPServer.IdleTimeout,
  }

  // блкирующий вызов запуска сервера
  if err := srv.ListenAndServe(); err != nil {
    log.Error("error starting server", sl.Err(err))
  }
  log.Error("stopping server")
}

const (
  envLocal = "local"
  envDev   = "dev"
  envProd  = "prod"
)

func setupLogger(env string) *slog.Logger {
  var log *slog.Logger

  switch env {
  case envLocal:
    //log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
    log = setupPrettySlog()
  case envDev:
    log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
  case envProd:
    log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
  }

  return log
}

func setupPrettySlog() *slog.Logger {
  opts := slogpretty.PrettyHandlerOptions{
    SlogOpts: &slog.HandlerOptions{
      Level: slog.LevelDebug,
    },
  }

  handler := opts.NewPrettyHandler(os.Stdout)

  return slog.New(handler)
}
