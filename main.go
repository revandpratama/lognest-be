package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"

	"github.com/revandpratama/lognest/cmd"
	"github.com/revandpratama/lognest/config"
	"github.com/revandpratama/lognest/internal/app"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type Server struct {
	shutdownCh chan os.Signal
	errCh      chan error
}

func NewServer() *Server {
	return &Server{
		shutdownCh: make(chan os.Signal, 1),
		errCh:      make(chan error),
	}
}

func main() {

	if err := config.LoadConfig(); err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	var rootCmd = &cobra.Command{
		Use:   "app",
		Short: "My app with subcommands",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info().Msg("starting server...")

			server := NewServer()
			server.Start()
		},
	}

	var migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info().Msg("Running migrations...")

			server := NewServer()
			server.Migrate()
			// run your migrations here
		},
	}

	var generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate modules",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info().Msg("Generating modules...")

			server := NewServer()
			server.GenerateModule(args[0])
			// run your migrations here
		},
	}

	rootCmd.AddCommand(migrateCmd, generateCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func (s *Server) Start() {

	signal.Notify(s.shutdownCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	apps, err := app.NewApp(
		app.WithDB(),
		app.WithRESTServer(),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create app")
	}

	select {
	case shutdown := <-s.shutdownCh:
		log.Info().Msgf("gracefully shutting down the app: %v", shutdown)
		if err := apps.Stop(); err != nil {
			log.Error().Err(err).Msgf("failed to stop app cleanly, cause: %v", err)
		}

		byeMsg := randomGoodbye()

		log.Info().Msgf("server stopped, %s!", byeMsg)
	case err := <-s.errCh:
		log.Error().Err(err).Msgf("failed to start server, cause: %v", err)
	}
}

func randomGoodbye() string {
	byes := []string{
		"Goodbye",                // English
		"Adiós",                  // Spanish
		"Au revoir",              // French
		"Auf Wiedersehen",        // German
		"Ciao",                   // Italian
		"Sayonara",               // Japanese
		"Annyeong",               // Korean
		"Do svidaniya",           // Russian
		"Selamat tinggal",        // Indonesian
		"Shukran, maʿ al-salāma", // Arabic
	}

	return byes[rand.Intn(len(byes))]
}

func (s *Server) Migrate() {
	apps, err := app.NewApp(
		app.WithDB(),
	)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to create app")
	}

	if err := cmd.EnsureSchema(apps.DB, config.ENV.LOGNEST_SCHEMA); err != nil {
		log.Fatal().Err(err).Msg("failed to ensure schema")
	}

	if err := cmd.MigrateDatabase(apps.DB); err != nil {
		log.Fatal().Err(err).Msg("failed to migrate database")
	}

	log.Info().Msg("database migrated successfully")

	if err := apps.Stop(); err != nil {
		log.Error().Err(err).Msgf("failed to stop app cleanly, cause: %v", err)
	}
}

func (s *Server) GenerateModule(moduleName string) {
	cmd.GenerateModule(moduleName)
}