package fun

type processOne[V any] func(worker *SliceWorker[V], pos int) error

type SliceWorker[V any] struct {
	data       []V
	processors []processOne[V]
}

func S[V any](data []V) *SliceWorker[V] {
	dst := make([]V, len(data))
	copy(dst, data)

	return &SliceWorker[V]{
		data: dst,
	}
}

func (s *SliceWorker[V]) Filter(check func(V) (bool, error)) *SliceWorker[V] {
	runCheck := func(worker *SliceWorker[V], pos int) error {
		for pos < len(worker.data) {
			satisfiesCondition, err := check(worker.data[pos])

			switch {
			case err != nil:
				return err

			case satisfiesCondition:
				return nil

			default:
				// Didn't satisfy the check, get it outta here and check the next one
				worker.data = append(worker.data[:pos], worker.data[pos+1:]...)
			}
		}

		return nil
	}

	s.processors = append(s.processors, runCheck)
	return s
}

func (s *SliceWorker[V]) Map(modify func(V) (V, error)) *SliceWorker[V] {
	runModify := func(worker *SliceWorker[V], pos int) error {
		updatedValue, err := modify(s.data[pos])
		if err != nil {
			return err
		}

		s.data[pos] = updatedValue

		return nil
	}

	s.processors = append(s.processors, runModify)
	return s
}

func (s *SliceWorker[V]) Finish() ([]V, error) {
	// Walk through our list once. Yay O(N)
	for i := 0; i < len(s.data); i++ {
		// run each processor for the current element
		for _, process := range s.processors {
			err := process(s, i)
			if err != nil {
				return []V{}, err
			}

			// In case the process function removed enough elements
			// to cause our index to be out of bounds (ex Filter could do that).
			// Note that this is okay; this is not an error.
			if i >= len(s.data) {
				return s.data, nil
			}
		}
	}

	return s.data, nil
}
