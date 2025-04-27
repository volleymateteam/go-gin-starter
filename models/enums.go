package models

type RoleEnum string

const (
	RoleSuperAdmin     RoleEnum = "super_admin"
	RoleAdmin          RoleEnum = "admin"
	RolePresident      RoleEnum = "president"
	RoleHeadCoach      RoleEnum = "head_coach"
	RoleAssistantCoach RoleEnum = "assistant_coach"
	RoleScoutman       RoleEnum = "scoutman"
	RolePlayer         RoleEnum = "player"
)

type GenderEnum string

const (
	GenderMale   GenderEnum = "male"
	GenderFemale GenderEnum = "female"
	GenderOther  GenderEnum = "other"
)
