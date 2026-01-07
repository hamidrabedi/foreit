package schema

// Type-specific field constructors.
// These return Field values that can be chained directly.

func Int64(name string) Field {
	return newField(name, TypeInt64)
}

func Int32(name string) Field {
	return newField(name, TypeInt32)
}

func String(name string) Field {
	return newField(name, TypeString)
}

func Text(name string) Field {
	return newField(name, TypeText)
}

func Email(name string) Field {
	return newField(name, TypeEmail)
}

func URL(name string) Field {
	return newField(name, TypeURL)
}

func Bool(name string) Field {
	return newField(name, TypeBool)
}

func Time(name string) Field {
	return newField(name, TypeTime)
}

func Date(name string) Field {
	return newField(name, TypeDate)
}

func DateTime(name string) Field {
	return newField(name, TypeDateTime)
}

func Float64(name string) Field {
	return newField(name, TypeFloat64)
}

func Float32(name string) Field {
	return newField(name, TypeFloat32)
}

func Decimal(name string) Field {
	return newField(name, TypeDecimal)
}

func JSON(name string) Field {
	return newField(name, TypeJSON)
}

func Bytes(name string) Field {
	return newField(name, TypeBytes)
}

func UUID(name string) Field {
	return newField(name, TypeUUID)
}
