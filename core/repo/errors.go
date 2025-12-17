package repo

import "errors"

var ErrMakeTempDir = errors.New("temp dir creation failed")
var ErrCloneRepo = errors.New("clone repository failed")
var ErrCopyRepo = errors.New("copy repository failed")
var ErrParseUrl = errors.New("parse providet url failed")
var ErrisValidHost = errors.New("--host is not a valid value; omit the flag to use the default value")
var ErrisValidUser = errors.New("--user is not a valid value")
var ErrisValidUrl = errors.New("--url is not a valid value")
