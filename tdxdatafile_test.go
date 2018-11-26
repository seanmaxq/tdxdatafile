package tdxdatafile

import (
	"fmt"
	"testing"
)

func TestAll(t *testing.T) {
	var oneRecord StockDataRaw
	oneRecord.StockType = "sz"
	fmt.Println(oneRecord.ToString)

}
