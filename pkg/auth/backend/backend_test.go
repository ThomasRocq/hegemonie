package hegemonie_auth_backend

import "testing"

func TestAuthBackendMem(t *testing.T) {
	var err error
	var users UserBackend
	var chars CharacterBackend

	users, err = ConnectUserBackend(":mem:")
	if err != nil {
		t.Fatal()
	}
	chars, err = ConnectCharacterBackend(":mem:")
	if err != nil {
		t.Fatal()
	}

	tab, err := users.List("", 10)
	if err != nil {
		t.Fatal()
	}
	if len(tab) != 0 {
		t.Fatal()
	}

	cTab, err := chars.List("", "", 10)
	if err != nil {
		t.Fatal()
	}
	if len(cTab) != 0 {
		t.Fatal()
	}
}
