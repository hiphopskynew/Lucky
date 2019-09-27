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

#### Status: 401

```json
{
    "error": {
        "message": "unauthorized access"
    }
}
```

#### Status: 401

```json
{
    "error": {
        "message": "token is invalid"
    }
}
```