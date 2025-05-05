package list

import (
	"encoding/json"
	"fmt"
)

type List[T any] struct {
	elements []T
}

func NewList[T any]() *List[T] {
	return &List[T]{elements: make([]T, 0)}
}

func AsList[T any](elements ...T) *List[T] {
	list := NewList[T]()
	if elements != nil && len(elements) > 0 {
		list.Add(elements...)
	}
	return list
}

func (l *List[T]) Add(element ...T) {
	l.elements = append(l.elements, element...)
}

func (l *List[T]) Remove(index int) error {
	if index < 0 || index >= l.Size() {
		return fmt.Errorf("index out of range")
	}
	l.elements = append(l.elements[:index], l.elements[index+1:]...)
	return nil
}

func (l *List[T]) RemoveIf(predicate func(T) bool) *List[T] {
	newArray := NewList[T]()
	for _, element := range l.Slice() {
		if !predicate(element) {
			newArray.Add(element)
		}
	}
	return newArray
}

func (l *List[T]) Get(index int) (T, error) {
	var zeroValue T
	if index < 0 || index >= l.Size() {
		return zeroValue, fmt.Errorf("index out of range")
	}
	return l.elements[index], nil
}

func (l *List[T]) Size() int {
	return len(l.elements)
}

func (l *List[T]) IsEmpty() bool {
	return l.Size() == 0
}

func (l *List[T]) Filter(predicate func(T) bool) *List[T] {
	filteredList := NewList[T]()
	for _, element := range l.Slice() {
		if predicate(element) {
			filteredList.Add(element)
		}
	}
	return filteredList
}

func (l *List[T]) ForEach(action func(T)) {
	for _, element := range l.Slice() {
		action(element)
	}
}

func (l *List[T]) Contains(element T, predicate func(T, T) bool) bool {
	for _, e := range l.Slice() {
		if predicate(e, element) {
			return true
		}
	}
	return false
}

func (l *List[T]) First() (T, error) {
	var zeroValue T
	if l.IsEmpty() {
		return zeroValue, fmt.Errorf("array is empty")
	}
	return l.elements[0], nil
}

func (l *List[T]) Last() (T, error) {
	var zeroValue T
	if l.IsEmpty() {
		return zeroValue, fmt.Errorf("array is empty")
	}
	return l.elements[l.Size()-1], nil
}

func (l *List[T]) Slice() []T {
	return l.elements
}

func Map[T, R any](arr *List[T], mapper func(T) R) *List[R] {
	mapped := NewList[R]()
	for _, element := range arr.Slice() {
		mapped.Add(mapper(element))
	}
	return mapped
}

func (l *List[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.elements)
}

func (l *List[T]) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &l.elements)
}

func (l *List[T]) Iter() <-chan T {
	ch := make(chan T)
	go func() {
		defer close(ch)
		for _, elem := range l.elements {
			ch <- elem
		}
	}()
	return ch
}
