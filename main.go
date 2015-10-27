package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"runtime"

	"github.com/golang/glog"
	"github.com/guregu/kami"
	"github.com/zenazn/goji/graceful"
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
	glog.Infoln("Starting server...")

	// load config
	d, err := ioutil.ReadFile(*configFile)
	if err != nil {
		glog.Fatalln("[ERROR] read config.yml", err)
	}

	cfg, err := config.Load(bytes.NewReader(d), *env)
	if err != nil {
		glog.Fatalln("[ERROR] config Load", err)
	}
	psAll := sphinx.NewSphinx(cfg.PSConfig, cpus)
	graceful.PostHook(func() {
		glog.Flush()
	})

	ctx := context.Background()
	ctx = sphinx.NewContext(ctx, psAll)
	kami.Context = ctx
	kami.Serve()
}
