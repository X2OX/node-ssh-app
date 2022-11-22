package nsa

import (
	"bytes"
	"fmt"
	"strings"

	"go.x2ox.com/THz"
	"go.x2ox.com/THz/render"
	"go.x2ox.com/sorbifolia/random"
)

func Router() *THz.THz {
	t := THz.New()

	api := t.Group("/api/v1")

	api.GET("/nsa/:os/:arch/register", func(c *THz.Context) {
		tunnel := NewTunnel(fmt.Sprintf("nsa-%s-%s-%s",
			c.Param(":os"), c.Param(":arch"),
			rd.RandString(5),
		))

		if tunnel == nil {
			c.Status(500)
			return
		}
		c.Status(200).Text(tunnel.Token)
	})

	t.GET("/", listAction)
	t.GET("/index.html", listAction)

	return t
}

func listAction(c *THz.Context) {
	tunnels := GetTunnel()
	arr := make([]string, 0, len(tunnels))

	for _, v := range tunnels {
		if !strings.HasPrefix(v.Name, "nsa-") {
			continue
		}

		arr = append(arr, v.Name)
	}

	c.Render(render.Data{
		ContentType: "text/html; charset=utf-8",
		Data:        GenHTML(arr),
	})
}

func GenHTML(arr []string) []byte {
	buf := new(bytes.Buffer)
	buf.WriteString(tplHead)

	for _, v := range arr {
		buf.WriteString("<li><a href=\"")
		buf.WriteString("https://")
		buf.WriteString(v)
		buf.WriteString(".")
		buf.WriteString(domain)
		buf.WriteString("\">")
		buf.WriteString(v)
		buf.WriteString("</a></li>\n")
	}
	buf.WriteString(tplFoot)
	return buf.Bytes()
}

var rd = random.Fast().SetRandBytes([]byte("abcdefghijklmnopqrstuvwxyz"))

const (
	tplHead = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport"
          content="width=device-width,minimum-scale=1,initial-scale=1,maximum-scale=5,viewport-fit=cover">
    <meta name="renderer" content="webkit">
    <title>NSA</title>

    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/github-markdown-css@5.1.0/github-markdown.min.css">
    <style>
        body {
            box-sizing: border-box;
            min-width: 200px;
            max-width: 980px;
            margin: 0 auto;
            padding: 45px;
            background-color: #F6F6F6;
        }

        #content {
            padding: 30px 60px 55px;
        }

        .markdown-body {
            background-color: #FFFFFF;
        }
    </style>
</head>
<body>
<article class="markdown-body"><div id="content">
<h1 id="title">NSA</h1>
<ul>

`
	tplFoot = `
</ul>
</div></article>
</body>
</html>
`
)
