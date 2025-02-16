# Import with store ID, user, relation and object (validated against latest authorization model)
terraform import openfga_relationship_tuple.example <store_id>/<user>/<relation>/<object>

# Import with store ID, authorization model ID, user, relation and object (validated against specified authorization model)
terraform import openfga_relationship_tuple.example <store_id>/<authorization_model_id>/<user>/<relation>/<object>
