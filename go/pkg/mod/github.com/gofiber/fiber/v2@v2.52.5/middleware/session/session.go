package session

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"

	"github.com/valyala/fasthttp"
)

type Session struct {
	mu         sync.RWMutex  // Mutex to protect non-data fields
	id         string        // session id
	fresh      bool          // if new session
	ctx        *fiber.Ctx    // fiber context
	config     *Store        // store configuration
	data       *data         // key value data
	byteBuffer *bytes.Buffer // byte buffer for the en- and decode
	exp        time.Duration // expiration of this session
}

var sessionPool = sync.Pool{
	New: func() interface{} {
		return new(Session)
	},
}

func acquireSession() *Session {
	s := sessionPool.Get().(*Session) //nolint:forcetypeassert,errcheck // We store nothing else in the pool
	if s.data == nil {
		s.data = acquireData()
	}
	if s.byteBuffer == nil {
		s.byteBuffer = new(bytes.Buffer)
	}
	s.fresh = true
	return s
}

func releaseSession(s *Session) {
	s.mu.Lock()
	s.id = ""
	s.exp = 0
	s.ctx = nil
	s.config = nil
	if s.data != nil {
		s.data.Reset()
	}
	if s.byteBuffer != nil {
		s.byteBuffer.Reset()
	}
	s.mu.Unlock()
	sessionPool.Put(s)
}

// Fresh is true if the current session is new
func (s *Session) Fresh() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.fresh
}

// ID returns the session id
func (s *Session) ID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.id
}

// Get will return the value
func (s *Session) Get(key string) interface{} {
	// Better safe than sorry
	if s.data == nil {
		return nil
	}
	return s.data.Get(key)
}

// Set will update or create a new key value
func (s *Session) Set(key string, val interface{}) {
	// Better safe than sorry
	if s.data == nil {
		return
	}
	s.data.Set(key, val)
}

// Delete will delete the value
func (s *Session) Delete(key string) {
	// Better safe than sorry
	if s.data == nil {
		return
	}
	s.data.Delete(key)
}

// Destroy will delete the session from Storage and expire session cookie
func (s *Session) Destroy() error {
	// Better safe than sorry
	if s.data == nil {
		return nil
	}

	// Reset local data
	s.data.Reset()

	s.mu.RLock()
	defer s.mu.RUnlock()

	// Use external Storage if exist
	if err := s.config.Storage.Delete(s.id); err != nil {
		return err
	}

	// Expire session
	s.delSession()
	return nil
}

// Regenerate generates a new session id and delete the old one from Storage
func (s *Session) Regenerate() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Delete old id from storage
	if err := s.config.Storage.Delete(s.id); err != nil {
		return err
	}

	// Generate a new session, and set session.fresh to true
	s.refresh()

	return nil
}

