package types

import (
	"fmt"
)

func (s ResourceSource) Validate() error {
	if s.URL == "" {
		return fmt.Errorf("URL is required in the source definition")
	}

	if s.Username == "" {
		return fmt.Errorf("Username is required in the source definition")
	}

	if s.Password == "" {
		return fmt.Errorf("Password is required in the source definition")
	}

	return nil
}
