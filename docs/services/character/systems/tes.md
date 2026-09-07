# The Electric State Characters

The `tes` character implementation creates player characters (`pc`) and non-player characters (`npc`). The submitted character data is stored as JSON in PostgreSQL.

## Request envelope

| Field       | Type    | Required | Description |
|-------------|---------|----------|-------------|
| `user_id`   | UUID    | Yes      | Owner of the character |
| `type`      | string  | Yes      | `pc` or `npc` |
| `rules`     | boolean | No       | Apply TES creation-rule validation when `true`; defaults to `false` |
| `character` | object  | Yes | Character sheet; integrity is always validated |

Both modes require a non-zero `user_id`, a valid character `type`, and a non-blank character name of at most 255 Unicode characters. `rules: false` allows a minimal sheet and custom creation values; it does not disable integrity checks.

## Character sheet

| Field            | Type       | Description |
|------------------|------------|-------------|
| `name`           | string     | Character name |
| `archetype`      | string     | TES archetype |
| `favourite_song` | string     | Favourite song |
| `description`    | string     | Free-form description |
| `stats`          | object     | `strength`, `agility`, `wits`, and `empathy` |
| `derivatives`    | object     | Derived `health` and `hope` |
| `bliss`          | object     | Current `bliss` and `permanent` bliss |
| `talents`        | string[]   | Talent identifiers |
| `dream`          | string     | Character dream |
| `flaw`           | string     | Character flaw |
| `gear`           | gear[]     | Starting inventory |
| `cash`           | integer    | Starting cash |
| `journey`        | object     | Journey `goal` and `threat` |
| `tension`        | tension[]  | Relations to other travellers |
| `conditions`     | object     | `injuries` and `traumas` arrays |
| `vehicle`        | object     | Group vehicle and shared gear |

## Validation summary

Validation consists of two independently composed policies.

**Integrity (always enabled):**

- Name must be non-blank and at most 255 Unicode characters.
- Base stats, health, hope, bliss, permanent bliss, and cash must be non-negative. Omitted numeric values default to zero.
- Every supplied item in `gear`, `vehicle.SharedGear`, and `vehicle.stats.gear` must have a non-blank name and code, a non-negative price, a supported type, and valid type-specific values (including damage/range, armor, and neurocaster fields).
- A minimal sheet containing only a name remains valid in free mode. Narrative fields, archetype/talent catalogs, and vehicle creation requirements belong to the creation policy.

**TES creation policy (only with `"rules": true`):**

- Every base stat must be between 2 and 6.
- `health` is `(strength + agility) / 2`, rounded up; the `tough` talent adds 2.
- `hope` is `(empathy + wits) / 2`, rounded up; the `dreamer` talent adds 2.
- Characters with a total stat score above 15 have one talent; characters at or below 15 have two distinct talents.
- Archetype, talents, starting cash, gear, journey, tension, vehicle stats, traits, and shared gear are validated against TES creation rules.
- Inventory must contain between one and four items and no more than one neurocaster.
- A vehicle must contain exactly three shared gear items.

Supported archetypes are `artist`, `criminal`, `devotee`, `doctor`, `dronePilot`, `investigator`, `outsider`, `runawayKid`, `scientist`, and `veteran`.

Gear uses a common object whose relevant fields depend on `type`: `gear`, `weapon`, `armor`, or `neurocaster`. Common fields are `name`, `code`, `type`, `price`, and `notes`; `bonus` belongs to ordinary gear and weapons. Specialized fields include weapon damage and range, armor values, or neurocaster processor, network, and graphics values.

### Typed item handling

The flat request shape is unchanged. The service converts each input item through a validating factory into one of four domain variants: equipment, weapon, armor, or neurocaster. The same factory is used for personal inventory, shared gear, and vehicle equipment, independently of `rules`.

Non-empty fields belonging to another variant are rejected with `422` and code `incompatible_field`, for example `character.gear[0].armor_level` on a weapon. Legacy irrelevant zero numbers, empty strings, and empty/null arrays are accepted and omitted from stored JSON. Stored items retain their `type` and only the fields of that variant. Existing database rows are not rewritten.

Collection limits (one to four starting items and at most one neurocaster) remain in the optional creation policy.

## Vehicle shape

The current API uses `passenger` (singular) and `SharedGear` (capitalized); clients must send these names exactly.

```json
{
  "vehicle_type": "4wdCar",
  "model": "Station wagon",
  "passenger": 5,
  "fuel": "gasoline",
  "description": "The group's vehicle",
  "stats": {
    "maneuverability": 1,
    "speed": 2,
    "hull": 3,
    "armor": 1,
    "traits": [
      {"name": "hull", "stat": "hull", "bonus": 1}
    ],
    "gear": []
  },
  "SharedGear": [
    {"name": "Toolbox", "code": "toolbox", "type": "gear", "bonus": 1, "price": 0, "notes": ""},
    {"name": "Radio", "code": "radio", "type": "gear", "bonus": 1, "price": 0, "notes": ""},
    {"name": "First aid kit", "code": "firstAid", "type": "gear", "bonus": 1, "price": 0, "notes": ""}
  ]
}
```

## Minimal transport example

Because TES rule validation covers the full sheet, the easiest way to call the endpoint is to keep the request in a `character.json` file:

```json
{
  "user_id": "d7a92c4c-7c65-41d8-946a-4b94d3e721f9",
  "type": "pc",
  "rules": false,
  "character": {
    "name": "Mira"
  }
}
```

```bash
curl -X POST "http://localhost:8081/character?system=tes" \
  -H "Content-Type: application/json" \
  --data @character.json
```
