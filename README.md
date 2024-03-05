# HTML table data extractor for Go

`htmltable` enables structured data extraction from HTML tables, requiring almost no external dependencies except `x/net/html`

## Installation

```bash
go get github.com/cel-edward/go-htmltable
```

# Usage
Pass an html string into `New()` or `NewFromString()`. 
`[]*Table` is returned, where Table.Data is of form `[][]string`.

rowspans and colspans are 'demerged', with the contained value copied into each spanned cell.

Cells with attribute `style="[...]display:none[...]"` are ignored.

Example html and results can be found in `parse_test.go`

# Notes

Strings values within returned tables are stripped of surrounding whitespace. 

Whitespace is inserted between multiple divs contained in a `<td>`. For example, if a `<td>` cell has two elements `<div>  text 1</div>` `<div>text2</div>` inside, the resulting text produced is `text 1 text 2`.

# Credits

This is a heavily modified fork of `github.com/nfx/go-htmltable`, designed for use with CEL algorithms. 

The main parsing algorithm has been completely rewritten as did not reliably function for our use cases, particularly with complex row/colspans.
Returned types are also adjusted.
