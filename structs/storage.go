package structs

type Storage interface {
	GetKeys() []string
	GetElement(key string) (interface{}, error)
	GetListElement(key string, index int) (string, error)
	GetDictionaryElement(key, keyInMap string) (string, error)

	PutOrUpdateString(key, value string) (string, bool)
	PutOrUpdateList(key string, value []string) ([]string, bool)
	PutOrUpdateDictionary(key string, value map[string]string) (map[string]string, bool)

	RemoveElement(key string)

	SetTTL(key string, keyTTL int)
}
