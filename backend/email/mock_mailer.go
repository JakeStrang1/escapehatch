package email

type MockMailer struct {
	SendParams SendParams
}

func NewMockMailer() *MockMailer {
	return &MockMailer{}
}

func (s *MockMailer) Send(params SendParams) error {
	s.SendParams = params
	return nil
}
