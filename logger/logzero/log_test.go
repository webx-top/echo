package logzero_test

import (
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
	logzero.Writer().Write([]byte(`test log message`))
}
