# VTM v5 System

Vampire: The Masquerade 5th Edition by White Wolf

- Uses d10 dice pools
- Successes on 6+
- Pairs of 10s grant critical successes (+2 extra successes per pair)
- Hunger dice can cause messy criticals or bestial failures

## Actions

- `roll` — Standard dice pool roll
- `reroll` — Willpower reroll (reroll specific main dice)
- `check` — Simple one-die check (rouse, frenzy, etc.)

---

## roll

Perform a standard VTM v5 dice pool roll.

### Request

| Field       | Type | Required | Constraints | Description                     |
|-------------|------|----------|-------------|---------------------------------|
| `attribute` | int  | Yes      | >= 1        | Attribute dots                  |
| `skill`     | int  | Yes      | >= 0        | Skill dots                      |
| `hunger`    | int  | Yes      | >= 0        | Current hunger level            |
| `target`    | int  | Yes      | >= 1        | Required successes (difficulty) |

### Response

| Field         | Type     | Description                                      |
|---------------|----------|--------------------------------------------------|
| `expression`  | string   | Dice notation (e.g., "5d10")                     |
| `main_roll`   | int[]    | Regular dice results                             |
| `hunger_roll` | int[]    | Hunger dice results                              |
| `successes`   | int      | Total successes (including critical bonus)       |
| `success`     | bool     | Whether roll met the target                      |
| `is_critical` | bool     | Whether a critical occurred                      |
| `crit_type`   | string   | Critical type (see below)                        |

### Critical Types

| Value              | Description                              |
|--------------------|------------------------------------------|
| `"none"`           | No critical event                        |
| `"critical"`       | Pair of 10s on regular dice only         |
| `"messy critical"` | Pair of 10s including hunger dice        |
| `"bestial failure"`| Failed roll with 1 on hunger dice        |

### Example

```bash
curl -X POST "http://localhost:8080/resolve?system=vtmv5&action=roll" \
  -H "Content-Type: application/json" \
  -d '{"attribute": 3, "skill": 2, "hunger": 2, "target": 3}'
```

```json
{
  "expression": "5d10",
  "main_roll": [10, 7, 3],
  "hunger_roll": [10, 4],
  "successes": 6,
  "success": true,
  "is_critical": true,
  "crit_type": "messy critical"
}
```

---

## reroll

Reroll specific dice from a previous roll (Willpower reroll). Only main dice can be rerolled, not hunger dice.

### Request

| Field          | Type  | Required | Constraints                             | Description                        |
|----------------|-------|----------|-----------------------------------------|------------------------------------|
| `main_roll`    | int[] | Yes      | Combined length >= 1                    | Previous main dice results         |
| `hunger_roll`  | int[] | Yes      | Combined length >= 1                    | Previous hunger dice results       |
| `reroll_index` | int[] | Yes      | length >= 1, indices < main_roll length | Indices of main_roll dice to reroll|
| `target`       | int   | Yes      | >= 1                                    | Required successes                 |

### Response

| Field              | Type     | Description                          |
|--------------------|----------|--------------------------------------|
| `expression`       | string   | Original dice notation               |
| `reroll_expression`| string   | Rerolled dice notation               |
| `main_roll`        | int[]    | Updated main dice results            |
| `hunger_roll`      | int[]    | Unchanged hunger dice results        |
| `successes`        | int      | New total successes                  |
| `success`          | bool     | Whether roll met the target          |
| `is_critical`      | bool     | Whether a critical occurred          |
| `crit_type`        | string   | Critical type (see roll action)      |

### Example

```bash
curl -X POST "http://localhost:8080/resolve?system=vtmv5&action=reroll" \
  -H "Content-Type: application/json" \
  -d '{
    "main_roll": [3, 5, 7],
    "hunger_roll": [10, 4],
    "reroll_index": [0, 1],
    "target": 3
  }'
```

```json
{
  "expression": "5d10",
  "reroll_expression": "2d10",
  "main_roll": [8, 10, 7],
  "hunger_roll": [10, 4],
  "successes": 6,
  "success": true,
  "is_critical": true,
  "crit_type": "messy critical"
}
```

---

## check

Simple one-die check (rouse check, frenzy check, etc.). Success on 6+.

### Request

Empty body required.

### Response

| Field        | Type   | Description              |
|--------------|--------|--------------------------|
| `expression` | string | Dice notation ("1d10")   |
| `result`     | int    | Die result               |
| `success`    | bool   | True if result >= 6      |

### Example

```bash
curl -X POST "http://localhost:8080/resolve?system=vtmv5&action=check" \
  -H "Content-Type: application/json" \
  -d '{}'
```

```json
{
  "expression": "1d10",
  "result": 7,
  "success": true
}
```
