package cnaErrors

import "errors"

var ErrDuplicateUser = errors.New("user already exists")
var ErrUserHasUnknownRole = errors.New("user has unknown role")
var ErrNoSuchUser = errors.New( "no such user" )
var ErrNoSuchRole = errors.New( "no such role" )
var ErrCannotAddExistingRole = errors.New( "cannot add existing role to app" )
var ErrUserAlreadyHasRole = errors.New( "user already has role" )
var ErrNoSuchApp = errors.New( "no such app" )
var ErrLoginOrPasswordWrong = errors.New("login or password is invalid" )
var ErrNoSuchSession = errors.New( "no such session" )
var ErrMigrationFailed = errors.New( "migration failed" )
var ErrDatabaseNotInitialised = errors.New( "database not initialised")
var ErrActionForbidden = errors.New( "action forbidden" );