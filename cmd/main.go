package main

import (
	"flag"
	"github.com/Newlooc/dt/pkg/apis"
	"github.com/Newlooc/dt/pkg/dtmanager"
	log "github.com/sirupsen/logrus"
	"time"
)

var (
	start string
	end   string
	code  string
	frq   int
	money float64
	rate  float64
	debug bool
)

func main() {
	startTime, err := time.Parse(apis.DateFormat, start)
	if err != nil {
		log.WithError(err).Fatalf("Date format should be YYYY-MM-DD. %s given.", start)
	}
	endTime, err := time.Parse(apis.DateFormat, end)
	if err != nil {
		log.WithError(err).Fatalf("Date format should be YYYY-MM-DD. %s given.", start)
	}

	dtm, err := dtmanager.NewDTManager(code, startTime, endTime, frq, money, rate)
	if err != nil {
		log.WithError(err).Fatal("Manager init with error.")
	}

	if err := dtm.Run(); err != nil {
		log.WithError(err).Fatal("Manager run with error.")
	}
}

func init() {
	flag.BoolVar(&debug, "debug", false, "debug")
	flag.StringVar(&start, "start", "", "YYYY-MM-DD")
	flag.StringVar(&end, "end", "", "YYYY-MM-DD")
	flag.StringVar(&code, "code", "", "code")
	flag.Float64Var(&money, "money", 100, "money")
	flag.IntVar(&frq, "frq", 1, "frq")
	flag.Float64Var(&rate, "rate", 0.15, "rate")
	flag.Parse()

	if debug {
		log.SetLevel(log.DebugLevel)
	}
}
