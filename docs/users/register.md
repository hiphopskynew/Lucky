## Register

### URL
- ```POST /api/users/{id}/register```

### Json Request
```json
{
	"email": "siggy@gmail.com",
	"password": "12345abcd"
	
}
```

### Json Response

#### Status: 201

```json
{
    "data": {
        "created_at": "2019-09-26T13:42:47.881144875+07:00",
        "email": "siggy@gmail.com",
        "id": "user:381e44bf-5ae9-4ec0-a571-6a3111c496b5",
        "status": "new",
        "token": "78832c7777704c6eae55d45b8eadb014f08abfdb628b4c66aec3345018c7c69441958f9feccb41738a9efe31ad2fe7f847c1a0f2359641fbbab1253a29ac88300f8feeb0dde04ca4ba4549adf517242d",
        "updated_at": "2019-09-26T13:42:47.881144952+07:00"
    }
}
```

### Json Error Cases

#### Status: 404

```json
{
    "error": {
        "message": "email already exist"
    }
}
```

#### Status: 400

##### Fields validate

- email
    - required
    - email format
    - max length 50
- password
    - required
    - min length 8
    - max length 140

```json
{
    "error": [
        {
            "key": "???",
            "messages": [
                "?????"
            ]
        }
    ]
}