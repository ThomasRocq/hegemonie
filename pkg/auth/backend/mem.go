// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_auth_backend

import (
	"github.com/google/uuid"
	"strings"
)

var (
	usersByID   = make(SetOfUsers, 0)
	usersByMail = make(map[string]string)
	charsByName = make(SetOfCharacters, 0)
)

type dbUsersMem struct{}
type dbCharactersMem struct{}

func registerMem() {
	RegisterUserConnector(func(endpoint string) (UserBackend, error) {
		if strings.HasPrefix(endpoint, ":mem:") {
			return &dbUsersMem{}, nil
		}
		return nil, ErrSkip
	})
	RegisterCharacterConnector(func(endpoint string) (CharacterBackend, error) {
		if strings.HasPrefix(endpoint, ":mem:") {
			return &dbCharactersMem{}, nil
		}
		return nil, ErrSkip
	})
}

func (db *dbUsersMem) List(marker string, max uint32) ([]User, error) {
	tab := make([]User, 0)
	for _, u := range usersByID.Slice(marker, max) {
		tab = append(tab, *u)
	}
	return tab, nil
}

func (db *dbUsersMem) Create(email string, pass []byte) (User, error) {
	if usersByMail[email] != "" {
		return User{}, ErrExists
	}
	u := User{
		ID:       uuid.New().String(),
		Name:     "NOT-SET",
		Email:    email,
		Password: pass,
		State:    UserStateActive,
	}
	usersByID.Add(&u)
	usersByMail[u.Email] = u.ID
	return u, nil
}

func (db *dbUsersMem) Show(ID string) (User, error) {
	u := usersByID.Get(ID)
	if u == nil {
		return User{}, ErrNotFound
	}
	return *u, nil
}

func (db *dbUsersMem) Promote(ID string) error {
	u := usersByID.Get(ID)
	if u == nil {
		return ErrNotFound
	}
	u.State = UserStateAdmin
	return nil
}

func (db *dbUsersMem) Demote(ID string) error {
	u := usersByID.Get(ID)
	if u == nil {
		return ErrNotFound
	}
	u.State = UserStateActive
	return nil
}

func (db *dbUsersMem) Suspend(ID string) error {
	u := usersByID.Get(ID)
	if u == nil {
		return ErrNotFound
	}
	u.State = UserStateSuspended
	return nil
}

func (db *dbUsersMem) Resume(ID string) error {
	u := usersByID.Get(ID)
	if u == nil {
		return ErrNotFound
	}
	u.State = UserStateActive
	return nil
}

func (db *dbUsersMem) Delete(ID string) error {
	u := usersByID.Get(ID)
	if u == nil {
		return ErrNotFound
	}
	u.State = UserStateDeleted
	return nil
}

func (db *dbUsersMem) Characters(userID, marker string, max uint32) ([]Character, error) {
	return nil, errNYI
}

func (db *dbCharactersMem) List(region, marker string, max uint32) ([]Character, error) {
	return nil, errNYI
}

func (db *dbCharactersMem) Create(userID, region, name string) (Character, error) {
	if !usersByID.Has(userID) {
		return Character{}, ErrNotFound
	}
	if charsByName.Has(region, name) {
		return Character{}, ErrExists
	}
	c := &Character{
		Name:   name,
		Region: region,
		State:  CharacterStateActive,
	}
	charsByName.Add(c)
	return *c, nil
}

func (db *dbCharactersMem) Show(region, name string) (Character, error) {
	return Character{}, errNYI
}

func (db *dbCharactersMem) Promote(region, name string) error {
	return errNYI
}

func (db *dbCharactersMem) Demote(region, name string) error {
	return errNYI
}

func (db *dbCharactersMem) Suspend(region, name string) error {
	return errNYI
}

func (db *dbCharactersMem) Resume(region, name string) error {
	return errNYI
}

func (db *dbCharactersMem) Delete(region, name string) error {
	return errNYI
}

//go:generate go run github.com/jfsmig/hegemonie/cmd/gen-set ./mem_auto.go hegemonie_auth_backend:SetOfUsers:*User           ID:string
//go:generate go run github.com/jfsmig/hegemonie/cmd/gen-set ./mem_auto.go hegemonie_auth_backend:SetOfCharacters:*Character Region:string Name:string
