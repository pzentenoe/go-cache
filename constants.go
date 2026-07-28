package cache

const (
	errOverflowFormat          = "overflow would occur for %s"
	errUnderflowFormat         = "underflow would occur for %s"
	errItemNotFoundFormat      = "item %s not found"
	errNotFloatTypeFormat      = "the value for %s does not have type float32 or float64"
	errGobRegistration         = "error registering item types with Gob library"
	errItemAlreadyExistsFormat = "item %s already exists"
	errItemDoesNotExistFormat  = "item %s doesn't exist"
	errNotIntegerOrFloatFormat = "the value for %s is not an integer or float"
	errTypeMismatchFormat      = "the value for %s does not match the type expected by this operation"
)
