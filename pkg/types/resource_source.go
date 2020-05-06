package types

import (
	"fmt"
)

func (s ResourceSource) Validate() error {
	if s.URL == "" {
		return fmt.Errorf("'url' is required in the source definition")
	}

	if s.Username == "" && s.APIToken == "" {
		return fmt.Errorf("'username' or 'api_token' are required in the source definition")
	}

	if s.Username != "" && s.Password == "" {
		return fmt.Errorf("'password' is required in the source definition when 'username' is present")
	}

	return nil
}
