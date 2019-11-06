package structs

type ValueType int

const (
	String ValueType = iota
	List
	Dictionary
)

func (t ValueType) String() string {
	return [...]string{"String", "List", "Dictionary"}[t]
}

type Storage interface {
	GetKeys() []string
	GetElement(key string) (interface{}, error)
	GetListElement(key string, index int) (string, error)
	GetDictionaryElement(key, internalKey string) (string, error)

	PutOrUpdateString(key, value string) (string, bool)
	PutOrUpdateList(key string, value []string) ([]string, bool)
	PutOrUpdateDictionary(key string, value map[string]string) (map[string]string, bool)

	RemoveElement(key string)

	SetTTL(key string, keyTTL uint64)
	GetType(key string) (ValueType, error)
}
