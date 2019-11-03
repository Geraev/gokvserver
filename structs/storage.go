package structs

type Storage interface {
	GetKeys() []string
	GetElement(key string) (interface{}, error)
	GetListElement(key string, index int) (string, error)
	GetDictionaryElement(key, keyInMap string) (string, error)

	PutOrUpdateString(key, value string) (string, error)
	PutOrUpdateList(key, value []string) ([]string, error)
	PutOrUpdateDictionary(key, value map[string]string) (map[string]string, error)

	SetTTL(key string, keyTTL int) (string, error)
}
