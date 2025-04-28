package models

type RoleEnum string
type GenderEnum string

const (
	RoleSuperAdmin     RoleEnum = "super_admin"
	RoleAdmin          RoleEnum = "admin"
	RolePresident      RoleEnum = "president"
	RoleHeadCoach      RoleEnum = "head_coach"
	RoleAssistantCoach RoleEnum = "assistant_coach"
	RoleScoutman       RoleEnum = "scoutman"
	RolePlayer         RoleEnum = "player"
)

const (
	GenderMale   GenderEnum = "male"
	GenderFemale GenderEnum = "female"
	GenderOther  GenderEnum = "other"
)

var AllRoles = []RoleEnum{
	RoleSuperAdmin,
	RoleAdmin,
	RolePresident,
	RoleHeadCoach,
	RoleAssistantCoach,
	RoleScoutman,
	RolePlayer,
}

var AllGenders = []GenderEnum{
	GenderMale,
	GenderFemale,
	GenderOther,
}

func IsValidRole(r RoleEnum) bool {
	switch r {
	case RoleSuperAdmin, RoleAdmin, RolePresident, RoleHeadCoach, RoleAssistantCoach, RoleScoutman, RolePlayer:
		return true
	default:
		return false
	}
}

func IsValidGender(g GenderEnum) bool {
	switch g {
	case GenderMale, GenderFemale, GenderOther:
		return true
	default:
		return false
	}
}
