package cnaErrors

import "errors"

var ErrDuplicateUser = errors.New("user already exists")
var ErrUserHasUnknownRole = errors.New("user has unknown role")
var ErrNoSuchUser = errors.New( "no such user" )
var ErrNoSuchRole = errors.New( "no such role" )
var ErrCannotAddExistingRole = errors.New( "cannot add existing role to app" )
var ErrUserAlreadyHasRole = errors.New( "user already has role" )
var ErrNoSuchApp = errors.New( "no such app" )
var ErrNoHydraAdminUrl = errors.New( "please set HYDRA_ADMIN_URL" )
var ErrUnexpectedHydraResponse = errors.New( "unexpected hydra response" )
var ErrLoginOrPasswordWrong = errors.New("login or password is invalid" )