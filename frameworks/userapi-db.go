package frameworks

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	cache "github.com/patrickmn/go-cache"
	"github.com/riomhaire/lightauth2/entities"
	"github.com/riomhaire/lightauth2/usecases"
	"gopkg.in/resty.v1"
)

// This file contains the code to interact with the userAPI as a user source
// Rather than the default csv
type UserAPIInteractor struct {
	registry      *usecases.Registry
	responseCache *cache.Cache
	useCache      bool
}

func NewUserAPIInteractor(registry *usecases.Registry) *UserAPIInteractor {
	d := UserAPIInteractor{}
	d.registry = registry
	// define cache if required
	if registry.Configuration.CacheTimeToLive > 0 {
		d.useCache = true
		ttl := time.Duration(registry.Configuration.CacheTimeToLive) * time.Second
		d.responseCache = cache.New(ttl, 2*ttl)
	}
	log.Printf("Reading UserAPI at %s, Caching Enabled %v for %v Seconds\n", registry.Configuration.UserAPIHost, d.useCache, registry.Configuration.CacheTimeToLive)
	return &d
}

func (d UserAPIInteractor) Lookup(username string) (entities.User, error) {
	// If we are using cache then check that 1st before the backend
	if d.useCache {

		cachedUser, found := d.responseCache.Get(username)
		if found {
			return cachedUser.(entities.User), nil
		}
	}

	target := fmt.Sprintf("%v/api/v1/user/account/%v", d.registry.Configuration.UserAPIHost, username)
	token := d.registry.Configuration.UserAPIKey
	d.registry.Logger.Log(usecases.Trace, fmt.Sprintf("Requesting '%v' from '%v'", username, target))
	// Assign Client Redirect Policy. Create one as per you need
	resty.SetRedirectPolicy(resty.FlexibleRedirectPolicy(15))
	resp, err := resty.R().SetHeader("Accept", "application/json").SetAuthToken(token).Get(target)
	user := entities.User{}
	if err != nil {
		d.registry.Logger.Log(usecases.Error, err.Error())
		return user, err
	}
	// Decode
	err = json.Unmarshal(resp.Body(), &user)
	if err != nil {
		d.registry.Logger.Log(usecases.Error, err.Error())
	}

	// If using cache remember
	if d.useCache {
		d.responseCache.Set(username, user, cache.DefaultExpiration)
	}
	return user, err
}

// We dont use this
func (d UserAPIInteractor) Create(user entities.User) error {
	return errors.New("Not Needed/Implemented For UserAPI")
}
