package helpers

import (
	"context"
	"math/rand"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type Export struct {
}
type ckey string

func (c Export) RandString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (c Export) GetUsername(ctx context.Context) string {
	if username, ok := ctx.Value(ckey("username")).(string); ok {
		return username
	}
	return ""

}
