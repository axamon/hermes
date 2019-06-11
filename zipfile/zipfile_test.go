package zipfile_test

import (
	"reflect"
	"testing"

	"github.com/golang/axamon/hermes/zipFile"
)

func ExampleReadAll() {

	zipFile.ReadAll("20190607_03_00_vodabr.cb.log.gz")

}

func TestReadAll(t *testing.T) {
	type args struct {
		zipFile string
	}
	tests := []struct {
		name        string
		args        args
		wantContent []byte
		wantErr     bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotContent, err := ReadAll(tt.args.zipFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotContent, tt.wantContent) {
				t.Errorf("ReadAll() = %v, want %v", gotContent, tt.wantContent)
			}
		})
	}
}
