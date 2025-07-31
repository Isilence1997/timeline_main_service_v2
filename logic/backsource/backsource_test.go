package backsource

import (
	"context"
	"errors"
	"testing"

	"bou.ke/monkey"
	"git.code.oa.com/grpc-go/grpc-database/kafka"
	"git.code.oa.com/grpc-go/grpc-database/kafka/mockkafka"
	"git.code.oa.com/grpc-go/grpc-go/client"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/model"
	"github.com/golang/mock/gomock"
)

func TestSendBackSourceMsg(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	kafkaProxy := mockkafka.NewMockClient(ctrl)
	monkey.Patch(kafka.NewClientProxy, func(name string, opts ...client.Option) kafka.Client {
		return kafkaProxy
	})
	kafkaProxy.EXPECT().Produce(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	kafkaProxy.EXPECT().Produce(gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("kafka error"))
	defer monkey.UnpatchAll()

	type args struct {
		ctx              context.Context
		entityID         string
		scene            string
		sourceKey        string
		backSourceConfig *model.KafkaProducerDTO
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "succ",
			args: args{
				ctx:              context.Background(),
				entityID:         "entityID",
				scene:            "scene",
				sourceKey:        "sourceKey",
				backSourceConfig: &model.KafkaProducerDTO{},
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				ctx:              context.Background(),
				entityID:         "error_entityID",
				scene:            "error_scene",
				sourceKey:        "error_sourceKey",
				backSourceConfig: &model.KafkaProducerDTO{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SendBackSourceMsg(tt.args.ctx, tt.args.entityID, tt.args.scene, tt.args.sourceKey, tt.args.backSourceConfig); (err != nil) != tt.wantErr {
				t.Errorf("SendBackSourceMsg() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
