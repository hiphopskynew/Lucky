## Get Profile

### URL
- ```GET /api/users/{id}/profile```

### Json Response

#### Status: 200

```json
{
    "data": {
        "address": "555 Thounsand Road",
        "date_of_birth": "12/09/1993",
        "first_name": "siggy",
        "id": "user:6b9d92da-4fee-4d34-96fd-0d072a4ea71f",
        "last_name": "sig",
        "user_id": "user:b49eebee-a0ff-4560-a86a-457fb93aef07"
    }
}
```

### Json Error Cases

#### Status: 404

```json
{
    "error": {
        "message": "profile not found"
    }
}
```