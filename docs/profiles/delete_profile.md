## Delete Profile

### URL
- ```DELETE /api/users/{id}/profiles/{id}```

### Json Response

#### Status: 200

```json
{
    "data": {
        "deleted": true
    }
}
```

#### Status: 200

```json
{
    "data": {
        "deleted": false
    }
}
```