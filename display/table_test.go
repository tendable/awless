package display

import (
	"bytes"
	"reflect"
	"testing"
)

func TestTableDisplay(t *testing.T) {
	table := NewTable([]*PropertyDisplayer{{Property: "c1"}, {Property: "C2"}, {Property: "c3"}})
	table.AddRow("v1.1", "v1.2", "v1.3")
	table.AddRow("v2.1", "v2.2", "v2.3", "v2.4")
	table.AddRow("v3.1", "v3.2")
	table.AddValue("c1", "v4.1")
	table.AddValue("C1", "v5.1")
	table.AddValue("c2", "v4.2")
	table.AddValue("c3", "v4.3")
	var print bytes.Buffer
	table.Fprint(&print)
	expected := `+------+------+------+
|  C1  |  C2  |  C3  |
+------+------+------+
| v1.1 | v1.2 | v1.3 |
| v2.1 | v2.2 | v2.3 |
| v3.1 | v3.2 |      |
| v4.1 | v4.2 | v4.3 |
| v5.1 |      |      |
+------+------+------+
`
	if got, want := print.String(), expected; got != want {
		t.Fatalf("got\n%s\nwant\n%s\n", got, want)
	}

	print.Reset()
	table.FprintColumns(&print, "c1", "C3")
	expected = `+------+------+
|  C1  |  C3  |
+------+------+
| v1.1 | v1.3 |
| v2.1 | v2.3 |
| v3.1 |      |
| v4.1 | v4.3 |
| v5.1 |      |
+------+------+
`
	if got, want := print.String(), expected; got != want {
		t.Fatalf("got\n%s\nwant\n%s\n", got, want)
	}

	print.Reset()
	table.FprintColumns(&print, "C1", "c4")
	expected = `+------+----+
|  C1  | C4 |
+------+----+
| v1.1 |    |
| v2.1 |    |
| v3.1 |    |
| v4.1 |    |
| v5.1 |    |
+------+----+
`
	if got, want := print.String(), expected; got != want {
		t.Fatalf("got\n%s\nwant\n%s\n", got, want)
	}

	if got, want := table.ColumnValues("c2"), []string{"v1.2", "v2.2", "v3.2", "v4.2"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("got\n%s\nwant\n%s\n", got, want)
	}

	print.Reset()
	table.FprintColumnValues(&print, "C3", " ")
	if got, want := print.String(), "v1.3 v2.3  v4.3\n"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestTableSpecialDisplays(t *testing.T) {
	table := NewTable([]*PropertyDisplayer{{Property: "c1", CollapseIdenticalValues: true}, {Property: "c2"}, {Property: "c3"}})
	table.AddRow("v1.1", "v1.2", "v1.3")
	table.AddRow("v1.1", "v2.2", "v2.3")
	table.AddRow("v1.1", "v3.2")
	table.AddRow("v4.1", "v4.2", "v4.3")
	table.AddRow("v4.1", "v4.2", "v5.3")
	var print bytes.Buffer
	table.Fprint(&print)
	expected := `+------+------+------+
|  C1  |  C2  |  C3  |
+------+------+------+
| v1.1 | v1.2 | v1.3 |
| //   | v2.2 | v2.3 |
| //   | v3.2 |      |
| v4.1 | v4.2 | v4.3 |
| //   | v4.2 | v5.3 |
+------+------+------+
`
	if got, want := print.String(), expected; got != want {
		t.Fatalf("got\n%s\nwant\n%s\n", got, want)
	}
	print.Reset()
	table.FprintColumnValues(&print, "c1", " ")
	if got, want := print.String(), "v1.1 v4.1\n"; got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestTableRanking(t *testing.T) {
	table := NewTable([]*PropertyDisplayer{{Property: " C1"}, {Property: "c2"}, {Property: "c3"}})
	table.AddRow("v1.1", "a1.2", "v1.3")
	table.AddRow("v5.1")
	table.AddRow("v3.1", "d3.2")
	table.AddRow("v4.1", "b4.2", "v2.3")
	table.AddRow("v2.1", "c2.2", "v1.3")
	var print bytes.Buffer
	table.Fprint(&print)
	expected := `+------+------+------+
|  C1  |  C2  |  C3  |
+------+------+------+
| v1.1 | a1.2 | v1.3 |
| v5.1 |      |      |
| v3.1 | d3.2 |      |
| v4.1 | b4.2 | v2.3 |
| v2.1 | c2.2 | v1.3 |
+------+------+------+
`
	if got, want := print.String(), expected; got != want {
		t.Fatalf("got\n%s\nwant\n%s\n", got, want)
	}

	print.Reset()
	table.SetSortBy("c1")
	table.Fprint(&print)
	expected = `+------+------+------+
| C1 ▲ |  C2  |  C3  |
+------+------+------+
| v1.1 | a1.2 | v1.3 |
| v2.1 | c2.2 | v1.3 |
| v3.1 | d3.2 |      |
| v4.1 | b4.2 | v2.3 |
| v5.1 |      |      |
+------+------+------+
`
	if got, want := print.String(), expected; got != want {
		t.Fatalf("got\n%s\nwant\n%s\n", got, want)
	}

	print.Reset()
	table.SetSortBy(" C2 ")
	table.Fprint(&print)
	expected = `+------+------+------+
|  C1  | C2 ▲ |  C3  |
+------+------+------+
| v5.1 |      |      |
| v1.1 | a1.2 | v1.3 |
| v4.1 | b4.2 | v2.3 |
| v2.1 | c2.2 | v1.3 |
| v3.1 | d3.2 |      |
+------+------+------+
`
	if got, want := print.String(), expected; got != want {
		t.Fatalf("got\n%s\nwant\n%s\n", got, want)
	}

	print.Reset()
	table.SetSortBy("c4", "c2") //c4 column does not exist
	table.Fprint(&print)
	if got, want := print.String(), expected; got != want {
		t.Fatalf("got\n%s\nwant\n%s\n", got, want)
	}

	print.Reset()
	table.SetSortBy("C3", "c1")
	table.Fprint(&print)
	expected = `+------+------+------+
|  C1  |  C2  | C3 ▲ |
+------+------+------+
| v3.1 | d3.2 |      |
| v5.1 |      |      |
| v1.1 | a1.2 | v1.3 |
| v2.1 | c2.2 | v1.3 |
| v4.1 | b4.2 | v2.3 |
+------+------+------+
`
	if got, want := print.String(), expected; got != want {
		t.Fatalf("got\n%s\nwant\n%s\n", got, want)
	}

}