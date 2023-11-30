package toolkit

import "os"

func (t *Tools) CreateDir(dir string) error {

	const mode = 0755

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, mode)
		if err != nil {
			return err
		}
	}
	return nil
}
