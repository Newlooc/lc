package dtmanager

import (
	"errors"
	"fmt"
	"github.com/Newlooc/dt/pkg/apis"
	"github.com/Newlooc/dt/pkg/output"
	"github.com/Newlooc/dt/pkg/parser"
	"github.com/Newlooc/dt/pkg/spider"
	log "github.com/sirupsen/logrus"
	"time"
)

var (
	errorInvalidDateRange = errors.New("Invalid data range")
	errorInvalidArgument  = errors.New("Invalid Arguments")
	urlTemplate           = "http://fund.eastmoney.com/data/FundInvestCaculator_AIPDatas.aspx?fcode=%s&sdate=%s&edate=%s&shdate=&round=%d&dtr=1&p=%.2f&je=%.2f&stype=1&needfirst=2&jsoncallback=FundDTSY.result"
	sleepDuration         = time.Millisecond * 300
)

type Manager struct {
	Config *ManagerConfig
	URLs   map[apis.URLConfig]string
	Parsed map[apis.URLConfig]*parser.DTMock
}

type ManagerConfig struct {
	Code          string
	Start         time.Time
	End           time.Time
	IntervalStart []time.Time
	IntervalEnd   []time.Time
	Frq           int
	Limit         time.Duration
	Money         float64
	Rate          float64
}

func NewDTManager(code string, start, end time.Time, frq int, money, rate float64) (*Manager, error) {
	if start.After(end) || start.Equal(end) {
		log.WithError(errorInvalidDateRange).Error("%s to %s.", start.Format(apis.DateFormat), end.Format(apis.DateFormat))
		return nil, errorInvalidDateRange
	}

	if money <= float64(0) || rate >= float64(100) || rate <= float64(0) {
		log.WithError(errorInvalidArgument).Errorf("code %s. start %s. end %s. frq %s. money %s. rate %s.", code, start, end, frq, money, rate)
		return nil, errorInvalidArgument
	}

	conf := &ManagerConfig{
		Code:  code,
		Start: start,
		End:   end,
		Frq:   frq,
		Money: money,
		Rate:  rate,
	}

	return &Manager{
		Config: conf,
		URLs:   make(map[apis.URLConfig]string),
		Parsed: make(map[apis.URLConfig]*parser.DTMock),
	}, nil
}

func (mr *Manager) Run() error {
	log.Info("Start gen interval")
	if err := mr.genInterval(); err != nil {
		log.WithError(err).Error("Failed to gen interval.")
		return err
	}

	log.Info("Start gen URLs")
	if err := mr.genURLs(); err != nil {
		log.WithError(err).Error("Failed to gen urls.")
		return err
	}
	log.Debugf("DEBUG %+v", mr.URLs)

	log.Info("Start visit URLs")
	visit := spider.NewVisit()
	var done, progress float64
	for urlConfig, url := range mr.URLs {
		_, b, err := visit.Do(url, false)
		progress = done / float64(len(mr.URLs)) * float64(100)
		log.Infof("[%.2f%%] Finished. %s-%s", progress, urlConfig.Start.Format(apis.DateFormat), urlConfig.End.Format(apis.DateFormat))
		log.Debugf("URL %s visited.", url)
		if err != nil {
			log.WithError(err).Errorf("Failed to run url %s.", url)
		}
		parsed := parser.NewDTMock()
		parsed.Parse(string(b))
		mr.Parsed[urlConfig] = parsed

		done++
		time.Sleep(sleepDuration)
	}

	// output
	log.Info("Start output")
	excel := output.NewExcel(mr.genXlsFileName(), mr.Config.Code)
	if err := excel.Write(mr.Parsed, mr.Config.IntervalEnd, mr.Config.IntervalStart); err != nil {
		log.WithError(err)
	}

	return nil
}

func (mr *Manager) genXlsFileName() string {
	return fmt.Sprintf("./xls/%s-%s-%s.xlsx", mr.Config.Code, mr.Config.Start.Format(apis.DateFormat), mr.Config.End.Format(apis.DateFormat))
}

func (mr *Manager) genInterval() error {
	cosr := mr.Config.Start
	for cosr.Before(mr.Config.End) {
		mr.Config.IntervalStart = append(mr.Config.IntervalStart, cosr)
		// Upgrade by month
		// TODO
		cosr = cosr.AddDate(1, 0, 0)
	}

	mr.Config.IntervalEnd = mr.Config.IntervalStart[1:]
	mr.Config.IntervalEnd = append(mr.Config.IntervalEnd, mr.Config.End)

	return nil
}

func (mr *Manager) genURLs() error {
	for _, end := range mr.Config.IntervalEnd {
		for _, start := range mr.Config.IntervalStart {
			if start.Before(end) {
				mr.URLs[apis.URLConfig{
					Start: start,
					End:   end,
				}] = fmt.Sprintf(urlTemplate, mr.Config.Code, start.Format(apis.DateFormat), end.Format(apis.DateFormat), mr.Config.Frq, mr.Config.Rate, mr.Config.Money)
			}
		}
	}
	return nil
}
