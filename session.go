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

// Options stores configuration for a session or session store.
// Fields are a subset of http.Cookie fields.
type SessionOptions struct {
	Engine string //Store Engine
	Name   string //Session Name
	*CookieOptions
}

// Wraps thinly gorilla-session methods.
// Session stores the values and optional configuration for a session.
type Session interface {
	// Get returns the session value associated to the given key.
	Get(key string) interface{}
	// Set sets the session value associated to the given key.
	Set(key string, val interface{}) Session
	SetId(id string) Session
	Id() string
	// Delete removes the session value associated to the given key.
	Delete(key string) Session
	// Clear deletes all values in the session.
	Clear() Session
	// AddFlash adds a flash message to the session.
	// A single variadic argument is accepted, and it is optional: it defines the flash key.
	// If not defined "_flash" is used by default.
	AddFlash(value interface{}, vars ...string) Session
	// Flashes returns a slice of flash messages from the session.
	// A single variadic argument is accepted, and it is optional: it defines the flash key.
	// If not defined "_flash" is used by default.
	Flashes(vars ...string) []interface{}

	Options(SessionOptions) Session

	// Save saves all sessions used during the current request.
	Save() error
}
