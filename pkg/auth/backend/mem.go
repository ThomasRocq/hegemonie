// Copyright (C) 2018-2020 	Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_auth_backend

import (
	"github.com/google/uuid"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"strings"
)

var (
	usersByID   = make(SetOfUsers, 0)
	usersByMail = make(map[string]string)
	charsByName = make(SetOfCharacters, 0)
)

type dbUsersMem struct{}
type dbCharactersMem struct{}

type charName struct {
	region, name string
}

type userMem struct {
	User
	roles SetOfCharacterNames
}

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
		tab = append(tab, u.User)
	}
	return tab, nil
}

func (db *dbUsersMem) Create(email string, pass []byte) (User, error) {
	if usersByMail[email] != "" {
		return User{}, ErrExists
	}
	u := &userMem{User: User{
		ID:       uuid.New().String(),
		Name:     "NOT-SET",
		Email:    email,
		Password: pass,
		State:    UserStateActive,
	}}
	usersByID.Add(u)
	usersByMail[u.Email] = u.ID
	return u.User, nil
}

func (db *dbUsersMem) Show(ID string) (User, error) {
	u := usersByID.Get(ID)
	if u == nil {
		return User{}, ErrNotFound
	}
	return u.User, nil
}

func (db *dbUsersMem) GetByMail(email string) (User, error) {
	ID := usersByMail[email]
	if ID == "" {
		return User{}, ErrNotFound
	}
	return db.Show(ID)
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

func (db *dbUsersMem) Characters(ID string) ([]Character, error) {
	u := usersByID.Get(ID)
	if u == nil {
		return nil, ErrNotFound
	}
	tab := make([]Character, 0)
	for _, cn := range u.roles {
		c := charsByName.Get(cn.region, cn.name)
		if c == nil {
			utils.Logger.Warn().Str("user", ID).Str("cname", cn.name).Str("region", cn.region).Msg("ghost character")
		} else {
			tab = append(tab, *c)
		}
	}
	return tab, nil
}

func (db *dbCharactersMem) List(region, marker string, max uint32) ([]Character, error) {
	tab := make([]Character, 0)
	for _, c := range charsByName.Slice(region, marker, max) {
		tab = append(tab, *c)
	}
	return tab, nil
}

func (db *dbCharactersMem) Create(userID, region, name string) (Character, error) {
	u := usersByID.Get(userID)
	if u == nil {
		return Character{}, ErrNotFound
	}
	if charsByName.Has(region, name) {
		return Character{}, ErrExists
	}
	c := &Character{
		Region: region,
		Name:   name,
		User:   userID,
		State:  CharacterStateActive,
	}
	charsByName.Add(c)
	u.roles.Add(charName{region: region, name: name})
	return *c, nil
}

func (db *dbCharactersMem) Show(region, name string) (Character, error) {
	c := charsByName.Get(region, name)
	if c == nil {
		return Character{}, ErrNotFound
	}
	return *c, nil
}

func (db *dbCharactersMem) Promote(region, name string) error {
	c := charsByName.Get(region, name)
	if c == nil {
		return ErrNotFound
	}
	c.State = CharacterStateSuper
	return nil
}

func (db *dbCharactersMem) Demote(region, name string) error {
	c := charsByName.Get(region, name)
	if c == nil {
		return ErrNotFound
	}
	c.State = CharacterStateActive
	return nil
}

func (db *dbCharactersMem) Suspend(region, name string) error {
	c := charsByName.Get(region, name)
	if c == nil {
		return ErrNotFound
	}
	c.State = CharacterStateSuspended
	return nil
}

func (db *dbCharactersMem) Resume(region, name string) error {
	c := charsByName.Get(region, name)
	if c == nil {
		return ErrNotFound
	}
	c.State = CharacterStateActive
	return nil
}

func (db *dbCharactersMem) Delete(region, name string) error {
	c := charsByName.Get(region, name)
	if c == nil {
		return ErrNotFound
	}
	c.State = CharacterStateDeleted
	return nil
}

//go:generate go run github.com/jfsmig/hegemonie/cmd/gen-set ./mem_auto.go hegemonie_auth_backend:SetOfUsers:*userMem ID:string
//go:generate go run github.com/jfsmig/hegemonie/cmd/gen-set ./mem_auto.go hegemonie_auth_backend:SetOfCharacters:*Character Region:string Name:string
//go:generate go run github.com/jfsmig/hegemonie/cmd/gen-set ./mem_auto.go hegemonie_auth_backend:SetOfCharacterNames:charName region:string name:string
