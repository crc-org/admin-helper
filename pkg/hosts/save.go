package hosts

import (
	"errors"
)

func (h *Hosts) saveHostsFile() error {
	file, err := OpenHostsFileAndDropPrivileges(h.File.Config.FilePath)
	if errors.Is(err, ErrPrivilegeDropUnsupported) {
		return h.File.SaveHostsFile()
	}
	if err != nil {
		return err
	}
	defer file.Close()

	// Truncate the file to 0 bytes to ensure it's empty before writing.
	if err := file.Truncate(0); err != nil {
		return err
	}
	// Seek to the beginning of the file to overwrite the existing contents.
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}
	dataBytes := []byte(h.File.RenderHostsFile())
	_, err = file.Write(dataBytes)
	if err != nil {
		return err
	}
	return file.Sync()
}
