//go:build go1.18

package echo

func MustGetValidated[T any](c Context) (T, error) {
	v, y := c.Internal().Get(`validated`).(*T)
	if y {
		return *v, nil
	}
	v = new(T)
	err := MustValidated[T](c, v)
	return *v, err
}

func MustValidated[T any](c Context, v *T) error {
	var valueDecoders BinderValueCustomDecoders
	var filters []FormDataFilter
	data := interface{}(v)
	if it, ok := data.(FiltersGetter); ok {
		filters = it.Filters(c)
	}
	if it, ok := data.(ValueDecodersGetter); ok {
		valueDecoders = it.ValueDecoders(c)
	}
	return c.MustBindAndValidateWithDecoder(v, valueDecoders, filters...)
}
