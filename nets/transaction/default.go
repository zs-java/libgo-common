package transaction

import (
	"errors"
	"github.com/zs-java/libgo-common/nets/packet"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

const (
	Waiting    uint32 = 1
	NotWaiting uint32 = 0
)

const DefaultCallbackCmd = -1

// ================= DefaultTransaction ==================

type DefaultTransaction struct {
	Id             int32
	Packet         *packet.Packet
	callbackPacket *packet.Packet
	wg             *sync.WaitGroup
	// isWait value see Waiting, NotWaiting
	isWait uint32
}

func NewTransaction(id int32, cmd int32, data []byte) Transaction {
	return NewDefaultTransaction(id, cmd, data)
}

func NewDefaultTransaction(id int32, cmd int32, data []byte) *DefaultTransaction {
	p := packet.NewPacketTransaction(cmd, id, data)
	var wg sync.WaitGroup
	wg.Add(1)
	return &DefaultTransaction{
		Id:     id,
		Packet: p,
		wg:     &wg,
	}
}

func (t *DefaultTransaction) GetId() int32 {
	return t.Id
}

func (t *DefaultTransaction) GetPacket() *packet.Packet {
	return t.Packet
}

func (t *DefaultTransaction) GetCallbackPacket() *packet.Packet {
	return t.callbackPacket
}

func (t *DefaultTransaction) Wait() error {
	if !atomic.CompareAndSwapUint32(&t.isWait, NotWaiting, Waiting) {
		return errors.New("WaitGroup is always wait")
	}
	t.isWait = Waiting
	t.wg.Wait()
	return nil
}

func (t *DefaultTransaction) ThenCallback(cb func(*packet.Packet)) {
	go func() {
		_ = t.Wait()
		cb(t.GetCallbackPacket())
	}()
}

func (t *DefaultTransaction) done(callbackPacket *packet.Packet) {
	t.callbackPacket = callbackPacket
	t.wg.Done()
}

// ================= DefaultManager ==================

type DefaultManager struct {
	CallbackCmd    int32
	transactionMap map[int32]*DefaultTransaction
	rd             *rand.Rand
	lock           sync.Mutex
}

func NewManager() Manager {
	return NewManagerCallbackCmd(DefaultCallbackCmd)
}

func NewManagerCallbackCmd(callbackCmd int32) Manager {
	return &DefaultManager{
		CallbackCmd:    callbackCmd,
		transactionMap: make(map[int32]*DefaultTransaction),
		rd:             rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (m *DefaultManager) CreateTransaction(cmd int32, data []byte) Transaction {
	m.lock.Lock()
	defer m.lock.Unlock()

	id := m.generateTransactionId()
	ts := NewDefaultTransaction(id, cmd, data)
	m.transactionMap[id] = ts
	return ts
}

func (m *DefaultManager) CreateCallbackPacket(transactionId int32, data []byte) *packet.Packet {
	return packet.NewPacketTransaction(m.CallbackCmd, transactionId, data)
}

func (m *DefaultManager) IsCallbackPacket(p *packet.Packet) bool {
	return m.CallbackCmd == p.Cmd
}

func (m *DefaultManager) DoneTransaction(p *packet.Packet) error {
	if !m.IsCallbackPacket(p) {
		return errors.New("packet is not callback cmd")
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	id := p.TransactionId
	ts, ok := m.transactionMap[id]
	if !ok {
		log.Printf("callback transactionId %d not found!\n", id)
		return nil
	}
	ts.done(p)
	delete(m.transactionMap, id)
	return nil
}

func (m *DefaultManager) generateTransactionId() int32 {
	for {
		id := m.rd.Int31()
		if _, ok := m.transactionMap[id]; !ok {
			return id
		}
	}
}
