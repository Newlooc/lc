package main

import (
	"fmt"
	"github.com/Newlooc/dt/cmd/data"
	"github.com/Newlooc/dt/pkg/parser"
	log "github.com/sirupsen/logrus"
)

func main() {

	dtmock := parser.NewDTMock()
	dtmock.Parse(data.D)
	fmt.Printf("%+v", dtmock)
}

func init() {

	log.SetLevel(log.DebugLevel)
}
