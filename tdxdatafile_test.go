package tdxdatafile

import (
	"fmt"
	"testing"
)

func TestAll(t *testing.T) {
	var oneRecord StockDataRaw
	oneRecord.StockType = "sz"
	fmt.Println(oneRecord.ToString())

}

func recordProc(record StockDataRaw) {
	fmt.Println(record.ToString())
}

func TestTDXFile(t *testing.T) {
	fmt.Println("StockSwiss")

	var cb TDXFileProcessControlBlock

	cb.Init("D:\\mxq\\stock\\data\\Index", recordProc)

	go cb.Receiver()

	cb.Traverse()

	cb.Waiting()
}
