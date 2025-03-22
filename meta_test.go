package echo_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
	. "github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware"
)

type MetaRequest struct {
	Name string `valid:"required"`
}

func NewMetaRequest() echo.MetaValidator {
	return &MetaRequest{}
}

func (m *MetaRequest) Methods() []string {
	return nil
}

func (m *MetaRequest) Filters(c echo.Context) []FormDataFilter {
	return nil
}

func (m *MetaRequest) ValueDecoders(c Context) BinderValueCustomDecoders {
	return nil
}

func TestEchoMeta(t *testing.T) {
	e := New()
	e.SetDebug(true)
	e.RouteDebug = true
	g := e.Group("/root").SetMetaKV(`parent`, `root`)
	routeName := `test.echo.meta`
	metaData := H{"version": 1.0, "data": H{"by": "handler"}}
	r := g.Get("/", e.MetaHandler(
		metaData,
		func(c Context) error {
			return c.JSON(c.Route().Meta)
		},
	))
	r.SetName(routeName)

	g2 := g.Group("/sub").SetMetaKV(`child`, `sub`)
	g2.Get("/", func(c Context) error {
		return c.JSON(c.Route().Meta)
	})

	e.RebuildRouter()

	assert.Equal(t, routeName, r.GetName())

	expectedMeta := H{
		`parent`: `root`, // group meta
	}
	expectedMeta.DeepMerge(metaData)
	assert.Equal(t, expectedMeta, r.GetMeta())

	var meta H

	for _, route := range e.Routes() {
		if route.Path == "/root/" {
			meta = route.Meta
		}
	}
	expected := H{
		"version": 1.0,
		"parent":  "root", // group meta
		"data": H{
			"by": "handler",
		},
	}
	assert.Equal(t, expected, meta)

	c, b := request(GET, "/root/", e)
	assert.Equal(t, http.StatusOK, c)
	expected2, _ := json.MarshalIndent(expected, "", "  ")
	assert.Equal(t, string(expected2), b)
	assert.Equal(t, `/root/`, e.URI(`test.echo.meta`))

	c, b = request(GET, "/root/sub/", e)
	assert.Equal(t, http.StatusOK, c)
	expected = H{
		"child":  "sub",  // group meta
		"parent": "root", // group meta
	}
	expected2, _ = json.MarshalIndent(expected, "", "  ")
	assert.Equal(t, string(expected2), b)
}

func TestEchoMetaRequestValidator(t *testing.T) {
	e := New()
	e.SetDebug(true)
	e.SetValidator(NewValidation())
	g := e.Group("/root")

	g.Post("/post", e.MakeHandler(
		func(c Context) error {
			data := GetValidated(c).(*MetaRequest)
			return c.String(data.Name)
		},
		NewMetaRequest,
	))
	e.RebuildRouter()

	c, b := request(POST, "/root/post", e, func(r *http.Request) {
		r.Form = url.Values{}
		r.Form.Add(`Name`, `OK`)
		r.Header.Set("Content-Type", echo.MIMEApplicationForm)
		r.Body = io.NopCloser(bytes.NewReader([]byte(r.Form.Encode())))
	})
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, `OK`, b)

	c, b = request(POST, "/root/post", e, func(r *http.Request) {
		r.Form = url.Values{}
		r.Header.Set("Content-Type", echo.MIMEApplicationForm)
		r.Body = io.NopCloser(bytes.NewReader([]byte(r.Form.Encode())))
	})
	assert.Equal(t, http.StatusInternalServerError, c)
	assert.Equal(t, `Name: Can not be empty`, b)

}

type testRequestData struct {
	Name string `valid:"required"`
}

type testResponseData struct {
	Name string `json:"name,omitempty"`
}

func TestEchoMetaRequestValidator2(t *testing.T) {
	e := New()
	e.Use(middleware.Validate(NewValidation))
	e.SetDebug(true)
	g := e.Group("/root")

	g.Post("/post2", e.MetaHandler(
		nil,
		func(c Context) error {
			data := GetValidated(c).(*testRequestData)
			return c.String(data.Name)
		},
		func() interface{} {
			return &testRequestData{}
		},
	))
	e.RebuildRouter()

	c, b := request(POST, "/root/post2", e, func(r *http.Request) {
		r.Form = url.Values{}
		r.Form.Add(`Name`, `OK`)
		r.Header.Set("Content-Type", echo.MIMEApplicationForm)
		r.Body = io.NopCloser(bytes.NewReader([]byte(r.Form.Encode())))
	})
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, `OK`, b)

	c, b = request(POST, "/root/post2", e, func(r *http.Request) {
		r.Form = url.Values{}
		r.Header.Set("Content-Type", echo.MIMEApplicationForm)
		r.Body = io.NopCloser(bytes.NewReader([]byte(r.Form.Encode())))
	})
	assert.Equal(t, http.StatusInternalServerError, c)
	assert.Equal(t, `Name: Can not be empty`, b)

}

