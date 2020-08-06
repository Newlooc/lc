package parser

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

const (
	dateformat = "2006-01-02"
)

type response struct {
	Info string `json:"info"`
}

type DTMock struct {
	Code              string
	CNName            string
	TotalRound        int64
	BaseInvestAmount  float64
	InvestType        string
	FinalAmount       float64
	ProfitPercent     float64
	MockCalculateDate time.Time
	dtRecords         []*dtRecord
}

type dtRecord struct {
	Date   time.Time
	Price  float64
	Amount float64
	Total  float64
}

var (
	ERRORUnexpectedRawData = errors.New("unexpected summary raw data")
)

func NewDTMock() *DTMock {
	return &DTMock{}
}

func (dt *DTMock) Parse(raw string) error {
	raw = strings.Trim(raw, " \n")
	uselessheadindex := strings.Index(raw, "{")
	if uselessheadindex == -1 {
		log.WithError(ERRORUnexpectedRawData).Error("raw data format error. %s", raw)
		return ERRORUnexpectedRawData
	}

	raw = raw[uselessheadindex:]
	uselesstailindex := strings.LastIndex(raw, ")")
	if uselesstailindex == -1 {
		log.WithError(ERRORUnexpectedRawData).Error("raw data format error. %s", raw)
		return ERRORUnexpectedRawData
	}
	raw = raw[:uselesstailindex]
	log.Debug(raw)

	response := &response{}
	if err := json.Unmarshal([]byte(raw), response); err != nil {
		log.WithError(err).Error("response should be jsonp body. %s.", raw)
		return err
	}

	dtmixitems := strings.Split(response.Info, " ")
	rawsummarys := dtmixitems[0]
	rawrecords := dtmixitems[1:]

	if err := dt.parseSummary(rawsummarys); err != nil {
		log.WithError(err).Error("failed to parse summary info.")
		return err
	}

	if err := dt.parseRecords(rawrecords); err != nil {
		log.WithError(err).Error("failed to parse dt record.")
		return err
	}

	return nil
}

func (dt *DTMock) parseRecords(raw []string) error {
	if len(raw) == 0 {
		return nil
	}
	if dt.dtRecords == nil {
		dt.dtRecords = make([]*dtRecord, 0, len(raw))
	}

	for _, recordstr := range raw {
		parsedrecord := &dtRecord{}
		recordpreparse := strings.Split(strings.TrimSpace(recordstr), "_")
		if len(recordpreparse) != 2 {
			log.WithError(ERRORUnexpectedRawData).Errorf("record unrecognised. %s.", recordstr)
			//return ERRORUnexpectedRawData
		} else {
			date, err := time.Parse(strings.TrimSpace(dateformat), recordpreparse[1])
			if err != nil {
				log.WithError(err).Errorf("failed to parse date, raw data: %s", recordstr)
				//return err
			}
			parsedrecord.Date = date
		}

		recordinfo := strings.Split(strings.TrimSpace(recordpreparse[0]), "~")
		if len(recordinfo) != 4 {
			log.WithError(ERRORUnexpectedRawData).Errorf("record info unrecognised. %s.", recordstr)
			//return ERRORUnexpectedRawData
		} else {
			price, err := strconv.ParseFloat(numstrnor(recordinfo[1]), 64)
			if err != nil {
				log.WithError(err).Error("failed to parse record info price. raw data: %s.", recordstr)
			}
			parsedrecord.Price = price

			amount, err := strconv.ParseFloat(numstrnor(recordinfo[2]), 64)
			if err != nil {
				log.WithError(err).Error("failed to parse record info amount. raw data: %s.", recordstr)
			}
			parsedrecord.Amount = amount

			total, err := strconv.ParseFloat(numstrnor(recordinfo[3]), 64)
			if err != nil {
				log.WithError(err).Error("failed to parse record info total. raw data: %s.", recordstr)
			}
			parsedrecord.Total = total
		}

		dt.dtRecords = append(dt.dtRecords, parsedrecord)
		log.Debugf("record parse successfully. %+v", parsedrecord)
	}

	return nil
}

func (dt *DTMock) parseSummary(raw string) error {

	summarypreparse := strings.Split(raw, "|")
	if len(summarypreparse) != 8 {
		log.WithError(ERRORUnexpectedRawData).Errorf("raw: %s.", raw)
		return ERRORUnexpectedRawData
	}

	code := strings.TrimSpace(summarypreparse[0])
	cnname := strings.TrimSpace(summarypreparse[1])
	investtype := strings.TrimSpace(summarypreparse[4])

	totalround, err := strconv.Atoi(numstrnor(summarypreparse[2][:len(summarypreparse[2])-3]))
	if err != nil {
		log.WithError(err).Errorf("failed to parse totalround, raw data: %s", raw)
		return err
	}

	baseinvestamount, err := strconv.ParseFloat(numstrnor(summarypreparse[3]), 64)
	if err != nil {
		log.WithError(err).Errorf("failed to parse baseinvestamount, raw data: %s", raw)
		return err
	}

	finalamount, err := strconv.ParseFloat(numstrnor(summarypreparse[5]), 64)
	if err != nil {
		log.WithError(err).Errorf("failed to parse finalamount, raw data: %s", raw)
		return err
	}

	profitpercent, err := strconv.ParseFloat(numstrnor(summarypreparse[6][:len(summarypreparse[6])-1]), 64)
	if err != nil {
		log.WithError(err).Errorf("failed to parse profitpercent, raw data: %s", raw)
		return err
	}

	mockcalculatedate, err := time.Parse(dateformat, strings.TrimSpace(summarypreparse[7]))
	if err != nil {
		log.WithError(err).Errorf("failed to parse mockcalculatedate, raw data: %s", raw)
		return err
	}

	dt.Code = code
	dt.CNName = cnname
	dt.TotalRound = int64(totalround)
	dt.BaseInvestAmount = baseinvestamount
	dt.InvestType = investtype
	dt.FinalAmount = finalamount
	dt.ProfitPercent = profitpercent
	dt.MockCalculateDate = mockcalculatedate

	return nil
}

func numstrnor(giv string) string {
	return strings.ReplaceAll(strings.Trim(giv, " \n"), ",", "")
}
