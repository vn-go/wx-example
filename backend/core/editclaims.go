package core

type EditClaims[T any, TKey any] struct {
	Data  T      `json:"data"`
	Key   TKey   `json:"-"`
	Token string `json:"token"`
}
