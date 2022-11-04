package echo_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
	. "github.com/webx-top/echo"
	"github.com/webx-top/echo/encoding/json"
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

func TestEchoMeta(t *testing.T) {
	e := New()
	e.SetDebug(true)
	g := e.Group("/root")

	g.Get("/", e.MetaHandler(
		H{"version": 1.0, "data": H{"by": "handler"}},
		func(c Context) error {
			return c.JSON(c.Route().Meta)
		},
	))

	e.RebuildRouter()

	var meta H

	for _, route := range e.Routes() {
		if route.Path == "/root/" {
			meta = route.Meta
		}
	}
	expected := H{
		"version": 1.0,
		"data": H{
			"by": "handler",
		},
	}
	assert.Equal(t, expected, meta)

	c, b := request(GET, "/root/", e)
	assert.Equal(t, http.StatusOK, c)
	expected2, _ := json.MarshalIndent(expected, "", "  ")
	assert.Equal(t, string(expected2), b)

}

func TestEchoMetaRequestValidator(t *testing.T) {
	e := New()
	e.SetDebug(true)
	e.SetValidator(NewValidation())
	g := e.Group("/root")

	g.Post("/post", e.MetaHandler(
		nil,
		func(c Context) error {
			data := c.Internal().Get(`validated`).(*MetaRequest)
			return c.String(data.Name)
		},
		NewMetaRequest,
	))
	e.RebuildRouter()

	c, b := request(POST, "/root/post", e, func(r *http.Request) {
		r.Form = url.Values{}
		r.Form.Add(`Name`, `OK`)
		r.Header.Set("Content-Type", echo.MIMEMultipartForm)
		r.Body = ioutil.NopCloser(bytes.NewReader([]byte(r.Form.Encode())))
	})
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, `OK`, b)

	c, b = request(POST, "/root/post", e, func(r *http.Request) {
		r.Form = url.Values{}
		r.Header.Set("Content-Type", echo.MIMEMultipartForm)
		r.Body = ioutil.NopCloser(bytes.NewReader([]byte(r.Form.Encode())))
	})
	assert.Equal(t, http.StatusInternalServerError, c)
	assert.Equal(t, `Name: Can not be empty`, b)

}

type testRequestData struct {
	Name string `valid:"required"`
}

func TestEchoMetaRequestValidator2(t *testing.T) {
	e := New()
	e.Use(middleware.Validate(NewValidation))
	e.SetDebug(true)
	g := e.Group("/root")

	g.Post("/post2", e.MetaHandler(
		nil,
		func(c Context) error {
			data := c.Internal().Get(`validated`).(*testRequestData)
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
		r.Header.Set("Content-Type", echo.MIMEMultipartForm)
		r.Body = ioutil.NopCloser(bytes.NewReader([]byte(r.Form.Encode())))
	})
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, `OK`, b)

	c, b = request(POST, "/root/post2", e, func(r *http.Request) {
		r.Form = url.Values{}
		r.Header.Set("Content-Type", echo.MIMEMultipartForm)
		r.Body = ioutil.NopCloser(bytes.NewReader([]byte(r.Form.Encode())))
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
			data := c.Internal().Get(`validated`).(*testRequestData)
			return c.String(data.Name)
		},
		&testRequestData{},
	))

	e.RebuildRouter()

	c, b := request(POST, "/root/post3", e, func(r *http.Request) {
		r.Form = url.Values{}
		r.Form.Add(`Name`, `OK`)
		r.Header.Set("Content-Type", echo.MIMEMultipartForm)
		r.Body = ioutil.NopCloser(bytes.NewReader([]byte(r.Form.Encode())))
	})
	assert.Equal(t, http.StatusOK, c)
	assert.Equal(t, `OK`, b)

	c, b = request(POST, "/root/post3", e, func(r *http.Request) {
		r.Form = url.Values{}
		r.Header.Set("Content-Type", echo.MIMEMultipartForm)
		r.Body = ioutil.NopCloser(bytes.NewReader([]byte(r.Form.Encode())))
	})
	assert.Equal(t, http.StatusInternalServerError, c)
	assert.Equal(t, `Name: Can not be empty`, b)
}
