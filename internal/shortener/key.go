package shortener

import (
	"encoding/base64"

	"github.com/google/uuid"
)

func generateKey(url string) string {
	id, _ := uuid.NewRandom()
	base64ID := base64.URLEncoding.EncodeToString([]byte(id.String()))
	return base64ID[:7]
}
