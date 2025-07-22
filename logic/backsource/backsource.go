package backsource

import (
	"context"
	"encoding/json"

	"git.code.oa.com/trpc-go/trpc-database/kafka"
	"git.code.oa.com/trpc-go/trpc-go/log"
	videoTimelineMsg "git.code.oa.com/trpcprotocol/component_plat/video_timeline_timeline_msg"
	"git.code.oa.com/v/main_logic/feeds/trpc_timeline_main_service_v2/model"
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
