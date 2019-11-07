package mapbased

import (
	"errors"
	"github.com/geraev/gokvserver/structs"
	"reflect"
	"sort"
	"sync"
	"testing"
)

var (
	storg = Storage{
		mx: &sync.RWMutex{},
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
)

func TestStorage_GetKeys(t *testing.T) {
	tests := []struct {
		name   string
		fields Storage
		want   []string
	}{
		{
			name:   "Testing GetElement",
			fields: storg,
			want: []string{
				"keyForStr1",
				"keyForStr2",
				"keyForList",
				"keyForDict",
			},
		},
		{
			name:   "Testing GetElement: return empty list",
			fields: Storage{
				mx: &sync.RWMutex{},
				data: map[string]interface{}{},
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				mx:   tt.fields.mx,
				data: tt.fields.data,
			}
			sort.Strings(tt.want)
			if got := s.GetKeys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_GetElement(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  Storage
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name:    "Testing GetElement: value type string",
			fields:  storg,
			args:    args{key: "keyForStr1"},
			want:    "ValueString_1",
			wantErr: false,
		},
		{
			name:    "Testing GetElement: value type string",
			fields:  storg,
			args:    args{key: "keyForStr2"},
			want:    "ValueString_2",
			wantErr: false,
		},
		{
			name:    "Testing GetElement: value type slice",
			fields:  storg,
			args:    args{key: "keyForList"},
			want:    []string{"new_string_1", "new_string_2"},
			wantErr: false,
		},
		{
			name:   "Testing GetElement: value type map",
			fields: storg,
			args:   args{key: "keyForDict"},
			want: map[string]string{
				"key_one": "value_one",
				"key_two": "value_two",
			},
			wantErr: false,
		},
		{
			name:   "Testing GetElement: failed type",
			fields: Storage{
				mx: &sync.RWMutex{},
				data: map[string]interface{}{
					"key_01": struct {
						Int int
						Str string
						Arr [5]int
					}{},
				},
			},
			args:   args{key: "key_01"},
			want: "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				mx:   tt.fields.mx,
				data: tt.fields.data,
			}
			got, err := s.GetElement(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetElement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetElement() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_GetListElement(t *testing.T) {
	type args struct {
		key   string
		index int
	}
	tests := []struct {
		name    string
		fields  Storage
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name:    "Testing GetListElement: index 1",
			fields:  storg,
			args:    args{key: "keyForList", index: 1},
			want:    "new_string_2",
			wantErr: false,
		},
		{
			name:    "Testing GetListElement: index out of range",
			fields:  storg,
			args:    args{key: "keyForList", index: -1},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Testing GetListElement: index out of range",
			fields:  storg,
			args:    args{key: "keyForList", index: 9999999},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Testing GetListElement: key not found",
			fields:  storg,
			args:    args{key: "failedKey", index: 1},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Testing GetListElement:  type error",
			fields:  storg,
			args:    args{key: "keyForDict", index: 1},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				mx:   tt.fields.mx,
				data: tt.fields.data,
			}
			got, err := s.GetListElement(tt.args.key, tt.args.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetListElement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetListElement() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_GetDictionaryElement(t *testing.T) {
	type args struct {
		key         string
		internalKey string
	}
	tests := []struct {
		name    string
		fields  Storage
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name:    "Testing GetDictionaryElement",
			fields:  storg,
			args:    args{key: "keyForDict", internalKey: "key_two"},
			want:    "value_two",
			wantErr: false,
		},
		{
			name:    "Testing GetDictionaryElement: key not found",
			fields:  storg,
			args:    args{key: "failedKey", internalKey: "failed_internal_key"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Testing GetDictionaryElement: internal key not found",
			fields:  storg,
			args:    args{key: "keyForDict", internalKey: "failed_internal_key"},
			want:    "",
			wantErr: true,
		},
		{
			name:    "Testing GetDictionaryElement: type error",
			fields:  storg,
			args:    args{key: "keyForList", internalKey: "failed_internal_key"},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				mx:   tt.fields.mx,
				data: tt.fields.data,
			}
			got, err := s.GetDictionaryElement(tt.args.key, tt.args.internalKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDictionaryElement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetDictionaryElement() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_PutOrUpdateDictionary(t *testing.T) {
	type args struct {
		key   string
		value map[string]string
	}
	tests := []struct {
		name            string
		fields          Storage
		args            args
		wantPreviousVal map[string]string
		wantIsUpdated   bool
	}{
		{
			name:            "Testing PutOrUpdateDictionary: put new value",
			fields:          storg,
			args:            args{key: "key0918635", value: map[string]string{"key0001": "Mars", "key0002": "Mercury", "key0003": "Neptune"}},
			wantPreviousVal: nil,
			wantIsUpdated:   false,
		},
		{
			name:            "Testing PutOrUpdateDictionary: update value",
			fields:          storg,
			args:            args{key: "keyForDict", value: map[string]string{"color": "Red", "size": "Large"}},
			wantPreviousVal: storg.data["keyForDict"].(map[string]string),
			wantIsUpdated:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				mx:   tt.fields.mx,
				data: tt.fields.data,
			}
			gotPreviousVal, gotIsUpdated := s.PutOrUpdateDictionary(tt.args.key, tt.args.value)
			if !reflect.DeepEqual(gotPreviousVal, tt.wantPreviousVal) {
				t.Errorf("PutOrUpdateDictionary() gotPreviousVal = %v, want %v", gotPreviousVal, tt.wantPreviousVal)
			}
			if gotIsUpdated != tt.wantIsUpdated {
				t.Errorf("PutOrUpdateDictionary() gotIsUpdated = %v, want %v", gotIsUpdated, tt.wantIsUpdated)
			}
		})
	}
}

func TestStorage_PutOrUpdateList(t *testing.T) {
	type args struct {
		key   string
		value []string
	}
	tests := []struct {
		name            string
		fields          Storage
		args            args
		wantPreviousVal []string
		wantIsUpdated   bool
	}{
		{
			name:            "Testing PutOrUpdateList: put new value",
			fields:          storg,
			args:            args{key: "list_077", value: []string{"Mars", "Mercury", "Neptune"}},
			wantPreviousVal: nil,
			wantIsUpdated:   false,
		},
		{
			name:            "Testing PutOrUpdateList: update value",
			fields:          storg,
			args:            args{key: "keyForList", value: []string{"Jean-Claude", "Van", "Damme"}},
			wantPreviousVal: storg.data["keyForList"].([]string),
			wantIsUpdated:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				mx:   tt.fields.mx,
				data: tt.fields.data,
			}
			sort.Strings(tt.wantPreviousVal)
			gotPreviousVal, gotIsUpdated := s.PutOrUpdateList(tt.args.key, tt.args.value)
			if !reflect.DeepEqual(gotPreviousVal, tt.wantPreviousVal) {
				t.Errorf("PutOrUpdateList() gotPreviousVal = %v, want %v", gotPreviousVal, tt.wantPreviousVal)
			}
			if gotIsUpdated != tt.wantIsUpdated {
				t.Errorf("PutOrUpdateList() gotIsUpdated = %v, want %v", gotIsUpdated, tt.wantIsUpdated)
			}
		})
	}
}

func TestStorage_PutOrUpdateString(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name            string
		fields          Storage
		args            args
		wantPreviousVal string
		wantIsUpdated   bool
	}{
		{
			name:            "Testing PutOrUpdateString: put new value",
			fields:          storg,
			args:            args{key: "str_19", value: "fruits"},
			wantPreviousVal: "",
			wantIsUpdated:   false,
		},
		{
			name:            "Testing PutOrUpdateString: update value",
			fields:          storg,
			args:            args{key: "keyForStr2", value: "gocache"},
			wantPreviousVal: storg.data["keyForStr2"].(string),
			wantIsUpdated:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				mx:   tt.fields.mx,
				data: tt.fields.data,
			}
			gotPreviousVal, gotIsUpdated := s.PutOrUpdateString(tt.args.key, tt.args.value)
			if gotPreviousVal != tt.wantPreviousVal {
				t.Errorf("PutOrUpdateString() gotPreviousVal = %v, want %v", gotPreviousVal, tt.wantPreviousVal)
			}
			if gotIsUpdated != tt.wantIsUpdated {
				t.Errorf("PutOrUpdateString() gotIsUpdated = %v, want %v", gotIsUpdated, tt.wantIsUpdated)
			}
		})
	}
}

func TestStorage_GetType(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  Storage
		args    args
		want    structs.ValueType
		wantErr error
	}{
		{
			name: "Testing GetType: return Dictionary",
			fields: storg,
			args: args{key: "keyForDict"},
			want: structs.Dictionary,
			wantErr: (error)(nil),
		},
		{
			name: "Testing GetType: return List",
			fields: storg,
			args: args{key: "keyForList"},
			want: structs.List,
			wantErr: (error)(nil),
		},
		{
			name: "Testing GetType: return String",
			fields: storg,
			args: args{key: "keyForStr2"},
			want: structs.String,
			wantErr: (error)(nil),
		},
		{
			name: "Testing GetType: key not found",
			fields: storg,
			args: args{key: "keyFailed"},
			want: 0,
			wantErr: errors.New("key not found"),
		},
		{
			name:   "Testing GetType: failed type",
			fields: Storage{
				mx: &sync.RWMutex{},
				data: map[string]interface{}{
					"key_01": struct {
						Int int
						Str string
						Arr [5]int
					}{},
				},
			},
			args:   args{key: "key_01"},
			want: 0,
			wantErr: errors.New("something wrong: type error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				mx:   tt.fields.mx,
				data: tt.fields.data,
			}
			got, err := s.GetType(tt.args.key)
			if err != tt.wantErr {
				t.Errorf("GetType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetType() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_RemoveElement(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields Storage
		args   args
	}{
		{
			name: "Testing RemoveElement",
			fields: storg,
			args: args{key: "keyForDict"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				mx:   tt.fields.mx,
				data: tt.fields.data,
			}
			s.RemoveElement(tt.args.key)

			_, err := s.GetElement(tt.args.key)
			if err == nil {
				t.Errorf("GetElement() got = %v, want %v", nil, err)
			}
		})
	}
}










/*func TestStorage_SetTTL(t *testing.T) {
	type fields struct {
		mx   *sync.RWMutex
		data map[string]interface{}
	}
	type args struct {
		key    string
		keyTTL uint64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				mx:   tt.fields.mx,
				data: tt.fields.data,
			}
		})
	}
}
*/

























