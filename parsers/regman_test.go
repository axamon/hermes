// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parsers_test

import (
	"context"
	"testing"

	"github.com/axamon/hermes/parsers"
)

func TestREGMAN(t *testing.T) {
	type args struct {
		ctx     context.Context
		logfile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Primo", args{ctx: context.TODO(), logfile: "testngasp.csv.gz"}, false},
		{"Fileinesistente", args{ctx: context.TODO(), logfile: "testngasp1.csv.gz"}, true}, // TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := parsers.REGMAN(tt.args.ctx, tt.args.logfile); (err != nil) != tt.wantErr {
				t.Errorf("REGMAN() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func BenchmarkREGMAN(b *testing.B) {
	for n := 0; n < b.N; n++ {
		parsers.REGMAN(context.TODO(), "testngasp.csv.gz")
	}
}
