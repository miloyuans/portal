package kcadmin

import (
	"fmt"
	"net/url"
)

func realmPath(realm string) string {
	return fmt.Sprintf("/%s", realm)
}

func clientsPath(realm string) string {
	return fmt.Sprintf("/%s/clients", realm)
}

func clientPath(realm, clientUUID string) string {
	return fmt.Sprintf("/%s/clients/%s", realm, clientUUID)
}

func clientServiceAccountPath(realm, clientUUID string) string {
	return fmt.Sprintf("/%s/clients/%s/service-account-user", realm, clientUUID)
}

func usersPath(realm string) string {
	return fmt.Sprintf("/%s/users", realm)
}

func userPath(realm, userID string) string {
	return fmt.Sprintf("/%s/users/%s", realm, userID)
}

func userRoleMappingsPath(realm, userID string) string {
	return fmt.Sprintf("/%s/users/%s/role-mappings", realm, userID)
}

func userEffectiveRealmRolesPath(realm, userID string) string {
	return fmt.Sprintf("/%s/users/%s/role-mappings/realm/composite", realm, userID)
}

func userEffectiveClientRolesPath(realm, userID, clientUUID string) string {
	return fmt.Sprintf("/%s/users/%s/role-mappings/clients/%s/composite", realm, userID, clientUUID)
}

func withQuery(path string, values url.Values) string {
	if values == nil || len(values) == 0 {
		return path
	}
	return path + "?" + values.Encode()
}
