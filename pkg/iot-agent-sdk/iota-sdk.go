package iotagentsdk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
)

const (
	urlBase        = "http://%v:%d"
	urlHealthcheck = urlBase + "/iot/about"
)

func (e ApiError) Error() string {
	return fmt.Sprintf("%s: %s", e.Name, e.Message)
}

func init() {
	logLvl := os.Getenv("LOG_LEVEL")
	if logLvl == "" {
		logLvl = "panic"
	}
	SetLogLevel(logLvl)
}

func SetLogLevel(ll string) {
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

func (i IoTA) Healthcheck() (*RespHealthcheck, error) {
	response, err := http.Get(fmt.Sprintf(urlHealthcheck, i.Host, i.Port))
	if err != nil {
		return nil, fmt.Errorf("Error while Healthcheck: %w", err)
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Error while Healthcheck: %w", err)
	}
	var respHealth RespHealthcheck
	json.Unmarshal(responseData, &respHealth)
	if respHealth.Version == "" {
		return nil, fmt.Errorf("Error healtchecking IoT-Agent, host: %s", i.Host)
	}
	log.Debug().Str("Response healthcheck", string(responseData)).Any("Healthcheck", respHealth).Send()
	return &respHealth, nil
}

func (i IoTA) GetAllServicePathsForService(service string) ([]string, error) {
	cgs, err := i.ListConfigGroups(FiwareService{service, "/*"})
	if err != nil {
		return nil, err
	}

	if cgs.Count == 0 {
		return nil, nil
	}

	servicePaths := []string{}
	for _, cg := range cgs.Services {
		if !slices.Contains(servicePaths, cg.ServicePath) {
			servicePaths = append(servicePaths, cg.ServicePath)
		}
	}

	return servicePaths, nil
}

func (i IoTA) Client() *http.Client {
	if i.client == nil {
		log.Debug().Msg("Creating http client")
		i.client = &http.Client{Timeout: 1 * time.Second}
	}
	return i.client
}
