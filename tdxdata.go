package tdxdatafile

import (
	"encoding/json"
	"errors"
	"time"
)

//
type StockData struct {
	StockType string
	StockCode string
	StockDate time.Time
	Open      float32
	High      float32
	Low       float32
	End       float32
	Volumn    int64
	Amount    float32
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

func (psdr *StockDataRaw) ToStockData() (*StockData, error) {
	if nil == psdr {
		return nil, errors.New("nil")
	}

	data := StockData{}

	data.StockType = psdr.StockType
	data.StockCode = psdr.StockCode
	stockdate, err1 := time.Parse("2006/01/02", psdr.StockDate)
	if nil == err1 {
		data.StockDate = stockdate
	}

	var f float32

	return &data, nil
}
