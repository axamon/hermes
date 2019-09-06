// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parsers_test

import (
	"context"
	"testing"

	"github.com/axamon/hermes/parsers"
)

func TestAVS(t *testing.T) {
	type args struct {
		ctx     context.Context
		logfile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Primo", args{ctx: context.TODO(), logfile: "testavs.csv.gz"}, false},
		{"Fileinesistente", args{ctx: context.TODO(), logfile: "testavs1.csv.gz"}, true},
		{"Tgucorta", args{ctx: context.TODO(), logfile: "testavstgusbagliata.csv.gz"}, false},
		{"filemalformato", args{ctx: context.TODO(), logfile: "testavsmalformato.csv.gz"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := parsers.AVS(tt.args.ctx, tt.args.logfile); (err != nil) != tt.wantErr {
				t.Errorf("AVS() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func BenchmarkAVS(b *testing.B) {
	for n := 0; n < b.N; n++ {

		parsers.AVS(context.TODO(), "testavs.csv.gz")
	}
}

func Test_gestisciMailMultiple(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want string
	}{

		{"Solo", args{line: `"fsdfa@fsdfa.it"`}, "fsdfa@fsdfa.it"},
		{"Primo", args{line: `"fsdfa@fsdfa.it | pipo@fsfs.it"`}, "fsdfa@fsdfa.it   pipo@fsfs.it"},
		{"Secondo", args{line: `"fsdfa@fsdfa.it | pipo@fsfs.it | jhjh@jljk.it"`}, "fsdfa@fsdfa.it   pipo@fsfs.it   jhjh@jljk.it"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parsers.GestisciMailMultiple(tt.args.line); got != tt.want {
				t.Errorf("gestisciMailMultiple() = %v, want %v", got, tt.want)
			}
		})
	}
}
