package application

import "backend/internal/users/domain"

type mockUserRepository struct {
	users          map[string]*domain.User
	byID           map[uint]*domain.User
	nextID         uint
	createErr      error
	updateErr      error
	findByEmailErr error
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:  make(map[string]*domain.User),
		byID:   make(map[uint]*domain.User),
		nextID: 1,
	}
}

func (m *mockUserRepository) Create(user *domain.User) error {
	if m.createErr != nil {
		return m.createErr
	}
	user.ID = m.nextID
	m.nextID++
	m.users[user.Email] = user
	m.byID[user.ID] = user
	return nil
}

func (m *mockUserRepository) FindByEmail(email string) (*domain.User, error) {
	if m.findByEmailErr != nil {
		return nil, m.findByEmailErr
	}
	if user, ok := m.users[email]; ok {
		return user, nil
	}
	return nil, domain.ErrUserNotFound
}

func (m *mockUserRepository) FindAll() ([]domain.User, error) {
	users := make([]domain.User, 0, len(m.users))
	for _, u := range m.users {
		users = append(users, *u)
	}
	return users, nil
}

func (m *mockUserRepository) FindAllByRole(role domain.UserRole) ([]domain.User, error) {
	users := make([]domain.User, 0)
	for _, u := range m.users {
		if u.GlobalRole == role {
			users = append(users, *u)
		}
	}
	return users, nil
}

func (m *mockUserRepository) FindByID(id uint) (*domain.User, error) {
	if user, ok := m.byID[id]; ok {
		return user, nil
	}
	return nil, domain.ErrUserNotFound
}

func (m *mockUserRepository) Update(user *domain.User) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	if _, ok := m.byID[user.ID]; !ok {
		return domain.ErrUserNotFound
	}

	for email, existingUser := range m.users {
		if existingUser.ID == user.ID && email != user.Email {
			delete(m.users, email)
		}
	}

	m.byID[user.ID] = user
	m.users[user.Email] = user
	return nil
}
