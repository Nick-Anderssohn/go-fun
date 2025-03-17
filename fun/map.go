package fun

import (
	"iter"
	"maps"
)

type MapWorker[K comparable, V any] struct {
	data       map[K]V
	processors []processOneMapElement[K, V]
	curKey     K
}

type processOneMapElement[K comparable, V any] func(
	next func() (K, bool),
) (iteratorIsOpen bool, err error)

func M[K comparable, V any](data map[K]V) *MapWorker[K, V] {
	dst := make(map[K]V, len(data))
	maps.Copy(dst, data)

	return &MapWorker[K, V]{
		data: dst,
	}
}

func (m *MapWorker[K, V]) Collect() (map[K]V, error) {
	keyIterator := maps.Keys(m.data)
	next, stop := iter.Pull(keyIterator)
	defer stop()

	iteratorIsOpen := true
	var err error

	// Walk through our map once, streaming each key/value through
	// the processors.
	for {
		m.curKey, iteratorIsOpen = next()
		if !iteratorIsOpen {
			return m.data, nil
		}

		// run each processor for the current element
		for _, process := range m.processors {
			iteratorIsOpen, err = process(next)

			switch {
			case err != nil:
				return nil, err

			// if the iterator was closed by the process function,
			// that means we are done processing elements
			case !iteratorIsOpen:
				return m.data, nil
			}
		}
	}
}

func (m *MapWorker[K, V]) Filter(check func(K, V) (bool, error)) *MapWorker[K, V] {
	runCheck := func(next func() (K, bool)) (bool, error) {
		iteratorIsOpen := true

		for iteratorIsOpen {
			satisfiesCondition, err := check(m.curKey, m.curVal())

			switch {
			case err != nil:
				return false, err

			case satisfiesCondition:
				return true, nil

			default:
				// Didn't satisfy the check, get it outta here and check the next one
				delete(m.data, m.curKey)
				m.curKey, iteratorIsOpen = next()
			}
		}

		return iteratorIsOpen, nil
	}

	m.processors = append(m.processors, runCheck)
	return m
}

func (m *MapWorker[K, V]) Map(transform func(K, V) (K, V, error)) *MapWorker[K, V] {
	runTransform := func(next func() (K, bool)) (bool, error) {
		updatedKey, updatedValue, err := transform(m.curKey, m.curVal())
		if err != nil {
			return true, err
		}

		m.data[updatedKey] = updatedValue

		if updatedKey != m.curKey {
			delete(m.data, m.curKey)
			m.curKey = updatedKey
		}

		return true, nil
	}

	m.processors = append(m.processors, runTransform)
	return m
}

func (m *MapWorker[K, V]) curVal() V {
	return m.data[m.curKey]
}
