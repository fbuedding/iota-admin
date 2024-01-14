package globals

import (
	"strings"

	"github.com/Netflix/go-env"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	Conf Config
)

type Config struct {
	LogLevel     string `env:"LOG_LEVEL,required=true"`
	CookieSecret string `env:"COOKIE_SECRET,required=true"`
	AppEnv       string `env:"APP_ENV,default=development"`
	BypassAuth   bool   `env:"BYPASS_AUTH"`
	IoTAHost     string `env:"IOTA_HOST,default=iot-agent"`
	IoTAPort     int    `env:"IOTA_PORT,default=4061"`
}

func init() {
	_, err := env.UnmarshalFromEnviron(&Conf)
	if err != nil {
		log.Fatal().Err(err).Msg("Error while reading envs")
	}

	setLogLevel(Conf.LogLevel)

}
func setLogLevel(ll string) {
	ll = strings.ToLower(ll)
	switch ll {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warning":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		log.Fatal().Msg("Log level need to be one of this: [TRACE DEBUG INFO WARNING ERROR FATAL PANIC]")
	}
}
