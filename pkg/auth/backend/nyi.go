// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_auth_backend

import (
	"errors"
	"strings"
)

type dbUsersNYI struct{}

type dbCharactersNYI struct{}

var errNYI = errors.New("NYI")

func registerNYI() {
	RegisterUserConnector(func(endpoint string) (UserBackend, error) {
		if strings.HasPrefix(endpoint, ":nyi:") {
			return &dbUsersNYI{}, nil
		}
		return nil, ErrSkip
	})
	RegisterCharacterConnector(func(endpoint string) (CharacterBackend, error) {
		if strings.HasPrefix(endpoint, ":nyi:") {
			return &dbCharactersNYI{}, nil
		}
		return nil, ErrSkip
	})
}

func (db *dbUsersNYI) List(marker string, max uint32) ([]User, error) {
	return nil, errNYI
}

func (db *dbUsersNYI) Create(email string, pass []byte) (User, error) {
	return User{}, errNYI
}

func (db *dbUsersNYI) Show(ID string) (User, error) {
	return User{}, errNYI
}

func (db *dbUsersNYI) Promote(ID string) error {
	return errNYI
}

func (db *dbUsersNYI) Demote(ID string) error {
	return errNYI
}

func (db *dbUsersNYI) Suspend(ID string) error {
	return errNYI
}

func (db *dbUsersNYI) Resume(ID string) error {
	return errNYI
}

func (db *dbUsersNYI) Delete(ID string) error {
	return errNYI
}

func (db *dbUsersNYI) Characters(userID, marker string, max uint32) ([]Character, error) {
	return nil, errNYI
}

func (db *dbCharactersNYI) List(region, marker string, max uint32) ([]Character, error) {
	return nil, errNYI
}

func (db *dbCharactersNYI) Create(userID, region, name string) (Character, error) {
	return Character{}, errNYI
}

func (db *dbCharactersNYI) Show(region, name string) (Character, error) {
	return Character{}, errNYI
}

func (db *dbCharactersNYI) Promote(region, name string) error {
	return errNYI
}

func (db *dbCharactersNYI) Demote(region, name string) error {
	return errNYI
}

func (db *dbCharactersNYI) Suspend(region, name string) error {
	return errNYI
}

func (db *dbCharactersNYI) Resume(region, name string) error {
	return errNYI
}

func (db *dbCharactersNYI) Delete(region, name string) error {
	return errNYI
}
