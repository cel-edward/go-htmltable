package htmltable

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestSimple(t *testing.T) {
	ts, err := NewFromString(testTable1)
	assertNoError(t, err)
	assertEqual(t, len(ts), 2)
}

func TestColspans(t *testing.T) {
	ts, err := NewFromString(testTable2)
	assertNoError(t, err)
	assertEqual(t, len(ts), 1)
	assertEqual(t, "Market capitalization change. [4]", (*ts[0])[2][5])
}

func TestComplex(t *testing.T) {
	ts, err := NewFromString(testTable3)
	assertNoError(t, err)
	assertEqual(t, len(ts), 1)
	for _, ss := range *ts[0] {
		assertEqual(t, len(ss), 24)
	}
	assertEqual(t, *ts[0], Table(testWant3))
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

const testTable1 = `<body>
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

const testTable2 = `<table>
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

const testTable3 = `
<table
  style="
    border-collapse: collapse;
    display: inline-table;
    margin-bottom: 5pt;
    vertical-align: text-bottom;
    width: 100%;
  "
>
  <tbody>
    <tr>
      <td style="width: 1%"></td>
      <td style="width: 40.566%"></td>
      <td style="width: 0.1%"></td>
      <td style="width: 1%"></td>
      <td style="width: 13.002%"></td>
      <td style="width: 0.1%"></td>
      <td style="width: 0.1%"></td>
      <td style="width: 0.441%"></td>
      <td style="width: 0.1%"></td>
      <td style="width: 1%"></td>
      <td style="width: 13.002%"></td>
      <td style="width: 0.1%"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td style="width: 0.1%"></td>
      <td style="width: 0.441%"></td>
      <td style="width: 0.1%"></td>
      <td style="width: 1%"></td>
      <td style="width: 13.002%"></td>
      <td style="width: 0.1%"></td>
      <td style="width: 0.1%"></td>
      <td style="width: 0.441%"></td>
      <td style="width: 0.1%"></td>
      <td style="width: 1%"></td>
      <td style="width: 13.005%"></td>
      <td style="width: 0.1%"></td>
    </tr>
    <tr>
      <td colspan="3" style="padding: 0px 1pt"></td>
      <td
        colspan="9"
        style="padding: 2px 1pt; text-align: center; vertical-align: middle"
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 700;
            line-height: 100%;
          "
          >Three Months Ended September 30,</span
        >
      </td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="padding: 0px 1pt"></td>
      <td
        colspan="9"
        style="padding: 2px 1pt; text-align: center; vertical-align: bottom"
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 700;
            line-height: 100%;
          "
          >Nine Months Ended September 30,</span
        >
      </td>
    </tr>
    <tr>
      <td colspan="3" style="padding: 0px 1pt"></td>
      <td
        colspan="3"
        style="
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 1pt;
          text-align: center;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 700;
            line-height: 100%;
          "
          >2022</span
        >
      </td>
      <td
        colspan="3"
        style="border-top: 1pt solid rgb(0, 0, 0); padding: 0px 1pt"
      ></td>
      <td
        colspan="3"
        style="
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 1pt;
          text-align: center;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 700;
            line-height: 100%;
          "
          >2023</span
        >
      </td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="padding: 0px 1pt"></td>
      <td
        colspan="3"
        style="
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 1pt;
          text-align: center;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 700;
            line-height: 100%;
          "
          >2022</span
        >
      </td>
      <td
        colspan="3"
        style="border-top: 1pt solid rgb(0, 0, 0); padding: 0px 1pt"
      ></td>
      <td
        colspan="3"
        style="
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 1pt;
          text-align: center;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 700;
            line-height: 100%;
          "
          >2023</span
        >
      </td>
    </tr>
    <tr>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
    </tr>
    <tr>
      <td
        colspan="3"
        style="
          background-color: rgb(204, 238, 255);
          padding: 2px 1pt;
          text-align: left;
          vertical-align: bottom;
        "
      >
        <div>
          <span
            style="
              color: rgb(0, 0, 0);
              font-family: 'Times New Roman', sans-serif;
              font-size: 10pt;
              font-weight: 400;
              line-height: 100%;
            "
            >Interest income </span
          ><span
            style="
              color: rgb(0, 0, 0);
              font-family: 'Times New Roman', sans-serif;
              font-size: 6.5pt;
              font-weight: 400;
              line-height: 100%;
              position: relative;
              top: -3.5pt;
              vertical-align: baseline;
            "
            >(1)</span
          >
        </div>
      </td>
      <td
        style="
          background-color: rgb(204, 238, 255);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 0px 2px 1pt;
          text-align: left;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          >$</span
        >
      </td>
      <td
        style="
          background-color: rgb(204, 238, 255);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          ><ix:nonfraction
            unitref="usd"
            contextref="c-7"
            decimals="-3"
            name="us-gaap:InterestIncomeOperating"
            format="ixt:num-dot-decimal"
            scale="3"
            id="f-448"
            >22,180</ix:nonfraction
          >&nbsp;</span
        >
      </td>
      <td
        style="
          background-color: rgb(204, 238, 255);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 1pt 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      ></td>
      <td
        colspan="3"
        style="background-color: rgb(204, 238, 255); padding: 0px 1pt"
      ></td>
      <td
        style="
          background-color: rgb(204, 238, 255);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 0px 2px 1pt;
          text-align: left;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          >$</span
        >
      </td>
      <td
        style="
          background-color: rgb(204, 238, 255);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          ><ix:nonfraction
            unitref="usd"
            contextref="c-8"
            decimals="-3"
            name="us-gaap:InterestIncomeOperating"
            format="ixt:num-dot-decimal"
            scale="3"
            id="f-449"
            >37,692</ix:nonfraction
          >&nbsp;</span
        >
      </td>
      <td
        style="
          background-color: rgb(204, 238, 255);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 1pt 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      ></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td
        colspan="3"
        style="background-color: rgb(204, 238, 255); padding: 0px 1pt"
      ></td>
      <td
        style="
          background-color: rgb(204, 238, 255);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 0px 2px 1pt;
          text-align: left;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          >$</span
        >
      </td>
      <td
        style="
          background-color: rgb(204, 238, 255);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          ><ix:nonfraction
            unitref="usd"
            contextref="c-9"
            decimals="-3"
            name="us-gaap:InterestIncomeOperating"
            format="ixt:num-dot-decimal"
            scale="3"
            id="f-450"
            >66,288</ix:nonfraction
          >&nbsp;</span
        >
      </td>
      <td
        style="
          background-color: rgb(204, 238, 255);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 1pt 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      ></td>
      <td
        colspan="3"
        style="background-color: rgb(204, 238, 255); padding: 0px 1pt"
      ></td>
      <td
        style="
          background-color: rgb(204, 238, 255);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 0px 2px 1pt;
          text-align: left;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          >$</span
        >
      </td>
      <td
        style="
          background-color: rgb(204, 238, 255);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          ><ix:nonfraction
            unitref="usd"
            contextref="c-1"
            decimals="-3"
            name="us-gaap:InterestIncomeOperating"
            format="ixt:num-dot-decimal"
            scale="3"
            id="f-451"
            >116,923</ix:nonfraction
          >&nbsp;</span
        >
      </td>
      <td
        style="
          background-color: rgb(204, 238, 255);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 1pt 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      ></td>
    </tr>
    <tr>
      <td
        colspan="3"
        style="
          background-color: rgb(255, 255, 255);
          padding: 2px 1pt;
          text-align: left;
          vertical-align: bottom;
        "
      >
        <div>
          <span
            style="
              color: rgb(0, 0, 0);
              font-family: 'Times New Roman', sans-serif;
              font-size: 10pt;
              font-weight: 400;
              line-height: 100%;
            "
            >Interest expense </span
          ><span
            style="
              color: rgb(0, 0, 0);
              font-family: 'Times New Roman', sans-serif;
              font-size: 6.5pt;
              font-weight: 400;
              line-height: 100%;
              position: relative;
              top: -3.5pt;
              vertical-align: baseline;
            "
            >(1)</span
          >
        </div>
      </td>
      <td
        colspan="2"
        style="
          background-color: rgb(255, 255, 255);
          padding: 2px 0px 2px 1pt;
          text-align: right;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          >(<ix:nonfraction
            unitref="usd"
            contextref="c-7"
            decimals="-3"
            name="us-gaap:InterestExpense"
            format="ixt:num-dot-decimal"
            scale="3"
            id="f-452"
            >3,050</ix:nonfraction
          >)</span
        >
      </td>
      <td
        style="
          background-color: rgb(255, 255, 255);
          padding: 2px 1pt 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      ></td>
      <td
        colspan="3"
        style="background-color: rgb(255, 255, 255); padding: 0px 1pt"
      ></td>
      <td
        colspan="2"
        style="
          background-color: rgb(255, 255, 255);
          padding: 2px 0px 2px 1pt;
          text-align: right;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          >(<ix:nonfraction
            unitref="usd"
            contextref="c-8"
            decimals="-3"
            name="us-gaap:InterestExpense"
            format="ixt:num-dot-decimal"
            scale="3"
            id="f-453"
            >9,414</ix:nonfraction
          >)</span
        >
      </td>
      <td
        style="
          background-color: rgb(255, 255, 255);
          padding: 2px 1pt 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      ></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td
        colspan="3"
        style="background-color: rgb(255, 255, 255); padding: 0px 1pt"
      ></td>
      <td
        colspan="2"
        style="
          background-color: rgb(255, 255, 255);
          padding: 2px 0px 2px 1pt;
          text-align: right;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          >(<ix:nonfraction
            unitref="usd"
            contextref="c-9"
            decimals="-3"
            name="us-gaap:InterestExpense"
            format="ixt:num-dot-decimal"
            scale="3"
            id="f-454"
            >6,322</ix:nonfraction
          >)</span
        >
      </td>
      <td
        style="
          background-color: rgb(255, 255, 255);
          padding: 2px 1pt 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      ></td>
      <td
        colspan="3"
        style="background-color: rgb(255, 255, 255); padding: 0px 1pt"
      ></td>
      <td
        colspan="2"
        style="
          background-color: rgb(255, 255, 255);
          padding: 2px 0px 2px 1pt;
          text-align: right;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          >(<ix:nonfraction
            unitref="usd"
            contextref="c-1"
            decimals="-3"
            name="us-gaap:InterestExpense"
            format="ixt:num-dot-decimal"
            scale="3"
            id="f-455"
            >20,828</ix:nonfraction
          >)</span
        >
      </td>
      <td
        style="
          background-color: rgb(255, 255, 255);
          padding: 2px 1pt 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      ></td>
    </tr>
    <tr>
      <td
        colspan="3"
        style="
          background-color: rgb(204, 238, 255);
          padding: 2px 1pt;
          text-align: left;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 115%;
          "
          >Fair value and other adjustments</span
        >
      </td>
      <td
        colspan="3"
        style="background-color: rgb(204, 238, 255); padding: 0px 1pt"
      ></td>
      <td
        colspan="3"
        style="background-color: rgb(204, 238, 255); padding: 0px 1pt"
      ></td>
      <td
        colspan="3"
        style="background-color: rgb(204, 238, 255); padding: 0px 1pt"
      ></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td
        colspan="3"
        style="background-color: rgb(204, 238, 255); padding: 0px 1pt"
      ></td>
      <td
        colspan="3"
        style="background-color: rgb(204, 238, 255); padding: 0px 1pt"
      ></td>
      <td
        colspan="3"
        style="background-color: rgb(204, 238, 255); padding: 0px 1pt"
      ></td>
      <td
        colspan="3"
        style="background-color: rgb(204, 238, 255); padding: 0px 1pt"
      ></td>
    </tr>
    <tr>
      <td
        colspan="3"
        style="
          background-color: rgb(255, 255, 255);
          padding: 2px 1pt 2px 13pt;
          text-align: left;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          >Unrealized loss, charge-offs, and other adjustments, net</span
        >
      </td>
      <td
        colspan="2"
        style="
          background-color: rgb(255, 255, 255);
          padding: 2px 0px 2px 1pt;
          text-align: right;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          >(<ix:nonfraction
            unitref="usd"
            contextref="c-7"
            decimals="-3"
            sign="-"
            name="upst:UnrealizedGainLossChargeOffsAndOtherFairValueAdjustmentsNet"
            format="ixt:num-dot-decimal"
            scale="3"
            id="f-456"
            >20,069</ix:nonfraction
          >)</span
        >
      </td>
      <td
        style="
          background-color: rgb(255, 255, 255);
          padding: 2px 1pt 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      ></td>
      <td
        colspan="3"
        style="background-color: rgb(255, 255, 255); padding: 0px 1pt"
      ></td>
      <td
        colspan="2"
        style="
          background-color: rgb(255, 255, 255);
          padding: 2px 0px 2px 1pt;
          text-align: right;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          >(<ix:nonfraction
            unitref="usd"
            contextref="c-8"
            decimals="-3"
            sign="-"
            name="upst:UnrealizedGainLossChargeOffsAndOtherFairValueAdjustmentsNet"
            format="ixt:num-dot-decimal"
            scale="3"
            id="f-457"
            >37,521</ix:nonfraction
          >)</span
        >
      </td>
      <td
        style="
          background-color: rgb(255, 255, 255);
          padding: 2px 1pt 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      ></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td
        colspan="3"
        style="background-color: rgb(255, 255, 255); padding: 0px 1pt"
      ></td>
      <td
        colspan="2"
        style="
          background-color: rgb(255, 255, 255);
          padding: 2px 0px 2px 1pt;
          text-align: right;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          >(<ix:nonfraction
            unitref="usd"
            contextref="c-9"
            decimals="-3"
            sign="-"
            name="upst:UnrealizedGainLossChargeOffsAndOtherFairValueAdjustmentsNet"
            format="ixt:num-dot-decimal"
            scale="3"
            id="f-458"
            >70,855</ix:nonfraction
          >)</span
        >
      </td>
      <td
        style="
          background-color: rgb(255, 255, 255);
          padding: 2px 1pt 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      ></td>
      <td
        colspan="3"
        style="background-color: rgb(255, 255, 255); padding: 0px 1pt"
      ></td>
      <td
        colspan="2"
        style="
          background-color: rgb(255, 255, 255);
          padding: 2px 0px 2px 1pt;
          text-align: right;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          >(<ix:nonfraction
            unitref="usd"
            contextref="c-1"
            decimals="-3"
            sign="-"
            name="upst:UnrealizedGainLossChargeOffsAndOtherFairValueAdjustmentsNet"
            format="ixt:num-dot-decimal"
            scale="3"
            id="f-459"
            >108,175</ix:nonfraction
          >)</span
        >
      </td>
      <td
        style="
          background-color: rgb(255, 255, 255);
          padding: 2px 1pt 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      ></td>
    </tr>
    <tr>
      <td
        colspan="3"
        style="
          background-color: rgb(204, 238, 255);
          padding: 2px 1pt 2px 13pt;
          text-align: left;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          >Realized loss on sale of loans, net</span
        >
      </td>
      <td
        colspan="2"
        style="
          background-color: rgb(204, 238, 255);
          padding: 2px 0px 2px 1pt;
          text-align: right;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          >(<ix:nonfraction
            unitref="usd"
            contextref="c-7"
            decimals="-3"
            sign="-"
            name="upst:RealizedGainLossOnTransferOfLoansNet"
            format="ixt:num-dot-decimal"
            scale="3"
            id="f-460"
            >21,176</ix:nonfraction
          >)</span
        >
      </td>
      <td
        style="
          background-color: rgb(204, 238, 255);
          padding: 2px 1pt 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      ></td>
      <td
        colspan="3"
        style="background-color: rgb(204, 238, 255); padding: 0px 1pt"
      ></td>
      <td
        colspan="2"
        style="
          background-color: rgb(204, 238, 255);
          padding: 2px 0px 2px 1pt;
          text-align: right;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          >(<ix:nonfraction
            unitref="usd"
            contextref="c-8"
            decimals="-3"
            sign="-"
            name="upst:RealizedGainLossOnTransferOfLoansNet"
            format="ixt:num-dot-decimal"
            scale="3"
            id="f-461"
            >2,955</ix:nonfraction
          >)</span
        >
      </td>
      <td
        style="
          background-color: rgb(204, 238, 255);
          padding: 2px 1pt 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      ></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td
        colspan="3"
        style="background-color: rgb(204, 238, 255); padding: 0px 1pt"
      ></td>
      <td
        colspan="2"
        style="
          background-color: rgb(204, 238, 255);
          padding: 2px 0px 2px 1pt;
          text-align: right;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          >(<ix:nonfraction
            unitref="usd"
            contextref="c-9"
            decimals="-3"
            sign="-"
            name="upst:RealizedGainLossOnTransferOfLoansNet"
            format="ixt:num-dot-decimal"
            scale="3"
            id="f-462"
            >45,255</ix:nonfraction
          >)</span
        >
      </td>
      <td
        style="
          background-color: rgb(204, 238, 255);
          padding: 2px 1pt 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      ></td>
      <td
        colspan="3"
        style="background-color: rgb(204, 238, 255); padding: 0px 1pt"
      ></td>
      <td
        colspan="2"
        style="
          background-color: rgb(204, 238, 255);
          padding: 2px 0px 2px 1pt;
          text-align: right;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          >(<ix:nonfraction
            unitref="usd"
            contextref="c-1"
            decimals="-3"
            sign="-"
            name="upst:RealizedGainLossOnTransferOfLoansNet"
            format="ixt:num-dot-decimal"
            scale="3"
            id="f-463"
            >22,255</ix:nonfraction
          >)</span
        >
      </td>
      <td
        style="
          background-color: rgb(204, 238, 255);
          padding: 2px 1pt 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      ></td>
    </tr>
    <tr>
      <td
        colspan="3"
        style="
          background-color: rgb(255, 255, 255);
          padding: 2px 1pt;
          text-align: left;
          vertical-align: bottom;
        "
      >
        <div>
          <span
            style="
              color: rgb(0, 0, 0);
              font-family: 'Times New Roman', sans-serif;
              font-size: 10pt;
              font-weight: 400;
              line-height: 100%;
            "
            >Total fair value and other adjustments, net </span
          ><span
            style="
              color: rgb(0, 0, 0);
              font-family: 'Times New Roman', sans-serif;
              font-size: 6.5pt;
              font-weight: 400;
              line-height: 100%;
              position: relative;
              top: -3.5pt;
              vertical-align: baseline;
            "
            >(1)</span
          >
        </div>
      </td>
      <td
        colspan="2"
        style="
          background-color: rgb(255, 255, 255);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 0px 2px 1pt;
          text-align: right;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          >(<ix:nonfraction
            unitref="usd"
            contextref="c-7"
            decimals="-3"
            name="upst:FairValueAndOtherAdjustmentsNet"
            format="ixt:num-dot-decimal"
            scale="3"
            id="f-464"
            >41,245</ix:nonfraction
          >)</span
        >
      </td>
      <td
        style="
          background-color: rgb(255, 255, 255);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 1pt 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      ></td>
      <td
        colspan="3"
        style="background-color: rgb(255, 255, 255); padding: 0px 1pt"
      ></td>
      <td
        colspan="2"
        style="
          background-color: rgb(255, 255, 255);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 0px 2px 1pt;
          text-align: right;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          >(<ix:nonfraction
            unitref="usd"
            contextref="c-8"
            decimals="-3"
            name="upst:FairValueAndOtherAdjustmentsNet"
            format="ixt:num-dot-decimal"
            scale="3"
            id="f-465"
            >40,476</ix:nonfraction
          >)</span
        >
      </td>
      <td
        style="
          background-color: rgb(255, 255, 255);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 1pt 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      ></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td
        colspan="3"
        style="background-color: rgb(255, 255, 255); padding: 0px 1pt"
      ></td>
      <td
        colspan="2"
        style="
          background-color: rgb(255, 255, 255);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 0px 2px 1pt;
          text-align: right;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          >(<ix:nonfraction
            unitref="usd"
            contextref="c-9"
            decimals="-3"
            name="upst:FairValueAndOtherAdjustmentsNet"
            format="ixt:num-dot-decimal"
            scale="3"
            id="f-466"
            >116,110</ix:nonfraction
          >)</span
        >
      </td>
      <td
        style="
          background-color: rgb(255, 255, 255);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 1pt 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      ></td>
      <td
        colspan="3"
        style="background-color: rgb(255, 255, 255); padding: 0px 1pt"
      ></td>
      <td
        colspan="2"
        style="
          background-color: rgb(255, 255, 255);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 0px 2px 1pt;
          text-align: right;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          >(<ix:nonfraction
            unitref="usd"
            contextref="c-1"
            decimals="-3"
            name="upst:FairValueAndOtherAdjustmentsNet"
            format="ixt:num-dot-decimal"
            scale="3"
            id="f-467"
            >130,430</ix:nonfraction
          >)</span
        >
      </td>
      <td
        style="
          background-color: rgb(255, 255, 255);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 1pt 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      ></td>
    </tr>
    <tr style="height: 15pt">
      <td
        colspan="3"
        style="background-color: rgb(204, 238, 255); padding: 0px 1pt"
      ></td>
      <td
        colspan="3"
        style="
          background-color: rgb(204, 238, 255);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 0px 1pt;
        "
      ></td>
      <td
        colspan="3"
        style="background-color: rgb(204, 238, 255); padding: 0px 1pt"
      ></td>
      <td
        colspan="3"
        style="
          background-color: rgb(204, 238, 255);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 0px 1pt;
        "
      ></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td
        colspan="3"
        style="background-color: rgb(204, 238, 255); padding: 0px 1pt"
      ></td>
      <td
        colspan="3"
        style="
          background-color: rgb(204, 238, 255);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 0px 1pt;
        "
      ></td>
      <td
        colspan="3"
        style="background-color: rgb(204, 238, 255); padding: 0px 1pt"
      ></td>
      <td
        colspan="3"
        style="
          background-color: rgb(204, 238, 255);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 0px 1pt;
        "
      ></td>
    </tr>
    <tr>
      <td
        colspan="3"
        style="
          background-color: rgb(255, 255, 255);
          padding: 2px 1pt;
          text-align: left;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          >Total interest income and fair value adjustments, net</span
        >
      </td>
      <td
        style="
          background-color: rgb(255, 255, 255);
          border-bottom: 3pt double rgb(0, 0, 0);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 0px 2px 1pt;
          text-align: left;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          >$</span
        >
      </td>
      <td
        style="
          background-color: rgb(255, 255, 255);
          border-bottom: 3pt double rgb(0, 0, 0);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          >(<ix:nonfraction
            unitref="usd"
            contextref="c-7"
            decimals="-3"
            sign="-"
            name="upst:InterestIncomeAndFairValueAdjustmentNet"
            format="ixt:num-dot-decimal"
            scale="3"
            id="f-468"
            >22,115</ix:nonfraction
          >)</span
        >
      </td>
      <td
        style="
          background-color: rgb(255, 255, 255);
          border-bottom: 3pt double rgb(0, 0, 0);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 1pt 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      ></td>
      <td
        colspan="3"
        style="background-color: rgb(255, 255, 255); padding: 0px 1pt"
      ></td>
      <td
        style="
          background-color: rgb(255, 255, 255);
          border-bottom: 3pt double rgb(0, 0, 0);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 0px 2px 1pt;
          text-align: left;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          >$</span
        >
      </td>
      <td
        style="
          background-color: rgb(255, 255, 255);
          border-bottom: 3pt double rgb(0, 0, 0);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          >(<ix:nonfraction
            unitref="usd"
            contextref="c-8"
            decimals="-3"
            sign="-"
            name="upst:InterestIncomeAndFairValueAdjustmentNet"
            format="ixt:num-dot-decimal"
            scale="3"
            id="f-469"
            >12,198</ix:nonfraction
          >)</span
        >
      </td>
      <td
        style="
          background-color: rgb(255, 255, 255);
          border-bottom: 3pt double rgb(0, 0, 0);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 1pt 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      ></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td colspan="3" style="display: none"></td>
      <td
        colspan="3"
        style="background-color: rgb(255, 255, 255); padding: 0px 1pt"
      ></td>
      <td
        style="
          background-color: rgb(255, 255, 255);
          border-bottom: 3pt double rgb(0, 0, 0);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 0px 2px 1pt;
          text-align: left;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          >$</span
        >
      </td>
      <td
        style="
          background-color: rgb(255, 255, 255);
          border-bottom: 3pt double rgb(0, 0, 0);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          >(<ix:nonfraction
            unitref="usd"
            contextref="c-9"
            decimals="-3"
            sign="-"
            name="upst:InterestIncomeAndFairValueAdjustmentNet"
            format="ixt:num-dot-decimal"
            scale="3"
            id="f-470"
            >56,144</ix:nonfraction
          >)</span
        >
      </td>
      <td
        style="
          background-color: rgb(255, 255, 255);
          border-bottom: 3pt double rgb(0, 0, 0);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 1pt 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      ></td>
      <td
        colspan="3"
        style="background-color: rgb(255, 255, 255); padding: 0px 1pt"
      ></td>
      <td
        style="
          background-color: rgb(255, 255, 255);
          border-bottom: 3pt double rgb(0, 0, 0);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 0px 2px 1pt;
          text-align: left;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          >$</span
        >
      </td>
      <td
        style="
          background-color: rgb(255, 255, 255);
          border-bottom: 3pt double rgb(0, 0, 0);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      >
        <span
          style="
            color: rgb(0, 0, 0);
            font-family: 'Times New Roman', sans-serif;
            font-size: 10pt;
            font-weight: 400;
            line-height: 100%;
          "
          >(<ix:nonfraction
            unitref="usd"
            contextref="c-1"
            decimals="-3"
            sign="-"
            name="upst:InterestIncomeAndFairValueAdjustmentNet"
            format="ixt:num-dot-decimal"
            scale="3"
            id="f-471"
            >34,335</ix:nonfraction
          >)</span
        >
      </td>
      <td
        style="
          background-color: rgb(255, 255, 255);
          border-bottom: 3pt double rgb(0, 0, 0);
          border-top: 1pt solid rgb(0, 0, 0);
          padding: 2px 1pt 2px 0px;
          text-align: right;
          vertical-align: bottom;
        "
      ></td>
    </tr>
  </tbody>
</table>
`

var testWant3 = [][]string{
	{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
	{"", "", "", "Three Months Ended September 30,", "Three Months Ended September 30,", "Three Months Ended September 30,", "Three Months Ended September 30,", "Three Months Ended September 30,", "Three Months Ended September 30,", "Three Months Ended September 30,", "Three Months Ended September 30,", "Three Months Ended September 30,", "", "", "", "Nine Months Ended September 30,", "Nine Months Ended September 30,", "Nine Months Ended September 30,", "Nine Months Ended September 30,", "Nine Months Ended September 30,", "Nine Months Ended September 30,", "Nine Months Ended September 30,", "Nine Months Ended September 30,", "Nine Months Ended September 30,"},
	{"", "", "", "2022", "2022", "2022", "", "", "", "2023", "2023", "2023", "", "", "", "2022", "2022", "2022", "", "", "", "2023", "2023", "2023"},
	{"Interest income (1)", "Interest income (1)", "Interest income (1)", "$", "22,180", "", "", "", "", "$", "37,692", "", "", "", "", "$", "66,288", "", "", "", "", "$", "116,923", ""},
	{"Interest expense (1)", "Interest expense (1)", "Interest expense (1)", "( 3,050 )", "( 3,050 )", "", "", "", "", "( 9,414 )", "( 9,414 )", "", "", "", "", "( 6,322 )", "( 6,322 )", "", "", "", "", "( 20,828 )", "( 20,828 )", ""},
	{"Fair value and other adjustments", "Fair value and other adjustments", "Fair value and other adjustments", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
	{"Unrealized loss, charge-offs, and other adjustments, net", "Unrealized loss, charge-offs, and other adjustments, net", "Unrealized loss, charge-offs, and other adjustments, net", "( 20,069 )", "( 20,069 )", "", "", "", "", "( 37,521 )", "( 37,521 )", "", "", "", "", "( 70,855 )", "( 70,855 )", "", "", "", "", "( 108,175 )", "( 108,175 )", ""},
	{"Realized loss on sale of loans, net", "Realized loss on sale of loans, net", "Realized loss on sale of loans, net", "( 21,176 )", "( 21,176 )", "", "", "", "", "( 2,955 )", "( 2,955 )", "", "", "", "", "( 45,255 )", "( 45,255 )", "", "", "", "", "( 22,255 )", "( 22,255 )", ""},
	{"Total fair value and other adjustments, net (1)", "Total fair value and other adjustments, net (1)", "Total fair value and other adjustments, net (1)", "( 41,245 )", "( 41,245 )", "", "", "", "", "( 40,476 )", "( 40,476 )", "", "", "", "", "( 116,110 )", "( 116,110 )", "", "", "", "", "( 130,430 )", "( 130,430 )", ""},
	{"", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "", ""},
	{"Total interest income and fair value adjustments, net", "Total interest income and fair value adjustments, net", "Total interest income and fair value adjustments, net", "$", "( 22,115 )", "", "", "", "", "$", "( 12,198 )", "", "", "", "", "$", "( 56,144 )", "", "", "", "", "$", "( 34,335 )", ""},
}
