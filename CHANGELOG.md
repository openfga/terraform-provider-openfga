## [0.5.0] - 2025-06-06

### Breaking Changes

- provider: Renamed `scopes` to `api_scopes`, `audience` to `api_audience` and `token_endpoint_url` to `api_token_issuer`. Changed environment variable names accordingly.

### Fixed

- goreleaser: Changed deprecated `archives.format` field to `archives.formats`.

## [0.4.0] - 2025-06-05

### Notes

- This is the first release as the **official Terraform provider** in the **OpenFGA organization**. Thank you to everyone who made this possible 🎉

### Security

- provider: Updated terraform provider SDK

## [0.3.2] - 2025-03-10

### Fixed

- data_source/authorization_model: Fixed nil pointer for non-existing latest authorization model
- data_source/\*_query: Added missing documentation

## [0.3.1] - 2025-02-27

### Fixed

- data_source/authorization_model_document: Fixed broken module file names

## [0.3.0] - 2025-02-27

### Added

- data_source/authorization_model_document: Added support for modular models

## [0.2.1] - 2025-02-22

### Fixed

- docs: Fixed missing provider attributes

## [0.2.0] - 2025-02-22

### Added

- provider: Added `scopes` and `audience` attributes

## [0.1.0] - 2025-02-19

### Added

- provider: Provider added
- resource/store: Resource added
- data_source/store: Data source added
- data_source/stores: Data source added
- resource/authorization_model Resource added
- data_source/authorization_model: Data source added
- data_source/authorization_models: Data source added
- data_source/authorization_model_document: Data source added
- resource/relationship_tuple Resource added
- data_source/relationship_tuple: Data source added
- data_source/relationship_tuples: Data source added
- data_source/check_query: Data source added
- data_source/list_objects_query: Data source added
- data_source/list_users_query: Data source added
