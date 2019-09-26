## Verify

### URL
- ```POST /api/users/{id}/verify```

### Json Request
```json
{
	"token": "40b269170df045f4a264b5341a327a3745fb3bf8e5c9475491031a6bfd65eb4ec9d54b8ff29f4fb58598de430c84382c2bf5b330f5a1490cb7b4ca09269b07f6350c9a2bb21a494daf1a6b7b3096a600"
}
```

### Json Response

#### Status: 200

```json
{
    "data": {
        "email": "siggy@gmail.com",
        "status": "verified"
    }
}
```

### Json Error Cases

#### Status: 404

```json
{
    "error": {
        "message": "email not found"
    }
}
```

#### Status: 400

```json
{
    "error": {
        "message": "token is invalid"
    }
}
```
