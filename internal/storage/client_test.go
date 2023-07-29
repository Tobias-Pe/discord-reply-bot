package storage

import (
	"github.com/Tobias-Pe/discord-reply-bot/internal/models"
	"github.com/alicebob/miniredis/v2"
	"reflect"
	"testing"
)

func TestAddElement(t *testing.T) {
	server := miniredis.RunT(t)
	InitClient(server.Addr())

	type args struct {
		key   models.MessageMatch
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Normal adding of a key-value", args: struct {
			key   models.MessageMatch
			value string
		}{key: models.MessageMatch{
			Message:      "Init",
			IsExactMatch: false,
		}, value: "Hello World"}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddElement(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("AddElement() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCheckConnection(t *testing.T) {
	server := miniredis.RunT(t)
	InitClient(server.Addr())

	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "Check current connection to test redis", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckConnection(); (err != nil) != tt.wantErr {
				t.Errorf("CheckConnection() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	server := miniredis.RunT(t)
	InitClient(server.Addr())

	testKey, testValue1 := addElement(t, false, "a", "b")
	_, testValue2 := addElement(t, false, "a", "b2")
	_, _ = addElement(t, true, "a", "b")

	type args struct {
		key models.MessageMatch
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{name: "Get pretest added val", args: struct{ key models.MessageMatch }{key: testKey}, want: []string{testValue1, testValue2}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAll(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAll() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllKeys(t *testing.T) {
	server := miniredis.RunT(t)
	InitClient(server.Addr())

	testKey, _ := addElement(t, false, "a", "b")
	testKey2, _ := addElement(t, true, "a", "b")

	tests := []struct {
		name    string
		want    []models.MessageMatch
		wantErr bool
	}{
		{
			name:    "Get pretest added key",
			want:    []models.MessageMatch{testKey, testKey2},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAllKeys()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllKeys() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllKeys() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitClient(t *testing.T) {
	server := miniredis.RunT(t)

	type args struct {
		redisUrl string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Normal init", args: struct{ redisUrl string }{redisUrl: server.Addr()}, wantErr: false},
		{name: "Fail init", args: struct{ redisUrl string }{redisUrl: "foo"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitClient(tt.args.redisUrl)
			err := CheckConnection()
			if (err != nil) != tt.wantErr {
				t.Errorf("InitClient() error = %v", err)
				return
			}
		})
	}
}

func TestRemoveElement(t *testing.T) {
	server := miniredis.RunT(t)
	InitClient(server.Addr())

	testKey, testValue := addElement(t, false, "a", "b")

	type args struct {
		key   models.MessageMatch
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "Normal removal", args: struct {
			key   models.MessageMatch
			value string
		}{key: testKey, value: testValue}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RemoveElement(tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("RemoveElement() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func addElement(t *testing.T, isExact bool, message, reply string) (models.MessageMatch, string) {
	testKey := models.MessageMatch{Message: message, IsExactMatch: isExact}
	testValue := reply
	err := AddElement(testKey, testValue)
	if err != nil {
		t.Errorf("Setup of test didnt work: %v", err)
	}

	return testKey, testValue
}
