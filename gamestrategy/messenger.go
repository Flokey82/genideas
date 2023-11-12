package gamestrategy

import "fmt"

type MsgReceiver interface {
	Receive(from int, message any)
}

type Messenger struct {
	Receivers []MsgReceiver
	ByID      map[int]MsgReceiver
}

func NewMessenger() *Messenger {
	return &Messenger{
		ByID: map[int]MsgReceiver{},
	}
}

func (m *Messenger) Register(id int, r MsgReceiver) {
	m.Receivers = append(m.Receivers, r)
	m.ByID[id] = r
}

func (m *Messenger) Broadcast(from int, message any) {
	m.Send(from, nil, message)
}

func (m *Messenger) Send(from int, to []int, message any) {
	if len(to) == 0 {
		// Broadcast to all receivers.
		// TODO: Only broadcast if the receiver is in range?
		for _, r := range m.Receivers {
			r.Receive(from, message)
		}
		return
	}

	// Send to specific receivers.
	for _, id := range to {
		if r, ok := m.ByID[id]; ok {
			r.Receive(from, message)
		}
	}
}

type Message struct {
	FromID int
	ToIDs  []int
}

func (m Message) String() string {
	return fmt.Sprintf("Message from %d to %v", m.FromID, m.ToIDs)
}

type MsgAttack struct {
	FromID  int   // Attacker
	ToID    int   // Defender
	AtCell  *Cell // Cell that was attacked
	Success bool  // Did the attack succeed?
}

type MsgBuild struct {
	FromID  int   // Player
	AtCell  *Cell // Cell that was built on
	Feature int64 // Feature that was built
}

type MsgExpand struct {
	FromID   int   // Player
	FromCell *Cell // Cell that was expanded from
	AtCell   *Cell // Cell that was expanded to
}

type MsgAbandon struct {
	FromID int   // Player
	AtCell *Cell // Cell that was abandoned
}
