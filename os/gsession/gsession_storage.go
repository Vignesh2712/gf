// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsession

import "time"

type Storage interface {
	// New creates a custom session id.
	// This function can be used for custom session creation.
	New(ttl time.Duration) (id string)

	// Get retrieves session value with given key.
	// It returns nil if the key does not exist in the session.
	Get(id string, key string) interface{}

	// GetMap retrieves all key-value pairs as map from storage.
	GetMap(id string) map[string]interface{}

	// GetSize retrieves the size of key-value pairs from storage.
	GetSize(id string) int

	// Set sets key-value session pair to the storage.
	// The parameter <ttl> specifies the TTL for the session id (not for the key-value pair).
	Set(id string, key string, value interface{}, ttl time.Duration) error

	// SetMap batch sets key-value session pairs with map to the storage.
	// The parameter <ttl> specifies the TTL for the session id(not for the key-value pair).
	SetMap(id string, data map[string]interface{}, ttl time.Duration) error

	// Remove deletes key with its value from storage.
	Remove(id string, key string) error

	// RemoveAll deletes all key-value pairs from storage.
	RemoveAll(id string) error

	// GetSession returns the session data as map for given session id.
	// The parameter <ttl> specifies the TTL for this session.
	// It returns nil if the TTL is exceeded.
	GetSession(id string, ttl time.Duration) map[string]interface{}

	// SetSession updates the data map for specified session id.
	// This function is called ever after session, which is changed dirty, is closed.
	// This copy all session data map from memory to storage.
	SetSession(id string, data map[string]interface{}, ttl time.Duration) error

	// UpdateTTL updates the TTL for specified session id.
	// This function is called ever after session, which is not dirty, is closed.
	UpdateTTL(id string, ttl time.Duration) error
}
