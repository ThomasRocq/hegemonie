// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_auth_backend

import (
	"bytes"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

const (
	CharacterStateUnknown   = iota
	CharacterStateActive    = iota
	CharacterStateSuper     = iota
	CharacterStateSuspended = iota
	CharacterStateDeleted   = iota
)

const (
	UserStateUnknown   = iota
	UserStateInvited   = iota
	UserStateActive    = iota
	UserStateAdmin     = iota
	UserStateSuspended = iota
	UserStateDeleted   = iota
)

type Character struct {
	Region string `json:"region"`
	Name   string `json:"name"`
	User   string `json:"user"`
	State  uint   `json:"state"`
}

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password []byte `json:"pwd"`
	State    uint   `json:"state"`
}

type CharacterConnector func(endpoint string) (CharacterBackend, error)

type UserConnector func(endpoint string) (UserBackend, error)

type UserBackend interface {
	// Returns a page of the list of User on the system.
	List(marker string, max uint32) ([]User, error)

	// Create a User bond to the given email and password.
	Create(email string, pass []byte) (User, error)

	// Present the details of the User.
	Show(ID string) (User, error)

	GetByMail(email string) (User, error)

	// Make the give User an administrator. An administrator has extended privileges on
	// all the regions.
	Promote(ID string) error

	// Remove the "admin" status of the given User.
	Demote(ID string) error

	// Suspend the user, so that it still appears in the List but it is not allowed to
	// authenticate.
	Suspend(ID string) error

	// Remove the suspended flag on the user.
	Resume(ID string) error

	// Mark the User as Deleted. It must not appear in the list of users,
	// must not be allowed to authenticate and lose all its privileges.
	Delete(ID string) error

	// List all the characters owned by the given user.
	Characters(ID string) ([]Character, error)
}

type CharacterBackend interface {
	// Returns a page of the list of Characters in the guven region.
	List(region, marker string, max uint32) ([]Character, error)

	// Create a Character owned by the given User, in the given Region.
	// The Character name must be unique in the Region
	Create(userID, region, name string) (Character, error)

	Show(region, name string) (Character, error)

	// Give the Character the privileges of a Game Master
	Promote(region, name string) error

	// Remove the Game Master privileges to the given Character
	Demote(region, name string) error

	Suspend(region, name string) error

	Resume(region, name string) error

	Delete(region, name string) error
}

var (
	ErrExists   = errors.New("user already exist")
	ErrNotFound = errors.New("user not found")
	ErrDenied   = errors.New("permission denied")
	ErrInternal = errors.New("internal error")
)

var (
	// Returned by a connector when the connection string does not match the
	// pattern expected
	ErrSkip              = errors.New("not suitable for the connector")
	ErrUnmanagedEndpoint = errors.New("the connection string matches no connector")
)

var (
	userConnectors      = make([]UserConnector, 0)
	characterConnectors = make([]CharacterConnector, 0)
)

func RegisterUserConnector(cb UserConnector) {
	userConnectors = append(userConnectors, cb)
}

func RegisterCharacterConnector(cb CharacterConnector) {
	characterConnectors = append(characterConnectors, cb)
}

func ConnectUserBackend(endpoint string) (UserBackend, error) {
	for _, cb := range userConnectors {
		backend, err := cb(endpoint)
		switch err {
		case nil:
			return backend, nil
		case ErrSkip:
			continue
		default:
			return nil, err
		}
	}
	return nil, ErrUnmanagedEndpoint
}

func ConnectCharacterBackend(endpoint string) (CharacterBackend, error) {
	for _, cb := range characterConnectors {
		backend, err := cb(endpoint)
		switch err {
		case nil:
			return backend, nil
		case ErrSkip:
			continue
		default:
			return nil, err
		}
	}
	return nil, ErrUnmanagedEndpoint
}

func (u User) Authenticate(pass []byte) bool {
	if bytes.HasPrefix(u.Password, []byte{':'}) {
		return bytes.Equal(u.Password[1:], pass)
	}
	return nil != bcrypt.CompareHashAndPassword(u.Password, pass)
}

func init() {
	registerMem()
}
