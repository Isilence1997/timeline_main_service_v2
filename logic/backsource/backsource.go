package backsource

import (
	"context"
	"encoding/json"

	"git.code.oa.com/grpc-go/grpc-database/kafka"
	"git.code.oa.com/grpc-go/grpc-go/log"
	videoTimelineMsg "git.code.oa.com/grpcprotocol/component_plat/video_timeline_timeline_msg"
	"git.code.oa.com/v/main_logic/feeds/grpc_timeline_main_service_v2/model"
)

// SendBackSourceMsg 发送回源消息
func SendBackSourceMsg(ctx context.Context, entityID, scene, sourceKey string,
	backSourceConfig *model.KafkaProducerDTO) error {
	kafkaProxy := kafka.NewClientProxy(
		backSourceConfig.ServiceName, backSourceConfig.KafkaProducerInfo2Options()...)

	msg := videoTimelineMsg.BackSourceMsg{
		EntityId:  entityID,
		Scene:     scene,
		SourceKey: sourceKey,
	}

	value, _ := json.Marshal(msg)
	if err := kafkaProxy.Produce(ctx, []byte(entityID), value); err != nil {
		log.Errorf("sendBackSourceMsg error, entityID:%s, msg:%s, err:%v", entityID, string(value), err)
		return err
	}
	return nil
}
