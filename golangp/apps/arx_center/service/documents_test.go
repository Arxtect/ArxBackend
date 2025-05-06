package service_test

import (
	"github.com/arxtect/ArxBackend/golangp/apps/arx_center/models"
	"github.com/arxtect/ArxBackend/golangp/apps/arx_center/service"
	"github.com/arxtect/ArxBackend/golangp/common/initializers"
	"reflect"
	"testing"
)

func Test_GetTagsByName(t *testing.T) {
	initializers.TestConnectDb()
	tags, err := service.GetTagsByName([]string{"天文", "机械"})
	if err != nil {
		t.Error(err)
	}
	t.Log(tags)
}

func TestGetTagsByName(t *testing.T) {
	type args struct {
		names []string
	}
	tests := []struct {
		name    string
		args    args
		want    []models.Tag
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				names: []string{"天文", "机械"},
			},
			want: []models.Tag{
				{
					Name: "天文",
				},
				{
					Name: "机械",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.GetTagsByName(tt.args.names)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTagsByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetTagsByName() got = %v, want %v", got, tt.want)
			}
		})
	}
}
