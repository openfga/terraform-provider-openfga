- add relationship_tuple (resources & data sources)
- add relationship_query_X (data sources) -> _check, _batch_check, _list_objects, _list_users, _expand
  ```
  data "openfga_relationship_query_check" "tom_can_read_file" {
    "user"     = "user:tom"
    "relation" = "can_read"
    "object"   = "file:dummy"
  }
  
  data.openfga_relationship_query.tom_can_read_file.result
  ```
- fix the upstream bug in openfga
- add `model` to `authorization_model_document`
- reduce code duplication (espically for JSON conversion)
- code review by Tom 