package echo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccept(t *testing.T) {
	a := NewAccepts(`application/vnd.example.v2+json, application/xhtml+xml, text/javascript, */*; q=0.01`).Advance()
	expected := &Accepts{
		Raw: "application/vnd.example.v2+json, application/xhtml+xml, text/javascript, */*; q=0.01",
		Type: []*Accept{
			{
				Raw:  "application/vnd.example.v2+json",
				Type: "application",
				Subtype: []string{
					"json",
				},
				Mime: "application/json",
				Vendor: []string{
					"example",
					"v2",
				},
			},
			{
				Raw:  "application/xhtml+xml",
				Type: "application",
				Subtype: []string{
					"xhtml",
					"xml",
				},
				Mime: "application/xml",
			},
			{
				Raw:  "text/javascript",
				Type: "text",
				Subtype: []string{
					"javascript",
				},
				Mime: "text/javascript",
			},
			{
				Raw:  "*/*",
				Type: "*",
				Subtype: []string{
					"*",
				},
				Mime: "*/*",
			},
		},
	}
	assert.Equal(t, Dump(expected, false), Dump(a, false))
	a = NewAccepts(`application/vnd.example.v2+json`).Advance()
	expected = &Accepts{
		Raw: "application/vnd.example.v2+json",
		Type: []*Accept{
			{
				Raw:  "application/vnd.example.v2+json",
				Type: "application",
				Subtype: []string{
					"json",
				},
				Mime: "application/json",
				Vendor: []string{
					"example",
					"v2",
				},
			},
		},
	}
	assert.Equal(t, Dump(expected, false), Dump(a, false))
}