func TestEchoMetaRequestValidator3(t *testing.T) {
	e := New()
	e.Use(middleware.Validate(NewValidation))
	e.SetDebug(true)
	g := e.Group("/root")

	g.Post("/post3", e.MetaHandler(
		nil,
		func(c Context) error {
			data := GetValidated(c).(*testRequestData)
			return c.String(data.Name)
		},
		&testRequestData{},
	))

	e.RebuildRouter()

	c, b := request(POST, "/root/post3", e, func(r *http.Request) {
		r.Form = url.Values{}
		r.Form.Add(`Name`, `OK`)
		r.Header.Set("Content-Type", echo.MIMEApplicationForm)
		r.Body = io.NopCloser(bytes.NewReader([]byte(r.Form.Encode())))
	})
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, `OK`, b)

	c, b = request(POST, "/root/post3", e, func(r *http.Request) {
		r.Form = url.Values{}
		r.Header.Set("Content-Type", echo.MIMEApplicationForm)
		r.Body = io.NopCloser(bytes.NewReader([]byte(r.Form.Encode())))
	})
	assert.Equal(t, http.StatusInternalServerError, c)
	assert.Equal(t, `Name: Can not be empty`, b)
}

func TestEchoMetaRequestValidatorX(t *testing.T) {
	e := New()
	e.Use(middleware.Validate(NewValidation))
	e.SetDebug(true)
	g := e.Group("/root")
	h := HandlerFuncWithArg[testRequestData, testResponseData](func(c Context, data *testRequestData) (testResponseData, error) {
		c.SetAuto(true)
		if data == nil {
			return testResponseData{}, ErrBadRequest
		}
		return testResponseData{Name: data.Name}, nil
	})
	reqOK := func(r *http.Request) {
		r.Form = url.Values{}
		r.Form.Add(`Name`, `OK`)
		r.Header.Set("Content-Type", echo.MIMEApplicationForm)
		r.Header.Set("Accept", echo.MIMEApplicationJSON)
		r.Body = io.NopCloser(bytes.NewReader([]byte(r.Form.Encode())))
	}
	reqInvalid := func(r *http.Request) {
		r.Form = url.Values{}
		r.Header.Set("Content-Type", echo.MIMEApplicationForm)
		r.Body = io.NopCloser(bytes.NewReader([]byte(r.Form.Encode())))
	}
	g.Post("/post", e.MakeHandler(h, &testRequestData{}))
	g.Route("GET,POST", "/post-invalid-arg-type", e.MetaHandler(
		nil,
		HandlerFuncWithArg[string, testResponseData](func(c Context, data *string) (testResponseData, error) {
			c.SetAuto(true)
			if data == nil {
				return testResponseData{}, ErrBadRequest
			}
			return testResponseData{Name: *data}, nil
		}),
		"POST",
	))
	g.Route("GET,POST", "/post-x", e.MakeHandler(h, "POST"))

	e.RebuildRouter()

	c, b := request(POST, "/root/post", e, reqInvalid)
	assert.Equal(t, http.StatusInternalServerError, c)
	assert.Equal(t, `Name: Can not be empty`, b)

	expected := `{
  "Code": 1,
  "State": "Success",
  "Info": null,
  "Data": {
    "name": "OK"
  }
}`
	c, b = request(POST, "/root/post", e, reqOK)
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, expected, b)

	c, b = request(GET, "/root/post-invalid-arg-type", e, reqOK)
	assert.Equal(t, http.StatusBadRequest, c)
	assert.Equal(t, `Bad Request`, b)

	c, b = request(POST, "/root/post-invalid-arg-type", e, reqOK)
	assert.Equal(t, http.StatusBadRequest, c)
	assert.Equal(t, `Bad Request`, b)

	c, b = request(GET, "/root/post-x", e, reqOK)
	assert.Equal(t, http.StatusBadRequest, c)
	assert.Equal(t, `Bad Request`, b)

	c, b = request(POST, "/root/post-x", e, reqOK)
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, expected, b)
}
