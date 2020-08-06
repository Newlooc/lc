package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

const (
	dateformat = "2006-01-02"
)

type response struct {
	info string `json:"info"`
}

type dtmock struct {
	code              string
	cnname            string
	totalround        int64
	baseinvestamount  float64
	investtype        string
	finalamount       float64
	profitpercent     float64
	mockcalculatedate time.time
	dtrecords         []*dtrecord
}

type dtrecord struct {
	date   time.time
	price  float64
	amount float64
	total  float64
}

var (
	errorunexpectedrawdata = errors.new("unexpected summary raw data")
)

func newdtmock() *dtmock {
	return &dtmock{}
}

func (dt *dtmock) parse(raw string) error {
	raw = strings.trim(raw, " \n")
	uselessheadindex := strings.index(raw, "{")
	if uselessheadindex == -1 {
		log.witherror(errorunexpectedrawdata).error("raw data format error. %s", raw)
		return errorunexpectedrawdata
	}

	raw = raw[uselessheadindex:]
	uselesstailindex := strings.lastindex(raw, ")")
	if uselesstailindex == -1 {
		log.witherror(errorunexpectedrawdata).error("raw data format error. %s", raw)
		return errorunexpectedrawdata
	}
	raw = raw[:uselesstailindex]
	log.debug(raw)

	response := &response{}
	if err := json.unmarshal([]byte(raw), response); err != nil {
		log.witherror(err).error("response should be jsonp body. %s.", raw)
		return err
	}

	dtmixitems := strings.split(response.info, " ")
	rawsummarys := dtmixitems[0]
	rawrecords := dtmixitems[1:]

	if err := dt.parsesummary(rawsummarys); err != nil {
		log.witherror(err).error("failed to parse summary info.")
		return err
	}

	if err := dt.parserecords(rawrecords); err != nil {
		log.witherror(err).error("failed to parse dt record.")
		return err
	}

	return nil
}

func (dt *dtmock) parserecords(raw []string) error {
	if len(raw) == 0 {
		return nil
	}
	if dt.dtrecords == nil {
		dt.dtrecords = make([]*dtrecord, 0, len(raw))
	}

	for _, recordstr := range raw {
		parsedrecord := &dtrecord{}
		recordpreparse := strings.split(strings.trimspace(recordstr), "_")
		if len(recordpreparse) != 2 {
			log.witherror(errorunexpectedrawdata).errorf("record unrecognised. %s.", recordstr)
			//return errorunexpectedrawdata
		} else {
			date, err := time.parse(strings.trimspace(dateformat), recordpreparse[1])
			if err != nil {
				log.witherror(err).errorf("failed to parse date, raw data: %s", recordstr)
				//return err
			}
			parsedrecord.date = date
		}

		recordinfo := strings.split(strings.trimspace(recordpreparse[0]), "~")
		if len(recordinfo) != 4 {
			log.witherror(errorunexpectedrawdata).errorf("record info unrecognised. %s.", recordstr)
			//return errorunexpectedrawdata
		} else {
			price, err := strconv.parsefloat(numstrnor(recordinfo[1]), 64)
			if err != nil {
				log.witherror(err).error("failed to parse record info price. raw data: %s.", recordstr)
			}
			parsedrecord.price = price

			amount, err := strconv.parsefloat(numstrnor(recordinfo[2]), 64)
			if err != nil {
				log.witherror(err).error("failed to parse record info amount. raw data: %s.", recordstr)
			}
			parsedrecord.amount = amount

			total, err := strconv.parsefloat(numstrnor(recordinfo[3]), 64)
			if err != nil {
				log.witherror(err).error("failed to parse record info total. raw data: %s.", recordstr)
			}
			parsedrecord.total = total
		}

		dt.dtrecords = append(dt.dtrecords, parsedrecord)
		log.debugf("record parse successfully. %+v", parsedrecord)
	}

	return nil
}

func (dt *dtmock) parsesummary(raw string) error {

	summarypreparse := strings.split(raw, "|")
	if len(summarypreparse) != 8 {
		log.witherror(errorunexpectedrawdata).errorf("raw: %s.", raw)
		return errorunexpectedrawdata
	}

	code := strings.trimspace(summarypreparse[0])
	cnname := strings.trimspace(summarypreparse[1])
	investtype := strings.trimspace(summarypreparse[4])

	totalround, err := strconv.atoi(numstrnor(summarypreparse[2][:len(summarypreparse[2])-3]))
	if err != nil {
		log.witherror(err).errorf("failed to parse totalround, raw data: %s", raw)
		return err
	}

	baseinvestamount, err := strconv.parsefloat(numstrnor(summarypreparse[3]), 64)
	if err != nil {
		log.witherror(err).errorf("failed to parse baseinvestamount, raw data: %s", raw)
		return err
	}

	finalamount, err := strconv.parsefloat(numstrnor(summarypreparse[5]), 64)
	if err != nil {
		log.witherror(err).errorf("failed to parse finalamount, raw data: %s", raw)
		return err
	}

	profitpercent, err := strconv.parsefloat(numstrnor(summarypreparse[6][:len(summarypreparse[6])-1]), 64)
	if err != nil {
		log.witherror(err).errorf("failed to parse profitpercent, raw data: %s", raw)
		return err
	}

	mockcalculatedate, err := time.parse(dateformat, strings.trimspace(summarypreparse[7]))
	if err != nil {
		log.witherror(err).errorf("failed to parse mockcalculatedate, raw data: %s", raw)
		return err
	}

	dt.code = code
	dt.cnname = cnname
	dt.totalround = int64(totalround)
	dt.baseinvestamount = baseinvestamount
	dt.investtype = investtype
	dt.finalamount = finalamount
	dt.profitpercent = profitpercent
	dt.mockcalculatedate = mockcalculatedate

	return nil
}

func numstrnor(giv string) string {
	return strings.replaceall(strings.trim(giv, " \n"), ",", "")
}
