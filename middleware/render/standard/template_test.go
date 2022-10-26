package standard

import (
	"bytes"
	"fmt"
	"html/template"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	tmpl := template.New(`test`)
	var now string
	funcMap := map[string]interface{}{
		`now`: func() template.HTML {
			now = time.Now().Format(time.RFC3339Nano)
			fmt.Println(now)
			return template.HTML(now)
		},
	}
	tmpl.Funcs(funcMap)
	_, err := tmpl.Parse(`{{now}}`)
	assert.NoError(t, err)
	for i := 0; i < 10; i++ {
		go func() {
			buf := bytes.NewBuffer(nil)
			err = tmpl.Execute(buf, nil)
			assert.NoError(t, err)
			//assert.Equal(t, now, buf.String())
			time.Sleep(time.Second * time.Duration(rand.Intn(5)+1))
		}()
	}
	time.Sleep(time.Second * 5)
}
