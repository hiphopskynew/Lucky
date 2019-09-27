## Login

### URL
- ```POST /api/users/{id}/login```

### Json Request
```json
{
	"email": "siggy@gmail.com",
	"password": "12345abcd"
	
}
```

### Json Response

#### Status: 200

```json
{
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjEsIm5hbWUiOiJcInVzZXI6YjQ5ZWViZWUtYTBmZi00NTYwLWE4NmEtNDU3ZmI5M2FlZjA3XCIiLCJleHAiOjE1Njk0ODQxODN9.AWnZnLGosVi0KTzAGD4simvtF9TBLxeVITc-6-nNwYI"
    }
}
```

### Json Error Cases

#### Status: 404

```json
{
    "error": {
        "message": "user not found"
    }
}
```

#### Status: 400

```json
{
    "error": {
        "message": "email or password is invalid"
    }
}
```