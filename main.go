package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"runtime"

	log "github.com/Sirupsen/logrus"
	"github.com/guregu/kami"
	"golang.org/x/net/context"

	"github.com/shunsukeaihara/sphinx-httpserver/config"
	"github.com/shunsukeaihara/sphinx-httpserver/sphinx"
)

var (
	env        = flag.String("e", "development", "app server enviroment")
	configFile = flag.String("c", "config.yml", "app server configuration file")
)

func main() {

	flag.Parse()
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stderr)
	if *env == "production" {
		log.SetLevel(log.WarnLevel)
	} else {
		log.SetLevel(log.DebugLevel)
	}
	log.Infoln("Starting server...")

	// load config
	d, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatalln("[ERROR] read config.yml", err)
	}

	cfg, err := config.Load(bytes.NewReader(d), *env)
	if err != nil {
		log.Fatalln("[ERROR] config Load", err)
	}
	psAll := sphinx.NewSphinx(cfg.PSConfig, cpus)

	ctx := context.Background()
	ctx = sphinx.NewContext(ctx, psAll)
	kami.Context = ctx
	kami.Serve()
}
