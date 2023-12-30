package main

import "errors"

type OptionalData[T any] interface {
	get() (T, error)
}
type None[T any] struct{}
type Some[T any] struct {data T}

func (None[T]) get() (T, error) {
	var data T
	return data, errors.New("no data, is none")
}

func (s Some[T]) get() (T, error) {
	return s.data, nil
}

func create_some[T any](data T) OptionalData[T] {
	return Some[T]{data}
}

func create_none[T any]() OptionalData[T] {
	return None[T]{}
}


type Option[T any] interface {
	is_none() bool
	is_some() bool
	find() OptionalData[T]
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

func (None[T]) find() OptionalData[T] {
	return create_none[T]()
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

func (s Some[T]) find() OptionalData[T] {
	return create_some[T](s.data)
}