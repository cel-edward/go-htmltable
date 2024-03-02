// htmltable enables structured data extraction from HTML tables and URLs
package htmltable

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// mock for tests
var htmlParse = html.Parse

// Page is the container for all tables parseable
type Page struct {
	Tables []*Table

	ctx      context.Context
	rowSpans []int
	colSpans []int
	row      []string
	rows     [][]string
	maxCols  int

	// current row
	colSpan []int
	rowSpan []int
	// all
	cSpans [][]int
	rSpans [][]int
}

// New returns an instance of the page with possibly more than one table
func New(ctx context.Context, r io.Reader) (*Page, error) {
	p := &Page{ctx: ctx}
	return p, p.init(r)
}

// NewFromString is same as New(ctx.Context, io.Reader), but from string
func NewFromString(r string) (*Page, error) {
	return New(context.Background(), strings.NewReader(r))
}

// NewFromResponse is same as New(ctx.Context, io.Reader), but from http.Response.
//
// In case of failure, returns `ResponseError`, that could be further inspected.
func NewFromResponse(resp *http.Response) (*Page, error) {
	p, err := New(resp.Request.Context(), resp.Body)
	if err != nil {
		return nil, err
	}
	return p, nil
}

// Len returns number of tables found on the page
func (p *Page) Len() int {
	return len(p.Tables)
}

func (p *Page) init(r io.Reader) error {
	root, err := htmlParse(r)
	if err != nil {
		return err
	}
	p.parse(root)
	p.finishTable()
	return nil
}

func (p *Page) parse(n *html.Node) {
	if n == nil {
		return
	}
	switch n.Data {
	case "td", "th":
		p.colSpan = append(p.colSpan, p.intAttrOr(n, "colspan", 1))
		p.rowSpan = append(p.rowSpan, p.intAttrOr(n, "rowspan", 1))
		var sb strings.Builder
		p.innerText(n, &sb)
		p.row = append(p.row, strings.TrimSpace(sb.String()))
		return
	case "tr":
		p.finishRow()
	case "table":
		p.finishTable()
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		p.parse(c)
	}
}

func (p *Page) intAttrOr(n *html.Node, attr string, default_ int) int {
	for _, a := range n.Attr {
		if a.Key != attr {
			continue
		}
		val, err := strconv.Atoi(a.Val)
		if err != nil {
			return default_
		}
		return val
	}
	return default_
}

func (p *Page) finishRow() {
	if len(p.row) == 0 {
		return
	}
	if len(p.row) > p.maxCols {
		p.maxCols = len(p.row)
	}
	p.rows = append(p.rows, p.row)
	p.cSpans = append(p.cSpans, p.colSpan)
	p.rSpans = append(p.rSpans, p.rowSpan)
	p.row = []string{}
	p.colSpan = []int{}
	p.rowSpan = []int{}
}

type cellSpan struct {
	BeginX, EndX int
	BeginY, EndY int
	Value        string
}

func (d *cellSpan) Match(x, y int) bool {
	if d.BeginX > x {
		return false
	}
	if d.EndX <= x {
		return false
	}
	if d.BeginY > y {
		return false
	}
	if d.EndY <= y {
		return false
	}
	return true
}

type spans []cellSpan

func (s spans) Value(x, y int) (string, bool) {
	for _, v := range s {
		if !v.Match(x, y) {
			continue
		}
		return v.Value, true
	}
	return "", false
}

func (p *Page) finishTable() {
	defer func() {
		if r := recover(); r != nil {
			firstRow := []string{}
			if len(p.rows) > 0 {
				firstRow = p.rows[0][:]
			}
			Logger(p.ctx, "unparsable table", "panic", fmt.Sprintf("%v", r), "firstRow", firstRow)
		}
		p.rows = [][]string{}
		p.colSpans = []int{}
		p.rowSpans = []int{}
		p.cSpans = [][]int{}
		p.rSpans = [][]int{}
		p.maxCols = 0
	}()
	p.finishRow()
	if len(p.rows) == 0 {
		return
	}

	rows := [][]string{}
	allSpans := spans{}
	rowSkips := 0

	for y := 0; y < len(p.rows); y++ { // rows cols addressable by x
		currentRow := []string{}
		skipRow := false
		j := 0 // p.rows cols addressable by j
		for x := 0; x < p.maxCols; x++ {
			value, ok := allSpans.Value(x, y)
			if ok {
				currentRow = append(currentRow, value)
				continue
			}
			if len(p.rSpans[y]) == j {
				break
			}
			rowSpan := p.rSpans[y][j]
			colSpan := p.cSpans[y][j]
			value = p.rows[y][j]
			if rowSpan > 1 || colSpan > 1 {
				allSpans = append(allSpans, cellSpan{
					BeginX: x,
					EndX:   x + colSpan,
					BeginY: y,
					EndY:   y + rowSpan,
					Value:  value,
				})
			}
			currentRow = append(currentRow, value)

			j++
		}
		if skipRow {
			rowSkips++
			y++
		}
		if len(currentRow) > p.maxCols {
			p.maxCols = len(currentRow)
		}
		rows = append(rows, currentRow)
	}
	Logger(p.ctx, "found table", "count", len(rows))
	p.Tables = append(p.Tables, &Table{
		Rows: rows,
	})
}

// CEL edit: this also adds a blank space between node texts
func (p *Page) innerText(n *html.Node, sb *strings.Builder) {
	if n.Type == html.TextNode {
		sb.WriteString(strings.TrimSpace(n.Data))
		sb.WriteString(" ")
		return
	}
	if n.FirstChild == nil {
		return
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		p.innerText(c, sb)
	}
}

// Table is the low-level representation of rows.
//
// Every cell string value is truncated of its whitespace.
type Table struct {
	// Rows holds slice of string slices
	Rows [][]string
}

func (table *Table) String() string {
	return fmt.Sprintf("Table (%d rows)", len(table.Rows))
}
