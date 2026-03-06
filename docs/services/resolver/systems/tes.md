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

| Field                     | Type     | Description                     |
|---------------------------|----------|---------------------------------|
| `record_id`               | string   | record id of this roll          |
| `payload.expression`      | string   | Dice notation (e.g., "6d6")     |
| `payload.attribute_rolls` | int[]    | Attribute + assist dice results |
| `payload.datagear_rolls`  | int[]    | Gear dice results               |
| `payload.successes`       | int      | Total 6s rolled                 |
| `payload.success`         | bool     | Whether roll met the target     |

### Example

```bash
curl -X POST "http://localhost:8080/resolve?system=tes&action=roll" \
  -H "Content-Type: application/json" \
  -d '{"attr": 4, "assist": 1, "gear": 2, "modificator": 0, "target": 2}'
```

```json
{
  "record_id": "90e8e1f0-c213-4fa3-9a64-7d7ce3115e04",
  "payload": {
    "expression": "7d6",
    "attribute_rolls": [6, 6, 1, 6, 3],
    "gear_rolls": [2, 4],
    "successes": 3,
    "success": true
  }
}
```

---

## push

Push a previous roll. Re-rolls all dice that aren't 1 or 6. Tracks hope loss and gear damage from new 1s.

### Request

| Field       | Type | Required | Constraints | Description                                          |
|-------------|------|----------|-------------|------------------------------------------------------|
| `record_id` | uuid | Yes      | not null    | record id of the roll that is stored in the database |

### Response

| Field                     | Type     | Description                    |
|---------------------------|----------|--------------------------------|
| `record_id`               | string   | record id of this roll         |
| `payload.expression`      | string   | Original dice notation         |
| `payload.push_expression` | string   | Pushed dice notation           |
| `payload.attribute_rolls` | int[]    | Updated attribute dice results |
| `payload.gear_rolls`      | int[]    | Updated gear dice results      |
| `payload.successes`       | int      | Total 6s rolled                |
| `payload.success`         | bool     | Whether roll met the target    |
| `payload.hope_losses`     | int      | Number of 1s on attribute dice |
| `payload.gear_damage`     | int      | Number of 1s on gear dice      |

### Example

```bash
curl -X POST "http://localhost:8080/resolve?system=tes&action=push" \
  -H "Content-Type: application/json" \
  -d '{"record_id": "90e8e1f0-c213-4fa3-9a64-7d7ce3115e04"}'
```

```json
{
  "record_id": "704d162c-dc92-404c-bb4b-2c48e22df290",
  "payload": {
    "expression": "7d6",
    "push_expression": "3d6",
    "attribute_rolls": [6, 6, 1, 6, 1],
    "gear_rolls": [3, 5],
    "successes": 3,
    "success": true,
    "hope_losses": 2,
    "gear_damage": 0
  }
}
```
