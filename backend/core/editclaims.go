package core

type DataContract[T any, TKey any] struct {
	Data   T      `json:"data"`
	Key    TKey   `json:"-"`
	Token  string `json:"token"`
	Status string `json:"status"`
}
