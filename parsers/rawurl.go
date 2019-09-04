package parsers


type RawURL struct {
	URL string
}

func (u RawURL) GetIDVideoteca() string {
	
	return extractIDVideoteca(u.URL)
}

func (u RawURL) GetManifestURL() string {

	result, _ := extractManifest(u.URL)

	return result 
}