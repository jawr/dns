API
===
The API route can be found at  `/api/v1/`

Searching
---------
Endpoints that return an array can be paginated using `limit` and `page` query parameters. Searching on object fields can be done by using the fields name, i.e. `?name=example`. Foreign objects should be refered to by their fk, i.e. `?tld=1`

Domain
------
```
{
	"name": "github", 
	"tld": {
		"id": 1, 
		"name": "biz"
	}, 
	"uuid": "5e79e57e-0e46-5e34-801d-6dc8e1c872f1"
}
```

| Method | Endpoint | Return | Description |
| -------- | ------ | ------- | --------- |
| GET | /domain  | Array | Return an array of Domain objects. Searchable. |
| GET | /domain/{uuid} | Object | Return an instance of a Domain by UUID. |

Record
------
```
{
	"added": "2014-12-13T15:48:24.391041+01:00", 
	"args": {
		"args": [
				"127.0.0.1"
			], 
		"ttl": 900
	}, 
	"domain": {
		"name": "google", 
		"tld": {
			"id": 1, 
			"name": "biz"
		}, 
		"uuid": "c36bee42-d2e9-51df-a694-b0dcdba886bf"
	}, 
	"name": "ns1", 
	"parse_date": "2014-06-22T00:00:00Z", 
	"type": {
		"id": 1, 
		"name": "a"
	}, 
	"uuid": "00148631-e798-5328-91ee-e4f1da1b74be"
}

```

| Method | Endpoint | Return | Description |
| -------- | ------ | ------- | --------- |
| GET | /record  | Array | Return an array of Record objects. Searchable. |
| GET | /record/{uuid} | Object | Return an instance of a Record by UUID. |

Record Type
-----------
```
{
	"id": 3, 
	"name": "ns"
}
```

| Method | Endpoint | Return | Description |
| -------- | ------ | ------- | --------- |
| GET | /record_type  | Array | Return an array of RecordType objects. Searchable. |
| GET | /record_type/{id} | Object | Return an instance of a RecordType by ID. |
| GET | /record_type/{name} | Object | Return an instance of a RecordType by Name. |


TLD
---
```
{
	"id": 1, 
	"name": "biz"
}
```

| Method | Endpoint | Return | Description |
| -------- | ------ | ------- | --------- |
| GET | /tld  | Array | Return an array of TLD objects. Searchable. |
| GET | /tld/{id} | Object | Return an instance of a TLD by ID. |
| GET | /tld/{name} | Object | Return an instance of a TLD by Name. |

Whois
-----
```
{
	"added": "2014-12-13T20:49:08.126472+01:00", 
	"data": "eyJzdGF0dXMiOiBbIm9rIl0sICJ1cGRh=",
	"domain": {
		"name": "httk", 
		"tld": {
			"id": 1, 
			"name": "biz"
		}, 
		"uuid": "df16d75f-7d86-51dd-9951-4b19e723a6d2"
	}, 
	"id": 1
}
```

| Method | Endpoint | Return | Description |
| -------- | ------ | ------- | --------- |
| GET | /whois  | Array | Return an array of Whois objects. Searchable. |
| GET | /whois/{id} | Object | Return an instance of a Whois by ID. |
