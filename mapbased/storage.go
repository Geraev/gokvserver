package mapbased

import (
	"errors"
	"sync"
	"time"
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

func (s *Storage) GetElement(key string) (interface{}, error) {
	s.mx.RLock()
	defer s.mx.RUnlock()

	val, ok := s.data[key]
	if !ok {
		return "", errors.New("key not found")
	}

	switch v := val.(type) {
	case string:
		return v, nil
	case []string:
		return v, nil
	case map[string]string:
		return v, nil
	default:
		return "", errors.New("something wrong: type error")
	}
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

// PutOrUpdateString добавляет либо обновляет значение ключа. Если ключь уже существовал, то перавым аргументом
// возвращает предыдущее значение ключа, а вторым аргументом возвращает true
 func (s *Storage) PutOrUpdateString(key, value string) (previousVal string, isUpdated bool) {
	s.mx.Lock()
	defer s.mx.Unlock()

	if val, ok := s.data[key]; ok {
		previousVal = val.(string)
		isUpdated = ok
	}
	s.data[key] = value
	return previousVal, isUpdated
}

// PutOrUpdateList добавляет либо обновляет значение ключа. Если ключь уже существовал, то перавым аргументом
// возвращает предыдущее значение ключа, а вторым аргументом возвращает true
func (s *Storage) PutOrUpdateList(key string, value []string) (previousVal []string, isUpdated bool) {
	s.mx.Lock()
	defer s.mx.Unlock()

	if val, ok := s.data[key]; ok {
		previousVal = val.([]string)
		isUpdated = ok
	}
	s.data[key] = value
	return previousVal, isUpdated
}

// PutOrUpdateDictionary добавляет либо обновляет значение ключа. Если ключь уже существовал, то перавым аргументом
// возвращает предыдущее значение ключа, а вторым аргументом возвращает true
func (s *Storage) PutOrUpdateDictionary(key string, value map[string]string) (previousVal map[string]string, isUpdated bool) {
	s.mx.Lock()
	defer s.mx.Unlock()

	if val, ok := s.data[key]; ok {
		previousVal = val.(map[string]string)
		isUpdated = ok
	}
	s.data[key] = value
	return previousVal, isUpdated
}

func (s *Storage) RemoveElement(key string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	delete(s.data, key)
	return
}

func (s *Storage) SetTTL(key string, keyTTL int) {
	if keyTTL <= 0 {
		return
	}
	time.AfterFunc(time.Millisecond*time.Duration(keyTTL), func() {
		s.mx.Lock()
		delete(s.data, key)
		s.mx.Unlock()
	})
	return
}
