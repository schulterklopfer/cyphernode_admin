package cnaErrors

import "errors"

var ErrDuplicateUser = errors.New("user already exists")
var ErrUserHasUnknownRole = errors.New("user has unknown role")
var ErrNoSuchUser = errors.New( "no such user" )
var ErrNoSuchRole = errors.New( "no such role" )
var ErrCannotAddExistingRole = errors.New( "cannot add existing role to app" )
var ErrUserAlreadyHasRole = errors.New( "user already has role" )
var ErrNoSuchApp = errors.New( "no such app" )
