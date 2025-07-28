package qbmodels

import "time"

type QBReponseToken struct {
	AccessToken          string    `json:"access_token" firestore:"access_token"`
	RefreshToken         string    `json:"refresh_token" firestore:"refresh_token"`
	ExpiresInSec         int       `json:"expires_in" firestore:"expires_in"`
	ObtainedAt           time.Time `json:"obtained_at" firestore:"obtained_at"`
	RefresTokenExpiresIn int       `json:"x_refresh_token_expires_in" firestore:"x_refresh_token_expires_in"`
	ExpiresAt            time.Time `json:"expires_at" firestore:"expires_at"`
	TokenType            string    `json:"token_type" firestore:"token_type"`
	State                string    `json:"state" firestore:"state"`
	RealmId              string    `json:"realmId" firestore:"realmId"`
	Scope                string    `json:"scope" firestore:"scope"`
}

func (r *QBReponseToken) IsExpired() bool {
	return time.Now().After(r.ExpiresAt)
}
func (r *QBReponseToken) SetObtainedAt() {
	r.ObtainedAt = time.Now()
}
func (r *QBReponseToken) SetExpiresAt() {
	r.ExpiresAt = time.Now().Add(time.Duration(r.ExpiresInSec) * time.Second)
}
func (r *QBReponseToken) SetRealmID(realmID string) {
	r.RealmId = realmID
}
func (r *QBReponseToken) SetState(state string) {
	r.State = state
}
