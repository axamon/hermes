package hasher

import (
	"os"
	"log"
	"crypto/sha256"
	"io"
	"fmt"
)

func Sum(file string) (hash string) {
f, err := os.Open(file)
if err != nil {
    log.Fatal(err)
}
defer f.Close()

h := sha256.New()
if _, err := io.Copy(h, f); err != nil {
  log.Fatal(err)
}

hash = fmt.Sprintf("%x", h.Sum(nil))

return hash
}