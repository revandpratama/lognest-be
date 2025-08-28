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

	var fresh bool
	var migrateCmd = &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		Run: func(cmd *cobra.Command, args []string) {
			log.Info().Msg("Running migrations...")

			if fresh {
				server := NewServer()
				server.Migrate("fresh")
			} else {
				server := NewServer()
				server.Migrate("")
			}

			// run your migrations here
		},
	}
	migrateCmd.Flags().BoolVarP(&fresh, "fresh", "f", false, "Run fresh migrations")

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
		app.WithAzureBlobStorage(),
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

func (s *Server) Migrate(flag string) {
	apps, err := app.NewApp(
		app.WithDB(),
	)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to create app")
	}

	if err := cmd.EnsureSchema(apps.DB, config.ENV.LOGNEST_SCHEMA); err != nil {
		log.Fatal().Err(err).Msg("failed to ensure schema")
	}

	if flag == "fresh" {
		if err := cmd.MigrateDatabaseFresh(apps.DB); err != nil {
			log.Fatal().Err(err).Msg("failed to migrate database fresh")
		}
	}

	if flag == "" {
		if err := cmd.MigrateDatabase(apps.DB); err != nil {
			log.Fatal().Err(err).Msg("failed to migrate database")
		}
	}

	log.Info().Msg("database migrated successfully")

	if err := apps.Stop(); err != nil {
		log.Error().Err(err).Msgf("failed to stop app cleanly, cause: %v", err)
	}
}

func (s *Server) GenerateModule(moduleName string) {
	cmd.GenerateModule(moduleName)
}
