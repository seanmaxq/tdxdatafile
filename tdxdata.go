package tdxdatafile

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

//
type StockData struct {
	StockType string
	StockCode string
	StockDate time.Time
	Open      float64
	High      float64
	Low       float64
	End       float64
	Volumn    int64
	Amount    float64
}

type StockDataRaw struct {
	StockType string
	StockCode string
	StockDate string
	Open      string
	High      string
	Low       string
	End       string
	Volumn    string
	Amount    string
}

func (psdr *StockDataRaw) ParseRawString(one string) error {
	one = strings.TrimSpace(one)
	ss := strings.Split(one, ",")

	if len(ss) == 7 {
		_, err := time.Parse("2006/01/02", ss[0])
		if nil != err {
			return errors.New("Invalid data(not started with a date)!")
		}
		psdr.StockDate = ss[0]
		psdr.Open = ss[1]
		psdr.High = ss[2]
		psdr.Low = ss[3]
		psdr.End = ss[4]
		psdr.Volumn = ss[5]
		psdr.Amount = ss[6]
	} else {
		return errors.New("Invalid data!")
	}
	return nil
}

func (psdr *StockDataRaw) ToString() (string, error) {
	if nil == psdr {
		return "", errors.New("nil")
	}
	tmp, err := json.Marshal(psdr)
	if nil == err {
		return string(tmp), nil
	}
	return "", err
}

func (psdr *StockDataRaw) ToStockData() (*StockData, string, error) {
	if nil == psdr {
		return nil, "", errors.New("nil")
	}

	strError := ""
	data := StockData{}

	data.StockType = psdr.StockType
	data.StockCode = psdr.StockCode
	stockdate, err := time.Parse("2006/01/02", psdr.StockDate)
	if nil == err {
		data.StockDate = stockdate
	} else {
		strError += err.Error()
		err = nil
	}

	var f float64
	f, err = strconv.ParseFloat(psdr.Open, 64)
	if nil == err {
		data.Open = f
	} else {
		strError += err.Error()
		err = nil
	}
	f, err = strconv.ParseFloat(psdr.High, 64)
	if nil == err {
		data.High = f
	} else {
		strError += err.Error()
		err = nil
	}
	f, err = strconv.ParseFloat(psdr.Low, 64)
	if nil == err {
		data.Low = f
	} else {
		strError += err.Error()
		err = nil
	}
	f, err = strconv.ParseFloat(psdr.End, 64)
	if nil == err {
		data.End = f
	} else {
		strError += err.Error()
		err = nil
	}
	var v int64
	v, err = strconv.ParseInt(psdr.Volumn, 10, 64)
	if nil == err {
		data.Volumn = v
	} else {
		strError += err.Error()
		err = nil
	}
	f, err = strconv.ParseFloat(psdr.Amount, 64)
	if nil == err {
		data.Amount = f
	} else {
		strError += err.Error()
		err = nil
	}

	return &data, strError, nil
}
