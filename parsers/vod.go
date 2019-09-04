package parsers



type Vod struct {
	RawURL      string
	IDVideoteca string
	ManifestURL string
}

func (v Vod) GetIDVideoteca() string {

	return extractIDVideoteca(v.RawURL)
}

func (v Vod) GetManifestURL() string {

	result, _ := extractManifest(v.RawURL)

	return result 
}
