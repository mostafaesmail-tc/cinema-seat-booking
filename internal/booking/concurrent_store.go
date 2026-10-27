package booking

import "sync"


type Concurrent_store struct {
	// seats --> booking
	// we are doing it in Concurrent first not in redis or postgresql

	// so we would go with a map

	bookings map[string]Booking
	sync.RWMutex
}


func NewConcurrentStore() *Concurrent_store {
	return &Concurrent_store{
		bookings:map[string]Booking{},
	}
}


func (s *Concurrent_store) Book(b Booking) error {
	s.Lock()
	defer s.Unlock()

	if _, exists := s.bookings[b.SeatID]; exists {
		return ErrSeatAlreadyBooked
	}

	s.bookings[b.SeatID] = b;

	return nil
}
func (s *Concurrent_store) ListBookings(movieId string) []Booking {

	s.RLock()
	defer s.RUnlock()

	var result []Booking

	for _, b := range s.bookings {
		if b.MovieID == movieId {
			result = append(result, b);
		}
	}

	return result;
}