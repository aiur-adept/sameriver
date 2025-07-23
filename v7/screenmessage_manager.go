package sameriver

type ScreenMessage interface {
	Position() [2]int
	Text() string
	Update(dt_ms int)
	IsActive() bool
}

type FixedScreenMessage struct {
	// the message to display
	Msg []string
	// how many lines of text to show before needing the
	// player to press space
	Lines int
	// age in milliseconds (used to scroll the text (if applicable)
	// and to set inactive)
	Age float64
}

type FloatingScreenMessage struct {
	// the message to display
	Msg string
	// the top-left corner of the box, where (0, 0) is
	// the bottom-left corner of the screen
	Position [2]int
	// how long the message should float for (in milliseconds)
	Duration float64
	// used to time the disappearance of the message
	Age float64
}

// responsible for spawning screen message entities
// managing their lifecycles, and destroying their resources
// when needed
type ScreenMessageManager struct {
	IDGen    IDGenerator
	messages map[int]ScreenMessage
}

func NewScreenMessageManager() *ScreenMessageManager {
	// arbitrary, can be tuned? Will grow
	capacity := 4
	return &ScreenMessageManager{
		IDGen:    NewIDGenerator(),
		messages: make(map[int]ScreenMessage, capacity),
	}
}

func (s *ScreenMessageManager) Update(dt_ms int) {
	for _, msg := range s.messages {
		msg.Update(dt_ms)
	}
}

func (s *ScreenMessageManager) Create(msg ScreenMessage) {
	s.messages[s.IDGen.Next()] = msg
}

func (s *ScreenMessageManager) Destroy(id int) {
	delete(s.messages, id)
}

func (s *ScreenMessageManager) Render() {
	for id, msg := range s.messages {
		if !msg.IsActive() {
			delete(s.messages, id)
		}
	}
}
