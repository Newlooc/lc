package main

import (
	"fmt"
	"github.com/Newlooc/dt/pkg/apis"
	"github.com/Newlooc/dt/pkg/dtmanager"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {

	start, _ := time.Parse(apis.DateFormat, "2010-01-01")
	end, _ := time.Parse(apis.DateFormat, "2020-01-01")

	dtm, err := dtmanager.NewDTManager("161005", start, end, 1, float64(100), 0.15)
	if err != nil {
		log.WithError(err).Fatal("INIT")
	}
	dtm.Run()
	if err != nil {
		log.WithError(err).Fatal("RUN")
	}
	fmt.Printf("%+v\n", dtm.Config.IntervalStart)
	fmt.Printf("%+v\n", dtm.Config.IntervalEnd)
}

func init() {

	log.SetLevel(log.DebugLevel)
}
