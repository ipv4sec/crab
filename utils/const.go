package utils

const (
	ErrBadRequestParam      = 10101
	ErrInternalServer       = 10102
	ErrBase64Decode         = 10103
	ErrJSONMarshal          = 10104
	//ErrYAMLMarshal          = 10105

	ErrClusterClientApply   = 20101
	ErrClusterClientPatch   = 20102

	ErrClusterGetConfigMap  = 20201
	ErrClusterSetConfigMap  = 20202
	ErrClusterUpdateConfigMap  = 20202

	ErrClusterGetDeployment = 20301

	ErrClusterGetSecret     = 20401
)
