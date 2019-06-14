package parsers

const seed = "vvkidtbcjujhgffbjnvrngvrinvufjkvljreucecvfcj"

type Utenti []Fruizioni

type Fruizioni struct {
	Hashfruizione map[string]bool
	Clientip      map[string]string
	Idvideoteca   map[string]string
	Idaps         map[string]string
	Edgeip        map[string]string
	Giorno        map[string]string
	Orario        map[string]string
	Details       map[string][]float64 `json:"-"`
}
