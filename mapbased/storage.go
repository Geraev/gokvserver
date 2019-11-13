package mapbased

import (
	"errors"
	"github.com/geraev/gokvserver/structs"
	"sort"
	"sync"
	"time"
)

type Storage struct {
	*sync.RWMutex
	data    map[string]interface{}
	expired map[string]uint64
	janitor *janitor
}

func NewStorage() *Storage {
	s := &Storage{
		RWMutex: new(sync.RWMutex),
		data:    make(map[string]interface{}),
		expired: make(map[string]uint64),
	}
	//S := &struct {
	//	*Storage
	//}{s}
	runJanitor(s, time.Millisecond*20)
	//runtime.SetFinalizer(S, stopJanitor)
	return s
}

func TestTestStorage() *Storage {
	return &Storage{
		RWMutex: new(sync.RWMutex),
		data: map[string]interface{}{
			"keyForStr1": "ValueString_1",
			"keyForStr2": "ValueString_2",
			"keyForList": []string{"new_string_1", "new_string_2"},
			"keyForDict": map[string]string{
				"key_one": "value_one",
				"key_two": "value_two",
			},
		},
	}
}

// GetKeys получение списка ключей
func (s *Storage) GetKeys() []string {
	s.RLock()
	defer s.RUnlock()

	if len(s.data) == 0 {
		return []string{}
	}

	result := make([]string, 0, len(s.data))
	for key := range s.data {
		result = append(result, key)
	}
	sort.Strings(result)
	return result
}

// GetElement получение элемента по ключу
func (s *Storage) GetElement(key string) (interface{}, error) {
	s.RLock()
	//defer s.RUnlock()

	val, ok := s.data[key]
	if !ok {
		s.RUnlock()
		return nil, errors.New("key not found")
	}

	switch v := val.(type) {
	case string, []string, map[string]string:
		s.RUnlock()
		return v, nil
	default:
		s.RUnlock()
		return "", errors.New("something wrong: type error")
	}
}

// GetListElement получение по индексу одного элемента из списка
func (s *Storage) GetListElement(key string, index int) (string, error) {
	s.RLock()
	defer s.RUnlock()

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

// GetDictionaryElement получение по ключу одного элемента из словаря
func (s *Storage) GetDictionaryElement(key, internalKey string) (string, error) {
	s.RLock()
	defer s.RUnlock()

	val, ok := s.data[key]
	if !ok {
		return "", errors.New("key not found")
	}

	v, ok := val.(map[string]string)
	if !ok {
		return "", errors.New("something wrong: type error")
	}

	item, ok := v[internalKey]
	if !ok {
		return "", errors.New("key not found")
	}

	return item, nil
}

// PutOrUpdateString добавление либо обновление значения ключа. Если ключь уже существовал, то перавым аргументом
// возвращается предыдущее значение ключа, а вторым аргументом возвращается true
func (s *Storage) PutOrUpdateString(key, value string) (previousVal string, isUpdated bool) {
	s.Lock()
	//defer s.Unlock()

	if val, ok := s.data[key]; ok {
		previousVal = val.(string)
		isUpdated = ok
	}
	s.data[key] = value
	s.Unlock()
	return previousVal, isUpdated
}

// PutOrUpdateList добавление либо обновление значения ключа. Если ключь уже существовал, то перавым аргументом
// возвращается предыдущее значение ключа, а вторым аргументом возвращается true
func (s *Storage) PutOrUpdateList(key string, value []string) (previousVal []string, isUpdated bool) {
	s.Lock()
	//defer s.Unlock()

	if val, ok := s.data[key]; ok {
		previousVal = val.([]string)
		isUpdated = ok
	}
	s.data[key] = value
	sort.Strings(previousVal)
	s.Unlock()
	return previousVal, isUpdated
}

// PutOrUpdateDictionary добавление либо обновление значения ключа. Если ключь уже существовал, то перавым аргументом
// возвращается предыдущее значение ключа, а вторым аргументом возвращается true
func (s *Storage) PutOrUpdateDictionary(key string, value map[string]string) (previousVal map[string]string, isUpdated bool) {
	s.Lock()
	//defer s.Unlock()

	if val, ok := s.data[key]; ok {
		previousVal = val.(map[string]string)
		isUpdated = ok
	}
	s.data[key] = value
	s.Unlock()
	return previousVal, isUpdated
}

// RemoveElement удаление элемента по ключу
func (s *Storage) RemoveElement(key string) {
	s.Lock()
	//defer s.Unlock()
	delete(s.data, key)
	s.Unlock()
	return
}

// SetTTL установка TTL для ключа и удаление элемента после по прошествии времени.
// TTL устанваливаетс в милисекундах
// Deprecated
func (s *Storage) SetTTL(key string, keyTTL uint64) {
	if keyTTL <= 0 {
		return
	}
	time.AfterFunc(time.Millisecond*time.Duration(keyTTL), func() {
		s.Lock()
		delete(s.data, key)
		s.Unlock()
	})
	return
}

// SetExpired установка TTL для ключа
func (s *Storage) SetExpired(key string, expired uint64) {
	if expired == 0 {
		return
	}
	s.Lock()
	//defer s.Unlock()
	e := time.Now().Add(time.Millisecond * time.Duration(expired)).UnixNano()
	s.expired[key] = uint64(e)
	s.Unlock()
	return
}

// DeleteExpired удаление просроченых кдючей
func (s *Storage) DeleteExpired() {
	s.Lock()
	//defer s.Unlock()

	now := time.Now().UnixNano()
	for key, expired := range s.expired {
		if uint64(now) >= expired {
			delete(s.data, key)
			delete(s.expired, key)
		}
	}
	s.Unlock()
}

func (s *Storage) GetType(key string) (structs.ValueType, error) {
	s.RLock()
	defer s.RUnlock()

	val, ok := s.data[key]
	if !ok {
		return 0, errors.New("key not found")
	}

	switch val.(type) {
	case string:
		return structs.String, nil
	case []string:
		return structs.List, nil
	case map[string]string:
		return structs.Dictionary, nil
	default:
		return 0, errors.New("something wrong: type error")
	}
}
