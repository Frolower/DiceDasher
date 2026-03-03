# TES System

The Electric State by Free League

- Uses d6 dice
- Successes only on 6
- Attribute and gear dice tracked separately
- Push mechanic re-rolls non-1, non-6 dice with consequences

## Actions

- `roll` — Standard dice pool roll
- `push` — Push a previous roll

---

## roll

Perform a standard TES dice pool roll.

### Request

| Field        | Type | Required | Constraints                   | Description                   |
|--------------|------|----------|-------------------------------|-------------------------------|
| `attr`       | int  | Yes      | >= 1                          | Attribute score               |
| `assist`     | int  | Yes      | 0-3                           | Assisting players count       |
| `gear`       | int  | Yes      | >= 0                          | Gear bonus dice               |
| `modificator`| int  | Yes      | total pool + modificator >= 1 | GM modifier (can be negative) |
| `target`     | int  | Yes      | 1 to (attr + assist + gear)   | Required successes            |

### Response

| Field            | Type     | Description                     |
|------------------|----------|---------------------------------|
| `expression`     | string   | Dice notation (e.g., "6d6")     |
| `attribute_rolls`| int[]    | Attribute + assist dice results |
| `gear_rolls`     | int[]    | Gear dice results               |
| `successes`      | int      | Total 6s rolled                 |
| `success`        | bool     | Whether roll met the target     |

### Example

```bash
curl -X POST "http://localhost:8080/resolve?system=tes&action=roll" \
  -H "Content-Type: application/json" \
  -d '{"attr": 4, "assist": 1, "gear": 2, "modificator": 0, "target": 2}'
```

```json
{
  "expression": "7d6",
  "attribute_rolls": [3, 6, 2, 5, 6],
  "gear_rolls": [4, 1],
  "successes": 2,
  "success": true
}
```

---

## push

Push a previous roll. Re-rolls all dice that aren't 1 or 6. Tracks hope loss and gear damage from new 1s.

### Request

| Field            | Type  | Required | Constraints                 | Description                |
|------------------|-------|----------|-----------------------------|----------------------------|
| `attribute_rolls`| int[] | Yes      | length >= 1                 | Previous attribute results |
| `gear_rolls`     | int[] | Yes      | —                           | Previous gear results      |
| `target`         | int   | Yes      | 0 to attribute_rolls length | Required successes         |

### Response

| Field            | Type     | Description                        |
|------------------|----------|------------------------------------|
| `expression`     | string   | Original dice notation             |
| `push_expression`| string   | Pushed dice notation               |
| `attribute_rolls`| int[]    | Updated attribute dice results     |
| `gear_rolls`     | int[]    | Updated gear dice results          |
| `successes`      | int      | Total 6s rolled                    |
| `success`        | bool     | Whether roll met the target        |
| `hope_losses`    | int      | Number of 1s on attribute dice     |
| `gear_damage`    | int      | Number of 1s on gear dice          |

### Example

```bash
curl -X POST "http://localhost:8080/resolve?system=tes&action=push" \
  -H "Content-Type: application/json" \
  -d '{
    "attribute_rolls": [3, 6, 2, 5, 6],
    "gear_rolls": [4, 1],
    "target": 3
  }'
```

```json
{
  "expression": "7d6",
  "push_expression": "4d6",
  "attribute_rolls": [1, 6, 6, 4, 6],
  "gear_rolls": [3, 1],
  "successes": 3,
  "success": true,
  "hope_losses": 1,
  "gear_damage": 1
}
```
