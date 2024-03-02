package rtccamgen

import (
	"container/list"
	"rtccam/rtccamlog"
	"sync"
)

type Generator interface {
	GenerateID() int64
	ReturnID(id int64)
}

func NewIDGenerator() *IDGenerator {
	idCounter := int64(0)
	return &IDGenerator{
		idCounter: &idCounter,
		queue:     list.New(),
	}
}

type IDGenerator struct {
	idCounter *int64
	queue     *list.List
	mutex     sync.Mutex
}

func (g *IDGenerator) GenerateID() int64 {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	if g.queue.Len() > 0 {
		element := g.queue.Front()
		g.queue.Remove(element)
		if id, ok := element.Value.(int64); ok {
			rtccamlog.Info().Any("Get in queue", id).Send()
			return id
		}
	}

	*g.idCounter++
	rtccamlog.Info().Any("New Generate Id", *g.idCounter).Send()
	return *g.idCounter
}

func (g *IDGenerator) ReturnID(id int64) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	rtccamlog.Info().Any("ReturnID", id).Send()
	g.queue.PushBack(id)
}
