package parsers


// RawURL rappresenta una qualsiasi URL
type RawURL struct {
	URL string
}

// GetIDVideoteca estrae l'id univoco di videoteca se prensente.
func (u RawURL) GetIDVideoteca() string {
	
	return extractIDVideoteca(u.URL)
}

// GetManifestURL restituisce la URL del manifest di riferimento.
// Sono supportati DASH e SS
func (u RawURL) GetManifestURL() string {

	result, _ := extractManifest(u.URL)

	return result 
}