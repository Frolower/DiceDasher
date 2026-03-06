# Generic System

Universal dice roller for NdM notation.

## Actions

- `roll` — Roll N dice of size M

---

## roll

Roll N dice of size M and return individual results with sum.

### Request

| Field    | Type | Required | Constraints | Description       |
|----------|------|----------|-------------|-------------------|
| `number` | int  | Yes      | >= 1        | Number of dice    |
| `size`   | int  | Yes      | >= 2        | Die size (faces)  |

### Response

| Field                | Type   | Description                   |
|----------------------|--------|-------------------------------|
| `record_id`          | string | record id of this roll        |
| `payload.expression` | string | Dice notation (e.g., "2d6")   |
| `payload.rolls`      | int[]  | Individual roll results       |
| `payload.sum`        | int    | Sum of all rolls              |

### Example

```bash
curl -X POST "http://localhost:8080/resolve?system=generic&action=roll" \
  -H "Content-Type: application/json" \
  -d '{"number": 3, "size": 6}'
```

```json
{
  "record_id": "23d5205e-dd90-4c7e-8347-25f7b15605f3",
  "payload": {
    "expression": "3d6",
    "rolls": [3, 1, 2],
    "sum": 6
  }
}
```
