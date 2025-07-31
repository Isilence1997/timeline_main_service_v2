package config

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"bou.ke/monkey"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/common/constant"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/dao/xvkj"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/model"
)

func TestGetFeedSceneConfig(t *testing.T) {
	monkey.Patch(xvkj.GetSceneConfigConfig, func(scene string) (*model.SceneConfig, error) {
		if scene == "error_scene" {
			return nil, fmt.Errorf("xvkj.GetSceneConfigConfig error,scene:%s", scene)
		}
		if scene == "scene_AvatarList" {
			return &model.SceneConfig{
				Scene:         scene,
				DiffusionType: constant.Read,
				SceneType:     constant.AvatarList,
			}, nil
		}
		if scene == "scene_Write" {
			return &model.SceneConfig{
				Scene:         scene,
				DiffusionType: constant.Write,
				SceneType:     constant.Feed,
			}, nil
		}
		return &model.SceneConfig{
			Scene:         scene,
			DiffusionType: constant.Read,
			SceneType:     constant.Feed,
			PageSize:      0,
			ReadConfigID:  "",
			WriteConfigID: "",
		}, nil
	})
	defer monkey.UnpatchAll()
	type args struct {
		ctx   context.Context
		scene string
	}
	tests := []struct {
		name    string
		args    args
		want    *model.SceneConfig
		wantErr bool
	}{
		{
			name: "scene error",
			args: args{
				ctx:   context.Background(),
				scene: "error_scene",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "SceneType error",
			args: args{
				ctx:   context.Background(),
				scene: "scene_AvatarList",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "DiffusionType error",
			args: args{
				ctx:   context.Background(),
				scene: "scene_Write",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "succ",
			args: args{
				ctx:   context.Background(),
				scene: "scene_succ",
			},
			want: &model.SceneConfig{
				Scene:         "scene_succ",
				DiffusionType: constant.Read,
				SceneType:     constant.Feed,
				PageSize:      0,
				ReadConfigID:  "",
				WriteConfigID: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFeedSceneConfig(tt.args.ctx, tt.args.scene)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFeedSceneConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFeedSceneConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAvatarListSceneConfig(t *testing.T) {
	monkey.Patch(xvkj.GetSceneConfigConfig, func(scene string) (*model.SceneConfig, error) {
		if scene == "error_scene" {
			return nil, fmt.Errorf("xvkj.GetSceneConfigConfig error,scene:%s", scene)
		}
		if scene == "scene_Feed" {
			return &model.SceneConfig{
				Scene:         scene,
				DiffusionType: constant.Read,
				SceneType:     constant.Feed,
			}, nil
		}
		if scene == "scene_Write" {
			return &model.SceneConfig{
				Scene:         scene,
				DiffusionType: constant.Write,
				SceneType:     constant.AvatarList,
			}, nil
		}
		return &model.SceneConfig{
			Scene:         scene,
			DiffusionType: constant.Read,
			SceneType:     constant.AvatarList,
			PageSize:      0,
			ReadConfigID:  "",
			WriteConfigID: "",
		}, nil
	})
	defer monkey.UnpatchAll()

	type args struct {
		ctx   context.Context
		scene string
	}
	tests := []struct {
		name    string
		args    args
		want    *model.SceneConfig
		wantErr bool
	}{
		{
			name: "scene error",
			args: args{
				ctx:   context.Background(),
				scene: "error_scene",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "SceneType error",
			args: args{
				ctx:   context.Background(),
				scene: "scene_Feed",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "DiffusionType error",
			args: args{
				ctx:   context.Background(),
				scene: "scene_Write",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "succ",
			args: args{
				ctx:   context.Background(),
				scene: "scene_succ",
			},
			want: &model.SceneConfig{
				Scene:         "scene_succ",
				DiffusionType: constant.Read,
				SceneType:     constant.AvatarList,
				PageSize:      0,
				ReadConfigID:  "",
				WriteConfigID: "",
			},
			wantErr: false,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetAvatarListSceneConfig(tt.args.ctx, tt.args.scene)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAvatarListSceneConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAvatarListSceneConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
