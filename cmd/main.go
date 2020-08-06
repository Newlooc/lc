package main

import (
	"fmt"
	"github.com/Newlooc/dt/pkg/apis"
	"time"
	//	"github.com/Newlooc/dt/cmd/data"
	"github.com/Newlooc/dt/pkg/dtmanager"
	//	"github.com/Newlooc/dt/pkg/parser"
	//"github.com/Newlooc/dt/pkg/spider"

	log "github.com/sirupsen/logrus"
)

func main() {

	//	dtmock := parser.NewDTMock()
	//	dtmock.Parse(data.D)
	//	fmt.Printf("%+v", dtmock)
	//	v := spider.NewVisit()
	//	h, b := v.Do("http://fund.eastmoney.com/data/FundInvestCaculator_AIPDatas.aspx?fcode=000118&sdate=2018-1-1&edate=2020-1-1&shdate=&round=1&dtr=1&p=0.15&je=100&stype=1&needfirst=2&jsoncallback=FundDTSY.result")
	//	fmt.Printf("%+v\n", h)
	//	fmt.Printf("%+v\n", string(b))

	start, _ := time.Parse(apis.DateFormat, "2018-02-05")
	end, _ := time.Parse(apis.DateFormat, "2020-09-17")

	dtm, err := dtmanager.NewDTManager("000175", start, end, 1, float64(100), 0.15)
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
