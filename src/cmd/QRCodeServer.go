// GoProject1 project main.go
package main

import (
	"fmt"
	"flag"
	"http"
	"io"
	"log"
	//"template"
	"old/template"
	"url"
	//"url"
)

var addr = flag.String("addr", ":1718", "http service address") // Q=17, R=18
var fmap = template.FormatterMap{
	"html":     template.HTMLFormatter,
	"url+html": UrlHtmlFormatter,
}
var templ = template.MustParse(templateStr, fmap)

func main() {
	flag.Parse()
	fmt.Println("Starting http server...")
	//http.Handle("/", http.HandlerFunc(QR))
	http.Handle("/", http.FileServer(http.Dir("/home/jtdeng")))
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func QR(w http.ResponseWriter, req *http.Request) {
	templ.Execute(w, req.FormValue("s"))
}

func UrlHtmlFormatter(w io.Writer, ft string, v ...interface{}) {
	fmt.Println("%v", v)
	template.HTMLEscape(w, []byte(url.QueryEscape(v[0].(string))))
}

const templateStr = `
<html>
<head>
<title>QR Link Generator</title>
</head>
<body>
{.section @}
<img src="http://chart.apis.google.com/chart?chs=300x300&cht=qr&choe=UTF-8&chl={@|url+html}"
/>
<br>
{@|html}
<br>
<br>
{.end}
<form action="/" name=f method="GET"><input maxLength=1024 size=70
name=s value="" title="Text to QR Encode"><input type=submit
value="Show QR" name=qr>
</form>
</body>
</html>
`
