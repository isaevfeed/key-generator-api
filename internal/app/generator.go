package app

import (
	"isaevfeed/internal/cache"
	"math"
	"math/rand"
	"strconv"
	"time"
)

const (
	KEY_NUM  = "number"
	KEY_TEXT = "text"
	KEY_ANY  = "any"
	LETTERS  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	NUMBERS  = "1234567890"
)

type Generator struct {
	keyType   string
	maxKeyLen float64
	cache     *cache.Cache
}

func NewGenerator(keyType string, maxKeyLen float64, cacheRedis *cache.Cache) *Generator {
	return &Generator{keyType: keyType, maxKeyLen: maxKeyLen, cache: cacheRedis}
}

func (g *Generator) GenerateKey() string {
	var keyRes string

	switch g.keyType {
	case KEY_NUM:
		keyRes = g.generateNumberKey()
		break
	case KEY_TEXT:
		keyRes = g.generateTextKey(false)
		break
	case KEY_ANY:
		keyRes = g.generateTextKey(true)
		break
	}

	if g.getCache(keyRes) != "" {
		return g.getCache(keyRes)
	}

	g.cacheResult(keyRes)

	return keyRes
}

func (g *Generator) cacheResult(keyRes string) {
	g.cache.Set(keyRes, keyRes)
}

func (g *Generator) getCache(keyRes string) string {
	return g.cache.Get(keyRes)
}

func (g *Generator) generateNumberKey() string {
	numLen := math.Pow(10, g.maxKeyLen)
	newTimer := rand.NewSource(time.Now().UnixNano())
	newRand := rand.New(newTimer)

	return strconv.Itoa(newRand.Intn(int(numLen)-1) + int(g.maxKeyLen))
}

func (g *Generator) generateTextKey(isAny bool) string {
	keyRes := make([]byte, int(g.maxKeyLen))
	TEMPL := LETTERS

	if isAny {
		TEMPL = LETTERS + NUMBERS
	}

	for i := range keyRes {
		keyRes[i] = TEMPL[rand.Intn(len(TEMPL))]
	}

	return string(keyRes)
}
