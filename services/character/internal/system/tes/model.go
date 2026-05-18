package tes

type createRequest struct {
	Type          string        `json:"type"`
	Rules         bool          `json:"rules"`
	CharacterList characterList `json:"character"`
}

type characterList struct {
	Name          string      `json:"name"`
	Archetype     string      `json:"archetype"`
	FavouriteSong string      `json:"favourite_song"`
	Description   string      `json:"description"`
	Stats         stats       `json:"stats"`
	Derivatives   derivatives `json:"derivatives"`
	Bliss         bliss       `json:"bliss"`
	Talents       []string    `json:"talents"`
	Dream         string      `json:"dream"`
	Flaw          string      `json:"flaw"`
	Inventory     []gear      `json:"gear"`
	Cash          int         `json:"cash"`
	Journey       journey     `json:"journey"`
	Tension       []tension   `json:"tension"`
	Conditions    conditions  `json:"conditions"`
	Vehicle       vehicle     `json:"vehicle"`
}

type stats struct {
	Strength int `json:"strength"`
	Agility  int `json:"agility"`
	Wits     int `json:"wits"`
	Empathy  int `json:"empathy"`
}

type derivatives struct {
	Health int `json:"health"`
	Hope   int `json:"hope"`
}

type bliss struct {
	Bliss     int `json:"bliss"`
	Permanent int `json:"permanent"`
}

type gear struct {
	// general gear
	Name  string `json:"name"`
	Code  string `json:"code"`
	Type  string `json:"type"`
	Bonus int    `json:"bonus"`
	Price int    `json:"price"`
	Notes string `json:"notes"`

	// weapon specific
	DamageValue int      `json:"damage_value"`
	DamageKind  string   `json:"damage_kind"`
	RangeMin    string   `json:"range_min"`
	RangeMax    string   `json:"range_max"`
	Tags        []string `json:"tags"`

	// explosive specific
	BlastPower int `json:"blast_power"`

	// armor specific
	ArmorLevel      int `json:"armor_level"`
	AgilityModifier int `json:"agility_modifier"`

	// neurocaster specific
	Processor int `json:"processor"`
	Network   int `json:"network"`
	Graphics  int `json:"graphics"`
}

type journey struct {
	Goal   string `json:"goal"`
	Threat string `json:"threat"`
}

type tension struct {
	TravellerName string `json:"traveller_name"`
	Tension       int    `json:"tension"`
}

type conditions struct {
	Injuries []string `json:"injuries"`
	Traumas  []string `json:"traumas"`
}

type vehicle struct {
	VehicleType string       `json:"vehicle_type"`
	Model       string       `json:"model"`
	Passengers  int          `json:"passenger"`
	Fuel        string       `json:"fuel"`
	Description string       `json:"description"`
	Stats       vehicleStats `json:"stats"`
	GroupGear   []gear       `json:"group_gear"`
}

type vehicleStats struct {
	Maneuverability int        `json:"maneuverability"`
	Speed           int        `json:"speed"`
	Hull            int        `json:"hull"`
	Armor           int        `json:"armor"`
	Traits          []carTrait `json:"traits"`
	Gear            []gear     `json:"gear"`
}

type carTrait struct {
	Name  string `json:"name"`
	Stat  string `json:"stat"`
	Bonus int    `json:"bonus"`
}
