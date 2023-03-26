package reqres

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/Fajar-Islami/go-simple-user-crud/internal/infrastructure/container"
	"github.com/Fajar-Islami/go-simple-user-crud/internal/pkg/repositories/rest"
	"github.com/Fajar-Islami/go-simple-user-crud/internal/utils"
)

func TestNewListUsersRepository(t *testing.T) {
	type args struct {
		URI  string
		Opts rest.Opts
	}
	tests := []struct {
		name string
		args args
		want ReqResAPI
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewReqResAPI(tt.args.URI, tt.args.Opts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewReqResAPI() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListUsersRepositoryImpl_GetListUser(t *testing.T) {
	type args struct {
		params ReqListUser
	}

	type structTest struct {
		name       string
		lur        *ReqResAPIImpl
		args       args
		wantResArr int
		wantErr    bool
	}

	cont := container.InitContainer("log", "reqresapi")

	tests := []structTest{{
		name: "Test 1",
		lur: &ReqResAPIImpl{
			URI: cont.ReqResAPI.URL,
			Opts: rest.Opts{
				Timeout:     time.Duration(cont.ReqResAPI.TimeOut * int(time.Minute)),
				Logger:      *cont.Logger,
				IsDebugging: cont.ReqResAPI.Debugging,
			},
		},
		args: args{
			params: ReqListUser{
				PerPage: 12,
				Page:    1,
			},
		},
		wantResArr: 12,
		wantErr:    false,
	}}
	for _, tt := range tests {
		ctx := context.Background()
		CtxTimeout := utils.EnvInt("apps_timeout")
		cntx, cancel := context.WithTimeout(ctx, time.Duration(time.Duration(CtxTimeout*int(time.Second))))
		defer cancel()

		t.Run(tt.name, func(t *testing.T) {

			gotRes, err := tt.lur.GetListUser(cntx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListUsersRepositoryImpl.GetListUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(gotRes.Data) != tt.wantResArr {
				t.Error("========ERRRO====")
				t.Errorf("Expected data got %d result get %d", len(gotRes.Data), tt.wantResArr)
			}

		})
	}
}
