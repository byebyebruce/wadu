package edgetts

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func TestCommunicate_stream(t *testing.T) {
	tts := New()
	f, err := os.OpenFile("a.mp3", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	ctx := context.Background()
	err = tts.SynthesisFile(ctx, "你好啊", "a.mp3")
	t.Log(err)
}

func Test_listVoices(t *testing.T) {
	got, err := ListVoices()
	if err != nil {
		t.Error(err)
	}
	for _, voice := range got {
		t.Log(voice)
		fmt.Println(voice)
	}
}

func TestVoicesManager_find(t *testing.T) {
	type args struct {
		attributes Voice
	}
	vm := &VoicesManager{}
	vm.create(nil)
	tests := []struct {
		name string
		vm   *VoicesManager
		args args
		want []Voice
	}{
		// TODO: Add test cases.
		{
			name: "test-1",
			vm:   vm,
			args: args{
				attributes: Voice{
					Locale: "zh-CN",
				},
			},
			want: []Voice{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.vm.find(tt.args.attributes)
			if len(got) <= 0 {
				t.Errorf("ListVoices() wantErr %v", tt.want)
				return
			}
			t.Logf("ListVoices() got %v", got)
		})
	}
}

func Test_getHeadersAndData(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]string
		want1   []byte
		wantErr bool
	}{
		{
			name: "test-1",
			args: args{
				data: "X-Timestamp:2022-01-01\r\nContent-Type:application/json; charset=utf-8\r\nPath:speech.config\r\n\r\n{\"context\":{\"synthesis\":{\"audio\":{\"metadataoptions\":{\"sentenceBoundaryEnabled\":false,\"wordBoundaryEnabled\":true},\"outputFormat\":\"audio-24khz-48kbitrate-mono-mp3\"}}}}",
			},
			want:    map[string]string{},
			want1:   []byte{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := getHeadersAndData(tt.args.data)
			t.Logf("%v \n%v \n", got, got1)
			if (err != nil) != tt.wantErr {
				t.Errorf("getHeadersAndData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("getHeadersAndData() got = %v, want %v", got, tt.want)
			// }
			// if !reflect.DeepEqual(got1, tt.want1) {
			// 	t.Errorf("getHeadersAndData() got1 = %v, want %v", got1, tt.want1)
			// }
		})
	}
}
