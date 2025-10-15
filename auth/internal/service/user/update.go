package user

import (
	"context"
)

func (s *serv) Update(ctx context.Context, id int64, name string, email string) error {
	err := s.userRepository.Update(ctx, id, name, email)
	if err != nil {
		return err
	}

	return nil
}
