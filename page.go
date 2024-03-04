// htmltable enables structured data extraction from HTML tables and URLs
package htmltable

import (
	"io"
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
type Table struct {
	Data [][]string
}

// cell is an internal structure for use in parsing, representing a <td> unit
type cell struct {
	Value   string
	RowSpan int
	ColSpan int
}

// row is an internal structure for use in parsing, representing a slice of cells
type row struct {
	Cells []cell
}

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
		var sb strings.Builder
		getInnerText(n, &sb)
		cell := cell{
			Value:   strings.TrimSpace(sb.String()),
			ColSpan: intAttrOr(n, "colspan", 1),
			RowSpan: intAttrOr(n, "rowspan", 1),
		}
		p.currentRow.Cells = append(p.currentRow.Cells, cell)
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

// intAttrOr returns the integer value of attribute attr for node n,
// returning defaultValue if integer parsing fails or attr is not found
func intAttrOr(n *html.Node, attr string, defaultValue int) int {
	for _, a := range n.Attr {
		if a.Key != attr {
			continue
		}
		val, err := strconv.Atoi(a.Val)
		if err != nil {
			return defaultValue
		}
		return val
	}
	return defaultValue
}

// finishRow handles the end of a <tr> block in the html, shifting the data into the parser's rows buffer
func (p *Parser) finishRow() {
	if len(p.currentRow.Cells) == 0 {
		return
	}
	if len(p.currentRow.Cells) > p.maxCols {
		p.maxCols = len(p.currentRow.Cells)
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

		for _, cell := range row.Cells {
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
	p.Tables = append(p.Tables, &Table{Data: tableData})

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
