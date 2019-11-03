package structs

type Storage interface {
	GetKeys() []string

	GetString(key string) (string, error)
	GetList(key string) ([]string, error)
	GetListElement(key string, index int) (string, error)
	GetDictionary(key string) (map[string]string, error)
	GetDictionaryElement(key, keyInMap string) (map[string]string, error)

	PutOrUpdateString(key, value string) (string, error)
	PutOrUpdateList(key, value []string) ([]string, error)
	PutOrUpdateDictionary(key, value map[string]string) (map[string]string, error)

	SetTTL(key string, keyTTL int) (string, error)
}
