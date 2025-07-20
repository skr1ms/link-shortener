package link

import (
	"crypto/rand"
	"linkshortener/internal/stats"
	"linkshortener/pkg/db"
	"math/big"

	"gorm.io/gorm"
)

type Link struct {
	gorm.Model
	OriginalURL string        `gorm:"not null"`
	Hash        string        `gorm:"not null;uniqueIndex:idx_hash;size:12"`
	Stats       []stats.Stats `gorm:"foreignKey:LinkId"`
}

func NewLink(url string) *Link {
	return &Link{
		OriginalURL: url,
		Hash:        GenerateHash(),
	}
}

var base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func GenerateHash() string {
	const hashLength = 12

	result := make([]byte, hashLength)
	charsetSize := big.NewInt(int64(len(base62Chars)))

	for i := 0; i < hashLength; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetSize)
		if err != nil {
			panic("failed to generate random number: " + err.Error())
		}
		result[i] = base62Chars[randomIndex.Int64()]
	}

	return string(result)
}

func CheckUniqueAndGenerateHash(db *db.Db) string {
	for {
		hash := GenerateHash()

		if err := db.First(&Link{}, "hash = ?", hash).Error; err != nil {
			return hash
		}
	}
}
