package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/russross/blackfriday/v2"
	"golang.org/x/net/html"
)

func main() {
	var cssFile string
	if len(os.Args) < 2 {
		fmt.Println("Please provide a file name and a optional CSS file")
		os.Exit(1)
	}

	fileName := os.Args[1]
	if len(os.Args) < 3 {
		cssFile = "style1.css"
	} else {
		cssFile = os.Args[2]
	}
	pdfFileName := strings.Replace(fileName, ".md", ".pdf", 1)

	input, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}
	css, err := os.ReadFile(cssFile)
	if err != nil {
		fmt.Println("Error reading CSS file:", err)
		os.Exit(1)
	}

	html := blackfriday.Run(input)
	fmt.Println(string(html))
	styledHtml := formatHtml(string(html), string(css))
	os.WriteFile("resume.html", []byte(styledHtml), 0644)
	pdf := htmlToPdf(styledHtml)
	writePdf(pdf, pdfFileName)

}

// Removes erroneous characters from the HTML string and wraps each heading and the section below it in a div
func formatHtml(htmlStr string, css string) string {
	htmlStr = removeErroneousChars(htmlStr)

	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		os.Exit(1)
	}

	var f func(*html.Node)
	var buffer bytes.Buffer
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "head" {
			style := &html.Node{
				Type: html.ElementNode,
				Data: "style",
			}
			style.AppendChild(&html.Node{
				Type: html.TextNode,
				Data: css,
			})
			n.AppendChild(style)
		}
		if n.Type == html.ElementNode && n.Data == "h2" {
			div := &html.Node{
				Type: html.ElementNode,
				Data: "div",
				Attr: []html.Attribute{
					{Key: "class", Val: strings.ToLower(n.FirstChild.Data)},
				},
			}
			n.Parent.InsertBefore(div, n)
			n.Parent.RemoveChild(n)
			div.AppendChild(n)
			next := n.Parent.NextSibling
			for next != nil {
				if next.Type == html.ElementNode && next.Data == "h2" {
					f(next)
					break
				}
				nextSibling := next.NextSibling
				next.Parent.RemoveChild(next)
				div.AppendChild(next)
				next = nextSibling
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	html.Render(&buffer, doc)

	return buffer.String()
}

// Converts the HTML string to a PDF
func htmlToPdf(html string) *wkhtmltopdf.PDFGenerator {
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		fmt.Println("Error creating generator:", err)
		os.Exit(1)
	}

	pdfg.MarginBottom.Set(0) // Set bottom margin to 0
	pdfg.MarginTop.Set(0)    // Set top margin to 0
	pdfg.MarginRight.Set(0)  // Set right margin to 0
	pdfg.MarginLeft.Set(0)   // Set left margin to 0
	page := wkhtmltopdf.NewPageReader(strings.NewReader(html))
	page.NoBackground.Set(false)

	pdfg.AddPage(page)

	err = pdfg.Create()
	if err != nil {
		fmt.Println("Error generating PDF:", err)
		os.Exit(1)
	}
	return pdfg
}

// Writes the PDF to a file
// Writes to directory of the input file
// Not necessarily the place the script is run from
func writePdf(pdf *wkhtmltopdf.PDFGenerator, fileName string) {
	err := pdf.WriteFile(fileName)
	if err != nil {
		fmt.Println("Error writing PDF:", err)
		os.Exit(1)
	}
}

// Removes non-breaking spaces and other erroneous characters from the string
func removeErroneousChars(str string) string {
	str = strings.ReplaceAll(str, "&rsquo;", "'") // when converted to html some single quotes become entities some didn't.
	return strings.Map(func(r rune) rune {
		// fmt.Printf("Rune: %c , %d, %q \n", r, r, r)

		if r >= 32 && r != 127 {
			if r == 160 {
				return 32
			}
			if r == '“' || r == '”' {
				return '"'
			}
			if r == '’' || r == '\'' { // this ain't a repeat of the above, this is a different character
				return '\''
			}
			if r == 8211 {
				return '-'
			}
			return r
		}
		return -1
	}, str)
}
