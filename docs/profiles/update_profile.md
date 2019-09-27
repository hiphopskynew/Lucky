## Update Profile

### URL
- ```PUT /api/users/{id}/profiles/{id}```

### Json Request
```json
{
	"first_name": "siggy",
	"last_name": "sig",
	"date_of_birth": "12/12/1993",
	"address": "555 Thoundsand Road"
}
```

### Json Response

#### Status: 200

```json
{
    "data": {
        "address": "555",
        "date_of_birth": "12/12/1993",
        "first_name": "siggy",
        "id": "prof:4ab8bfb0-1abc-4595-ac02-22e1676e7c38",
        "last_name": "sig",
        "user_id": ""
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

#### Status: 400

```json
{
    "error": {
        "message": "date of birth does not match"
    }
}
```

#### Status: 404

```json
{
    "error": {
        "message": "profile does not exist"
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