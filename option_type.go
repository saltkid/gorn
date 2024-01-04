package main

import "errors"

type None[T any] struct{}
type Some[T any] struct{ data T }

type Option[T any] interface {
	IsNone() bool
	IsSome() bool
	Get() (T, error)
}

func none[T any]() Option[T] {
	return None[T]{}
}

func (None[T]) IsNone() bool {
	return true
}

func (None[T]) IsSome() bool {
	return false
}

func (None[T]) Get() (T, error) {
	var data T
	return data, errors.New("no data, is none")
}

func some[T any](data T) Option[T] {
	return Some[T]{data}
}

func (Some[T]) IsNone() bool {
	return false
}

func (Some[T]) IsSome() bool {
	return true
}

func (s Some[T]) Get() (T, error) {
	return s.data, nil
}
