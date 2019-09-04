// Copyright 2019 Alberto Bregliano. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parsers



// Vod è una struttura per archiviare i dati estratti.
type Vod struct {
	RawURL      string
	IDVideoteca string
	ManifestURL string
}

// GetIDVideoteca estrae l'id univoco di videoteca se prensente.
func (v Vod) GetIDVideoteca() string {

	return extractIDVideoteca(v.RawURL)
}


// GetManifestURL restituisce la URL del manifest di riferimento.
// Sono supportati DASH e SS
func (v Vod) GetManifestURL() string {

	result, _ := extractManifest(v.RawURL)

	return result 
}
