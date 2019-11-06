package mapbased

import (
	"reflect"
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

/*func TestNewStorage(t *testing.T) {
	tests := []struct {
		name string
		want *Storage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStorage(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_GetDictionaryElement(t *testing.T) {
	type fields struct {
		mx   *sync.RWMutex
		data map[string]interface{}
	}
	type args struct {
		key      string
		keyInMap string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				mx:   tt.fields.mx,
				data: tt.fields.data,
			}
			got, err := s.GetDictionaryElement(tt.args.key, tt.args.keyInMap)
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
*/
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

/*func TestStorage_GetKeys(t *testing.T) {
	type fields struct {
		mx   *sync.RWMutex
		data map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				mx:   tt.fields.mx,
				data: tt.fields.data,
			}
			if got := s.GetKeys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_GetListElement(t *testing.T) {
	type fields struct {
		mx   *sync.RWMutex
		data map[string]interface{}
	}
	type args struct {
		key   string
		index int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
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

func TestStorage_PutOrUpdateDictionary(t *testing.T) {
	type fields struct {
		mx   *sync.RWMutex
		data map[string]interface{}
	}
	type args struct {
		key   string
		value map[string]string
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantPreviousVal map[string]string
		wantIsUpdated   bool
	}{
		// TODO: Add test cases.
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
	type fields struct {
		mx   *sync.RWMutex
		data map[string]interface{}
	}
	type args struct {
		key   string
		value []string
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantPreviousVal []string
		wantIsUpdated   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{
				mx:   tt.fields.mx,
				data: tt.fields.data,
			}
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
	type fields struct {
		mx   *sync.RWMutex
		data map[string]interface{}
	}
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		wantPreviousVal string
		wantIsUpdated   bool
	}{
		// TODO: Add test cases.
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

func TestStorage_RemoveElement(t *testing.T) {
	type fields struct {
		mx   *sync.RWMutex
		data map[string]interface{}
	}
	type args struct {
		key string
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

func TestStorage_SetTTL(t *testing.T) {
	type fields struct {
		mx   *sync.RWMutex
		data map[string]interface{}
	}
	type args struct {
		key    string
		keyTTL int
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
}*/
