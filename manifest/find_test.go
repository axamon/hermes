package manifest

import "testing"

func TestFind(t *testing.T) {
	type args struct {
		rawurl string
	}
	tests := []struct {
		name            string
		args            args
		wantUrlmanifest string
		wantErr         bool
	}{
		{"Vod", args{rawurl: "http://vodabr.cb.ticdn.it/videoteca2/V3/Film/2012/11/50290465/SS/20085566/20085566_HD.ism/QualityLevels(192000)/Fragments(audio_ita_2=71249710000)"}, "http://vodabr.cb.ticdn.it/videoteca2/V3/Film/2012/11/50290465/SS/20085566/20085566_HD.ism/Manifest", false},
		{"URLEncoded", args{rawurl: "http%3A%2F%2Fvodabr.cb.ticdn.it%2Fvideoteca2%2FV3%2FFilm%2F2017%2F06%2F50670127%2FSS%2F11473278%2F11473278_SD.ism%2FManifest%23https%3A%2F%2Flicense.cubovision.it%2FLicense%2Frightsmanager.asmx"}, "http://vodabr.cb.ticdn.it/videoteca2/V3/Film/2017/06/50670127/SS/11473278/11473278_SD.ism/Manifest", false },
		{"dash", args{rawurl: "https://voddashhttps.cb.ticdn.it/videoteca2/V3/Film/2016/05/50565607/DASH_H265/20083913/20083913_HD_ANSN.mpd"}, "https://voddashhttps.cb.ticdn.it/videoteca2/V3/Film/2016/05/50565607/DASH_H265/20083913/20083913_HD_ANSN.mpd", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUrlmanifest, err := Find(tt.args.rawurl)
			if (err != nil) != tt.wantErr {
				t.Errorf("Find() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUrlmanifest != tt.wantUrlmanifest {
				t.Errorf("Find() = %v, want %v", gotUrlmanifest, tt.wantUrlmanifest)
			}
		})
	}
}
