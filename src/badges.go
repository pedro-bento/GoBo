package main

import "strings"

type badge int

const (
	badgeGod         badge = -1
	badgeBroadcaster badge = 0
	badgeModerator   badge = 1
	badgeNone        badge = 3
)

func badgeFromString(str string) badge {
	if strings.Contains(str, "broadcaster") {
		return badgeBroadcaster
	} else if strings.Contains(str, "moderator") {
		return badgeModerator
	} else if strings.Contains(str, "god") {
		return badgeGod
	}

	return badgeNone
}

func asPerm(perm, userBadge badge) bool {
	if perm == badgeGod {
		return false
	} else if perm == badgeNone {
		return true
	} else if perm == badgeModerator {
		return userBadge == badgeModerator || userBadge == badgeBroadcaster
	} else if perm == badgeBroadcaster {
		return userBadge == badgeBroadcaster
	}

	return false
}
