package models

// --- Role ---
type RoleEnum string

const (
	RoleSuperAdmin     RoleEnum = "super_admin"
	RoleAdmin          RoleEnum = "admin"
	RolePresident      RoleEnum = "president"
	RoleManager        RoleEnum = "manager"
	RoleHeadCoach      RoleEnum = "head_coach"
	RoleAssistantCoach RoleEnum = "assistant_coach"
	RoleScoutman       RoleEnum = "scoutman"
	RolePlayer         RoleEnum = "player"
	RoleAgent          RoleEnum = "agent"
	RoleGuest          RoleEnum = "guest"
)

// --- Role Pemissions ---
var RolePermissions = map[RoleEnum][]string{
	RoleSuperAdmin: {"all"},
	RoleAdmin: {
		"manage_users",
		"manage_teams",
		"manage_matches",
		"upload_video",
		"upload_scout",
		"manage_season",
		"manage_waitlist",
		"manage_roles",
		"view_audit_logs",
		"manage_permissions",
		"manage_notifications",
		"manage_settings",
		"manage_reports",
		"manage_logs",
		"manage_feedback",
		"manage_subscriptions",
		"manage_payments",
		"manage_tickets",
		"manage_events",
		"manage_promotions",
		"manage_partners",
		"manage_sponsors",
		"manage_merchandise",
		"manage_analytics",
		"manage_marketing",
		"manage_content",
		"manage_social_media",
		"manage_website",
		"manage_app",
		"manage_integration"},
	RolePresident: {
		"view_team",
		"approve_join_request",
		"view_match",
		"manage_team_roster",
		"schedule_practices",
		"assign_roles",
	},
	RoleManager: {
		"view_team",
		"approve_join_request",
		"view_match",
		"manage_team_roster",
		"schedule_practices",
		"assign_roles",
	},
	RoleHeadCoach: {
		"view_team",
		"upload_video",
		"upload_scout",
		"view_match",
		"view_scout_data",
	},
	RoleAssistantCoach: {
		"view_team",
		"view_match",
		"view_scout_data",
	},
	RoleScoutman: {
		"upload_video",
		"upload_scout",
		"view_match",
	},
	RolePlayer: {
		"view_own_stats",
		"view_match",
	},
	RoleAgent: {
		"view_player_profiles",
	},
	RoleGuest: {
		"view_public_stats",
	},
}

// --- Gender ---
type GenderEnum string

const (
	GenderMale   GenderEnum = "male"
	GenderFemale GenderEnum = "female"
	GenderOther  GenderEnum = "other"
)

// --- Season Name ---
type SeasonNameEnum string

const (
	SeasonBundesliga           SeasonNameEnum = "Bundesliga"
	SeasonPlusLiga             SeasonNameEnum = "PlusLiga"
	SeasonSerieA               SeasonNameEnum = "Serie A"
	SeasonSuperLega            SeasonNameEnum = "SuperLega"
	SeasonLigueA               SeasonNameEnum = "Ligue A"
	SeasonChampionsLeague      SeasonNameEnum = "Champions League"
	SeasonVNL                  SeasonNameEnum = "Volleyball Nations League"
	SeasonWorldChampionship    SeasonNameEnum = "World Championship"
	SeasonEuropeanChampionship SeasonNameEnum = "European Championship"
)

// --- Season Type ---
type SeasonTypeEnum string

const (
	SeasonTypeLeague     SeasonTypeEnum = "League"
	SeasonTypeTournament SeasonTypeEnum = "Tournament"
)

// --- Country ---
type CountryEnum string

const (
	CountryGermany CountryEnum = "Germany"
	CountryItaly   CountryEnum = "Italy"
	CountryFrance  CountryEnum = "France"
	CountryPoland  CountryEnum = "Poland"
)

// --- Round ---
type RoundEnum string

const (
	RoundFirstRound        RoundEnum = "First Round"
	RoundSecondRound       RoundEnum = "Second Round"
	RoundGroupStage        RoundEnum = "Group Stage"
	RoundPlayouts          RoundEnum = "Playouts"
	RoundPlayoffs          RoundEnum = "Playoffs"
	RoundQuarterfinals     RoundEnum = "Quarterfinals"
	RoundSemifinals        RoundEnum = "Semifinals"
	RoundFinals            RoundEnum = "Finals"
	RoundThirdPlacePlayoff RoundEnum = "Third Place Playoff"
	RoundSuperFinal        RoundEnum = "Super Final"
)

// --- Validations ---
func IsValidRole(r RoleEnum) bool {
	switch r {
	case RoleSuperAdmin, RoleAdmin, RolePresident, RoleHeadCoach, RoleAssistantCoach, RoleScoutman, RolePlayer, RoleAgent, RoleGuest:
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

func IsValidRound(r RoundEnum) bool {
	switch r {
	case RoundFirstRound, RoundSecondRound, RoundGroupStage, RoundPlayouts, RoundPlayoffs,
		RoundQuarterfinals, RoundSemifinals, RoundFinals, RoundThirdPlacePlayoff, RoundSuperFinal:
		return true
	default:
		return false
	}
}

func IsValidCountry(c CountryEnum) bool {
	switch c {
	case CountryGermany, CountryItaly, CountryFrance, CountryPoland:
		return true
	default:
		return false
	}
}

func IsValidSeasonName(s SeasonNameEnum) bool {
	switch s {
	case SeasonBundesliga, SeasonPlusLiga, SeasonSerieA, SeasonSuperLega,
		SeasonLigueA, SeasonChampionsLeague, SeasonVNL,
		SeasonWorldChampionship, SeasonEuropeanChampionship:
		return true
	default:
		return false
	}
}
