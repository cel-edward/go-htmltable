package htmltable

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

const fixture = `<body>
<h1>foo</h2>
<table>
	<tr><td>a</td><td>b</td></tr>
	<tr><td> 1 </td><td>2</td></tr>
	<tr><td>3  </td><td>4   </td></tr>
</table>
<h1>bar</h2>
<table>
	<tr><th>b</th><th>c</th><th>d</th></tr>
	<tr><td>1</td><td>2</td><td>5</td></tr>
	<tr><td>3</td><td>4</td><td>6</td></tr>
</table>
</body>`

func TestFindsAllTables(t *testing.T) {
	ts, err := NewFromString(fixture)
	assertNoError(t, err)
	assertEqual(t, len(ts), 2)
}

// added public domain data from https://en.wikipedia.org/wiki/List_of_S&P_500_companies
const fixtureColspans = `<table>
	<thead>
		<tr>
			<th rowspan="2">Date</th>
			<th colspan="2">Added</th>
			<th colspan="2">Removed</th>
			<th rowspan="2">Reason</th>
		</tr>
		<tr>
			<th rowspan="@#$%^&">Ticker</th>
			<th>Security</th>
			<th>Ticker</th>
			<th>Security</th>
		</tr>
	</thead>
	<tbody>
		<tr>
			<td>June 21, 2022</td>
			<td>KDP</td>
			<td><a href="/wiki/Keurig_Dr_Pepper" title="Keurig Dr Pepper">Keurig Dr Pepper</a></td>
			<td>UA/UAA</td>
			<td><a href="/wiki/Under_Armour" title="Under Armour">Under Armour</a></td>
			<td>Market capitalization change.<sup id="cite_ref-sp20220603_4-0" class="reference"><a href="#cite_note-sp20220603-4">[4]</a></sup></td>
		</tr>
		<tr>
			<td>June 21, 2022</td>
			<td>ON</td>
			<td><a href="/wiki/ON_Semiconductor" class="mw-redirect" title="ON Semiconductor">ON Semiconductor</a></td>
			<td>IPGP</td>
			<td><a href="/wiki/IPG_Photonics" title="IPG Photonics">IPG Photonics</a></td>
			<td>Market capitalization change.<sup id="cite_ref-sp20220603_4-1" class="reference"><a href="#cite_note-sp20220603-4">[4]</a></sup></td>
		</tr>
	</tbody>
</table>`

func TestFindsWithColspans(t *testing.T) {
	ts, err := NewFromString(fixtureColspans)
	assertNoError(t, err)
	assertEqual(t, len(ts), 1)
	assertEqual(t, "Market capitalization change. [4]", ts[0].Data[2][5])
}

func TestInitFails(t *testing.T) {
	prev := htmlParse
	t.Cleanup(func() {
		htmlParse = prev
	})
	htmlParse = func(r io.Reader) (*html.Node, error) {
		return nil, fmt.Errorf("nope")
	}
	_, err := New(strings.NewReader(".."))

	assertEqualError(t, err, "nope")
}
