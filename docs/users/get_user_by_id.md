## Get User By ID

### URL
- ```GET /api/users/{id}```

### Json Response

#### Status: 200

```json
{
    "data": {
        "created_at": "2019-09-26T06:44:59Z",
        "email": "siggy@gmail.com",
        "id": "user:b49eebee-a0ff-4560-a86a-457fb93aef07",
        "password": "123",
        "status": "verified",
        "updated_at": "2019-09-26T06:44:59Z"
    }
}
```

### Json Error Cases

#### Status: 404

```json
{
    "error": {
        "message": "user does not exist"
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