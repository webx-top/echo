package logzero_test

import (
	"log"
	"os"
	"testing"

	"github.com/webx-top/echo/logger/logzero"
)

func TestLog(t *testing.T) {
	logzero.Debug(`test log message`)
	logzero.Error(`test log message`)
	logzero.Warn(`test log message`)
	logzero.Info(`test log message`)
	//logzero.Fatal(`test log message`)
	logzero.GetLogger(`test`).Warn(`test log message`)
	log.SetOutput(logzero.Default())
	log.Println(`test log message`)
	log.SetOutput(logzero.GetLogger(`test`))
	log.Println(`test log message(test)`)
	f, err := os.OpenFile(`./app.log`, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	log.SetOutput(logzero.GetLogger(`test2`, f))
	log.Println(`test log message(test)`)
}
