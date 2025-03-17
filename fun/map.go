package fun

import (
	"iter"
	"maps"
)

type MapStream[K comparable, V any] struct {
	data       map[K]V
	processors []processMapElement[K, V]

	next func() (K, bool)
	stop func()
}

type processMapElement[K comparable, V any] func(
	k K,
	v V,
) (updatedKey K, updatedValue V, iteratorIsOpen bool, err error)

func NewMapStream[K comparable, V any](data map[K]V) *MapStream[K, V] {
	keyIterator := maps.Keys(data)
	next, stop := iter.Pull(keyIterator)

	return &MapStream[K, V]{
		data: data,
		next: next,
		stop: stop,
	}
}

func (m *MapStream[K, V]) Collect() (map[K]V, error) {
	var err error
	var curKey K
	iteratorIsOpen := true
	dst := map[K]V{}

	// Walk through our map once, streaming each key/value through the processors.
	for iteratorIsOpen {
		curKey, iteratorIsOpen = m.next()
		if !iteratorIsOpen {
			return dst, nil
		}

		// run each processor for the current element
		curVal := m.data[curKey]
		for _, process := range m.processors {
			curKey, curVal, iteratorIsOpen, err = process(curKey, curVal)

			switch {
			case err != nil:
				return nil, err

			// if the iterator was closed by the process function,
			// that means we are done processing elements
			case !iteratorIsOpen:
				return dst, nil
			}
		}

		dst[curKey] = curVal
	}

	return dst, nil
}

func (m *MapStream[K, V]) Filter(check func(K, V) (bool, error)) *MapStream[K, V] {
	runCheck := func(k K, v V) (K, V, bool, error) {
		iteratorIsOpen := true

		for iteratorIsOpen {
			satisfiesCondition, err := check(k, v)

			switch {
			case err != nil:
				return k, v, false, err

			case satisfiesCondition:
				return k, v, true, nil

			default:
				// Didn't satisfy the check, move on and check the next one
				k, iteratorIsOpen = m.next()
				v = m.data[k]
			}
		}

		return k, v, iteratorIsOpen, nil
	}

	m.processors = append(m.processors, runCheck)
	return m
}

func (m *MapStream[K, V]) Map(transform func(K, V) (K, V, error)) *MapStream[K, V] {
	runTransform := func(k K, v V) (K, V, bool, error) {
		updatedKey, updatedValue, err := transform(k, v)
		if err != nil {
			return k, v, true, err
		}

		return updatedKey, updatedValue, true, nil
	}

	m.processors = append(m.processors, runTransform)
	return m
}
