package main

import (
	"explorer/crawler"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
}

func loadViper(envPrefix, configPath string) (*viper.Viper, error) {
	if len(envPrefix) > 0 {
		viper.SetEnvPrefix(envPrefix)
	}

	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		if err, ok := err.(*os.PathError); !ok {
			return nil, err
		}
	}
	return viper.GetViper(), nil
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
	configPath, err := filepath.Abs(os.Args[1])
	if err != nil {
		panic(err)
	}

	if configPath == "" {
		panic(`must had config file`)
	}

	fmt.Printf("config file path:%s.\n", configPath)

	_, err = loadViper("", configPath)
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

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"

	}
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"https://*.foo.com"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true}))
	beego.Run()
}
