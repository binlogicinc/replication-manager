// replication-manager - Replication Manager Monitoring and CLI for MariaDB and MySQL
// Copyright 2017 Signal 18 SARL
// Authors: Guillaume Lefranc <guillaume@signal18.io>
//          Stephane Varoqui  <svaroqui@gmail.com>
// This source code is licensed under the GNU General Public License, version 3.

package ogórek

import (
	"bytes"
	"reflect"
	"testing"
)

func TestEncode(t *testing.T) {

	type foo struct {
		Foo string
		Bar int32
	}

	tests := []struct {
		name   string
		input  interface{}
		output interface{}
	}{
		{
			"graphite message",
			[]interface{}{map[interface{}]interface{}{"values": []interface{}{float64(473), float64(497), float64(540), float64(1497), float64(1808), float64(1890), float64(2013), float64(1821), float64(1847), float64(2176), float64(2156), float64(1250), float64(2055), float64(1570), None{}, None{}}, "start": int64(1383782400), "step": int64(86400), "end": int64(1385164800), "name": "ZZZZ.UUUUUUUU.CCCCCCCC.MMMMMMMM.XXXXXXXXX.TTT"}},
			nil,
		},
		{
			"small types",
			[]interface{}{int64(0), int64(1), int64(258), int64(65537), false, true},
			nil,
		},
		{
			"array of struct types",
			[]foo{{"Qux", 4}},
			[]interface{}{map[interface{}]interface{}{"Foo": "Qux", "Bar": int64(4)}},
		},
	}

	for _, tt := range tests {
		p := &bytes.Buffer{}
		e := NewEncoder(p)
		e.Encode(tt.input)

		d := NewDecoder(bytes.NewReader(p.Bytes()))
		output, _ := d.Decode()

		want := tt.output
		if want == nil {
			want = tt.input
		}

		if !reflect.DeepEqual(want, output) {
			t.Errorf("%s: got\n%q\n expected\n%q", tt.name, output, want)
		}

	}
}

type testMarshalPickle struct {
}

func (p *testMarshalPickle) MarshalPickle() (text []byte, err error) {
	return []byte("(lp0\nI1\naI2\na"), nil
}
func TestMarshalPickle(t *testing.T) {
	testData := map[string]interface{}{
		"key": &testMarshalPickle{},
	}

	p := &bytes.Buffer{}
	e := NewEncoder(p)
	e.Encode(testData)

	if p.String() != "}(U\x03key(lp0\nI1\naI2\nau." {
		t.FailNow()
	}
}
