package config

import (
	"context"
	"fmt"
	"testing"

	"bou.ke/monkey"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/dao/xvkj"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/model"
)

func TestCheckAccessReq(t *testing.T) {
	monkey.Patch(xvkj.GetAccessConfig, func(scene string) (*model.AccessConfig, error) {
		if scene == "error_scene" {
			return nil, fmt.Errorf("(xvkj.GetAccessConfig error, scene:%+v", scene)
		}

		if scene == "appkey_enable" {
			return &model.AccessConfig{
				Appkey:       "appkey",
				AppkeyEnable: true,
				Scene:        scene,
			}, nil
		}

		return &model.AccessConfig{
			Appkey:       "appkey",
			AppkeyEnable: false,
			Scene:        scene,
		}, nil
	})
	type args struct {
		ctx    context.Context
		scene  string
		appkey string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "scene error",
			args: args{
				ctx:    context.Background(),
				scene:  "error_scene",
				appkey: "error_appkey",
			},
			wantErr: true,
		},
		{
			name: "AppkeyEnable false",
			args: args{
				ctx:    context.Background(),
				scene:  "appkey_unable",
				appkey: "",
			},
			wantErr: false,
		},
		{
			name: "AppkeyEnable true succ",
			args: args{
				ctx:    context.Background(),
				scene:  "appkey_enable",
				appkey: "appkey",
			},
			wantErr: false,
		},
		{
			name: "AppkeyEnable true error",
			args: args{
				ctx:    context.Background(),
				scene:  "appkey_enable",
				appkey: "appkey_error",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckAccessReq(tt.args.ctx, tt.args.scene, tt.args.appkey); (err != nil) != tt.wantErr {
				t.Errorf("CheckAccessReq() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
