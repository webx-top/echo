package subdomains

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
)

func TestSortHosts(t *testing.T) {
	a := New()
	e := echo.New()
	a.Add(`frontend`, e)
	e2 := echo.New()
	e2.SetPrefix(`/admin`)
	a.Add(`backend`, e2)
	a.sortHosts()
	assert.Equal(t, []string{`backend`, `frontend`}, a.Hosts[``])
}
