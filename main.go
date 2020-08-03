package main

import (
	"context"
	"errors"
	"explorer/crawler"
	"explorer/db"
	"explorer/handler"
	"explorer/utils"
	"explorer/web"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type hooker struct{}

func start(_ *hooker, _ web.Daemon, _ *web.HTTPHandler) {}
func NewHooker(lc fx.Lifecycle, srv web.Daemon) *hooker {

	var check = func(d web.Daemon, errChan chan error) {
		err := d.Start()
		if err == nil {
			return
		}

		select {
		case errChan <- err:
		default:
		}
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			errChan := make(chan error)

			go check(srv, errChan)

			select {
			case e := <-errChan:
				return errors.New(`start fail: ` + e.Error())
			case <-time.After(time.Second * 2):
				return nil
			}
		},

		OnStop: func(ctx context.Context) error {
			srv.Stop()
			time.Sleep(time.Second)
			return nil
		},
	})
	return &hooker{}
}

func setupZeroLogFieldName() {
	zerolog.TimestampFieldName = "ts"
	zerolog.LevelFieldName = "l"
	zerolog.MessageFieldName = "m"
	zerolog.CallerFieldName = "c"
	//zerolog.TimeFieldFormat = "0102-15:04:05Z07"
	zerolog.ErrorFieldName = "e"
	zerolog.CallerMarshalFunc = func(file string, line int) string {
		return filepath.Base(file) + ":" + strconv.Itoa(line)
	}
	l := log.With().Caller().Logger()
	log.Logger = l
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	configPath := os.Getenv(`explorer_config_path`)

	_, err := utils.LoadViper("", configPath)
	if err != nil {
		panic(err)
	}

	setupZeroLogFieldName()
	zerologLevel := viper.GetString(`LOG.Level`)
	fmt.Printf("current log level:%s.\n", zerologLevel)

	level, _ := zerolog.ParseLevel(zerologLevel)
	if level == zerolog.NoLevel {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	crawler.OnStart()

	mgoStore := db.NewMongoStore()

	engine := web.NewEngine()

	app := fx.New(
		fx.Provide(func() *gin.Engine { return engine }),
		fx.Provide(func() db.MgoOperator { return mgoStore }),
		fx.Provide(handler.NewBaseHandler),
		fx.Provide(handler.NewDelegationsHandler),
		fx.Provide(handler.NewValidatorHandler),
		fx.Provide(handler.NewTokensHandler),
		fx.Provide(handler.NewAccountHandler),
		fx.Provide(handler.NewTransactionHandler),
		fx.Provide(handler.NewBlockHandler),
		fx.Provide(web.NewHTTPHandler),
		fx.Provide(web.NewServer),
		fx.Provide(NewHooker),
		fx.Invoke(start),
	)
	app.Run()
}
