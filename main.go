package main

import (
	"bufio"
	"flag"

	"context"
	"net/http"
	"os"
	"time"

	"github.com/Latezly/nyaa_go/config"
	"github.com/Latezly/nyaa_go/controllers"
	"github.com/Latezly/nyaa_go/controllers/databasedumps"
	"github.com/Latezly/nyaa_go/models"
	"github.com/Latezly/nyaa_go/utils/cookies"
	"github.com/Latezly/nyaa_go/utils/log"
	"github.com/Latezly/nyaa_go/utils/publicSettings"
	"github.com/Latezly/nyaa_go/utils/search"
	"github.com/Latezly/nyaa_go/utils/signals"
)

var buildversion string

// RunServer runs webapp mainloop
func RunServer() {
	conf := config.Get()
	// TODO Use config from cli
	os.Mkdir(databasedumpsController.DatabaseDumpPath, 0700)
	// TODO Use config from cli
	os.Mkdir(databasedumpsController.GPGPublicKeyPath, 0700)

	http.Handle("/", controllers.CSRFRouter)

	// Set up server,
	srv := &http.Server{
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	l, err := CreateHTTPListener(conf)
	log.CheckError(err)
	if err != nil {
		return
	}
	log.Infof("listening on %s", l.Addr())
	// http.Server.Shutdown closes associated listeners/clients.
	// context.Background causes srv to indefinitely try to
	// gracefully shutdown. add a timeout if this becomes a problem.
	signals.OnInterrupt(func() {
		srv.Shutdown(context.Background())
	})
	err = srv.Serve(l)
	if err == nil {
		log.Panic("http.Server.Serve never returns nil")
	}
	if err == http.ErrServerClosed {
		return
	}
	log.CheckError(err)
}

func main() {
	if buildversion != "" {
		config.Get().Build = buildversion
	} else {
		config.Get().Build = "unknown"
	}
	defaults := flag.Bool("print-defaults", false, "print the default configuration file on stdout")
	callback := config.BindFlags()
	flag.Parse()
	if *defaults {
		stdout := bufio.NewWriter(os.Stdout)
		err := config.Get().Pretty(stdout)
		if err != nil {
			log.Fatal(err.Error())
		}
		err = stdout.Flush()
		if err != nil {
			log.Fatal(err.Error())
		}
		os.Exit(0)
	} else {
		callback()
		var err error
		models.ORM, err = models.GormInit(models.DefaultLogger)
		if err != nil {
			log.Fatal(err.Error())
		}
		if config.Get().Search.EnableElasticSearch {
			log.Info("ES Enabled in Config")
			models.ElasticSearchClient, _ = models.ElasticSearchInit()
		} else {
			log.Info("ES is disabled in Config")
		}
		err = publicSettings.InitI18n(config.Get().I18n, cookies.NewCurrentUserRetriever())
		if err != nil {
			log.Fatal(err.Error())
		}
		err = search.Configure(&config.Get().Search)
		if err != nil {
			log.Fatal(err.Error())
		}
		signals.Handle()
		if len(config.Get().Torrents.FileStorage) > 0 {
			err := os.MkdirAll(config.Get().Torrents.FileStorage, 0700)
			if err != nil {
				log.Fatal(err.Error())
			}
		}
		RunServer()
	}
}
