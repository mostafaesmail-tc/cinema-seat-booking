package booking


type Memory_store struct {
	// seats --> booking
	// we are doing it in memory first not in redis or postgresql

	// so we would go with a map

	bookings map[string]Booking
}


func NewMemoryStore() *Memory_store {
	return &Memory_store{
		bookings:map[string]Booking{},
	}
}


func (s *Memory_store) Book(b Booking) error {

	if _, exists := s.bookings[b.SeatID]; exists {
		return ErrSeatAlreadyBooked
	}

	s.bookings[b.SeatID] = b;

	return nil
}
func (s *Memory_store) ListBookings(movieId string) []Booking {

	var result []Booking

	for _, b := range s.bookings {
		if b.MovieID == movieId {
			result = append(result, b);
		}
	}

	return result;
}