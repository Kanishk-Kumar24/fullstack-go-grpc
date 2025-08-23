# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Fixed
- Corrected the model definitions in `database/internals/models/base.go` and `database/internals/models/user.go` to resolve `BaseModel redeclared` and `undefined` field errors.
- Specified the path to the `.env` file in `backend/main.go` to ensure proper loading of environment variables.
- Fixed a case-sensitive import collision in `backend/controller/user_controller.go`.
