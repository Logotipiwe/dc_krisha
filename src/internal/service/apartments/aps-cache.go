package apartments

import (
	"krisha/src/internal/domain"
)

type ApsCacheService struct {
	cache map[string]*domain.Ap
}

func NewApsCacheService() *ApsCacheService {
	return &ApsCacheService{
		make(map[string]*domain.Ap),
	}
}

func (s *ApsCacheService) AddToCache(aps map[string]*domain.Ap) {
	for id, val := range aps {
		existing := s.cache[id]
		if existing == nil {
			s.cache[id] = val
		}
	}
}

func (s *ApsCacheService) IsInCache(id string) bool {
	_, is := s.cache[id]
	return is
}
