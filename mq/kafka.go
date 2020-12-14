package chxlib

import (
	"errors"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"sync"
	"time"
)

var kafkaLock sync.Mutex
var kafkaMap map[string]*Kafka = nil

type Kafka struct {
	Host          string
	Consumer      sarama.Consumer
	Producer      sarama.AsyncProducer
	client        sarama.Client
	offsetMgr     sarama.OffsetManager
	partitionMgrs map[string]map[int32]sarama.PartitionOffsetManager
	partitionComs map[string]map[int32]sarama.PartitionConsumer
}

func Kaf(mod string) *Kafka {
	if nil == kafkaMap {
		logs.Error("xdb.Kaf kafkaMap is nil")
		return nil
	}
	kafkaLock.Lock()
	defer kafkaLock.Unlock()
	conn := kafkaMap[mod]
	if nil == conn {
		conn = NewKafka(mod)
		kafkaMap[mod] = conn
	}
	return conn
}

func NewKafka(mod string) *Kafka {
	k := new(Kafka)
	k.Host = getConfString(mod, "host")
	offset := getConfBool(mod, "offset")
	consumer := getConfString(mod, "consumer")

	var err error
	config := sarama.NewConfig()
	config.Version = sarama.V1_1_1_0
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = time.Second
	k.client, err = sarama.NewClient([]string{k.Host}, config)
	if err != nil {
		logs.Error(err.Error())
		return nil
	}
	k.Consumer, err = sarama.NewConsumerFromClient(k.client)
	if nil != err {
		logs.Error(err.Error())
		return nil
	}
	k.Producer, err = sarama.NewAsyncProducerFromClient(k.client)
	if nil != err {
		logs.Error(err.Error())
		return nil
	}
	k.offsetMgr, err = sarama.NewOffsetManagerFromClient(consumer, k.client)
	if nil != err {
		logs.Error(err.Error())
		return nil
	}

	topics, err := k.Consumer.Topics()
	if nil != err {
		logs.Error(err.Error())
		return nil
	}
	k.partitionMgrs = make(map[string]map[int32]sarama.PartitionOffsetManager)
	k.partitionComs = make(map[string]map[int32]sarama.PartitionConsumer)
	for _, topic := range topics {
		k.partitionMgrs[topic] = make(map[int32]sarama.PartitionOffsetManager)
		k.partitionComs[topic] = make(map[int32]sarama.PartitionConsumer)
		pids, err := k.Consumer.Partitions(topic)
		if nil != err {
			logs.Error(err.Error())
			return nil
		}
		for i := 0; i < len(pids); i++ {
			pom, err := k.offsetMgr.ManagePartition(topic, pids[i])
			if nil != err {
				logs.Error("k.offsetMgr.ManagePartition error %s topic=%s pid=%d", err.Error(), topic, i)
				return nil
			}
			k.partitionMgrs[topic][pids[i]] = pom
			next := sarama.OffsetNewest
			if offset {
				next, _ = pom.NextOffset()
			}
			po, err := k.Consumer.ConsumePartition(topic, pids[i], next)
			if nil != err {
				logs.Error("k.Consumer.ConsumePartition error %s topic=%s pid=%d offset=%d", err.Error(), topic, i, next)
				return nil
			}
			k.partitionComs[topic][pids[i]] = po
		}
	}
	logs.Info("connect kafka ok host:%s offset:%v", k.Host, offset)
	return k
}

//获取当前topic下的所有partitions
func (k *Kafka) Partitions(topic string) []int32 {
	ids, err := k.Consumer.Partitions(topic)
	if nil != err {
		return []int32{}
	}
	return ids
}

//offset提高
func (k *Kafka) MarkOffset(topic string, pid int32, offset int64) {
	if nil != k.partitionMgrs[topic] {
		if nil != k.partitionMgrs[topic][pid] {
			k.partitionMgrs[topic][pid].MarkOffset(offset, fmt.Sprintf("partition:%d", pid))
		}
	}
}

//写入一个消息
func (k *Kafka) PutJob(topic, msg string) bool {
	select {
	case k.Producer.Input() <- &sarama.ProducerMessage{Topic: topic, Key: nil, Value: sarama.StringEncoder(msg)}:
		return true
	case err := <-k.Producer.Errors():
		logs.Error("Failed to produce message", err)
	}
	return false
}

//获取一个消息，阻塞线程
func (k *Kafka) Reserve(topic string, pid int32) (*sarama.ConsumerMessage, error) {
	if nil == k.partitionComs[topic] {
		return nil, errors.New("k.partitionComs[topic] is nil")
	}
	pcm := k.partitionComs[topic][pid]
	if nil == pcm {
		return nil, errors.New("k.partitionComs[topic][pid] is nil")
	}
	select {
	case msg := <-pcm.Messages():
		return msg, nil
	case err := <-pcm.Errors():
		return nil, err
	}
}

//释放所有资源
func (k *Kafka) Close() {
	for _, v := range k.partitionMgrs {
		for _, pom := range v {
			if err := pom.Close(); nil != err {
				logs.Error(err.Error())
			}
		}
	}
	for _, v := range k.partitionComs {
		for _, pcm := range v {
			if err := pcm.Close(); nil != err {
				logs.Error(err.Error())
			}
		}
	}
	if err := k.Consumer.Close(); nil != err {
		logs.Error(err.Error())
	}
	if err := k.offsetMgr.Close(); nil != err {
		logs.Error(err.Error())
	}
	if err := k.client.Close(); nil != err {
		logs.Error(err.Error())
	}
}
