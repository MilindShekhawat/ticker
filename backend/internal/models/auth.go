package models

import "fmt"

func SignupWithSession(email, password, name, ip, userAgent string) (*User, *Session, error) {
	userCount, err := CountUsers()
	if err != nil {
		return nil, nil, err
	}

	if userCount == 0 {
		return bootstrapFirstUserWithSession(email, password, name, ip, userAgent)
	}

	user, err := CreateUser(email, password, name)
	if err != nil {
		return nil, nil, err
	}

	session, err := CreateSession(user.ID, ip, userAgent)
	if err != nil {
		return nil, nil, err
	}

	return user, session, nil
}

func LoginWithSession(email, password, ip, userAgent string) (*User, *Session, error) {
	user, err := AuthenticateUser(email, password)
	if err != nil {
		return nil, nil, err
	}

	session, err := CreateSession(user.ID, ip, userAgent)
	if err != nil {
		return nil, nil, err
	}

	return user, session, nil
}

func bootstrapFirstUserWithSession(email, password, name, ip, userAgent string) (*User, *Session, error) {
	user, err := CreateUser(email, password, name)
	if err != nil {
		return nil, nil, err
	}

	session, err := CreateSession(user.ID, ip, userAgent)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create bootstrap session: %w", err)
	}

	return user, session, nil
}
