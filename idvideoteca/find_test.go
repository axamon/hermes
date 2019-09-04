package idvideoteca

import (
	"fmt"
	"strings"
	"testing"
)

func TestFind(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name            string
		args            args
		wantIdvideoteca string
		wantErr         bool
	}{
		{"Sringa vuota", args{s: ""}, "", false},
		{"SingleTitle", args{s: "http://vodabr.cb.ticdn.it/videoteca2/V3/SingleTitle/2019/05/60000243/SS/20089777/20089777_HD.ism/Manifest"}, "60000243", false},
		{"SS", args{s: "http://vodabr.cb.ticdn.it/videoteca2/V3/Film/2014/09/50434361/SS/20086428/20086428_HD.ism/Manifest"}, "50434361", false},
		{"Urlencoded", args{s: "http%3A%2F%2Fvodabr.cb.ticdn.it%2Fvideoteca2%2FV3%2FFilm%2F2017%2F06%2F50670127%2FSS%2F11473278%2F11473278_SD.ism%2FManifest%23https%3A%2F%2Flicense.cubovision.it%2FLicense%2Frightsmanager.asmx"}, "50670127", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIdvideoteca, err := Find(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotIdvideoteca != tt.wantIdvideoteca {
				t.Errorf("Find() = %v, want %v", gotIdvideoteca, tt.wantIdvideoteca)
			}
		})
	}
}

var str = `
http://vodabr.cb.ticdn.it/videoteca2/V3/SingleTitle/2019/05/60000242/SS/20089785/20089785_HD.ism/Manifest
http://vodabr.cb.ticdn.it/videoteca2/V3/SingleTitle/2019/05/60000241/SS/20089779/20089779_HD.ism/Manifest
http://vodabr.cb.ticdn.it/videoteca2/V3/Film/2019/01/50734100/SS/11483192/11483192_HD.ism/Manifest
0000000B;000774571385;2019-06-07 22:57:26.534;0000000B;dwt765ti;cubo;82.55.223.125;;7000;;;;;;;;;SS_QUALITY;;;7000;7000;7000;;;;The Handmaid's Tale;VoD;http://vodabr.cb.ticdn.it/videoteca2/V3/SingleTitle/2019/05/60000243/SS/20089777/20089777_HD.ism/Manifest;systemInput;7.0-4.2.9.1-2018.265;wifi;;2019-06-07 22:56:20;;5cfaebd635172f2ca79a29e1;;10.16.6;2019-06-07 22:57:26.534;TIM;;;;;;;;;;;
0000000B;000774571385;2019-06-07 22:42:46.398;0000000B;dwt765ti;cubo;82.55.223.125;;7000;;;;;;;;;SS_QUALITY;;;7000;7000;7000;;;;The Handmaid's Tale;VoD;http://vodabr.cb.ticdn.it/videoteca2/V3/SingleTitle/2019/05/60000243/SS/20089777/20089777_HD.ism/Manifest;systemInput;7.0-4.2.9.1-2018.265;wifi;;2019-06-07 22:41:20;;5cfae86635172f2ca79a27f7;;10.16.6;2019-06-07 22:42:46.398;TIM;;;;;;;;;;;
0000000B;000774571385;2019-06-07 22:52:26.421;0000000B;dwt765ti;cubo;82.55.223.125;;7000;;;;;;;;;SS_QUALITY;;;7000;7000;7000;;;;The Handmaid's Tale;VoD;http://vodabr.cb.ticdn.it/videoteca2/V3/SingleTitle/2019/05/60000243/SS/20089777/20089777_HD.ism/Manifest;systemInput;7.0-4.2.9.1-2018.265;wifi;;2019-06-07 22:51:20;;5cfaeaaa35172f2ca79a2943;;10.16.6;2019-06-07 22:52:26.421;TIM;;;;;;;;;;;
0000000B;000774571385;2019-06-07 22:27:26.039;0000000B;dwt765ti;cubo;82.55.223.125;;;;;;;;;;;PLAY;;;;;;;;;The Handmaid's Tale;VoD;http://vodabr.cb.ticdn.it/videoteca2/V3/SingleTitle/2019/05/60000243/SS/20089777/20089777_HD.ism/Manifest;systemInput;7.0-4.2.9.1-2018.265;wifi;;2019-06-07 22:26:20;;5cfae4d035172f2ca79a253b;;10.16.6;2019-06-07 22:27:26.039;TIM;;;;;;;;;;;
0000000B;000774571385;2019-06-07 22:25:06.548;0000000B;dwt765ti;cubo;82.55.223.125;;6939;;;;;;;;;SS_QUALITY;;;7000;7000;3000;;;;The Handmaid's Tale;VoD;http://vodabr.cb.ticdn.it/videoteca2/V3/SingleTitle/2019/05/60000242/SS/20089785/20089785_HD.ism/Manifest;systemInput;7.0-4.2.9.1-2018.265;wifi;;2019-06-07 22:24:00;;5cfae44535172f2ca79a24e1;;10.16.6;2019-06-07 22:25:06.548;TIM;;;;;;;;;;;
0000000B;000774571385;2019-06-07 23:12:26.382;0000000B;dwt765ti;cubo;82.55.223.125;;7000;;;;;;;;;SS_QUALITY;;;7000;7000;7000;;;;The Handmaid's Tale;VoD;http://vodabr.cb.ticdn.it/videoteca2/V3/SingleTitle/2019/05/60000243/SS/20089777/20089777_HD.ism/Manifest;systemInput;7.0-4.2.9.1-2018.265;wifi;;2019-06-07 23:11:20;;5cfaef5a35172f2ca79a2b75;;10.16.6;2019-06-07 23:12:26.382;TIM;;;;;;;;;;;
http%3A%2F%2Fvodabr.cb.ticdn.it%2Fvideoteca2%2FV3%2FFilm%2F2017%2F06%2F50670127%2FSS%2F11473278%2F11473278_SD.ism%2FManifest%23https%3A%2F%2Flicense.cubovision.it%2FLicense%2Frightsmanager.asmx`

var elements = strings.Split(str, "\n")

func ExampleFind() {

	for _, element := range elements {
		//fmt.Println(element)
		idv, err := Find(element)
		if err != nil {
			idv = "NON DISPOBINILE"
		}
		fmt.Println(idv)
	}
	// Output:
	// 60000242
	// 60000241
	// 50734100
	// 60000243
	// 60000243
	// 60000243
	// 60000243
	// 60000242
	// 60000243
	// 50670127
}

func BenchmarkFind(b *testing.B) {
	for n := 0; n < b.N; n++ {
		for _, element := range elements {
			//fmt.Println(element)
			Find(element)
		}
	}
}
