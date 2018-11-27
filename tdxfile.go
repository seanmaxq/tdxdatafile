package tdxdatafile

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	PATH_SEPARATOR = GetPathSeparator()
)

func GetPathSeparator() string {
	if runtime.GOOS == "windows" {
		return "\\"
	}
	return "/"
}

func ProcessFileName(pathfilename, sep string) (stocktype, stockcode string) {
	stocktype = ""
	stockcode = "000000"

	ss1 := strings.Split(pathfilename, sep)
	filename := ss1[len(ss1)-1]

	ss2 := strings.Split(filename, ".")

	ss3 := strings.Split(ss2[0], "#")
	stocktype = ss3[0]
	stockcode = ss3[1]
	return
}

func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

type TDXFileProcessControlBlock struct {
	path                string
	recordChan          chan StockDataRaw
	funcRecordProcessor FuncRecordProcessor
	doneChan            chan int // no more records
	quitChan            chan int // receiver is closed
}

type FuncRecordProcessor func(record StockDataRaw)

func (pTDXFilePCB *TDXFileProcessControlBlock) Init(path string, funcRecordProcessor FuncRecordProcessor) {
	// build channels
	pTDXFilePCB.recordChan = make(chan StockDataRaw)
	pTDXFilePCB.doneChan = make(chan int)
	pTDXFilePCB.quitChan = make(chan int)

	pTDXFilePCB.path = path
	pTDXFilePCB.funcRecordProcessor = funcRecordProcessor
}

func (pTDXFilePCB *TDXFileProcessControlBlock) SetRecordProcessor(funcRecordProcessor FuncRecordProcessor) {
	pTDXFilePCB.funcRecordProcessor = funcRecordProcessor
}

func (pTDXFilePCB *TDXFileProcessControlBlock) Receiver() {
	defer func() {
		pTDXFilePCB.quitChan <- 0
	}()

	nRecords := 0

	for {
		select {
		case <-pTDXFilePCB.doneChan:
			fmt.Printf("[%s]Receiver() ended. Received %d.\n", pTDXFilePCB.path, nRecords)
			return

		case strOneRecord := <-pTDXFilePCB.recordChan:
			nRecords++
			if nil != pTDXFilePCB.funcRecordProcessor {
				pTDXFilePCB.funcRecordProcessor(strOneRecord)
			}
		}
	}
}

func (pTDXFilePCB *TDXFileProcessControlBlock) Waiting() {
	<-pTDXFilePCB.quitChan
}

func (pTDXFilePCB *TDXFileProcessControlBlock) Traverse() error {
	defer func() {
		pTDXFilePCB.doneChan <- 0
	}()

	if nil == pTDXFilePCB {
		return errors.New("nil parameter!")
	}
	count := 0

	processTDXFile := func(pathfilename string) error {
		stockType, stockCode := ProcessFileName(pathfilename, PATH_SEPARATOR)
		var fd *os.File
		var err1 error
		if CheckFileIsExist(pathfilename) { //如果文件存在
			fd, err1 = os.OpenFile(pathfilename, os.O_RDONLY, os.ModePerm) //打开文件
			if nil != err1 {
				return err1
			}
			buf := bufio.NewReader(fd)
			defer fd.Close()

			// GBK
			//decoder := mahonia.NewDecoder("gbk")
			//if decoder == nil {
			//	fmt.Println("编码不存在!")
			//}
			count := 0
			wrLnCnt := 0
			for {
				line, err33 := buf.ReadString('\n')
				if err33 != nil {
					if err33 == io.EOF {
						break
					}
					fmt.Println(err33)
					break
				}
				count++
				line = strings.TrimSpace(line)
				ss := strings.Split(line, ",")

				if len(ss) == 7 {
					_, err31 := time.Parse("2006/01/02", ss[0])
					if nil != err31 {
						fmt.Println(err31)
						continue
					}

					rawData := &StockDataRaw{}
					rawData.StockType = stockType
					rawData.StockCode = stockCode
					rawData.StockDate = ss[0]
					rawData.Open = ss[1]
					rawData.High = ss[2]
					rawData.Low = ss[3]
					rawData.End = ss[4]
					rawData.Volumn = ss[5]
					rawData.Amount = ss[6]

					pTDXFilePCB.recordChan <- *rawData
					wrLnCnt++
				}
			}
			//fmt.Printf("count:%d, wrLnCnt:%d\n", count, wrLnCnt)
		}

		return nil
	}

	listFunc := func(pathfile string, f os.FileInfo, err error) error {
		if f == nil {
			fmt.Println(err)
			return err
		}
		if f.IsDir() {
			fmt.Println("Directory:" + pathfile)
		} else {
			processTDXFile(pathfile)
		}
		count++

		return nil
	}

	//var strRet string
	err := filepath.Walk(pTDXFilePCB.path, listFunc) //

	if err != nil {
		fmt.Printf("filepath.Walk() returned %v\n", err)
	} else {
		fmt.Println(count)
	}

	return err
}
