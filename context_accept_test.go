package echo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccept(t *testing.T) {
	actual := `application/vnd.example.v2+json, application/xhtml+xml, text/javascript, */*; q=0.01`
	a := NewAccepts(actual).Advance()
	expected := &Accepts{
		Raw: actual,
		Accepts: []*AcceptQuality{
			{
				Quality: 0.01,
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
			},
		},
	}
	assert.Equal(t, Dump(expected, false), Dump(a, false))
	actual = `application/vnd.example.v2+json`
	a = NewAccepts(actual).Advance()
	expected = &Accepts{
		Raw: actual,
		Accepts: []*AcceptQuality{
			{
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
			},
		},
	}
	assert.Equal(t, Dump(expected, false), Dump(a, false))
}

func TestAccept2(t *testing.T) {
	actual := `text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8`
	a := NewAccepts(actual).Advance()
	expected := &Accepts{
		Raw: actual,
		Accepts: []*AcceptQuality{
			{
				Quality: 0.9,
				Type: []*Accept{
					{
						Raw:  "text/html",
						Type: "text",
						Subtype: []string{
							"html",
						},
						Mime: "text/html",
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
						Raw:  "application/xml",
						Type: "application",
						Subtype: []string{
							"xml",
						},
						Mime: "application/xml",
					},
				},
			},
			{
				Quality: 0.8,
				Type: []*Accept{
					{
						Raw:  "image/webp",
						Type: "image",
						Subtype: []string{
							"webp",
						},
						Mime: "image/webp",
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
			},
		},
	}
	assert.Equal(t, Dump(expected, false), Dump(a, false))
}
