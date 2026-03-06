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

| Field                  | Type     | Description                                |
|------------------------|----------|--------------------------------------------|
| `record_id`            | string   | record id of this roll                     |
| `payload.expression`   | string   | Dice notation (e.g., "5d10")               |
| `payload.main_roll`    | int[]    | Regular dice results                       |
| `payload.hunger_roll`  | int[]    | Hunger dice results                        |
| `payload.successes`    | int      | Total successes (including critical bonus) |
| `payload.success`      | bool     | Whether roll met the target                |
| `payload.is_critical`  | bool     | Whether a critical occurred                |
| `payload.crit_type`    | string   | Critical type (see below)                  |

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
  "record_id": "923a67f6-a421-41cf-b323-f71c0e7caffc",
  "payload": {
    "expression": "5d10",
    "main_roll": [7, 4, 10],
    "hunger_roll": [10, 2],
    "successes": 5,
    "success": true,
    "is_critical": true,
    "crit_type": "messy critical"
  }
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

| Field                       | Type   | Description                     |
|-----------------------------|--------|---------------------------------|
| `record_id`                 | string | record id of this roll          |
| `payload.expression`        | string | Original dice notation          |
| `payload.reroll_expression` | string | Rerolled dice notation          |
| `payload.main_roll`         | int[]  | Updated main dice results       |
| `payload.hunger_roll`       | int[]  | Unchanged hunger dice results   |
| `payload.successes`         | int    | New total successes             |
| `payload.success`           | bool   | Whether roll met the target     |
| `payload.is_critical`       | bool   | Whether a critical occurred     |
| `payload.crit_type`         | string | Critical type (see roll action) |

### Example

```bash
curl -X POST "http://localhost:8080/resolve?system=vtmv5&action=reroll" \
  -H "Content-Type: application/json" \
  -d '{
    "record_id": "923a67f6-a421-41cf-b323-f71c0e7caffc",
    "reroll_index": [0, 1]
  }'
```

```json
{
  "record_id": "17ea1c08-78f3-4ebc-b69e-ec344a5bf3f3",
  "payload": {
    "expression": "5d10",
    "reroll_expression": "0d10",
    "main_roll": [7, 4, 10],
    "hunger_roll": [10, 2],
    "successes": 5,
    "success": true,
    "is_critical": true,
    "crit_type": "messy critical"
  }
}
```

---

## check

Simple one-die check (rouse check, frenzy check, etc.). Success on 6+.

### Request

Empty body required.

### Response

| Field                | Type   | Description           |
|----------------------|--------|-----------------------|
| `record_id`          | string | record id of this roll|
| `payload.expression` | string | Dice notation ("1d10") |
| `payload.result`     | int    | Die result            |
| `payload.success`    | bool   | True if result >= 6   |

### Example

```bash
curl -X POST "http://localhost:8080/resolve?system=vtmv5&action=check" \
  -H "Content-Type: application/json" \
  -d '{}'
```

```json
{
  "record_id": "dedd2ad9-8b72-4f58-a022-bb41be17f6dd",
  "payload": {
    "expression": "1d10",
    "result": 8,
    "success": true
  }
}
```
