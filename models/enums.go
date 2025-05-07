package models

// --- Role ---
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

func IsValidRound(r RoundEnum) bool {
	switch r {
	case RoundFirstRound, RoundSecondRound, RoundGroupStage, RoundPlayouts, RoundPlayoffs,
		RoundQuarterfinals, RoundSemifinals, RoundFinals, RoundThirdPlacePlayoff, RoundSuperFinal:
		return true
	default:
		return false
	}
}
