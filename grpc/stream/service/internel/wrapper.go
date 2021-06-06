package internel

import (
	"context"
	"fmt"
	"github.com/hxzhouh/go_example/grpc/stream/stream"
	"go.uber.org/zap"
	"io"
	"log"
	"time"
)

type StreamServiceImpl struct{}
type messageResponseReader func(context.Context, int32) (*stream.StreamReply, error)
type messageResponseWriter func(reply *stream.StreamRequest) error
type BitStream stream.EchoStream_SayHelloStreamServer

var heartbeatTimeout = 60 * time.Second

type StreamWrapper struct {
	context   context.Context
	msg       stream.EchoStream_SayHelloStreamServer //
	msgReader messageResponseReader
	msgWriter messageResponseWriter
	sendChl   chan stream.StreamReply
	BitStream
	heartbeatTimeout    *time.Timer //超时定时器
	errHeartbeatTimeout bool
	cancelFn            context.CancelFunc // 出现错误 关闭整个 双向流
	logger              *zap.Logger
}

func newStreamWrapper(s BitStream) *StreamWrapper {
	heartbeatTimeout := time.NewTimer(heartbeatTimeout)
	ctx, cancel := context.WithCancel(s.Context())
	logger, _ := zap.NewDevelopment()
	logger.Info("create newStreamWrapper")
	return &StreamWrapper{
		BitStream:        s,
		heartbeatTimeout: heartbeatTimeout,
		sendChl:          make(chan stream.StreamReply, 10),
		cancelFn:         cancel,
		context:          ctx,
		logger:           logger,
	}
}

func (t *StreamWrapper) sender() {
	for {
		select {
		case value := <-t.sendChl: //等待数据发送
			err := t.Send(&value)
			if err != nil {
				t.logger.Error("Send message error", zap.Error(err))
				t.cancelFn()
				break
			}
		case <-t.context.Done():
			break
		}
	}
}

// 接收函数
func (t *StreamWrapper) recv() {
	type result struct {
		msg *stream.StreamRequest
		err error
	}
	recvBuf := make(chan *result, 1)
	go func(w *StreamWrapper, buf chan *result) {
		for {
			msg, err := t.Recv() //这里会阻塞
			if err != nil {
				t.logger.Warn("recv stream error", zap.Error(err))
				recvBuf <- &result{
					msg: msg,
					err: err,
				}
				break
			}
			// todo 优化
			select {
			case <-t.Context().Done():
				t.cancelFn() // 退出函数
				break
			default:
				recvBuf <- &result{
					msg: msg,
					err: err,
				}
			}
		}
	}(t, recvBuf)

	for {
		var (
			msg *stream.StreamRequest
			err error
		)
		select {
		case <-t.context.Done():
			err = t.context.Err()
		case ret := <-recvBuf: // 收到 客户端的消息，或者是心跳
			if ret.err != nil {
				log.Println(err)
				break
			}
			t.sendChl <- stream.StreamReply{Ok: true}
			err = ret.err
			msg = ret.msg
			if ok := t.heartbeatTimeout.Reset(heartbeatTimeout); !ok {
				err = fmt.Errorf("reset heartbeattime ticker error")
				break
			}
			log.Println(msg.Body.Name)
		case <-t.heartbeatTimeout.C:
			err = fmt.Errorf("heart beat timeout")

		}
		if err == io.EOF { // 客户端断开
			// 下线逻辑
			log.Println("客户端离线")
			break
		}
		if err == fmt.Errorf("heart beat timeout") {
			t.errHeartbeatTimeout = true
			break
		}
	}
}

func (t *StreamWrapper) start() error {
	t.context, t.cancelFn = context.WithCancel(t.Context())
	var err error
	go t.sender()
	go t.recv()
	err = connectMap.Add("test", t)
	if err != nil {
		return err
	}
	select {
	case <-t.context.Done(): // 等待结束
		err = t.context.Err()
		if t.errHeartbeatTimeout {
			err = fmt.Errorf("timeout")
		}
	}
	t.heartbeatTimeout.Stop()
	return err
}

func (s StreamServiceImpl) SayHelloStream(server stream.EchoStream_SayHelloStreamServer) error {
	streamWrapper := newStreamWrapper(server) // 新建一个双向流连接
	err := streamWrapper.start()              //阻塞在这里
	if err != nil {
		streamWrapper.logger.Warn("stream 断开", zap.Error(err))
	}
	return nil
}
