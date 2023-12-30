package main

import "errors"

type None[T any] struct{}
type Some[T any] struct {data T}

type Option[T any] interface {
	is_none() bool
	is_some() bool
	get() (T, error)
}

func none[T any]() Option[T] {
	return None[T]{}
}

func (None[T]) is_none() bool {
	return true
}

func (None[T]) is_some() bool {
	return false
}

func (None[T]) get() (T, error) {
	var data T
	return data, errors.New("no data, is none")
}

func some[T any](data T) Option[T] {
	return Some[T]{data}
}

func (Some[T]) is_none() bool {
	return false
}

func (Some[T]) is_some() bool {
	return true
}

func (s Some[T]) get() (T, error) {
	return s.data, nil
}