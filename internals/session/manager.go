package session

import "github.com/google/uuid"

type Manager struct {
	currentId string
}

func (m *Manager) CurrentId() string {
	if m.currentId == "" {
		m.currentId = uuid.NewString()
	}

	return m.currentId
}
