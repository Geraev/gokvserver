package mapbased

import (
	"errors"
	"sync"
)

type Storage struct {
	mx   *sync.RWMutex
	data map[string]interface{}
}

func NewStorage() *Storage {
	return &Storage{
		mx:   new(sync.RWMutex),
		data: make(map[string]interface{}),
	}
}

func (s *Storage) GetKeys() []string {
	s.mx.RLock()
	defer s.mx.RUnlock()

	if len(s.data) == 0 {
		return []string{}
	}

	result := make([]string, 0, len(s.data))
	for key := range s.data {
		result = append(result, key)
	}

	return result
}

// TODO
func (s *Storage) GetString(key string) (string, error) {
	panic("implement me")
}

// TODO
func (s *Storage) GetList(key string) ([]string, error) {
	panic("implement me")
}

func (s *Storage) GetListElement(key string, index int) (string, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	if index < 0 {
		return "", errors.New("index out of range")
	}

	val, ok := s.data[key]
	if !ok {
		return "", errors.New("key not found")
	}

	v, ok := val.([]string)
	if !ok {
		return "", errors.New("something wrong: type error")
	}

	if index >= len(v) {
		return "", errors.New("index out of range")
	}

	return v[index], nil
}

// TODO
func (s *Storage) GetDictionary(key string) (map[string]string, error) {
	panic("implement me")
}

func (s *Storage) GetDictionaryElement(key, keyInMap string) (string, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	val, ok := s.data[key]
	if !ok {
		return "", errors.New("key not found")
	}

	v, ok := val.(map[string]string)
	if !ok {
		return "", errors.New("something wrong: type error")
	}

	item, ok := v[keyInMap]
	if !ok {
		return "",errors.New("key not found")
	}

	return item, nil
}

// TODO
func (s *Storage) PutOrUpdateString(key, value string) (string, error) {
	panic("implement me")
}

// TODO
func (s *Storage) PutOrUpdateList(key, value []string) ([]string, error) {
	panic("implement me")
}

// TODO
func (s *Storage) PutOrUpdateDictionary(key, value map[string]string) (map[string]string, error) {
	panic("implement me")
}

// TODO
func (s *Storage) SetTTL(key string, keyTTL int) (string, error) {
	panic("implement me")
}
