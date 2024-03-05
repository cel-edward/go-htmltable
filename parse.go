// htmltable enables structured data extraction from HTML tables and URLs
package htmltable

import (
	"io"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// mock for tests
var htmlParse = html.Parse

type Parser struct {
	Tables     []*Table
	currentRow row
	rows       []row
	maxCols    int
}

// Table contains the 2D slice of string data parsed from html.
//
// Each string value is stripped of whitespace.
type Table [][]string

// cell is an internal structure for use in parsing, representing a <td> unit
type cell struct {
	Value   string
	RowSpan int
	ColSpan int
}

// row is an internal structure for use in parsing, representing a slice of cells
type row []cell

// New returns an instance of the page with possibly more than one table
func New(r io.Reader) ([]*Table, error) {
	parser := Parser{
		Tables: nil,
	}
	err := parser.parse(r)
	if err != nil {
		return nil, err
	}
	return parser.Tables, nil
}

// NewFromString is same as New(ctx.Context, io.Reader), but from string
func NewFromString(r string) ([]*Table, error) {
	return New(strings.NewReader(r))
}

func (p *Parser) parse(r io.Reader) error {
	root, err := htmlParse(r)
	if err != nil {
		return err
	}
	p.traverse(root)
	p.finishTable()
	return nil
}

// traverse recursively walks the node and its children, handling table node elements
func (p *Parser) traverse(n *html.Node) {
	if n == nil {
		return
	}
	switch n.Data {
	case "td", "th":
		rowspan, colspan, isDisplayNone := getAttributes(n)
		if isDisplayNone {
			return
		}
		var sb strings.Builder
		getInnerText(n, &sb)
		cell := cell{
			Value:   strings.TrimSpace(sb.String()),
			ColSpan: colspan,
			RowSpan: rowspan,
		}
		p.currentRow = append(p.currentRow, cell)
		return
	case "tr":
		p.finishRow()
	case "table":
		p.finishTable()
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		p.traverse(c)
	}
}

// getAttributes returns attributes for node n that are relevant for parsing,
// namely rowspan, colspan, and display:none
//
// If not found, defaults returned are row/colspan = 1 and isDisplayNone = false
func getAttributes(n *html.Node) (rowspan int, colspan int, isDisplayNone bool) {
	colspan = 1
	rowspan = 1
	isDisplayNone = false
	displayNoneRegexp := regexp.MustCompile(`display:\s*none`)
	for _, a := range n.Attr {
		key := strings.ToLower(a.Key)
		if key == "colspan" {
			val, err := strconv.Atoi(a.Val)
			if err == nil {
				colspan = val
			}
		} else if key == "rowspan" {
			val, err := strconv.Atoi(a.Val)
			if err == nil {
				rowspan = val
			}
		} else if key == "style" {
			if displayNoneRegexp.Match([]byte(a.Val)) {
				isDisplayNone = true
			}
		}
	}
	return rowspan, colspan, isDisplayNone
}

// finishRow handles the end of a <tr> block in the html, shifting the data into the parser's rows buffer
func (p *Parser) finishRow() {
	if len(p.currentRow) == 0 {
		return
	}
	if len(p.currentRow) > p.maxCols {
		p.maxCols = len(p.currentRow)
	}
	p.rows = append(p.rows, p.currentRow)
	p.currentRow = row{}
}

// finishTable handles the end of a <table> block in the html.
// The string representation is calculated (handling row/colspans)
// and the data is appended to parser.Tables
func (p *Parser) finishTable() {

	p.finishRow()
	if len(p.rows) == 0 {
		return
	}

	tableData := [][]string{}

	// carryover handles row spans > 1, by keeping track of cells that need to be handled
	// in subsequent rows during the main row loop
	type carryover struct {
		Value   string // value of cell
		Index   int    // column index of the carryover
		RowSpan int    // how many spans left still in the carryover
	}
	var rowCarryover []carryover

	for _, row := range p.rows {
		var rowData []string
		nextRowCarryover := []carryover{}
		currentIndex := 0

		for _, cell := range row {
			// if there is carryover from prior rows, make sure they're accounted for
			// e.g. if index 3 had a carryover this needs to come the cell that otherwise is the fourth element in the row
			for len(rowCarryover) > 0 && rowCarryover[0].Index <= currentIndex {
				// pop out the first item
				co := rowCarryover[0]
				if len(rowCarryover) > 1 {
					rowCarryover = rowCarryover[1:]
				} else {
					rowCarryover = nil
				}

				rowData = append(rowData, co.Value)
				// add to the next row if there is still additional rowspan
				if co.RowSpan > 1 {
					nextRowCarryover = append(nextRowCarryover, carryover{
						Value:   co.Value,
						RowSpan: co.RowSpan - 1,
						Index:   co.Index,
					})
				}
			}

			// now we can start copying values from the current cell
			for i := 0; i < cell.ColSpan; i++ {
				rowData = append(rowData, cell.Value)
				// account for rowspan into subsequent rows
				if cell.RowSpan > 1 {
					co := carryover{
						Value:   cell.Value,
						RowSpan: cell.RowSpan - 1,
						Index:   currentIndex,
					}
					nextRowCarryover = append(nextRowCarryover, co)
				}
				currentIndex++
			}
		}

		// this is for any columns that only exist at the bottom due to rowspan (i.e. no <td> standalone)
		for _, co := range rowCarryover {
			rowData = append(rowData, co.Value)
			// add to the next row if there is still additional rowspan
			if co.RowSpan > 1 {
				nextRowCarryover = append(nextRowCarryover, carryover{
					Value:   co.Value,
					RowSpan: co.RowSpan - 1,
					Index:   co.Index,
				})
			}
		}

		tableData = append(tableData, rowData)
		rowCarryover = nextRowCarryover
	}
	newTable := Table(tableData)
	p.Tables = append(p.Tables, &newTable)

	p.maxCols = 0
	p.currentRow = row{}
	p.rows = nil
}

// getInnerText retrieves any text from child nodes of n and adds it to sb.
// Texts from different nodes will have a whitespace inserted between.
func getInnerText(n *html.Node, sb *strings.Builder) {
	if n.Type == html.TextNode {
		sb.WriteString(strings.TrimSpace(n.Data))
		sb.WriteString(" ")
		return
	}
	if n.FirstChild == nil {
		return
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		getInnerText(c, sb)
	}
}

// type cellSpan struct {
// 	BeginX, EndX int
// 	BeginY, EndY int
// 	Value        string
// }

// func (d *cellSpan) Match(x, y int) bool {
// 	if d.BeginX > x {
// 		return false
// 	}
// 	if d.EndX <= x {
// 		return false
// 	}
// 	if d.BeginY > y {
// 		return false
// 	}
// 	if d.EndY <= y {
// 		return false
// 	}
// 	return true
// }

// type spans []cellSpan

// func (s spans) Value(x, y int) (string, bool) {
// 	for _, v := range s {
// 		if !v.Match(x, y) {
// 			continue
// 		}
// 		return v.Value, true
// 	}
// 	return "", false
// }
