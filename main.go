package main

import (
	"explorer/crawler"
	"fmt"
	"github.com/rs/zerolog/log"
	"path/filepath"
	"strconv"

	_ "explorer/routers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

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

//go:generate sh -c "echo 'package routers; import \"github.com/astaxie/beego\"; func init() {beego.BConfig.RunMode = beego.DEV}' > routers/0.go"
//go:generate sh -c "echo 'package routers; import \"os\"; func init() {os.Exit(0)}' > routers/z.go"
//go:generate go run $GOFILE conf/config.toml
//go:generate sh -c "rm routers/0.go routers/z.go"
