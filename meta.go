/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

package echo

type MetaValidator interface {
	MethodGetter
	FiltersGetter
	ValueDecodersGetter
}

type FiltersGetter interface {
	Filters(Context) []FormDataFilter
}

type MethodGetter interface {
	Methods() []string
}

type RequestValidator func() MetaValidator

func NewBaseRequestValidator(data interface{}, method ...string) *BaseRequestValidator {
	if len(method) == 0 {
		if mt, ok := data.(MethodGetter); ok {
			method = mt.Methods()
		}
	}
	return &BaseRequestValidator{data: data, methods: method}
}

type BaseRequestValidator struct {
	methods []string
	data    interface{}
}

func (b *BaseRequestValidator) SetStruct(data interface{}) *BaseRequestValidator {
	b.data = data
	return b
}

func (b *BaseRequestValidator) Methods() []string {
	return b.methods
}

func (b *BaseRequestValidator) Filters(c Context) []FormDataFilter {
	if gt, ok := b.data.(FiltersGetter); ok {
		return gt.Filters(c)
	}
	return nil
}

func (b *BaseRequestValidator) ValueDecoders(c Context) BinderValueCustomDecoders {
	if gt, ok := b.data.(ValueDecodersGetter); ok {
		return gt.ValueDecoders(c)
	}
	return nil
}

type MetaHandler struct {
	meta    H
	request RequestValidator
	Handler
}

func (m *MetaHandler) Name() string {
	if v, y := m.Handler.(Name); y {
		return v.Name()
	}
	return HandlerName(m.Handler)
}

func (m *MetaHandler) Meta() H {
	return m.meta
}

func (m *MetaHandler) Handle(c Context) error {
	if m.request == nil {
		return m.Handler.Handle(c)
	}
	recv := m.request()
	methods := recv.Methods()
	var data interface{}
	if bs, ok := recv.(*BaseRequestValidator); ok {
		data = bs.data
	} else {
		data = recv
	}
	if len(methods) > 0 && !InSliceFold(c.Method(), methods) {
		return m.Handler.Handle(c)
	}
	if err := c.MustBindAndValidateWithDecoder(data, recv.ValueDecoders(c), recv.Filters(c)...); err != nil {
		return err
	}
	c.Internal().Set(`validated`, data)
	return m.Handler.Handle(c)
}

func GetValidated(c Context, defaults ...interface{}) interface{} {
	return c.Internal().Get(`validated`, defaults...)
}

func MustGetValidated[T any](c Context) (*T, error) {
	v, y := c.Internal().Get(`validated`).(*T)
	if y {
		return v, nil
	}
	v = new(T)
	err := MustValidated(c, v)
	return v, err
}

func MustValidated[T any](c Context, v T) error {
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
