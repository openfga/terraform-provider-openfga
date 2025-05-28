## 0.4.0 (May 28, 2025)

NOTES:

* This is the first release as the **official Terraform provider** in the **OpenFGA organization**. Thank you to everyone who made this possible 🎉

SECURITY:

* provider: Updated terraform provider SDK

## 0.3.2 (March 10, 2025)

BUG FIXES:

* data_source/authorization_model: Fixed nil pointer for non-existing latest authorization model
* data_source/\*_query: Added missing documentation

## 0.3.1 (February 27, 2025)

BUG FIXES:

* data_source/authorization_model_document: Fixed broken module file names

## 0.3.0 (February 27, 2025)

FEATURES:

* data_source/authorization_model_document: Added support for modular models

## 0.2.1 (February 22, 2025)

BUG FIXES:

* docs: Fixed missing provider attributes

## 0.2.0 (February 22, 2025)

FEATURES:

* provider: Added `scopes` and `audience` attributes

## 0.1.0 (February 19, 2025)

FEATURES:

* provider: Provider added
* resource/store: Resource added
* data_source/store: Data source added
* data_source/stores: Data source added
* resource/authorization_model Resource added
* data_source/authorization_model: Data source added
* data_source/authorization_models: Data source added
* data_source/authorization_model_document: Data source added
* resource/relationship_tuple Resource added
* data_source/relationship_tuple: Data source added
* data_source/relationship_tuples: Data source added
* data_source/check_query: Data source added
* data_source/list_objects_query: Data source added
* data_source/list_users_query: Data source added
