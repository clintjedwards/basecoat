package service

import (
	"os"
	"strings"

	"github.com/clintjedwards/basecoat/internal/app"
	"github.com/clintjedwards/basecoat/internal/cmd/cl"
	"github.com/clintjedwards/basecoat/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

var cmdServiceStart = &cobra.Command{
	Use:   "start",
	Short: "Start the Basecoat GRPC/HTTP combined server",
	Long: `Start the Basecoat GRPC/HTTP combined server.

Basecoat runs as a GRPC backend combined with GRPC-WEB/HTTP. Running this command attempts to start the long
running service. This command will block and only gracefully stop on SIGINT or SIGTERM signals

## Configuration

### List of Environment Variables

` + strings.Join(config.GetAPIEnvVars(), "\n"),
	RunE: serverStart,
}

func init() {
	cmdServiceStart.Flags().BoolP("dev-mode", "d", false, "Alters several feature flags such that development is easy. "+
		"This is not to be used in production and may turn off features that are useful for even development like authentication")
	CmdService.AddCommand(cmdServiceStart)
}

func serverStart(cmd *cobra.Command, _ []string) error {
	cl.State.Fmt.Finish()

	configPath, _ := cmd.Flags().GetString("config")
	devMode, _ := cmd.Flags().GetBool("dev-mode")
	conf, err := config.InitAPIConfig(configPath, true, devMode)
	if err != nil {
		log.Fatal().Err(err).Msg("error in config initialization")
	}

	setupLogging(conf.LogLevel, conf.Development.PrettyLogging)
	app.StartServices(conf)

	return nil
}

func setupLogging(loglevel string, pretty bool) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.With().Caller().Logger()
	zerolog.SetGlobalLevel(parseLogLevel(loglevel))
	if pretty {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

func parseLogLevel(loglevel string) zerolog.Level {
	switch loglevel {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	default:
		log.Error().Msgf("loglevel %s not recognized; defaulting to debug", loglevel)
		return zerolog.DebugLevel
	}
}
