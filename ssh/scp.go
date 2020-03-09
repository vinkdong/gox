package ssh

import (
	"github.com/tmc/scp"
	"github.com/vinkdong/gox/log"
	"io/ioutil"
	"os"
)

func (s *SSH) WriteBytesToRemoteHost(data []byte, filename string) error {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "gox.*")
	if err != nil {
		return err
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			log.Error(err)
		}
	}()

	if _, err = tmpFile.Write(data); err != nil {
		return err
	}

	session, err := s.newSession()
	if err != nil {
		return err
	}

	session.Stderr = s.Stderr
	session.Stdout = s.Stdout
	return scp.CopyPath(tmpFile.Name(), filename, session)
}