// Reset generates a new session id, deletes the old one from storage, and resets the associated data
func (s *Session) Reset() error {
	// Reset local data
	if s.data != nil {
		s.data.Reset()
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Reset byte buffer
	if s.byteBuffer != nil {
		s.byteBuffer.Reset()
	}
	// Reset expiration
	s.exp = 0

	// Delete old id from storage
	if err := s.config.Storage.Delete(s.id); err != nil {
		return err
	}

	// Expire session
	s.delSession()

	// Generate a new session, and set session.fresh to true
	s.refresh()

	return nil
}

// refresh generates a new session, and set session.fresh to be true
func (s *Session) refresh() {
	s.id = s.config.KeyGenerator()
	s.fresh = true
}

// Save will update the storage and client cookie
//
// sess.Save() will save the session data to the storage and update the
// client cookie, and it will release the session after saving.
//
// It's not safe to use the session after calling Save().
func (s *Session) Save() error {
	// Better safe than sorry
	if s.data == nil {
		return nil
	}

	s.mu.Lock()

	// Check if session has your own expiration, otherwise use default value
	if s.exp <= 0 {
		s.exp = s.config.Expiration
	}

	// Update client cookie
	s.setSession()

	// Convert data to bytes
	encCache := gob.NewEncoder(s.byteBuffer)
	err := encCache.Encode(&s.data.Data)
	if err != nil {
		return fmt.Errorf("failed to encode data: %w", err)
	}

	// Copy the data in buffer
	encodedBytes := make([]byte, s.byteBuffer.Len())
	copy(encodedBytes, s.byteBuffer.Bytes())

	// Pass copied bytes with session id to provider
	if err := s.config.Storage.Set(s.id, encodedBytes, s.exp); err != nil {
		return err
	}

	s.mu.Unlock()

	// Release session
	// TODO: It's not safe to use the Session after calling Save()
	releaseSession(s)

	return nil
}

// Keys will retrieve all keys in current session
func (s *Session) Keys() []string {
	if s.data == nil {
		return []string{}
	}
	return s.data.Keys()
}

// SetExpiry sets a specific expiration for this session
func (s *Session) SetExpiry(exp time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.exp = exp
}

func (s *Session) setSession() {
	if s.config.source == SourceHeader {
		s.ctx.Request().Header.SetBytesV(s.config.sessionName, []byte(s.id))
		s.ctx.Response().Header.SetBytesV(s.config.sessionName, []byte(s.id))
	} else {
		fcookie := fasthttp.AcquireCookie()
		fcookie.SetKey(s.config.sessionName)
		fcookie.SetValue(s.id)
		fcookie.SetPath(s.config.CookiePath)
		fcookie.SetDomain(s.config.CookieDomain)
		// Cookies are also session cookies if they do not specify the Expires or Max-Age attribute.
		// refer: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie
		if !s.config.CookieSessionOnly {
			fcookie.SetMaxAge(int(s.exp.Seconds()))
			fcookie.SetExpire(time.Now().Add(s.exp))
		}
		fcookie.SetSecure(s.config.CookieSecure)
		fcookie.SetHTTPOnly(s.config.CookieHTTPOnly)

		switch utils.ToLower(s.config.CookieSameSite) {
		case "strict":
			fcookie.SetSameSite(fasthttp.CookieSameSiteStrictMode)
		case "none":
			fcookie.SetSameSite(fasthttp.CookieSameSiteNoneMode)
		default:
			fcookie.SetSameSite(fasthttp.CookieSameSiteLaxMode)
		}
		s.ctx.Response().Header.SetCookie(fcookie)
		fasthttp.ReleaseCookie(fcookie)
	}
}

func (s *Session) delSession() {
	if s.config.source == SourceHeader {
		s.ctx.Request().Header.Del(s.config.sessionName)
		s.ctx.Response().Header.Del(s.config.sessionName)
	} else {
		s.ctx.Request().Header.DelCookie(s.config.sessionName)
		s.ctx.Response().Header.DelCookie(s.config.sessionName)

		fcookie := fasthttp.AcquireCookie()
		fcookie.SetKey(s.config.sessionName)
		fcookie.SetPath(s.config.CookiePath)
		fcookie.SetDomain(s.config.CookieDomain)
		fcookie.SetMaxAge(-1)
		fcookie.SetExpire(time.Now().Add(-1 * time.Minute))
		fcookie.SetSecure(s.config.CookieSecure)
		fcookie.SetHTTPOnly(s.config.CookieHTTPOnly)

		switch utils.ToLower(s.config.CookieSameSite) {
		case "strict":
			fcookie.SetSameSite(fasthttp.CookieSameSiteStrictMode)
		case "none":
			fcookie.SetSameSite(fasthttp.CookieSameSiteNoneMode)
		default:
			fcookie.SetSameSite(fasthttp.CookieSameSiteLaxMode)
		}

		s.ctx.Response().Header.SetCookie(fcookie)
		fasthttp.ReleaseCookie(fcookie)
	}
}

// decodeSessionData decodes the session data from raw bytes.
func (s *Session) decodeSessionData(rawData []byte) error {
	_, _ = s.byteBuffer.Write(rawData) //nolint:errcheck // This will never fail
	encCache := gob.NewDecoder(s.byteBuffer)
	if err := encCache.Decode(&s.data.Data); err != nil {
		return fmt.Errorf("failed to decode session data: %w", err)
	}
	return nil
}
