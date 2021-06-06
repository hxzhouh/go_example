package internel

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var connectMap ConnectMap

type ConnectMap struct {
	connectMap map[string]*StreamWrapper
}

func init() {
	connectMap = ConnectMap{connectMap: make(map[string]*StreamWrapper)}
}

func (t *ConnectMap) Add(id string, value *StreamWrapper) error {
	if _, ok := t.connectMap[id]; ok {
		return status.Error(codes.AlreadyExists, "clientId is online")
	}
	t.connectMap[id] = value
	return nil
}

func (t *ConnectMap) Del(id string) error {
	if _, ok := t.connectMap[id]; !ok {
		return status.Error(codes.NotFound, "clientId is not online")
	}
	delete(t.connectMap, id)
	return nil
}

func (t *ConnectMap) Get(id string) (error, *StreamWrapper) {
	if _, ok := t.connectMap[id]; !ok {
		return status.Error(codes.NotFound, "clientId is not online"), nil
	}
	return nil, t.connectMap[id]
}
