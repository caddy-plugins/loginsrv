package htpasswd

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/caddy-plugins/loginsrv/login"
	. "github.com/stretchr/testify/assert"
)

func TestSetupOneFile(t *testing.T) {
	p, exist := login.GetProvider(ProviderName)
	True(t, exist)
	NotNil(t, p)

	files := writeTmpfile(testfile)
	backend, err := p(map[string]string{
		"file": files[0],
	})

	NoError(t, err)
	Equal(t,
		[]File{File{files[0], modTime(files[0])}},
		backend.(*Backend).auth.filenames)
}

func TestSetupTwoFiles(t *testing.T) {
	p, exist := login.GetProvider(ProviderName)
	True(t, exist)
	NotNil(t, p)

	filenames := writeTmpfile(testfile, testfile)

	var morphed []File
	for _, curFile := range filenames {
		morphed = append(morphed, File{curFile, modTime(curFile)})
	}
	backend, err := p(map[string]string{
		"file": strings.Join(filenames, ";"),
	})

	NoError(t, err)
	Equal(t,
		morphed,
		backend.(*Backend).auth.filenames)
}

func TestSetup_Error(t *testing.T) {
	p, exist := login.GetProvider(ProviderName)
	True(t, exist)
	NotNil(t, p)

	_, err := p(map[string]string{})
	Error(t, err)
}

func TestSimpleBackend_Authenticate(t *testing.T) {
	backend, err := NewBackend(writeTmpfile(testfile))
	NoError(t, err)

	authenticated, userInfo, err := backend.Authenticate("bob-bcrypt", "secret")
	True(t, authenticated)
	Equal(t, "bob-bcrypt", userInfo.Subject)
	NoError(t, err)

	authenticated, userInfo, err = backend.Authenticate("bob-bcrypt", "fooo")
	False(t, authenticated)
	Equal(t, "", userInfo.Subject)
	NoError(t, err)

	authenticated, userInfo, err = backend.Authenticate("", "")
	False(t, authenticated)
	Equal(t, "", userInfo.Subject)
	NoError(t, err)
}

func modTime(f string) time.Time {
	fileInfo, err := os.Stat(f)
	if err != nil {
		panic(err)
	}
	return fileInfo.ModTime()
}
