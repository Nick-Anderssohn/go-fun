package fun

type SliceStream[V any] struct {
	data       []V
	processors []processSliceElement[V]

	pos int
}

type processSliceElement[V any] func(v V) (V, error)

func NewSliceStream[V any](data []V) *SliceStream[V] {
	return &SliceStream[V]{
		data: data,
	}
}

func (s *SliceStream[V]) Collect() ([]V, error) {
	dst := []V{}

	// Walk through our list once, streaming each value through
	// the processors.
	for ; s.pos < len(s.data); s.pos++ {
		var err error
		v := s.data[s.pos]

		// run each processor for the current element
		for _, process := range s.processors {
			v, err = process(v)

			switch {
			case err == errEndOfStream:
				return dst, nil

			case err != nil:
				return []V{}, err
			}
		}

		dst = append(dst, v)
	}

	return dst, nil
}

func (s *SliceStream[V]) Filter(check func(V) (bool, error)) *SliceStream[V] {
	runCheck := func(v V) (V, error) {
		for s.pos < len(s.data) {
			satisfiesCondition, err := check(v)

			switch {
			case err != nil:
				return v, err

			case satisfiesCondition:
				return v, nil

			default:
				// We'll try the next one
				s.pos++

				if s.pos < len(s.data) {
					v = s.data[s.pos]
				}
			}
		}

		return v, errEndOfStream
	}

	s.processors = append(s.processors, runCheck)
	return s
}

func (s *SliceStream[V]) Map(transform func(V) (V, error)) *SliceStream[V] {
	runTransform := func(v V) (V, error) {
		updatedValue, err := transform(v)
		if err != nil {
			return v, err
		}

		return updatedValue, err
	}

	s.processors = append(s.processors, runTransform)
	return s
}
