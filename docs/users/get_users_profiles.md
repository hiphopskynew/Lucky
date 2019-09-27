## Get User By ID

### URL
- ```GET /api/users/{id}```

### Json Response

#### Status: 200

```json
{
    "data": [
        {
            "address": "555",
            "created_at": "2019-09-26T06:44:59Z",
            "date_of_birth": "12/12/1993",
            "email": "siggy@gmail.com",
            "first_name": "siggy",
            "id": "user:b49eebee-a0ff-4560-a86a-457fb93aef07",
            "last_name": "sig",
            "profile_id": "user:4ab8bfb0-1abc-4595-ac02-22e1676e7c38",
            "status": "verified",
            "updated_at": "2019-09-26T06:44:59Z"
        },
        {
            "address": "",
            "created_at": "2019-09-27T06:54:15Z",
            "date_of_birth": "",
            "email": "siggy2@gmail.com",
            "first_name": "",
            "id": "user:108fceb5-3bac-46b2-9ec0-2b519f9845a7",
            "last_name": "",
            "profile_id": "",
            "status": "verified",
            "updated_at": "2019-09-27T06:54:15Z"
        }
    ]
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