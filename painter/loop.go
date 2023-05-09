package painter

import (
	"image"

	"golang.org/x/exp/shiny/screen"
)

// Receiver отримує текстуру, яка була підготовлена в результаті виконання команд у циелі подій.
type Receiver interface {
	Update(t screen.Texture)
}

// Loop реалізує цикл подій для формування текстури отриманої через виконання операцій отриманих з внутрішньої черги.
type Loop struct {
	Receiver Receiver

	next screen.Texture // текстура, яка зараз формується
	prev screen.Texture // текстура, яка була відправленя останнього разу у Receiver

	MsgQueue messageQueue
}

var size = image.Pt(800, 800)

// Start запускає цикл подій. Цей метод потрібно запустити до того, як викликати на ньому будь-які інші методи.
func (l *Loop) Start(s screen.Screen) {
	l.next, _ = s.NewTexture(size)
	l.prev, _ = s.NewTexture(size)

	l.MsgQueue = messageQueue{} // ініціалізація черги подій

	go func() { //запуск рутини обробки повідомлень у черзі подій.
		for {
			op := l.MsgQueue.pull()
			if op == nil {
				continue
			}
			isUpdate := op.Do(l.next)
			if isUpdate {
				l.Receiver.Update(l.next)
				l.next, l.prev = l.prev, l.next
			}
		}
	}()
}

// Post додає нову операцію у внутрішню чергу.
func (l *Loop) Post(op Operation) {
	if op == nil {
		return
	}
	l.MsgQueue.push(op)
}

// StopAndWait сигналізує
func (l *Loop) StopAndWait() {

}

// черга повідомлень.
type messageQueue struct {
	Queue []Operation
}

func (mq *messageQueue) push(op Operation) {
	mq.Queue = append(mq.Queue, op)
}

func (mq *messageQueue) pull() Operation {
	if len(mq.Queue) == 0 {
		return nil
	}

	op := mq.Queue[0]
	mq.Queue = mq.Queue[1:]
	return op
}
