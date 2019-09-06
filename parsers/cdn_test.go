// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parsers_test

import (
	"github.com/axamon/hermes/parsers"
	"context"
	"testing"
)

func TestCDN(t *testing.T) {
	type args struct {
		ctx            context.Context
		logfile        string
		maxNumRoutines int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Primo", args{ctx: context.TODO(), logfile: "testcdn.csv.gz", maxNumRoutines: 10}, false},
		{"Inesistente", args{ctx: context.TODO(), logfile: "testcdn1.csv.gz", maxNumRoutines: 1000}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := parsers.CDN(tt.args.ctx, tt.args.logfile, tt.args.maxNumRoutines); (err != nil) != tt.wantErr {
				t.Errorf("CDN() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func benchmarkCDN(numofgoroutins int, b *testing.B) {
	for n := 0; n < b.N; n++ {
		parsers.CDN(context.TODO(), "testcdn.csv.gz", numofgoroutins)
	}
}

func BenchmarkCDN10(b *testing.B)  {benchmarkCDN(10, b)}
func BenchmarkCDN100(b *testing.B)  {benchmarkCDN(100, b)}
func BenchmarkCDN1000(b *testing.B)  {benchmarkCDN(1000, b)}
func BenchmarkCDN10000(b *testing.B)  {benchmarkCDN(10000, b)}
