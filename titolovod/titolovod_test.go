package titolovod_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/axamon/hermes/titolovod"
)

const testjson = `[{"resultCode":"OK","systemTime":1568050273,"errors":[],"resultObj":{"id":"50649232","containers":[{"actions":[{"key":"onClick","uri":"/DETAILScontentId=50649232&type=MOVIE&deviceType=ALL&serviceName=ALL","targetType":"MOVIE"}],"layout":"CONTENT_DETAILS","metadata":{"shortDescription":"Martino (Alessandro Siani) è uno scansafatiche che vive in Svizzera presso la sorella Caterina che fa la domestica in casa del dottor Guglielmo Gioia, una specie di guru motivazionale per persone in crisi. Fino a quando la donna per un incidente...","parentalControlLevel":"1","language":"","contentId":"50649232","expirationDate":1575327540,"title":"Mister Felicità","contentAnalyticsType":"svod","type":"Vod","longDescription":"Martino (Alessandro Siani) è uno scansafatiche che vive in Svizzera presso la sorella Caterina che lavora come donna delle pulizie in casa del dottor Guglielmo Gioia (Diego Abatantuono), una specie di guru motivazionale per persone in crisi. Ma quando la donna ha un incidente che non le permette più di lavorare, Martino è costretto a sostituirla nelle faccende domestiche in casa Gioia, dove comincerà anche a spacciarsi per l'assistente del capo.","year":"2017","contentProvider":"raicinema","duration":5399,"genre":"Commedia","supportedDevices":["ANDROIDSMARTPHONE","ANDROIDTABLET","CONNECTEDTV_SS","CUBO","IPAD","IPHONE"],"categoriesThemeArea":["CINEMA"],"rating":"3.0","badge":"","offerTypeLabel":null,"imageUrl":"http://images2.timvision.it/videoteca2/VCMS/VOD/2018/10/50649232/1_50649232__CA__23__WL___CA_23_WL_MED.jpg","imageUrlOther":"http://images2.timvision.it/videoteca2/VCMS/VOD/2018/10/50649232/1_50649232__CA__43__NL___CA_43_NL_LOW.jpg","bgImageUrl":"http://images2.timvision.it/videoteca2/VCMS/VOD/2018/10/50649232/1_50649232__BG_01__169__NL___NL___BG_HIGH.jpg","displayExpirationDate":"02/12/19","directors":["Alessandro Siani"],"actors":["Alessandro Siani","Diego Abatantuono","Carla Signoris","Elena Cucci","Cristina Dell'Anna","Yari Gugliucci","Ernesto Mahieux","Pippo Lorusso"],"broadcastDisplayDate":"","objectSubtype":"MOVIE","objectType":"VIDEO","titleBrief":"","videoType":["SD"],"playback":{"audioLanguages":[{"id":"0","name":"ITA","isPreferred":"Y"}],"subtitles":[{"id":"0","name":"ITA","isPreferred":"N"}]},"visionPercentage":"0","contentOptions":["SUBTITLED","ISDOWNLOADABLE"]}},{"layout":"RIGHTS","actions":[],"retrieveItems":{"uri":"/RIGHTScontentId=50649232&deviceType=ALL&serviceName=ALL&type=MOVIE","type":"REMOTE"},"metadata":{}},{"layout":"SMALL_CARDS","actions":[],"retrieveItems":{"uri":"/TRAY/RECOMdeviceType=ALL&serviceName=ALL&dataSet=RICH&recomType=MORE_LIKE_THIS&contentId=50649232&maxResults=25&category=CINEMA","type":"REMOTE"},"metadata":{"label":"GUARDA ANCHE","imageUrl":null}},{"layout":"SMALL_CARDS","actions":[],"retrieveItems":{"uri":"/TRAY/CELEBRITIESfrom=0&to=50&deviceType=ALL&serviceName=ALL&contentId=50649232","type":"REMOTE"},"metadata":{"label":"ATTORI E REGISTI","imageUrl":null}}],"total":4},"contentId":"50649232"}]`

var testServer = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
	//	res.WriteHeader(http.StatusOK)
	res.Write([]byte(testjson))
}))

func TestElaboraResp(t *testing.T) {

	type args struct {
		bodyBytes []byte
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
		wantErr    bool
	}{
		{"Secondo", args{bodyBytes: []byte(testjson)}, "Mister Felicità", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := titolovod.ElaboraResp(tt.args.bodyBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("ElaboraResp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !(gotResult[0].ResultObj.Containers[0].Metadata.Title == tt.wantResult) {
				t.Errorf("ElaboraResp() = %v, want %v", gotResult[0].ResultObj.Containers[0].Metadata.Title, tt.wantResult)
			}

		})
	}
}
