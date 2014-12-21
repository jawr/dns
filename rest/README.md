API
===
The API route can be found at  `/api/v1/`.

Searching
---------
Endpoints that return an array can be paginated using `limit` and `page` query parameters. Searching on object fields can be done by using the fields name, i.e. `?name=example`. Foreign objects are refered to by their name, i.e. `tld=com` or `tld=co.uk``.

By default Search endpoints are limited to 15 results with a maximum of 50 results. Pagination can be done using the `page` parameter, i.e. `?limit=50&page=1` -> `?limit=50&page=10`.

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
| GET | /domain/{uuid}/whois | Array | Return an array of Whois objects filtered by Domain. Searchable. |
| GET | /domain/{uuid}/records | Array | Return an array of Records objects filtered by Domain. Searchable. |
| POST | /domain/{uuid}/whois/ | Array | Creates a new Whois instance. Returns Domain's Whois records. |
| GET | /domain/{name} | Object | Return an instance of a Domain by Name.TLD, i.e. 'google.biz'. |
| GET | /domain/{name}/whois | Array | Return an array of Record objects filtered by Domain. Searchable. |
| GET | /domain/{name}/records | Array | Return an array of Records objects filtered by Domain. Searchable. |
| GET | /domain/query/emails | Array | Return an array of Domain objects filtered by related Whois document emails. |

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
	id: 1,
	domain: {
		uuid: "df1c4cea-6bb8-5e21-a467-d2f180042de0",
		name: "taqaglobal",
		tld: {
			id: 43,
			name: "biz"
		}
	},
	data: {
		id: [
			"D48586522-BIZ"
		],
		status: [
			"clientTransferProhibited"
		],
		registrar: [
			"Network Solutions INC."
		],
		nameservers: [
			"ns11.worldnic.com",
			"ns12.worldnic.com"
		],
		updated_date: [
			"2012-01-23T12:43:43",
			"2014-12-21T18:55:01"
		],
		creation_date: [
			"2012-01-23T12:42:36"
		],
		expiration_date: [
			"2017-01-22T23:59:59"
		]
	},
	raw: [
		"Domain Name: TAQAGLOBAL.BIZ Domain ID: D48586522-BIZ Sponsoring Registrar: NETWORK SOLUTIONS INC. Sponsoring Registrar IANA ID: 2 Registrar URL (registration services): whois.biz Domain Status: clientTransferProhibited Registrant ID: 43438924V Registrant Name: Perfect Privacy, LLC Registrant Address1: 12808 Gran Bay Parkway West Registrant Address2: care of Network Solutions Registrant City: Jacksonville Registrant State/Province: FL Registrant Postal Code: 32258 Registrant Country: United States Registrant Country Code: US Registrant Phone Number: +1.5707088780 Registrant Email: c872446c4kx@networksolutionsprivateregistration.com Administrative Contact ID: 43438924V Administrative Contact Name: Perfect Privacy, LLC Administrative Contact Address1: 12808 Gran Bay Parkway West Administrative Contact Address2: care of Network Solutions Administrative Contact City: Jacksonville Administrative Contact State/Province: FL Administrative Contact Postal Code: 32258 Administrative Contact Country: United States Administrative Contact Country Code: US Administrative Contact Phone Number: +1.5707088780 Administrative Contact Email: c872446c4kx@networksolutionsprivateregistration.com Billing Contact ID: 41251237V Billing Contact Name: Perfect Privacy, LLC Billing Contact Organization: TAQA Billing Contact Address1: 12808 Gran Bay Parkway West Billing Contact Address2: care of Network Solutions Billing Contact City: Jacksonville Billing Contact State/Province: FL Billing Contact Postal Code: 32258 Billing Contact Country: United States Billing Contact Country Code: US Billing Contact Phone Number: +1.5707088780 Billing Contact Email: ha9729574er@networksolutionsprivateregistration.com Technical Contact ID: 43438924V Technical Contact Name: Perfect Privacy, LLC Technical Contact Address1: 12808 Gran Bay Parkway West Technical Contact Address2: care of Network Solutions Technical Contact City: Jacksonville Technical Contact State/Province: FL Technical Contact Postal Code: 32258 Technical Contact Country: United States Technical Contact Country Code: US Technical Contact Phone Number: +1.5707088780 Technical Contact Email: c872446c4kx@networksolutionsprivateregistration.com Name Server: NS11.WORLDNIC.COM Name Server: NS12.WORLDNIC.COM Created by Registrar: NETWORK SOLUTIONS INC. Last Updated by Registrar: NETWORK SOLUTIONS INC. Domain Registration Date: Mon Jan 23 12:42:36 GMT 2012 Domain Expiration Date: Sun Jan 22 23:59:59 GMT 2017 Domain Last Updated Date: Mon Jan 23 12:43:43 GMT 2012 >>>> Whois database was last updated on: Sun Dec 21 18:55:01 GMT 2014 <<<< NeuStar, Inc., the Registry Operator for .BIZ, has collected this information for the WHOIS database through an ICANN-Accredited Registrar. This information is provided to you for informational purposes only and is designed to assist persons in determining contents of a domain name registration record in the NeuStar registry database. NeuStar makes this information available to you "as is" and does not guarantee its accuracy. By submitting a WHOIS query, you agree that you will use this data only for lawful purposes and that, under no circumstances will you use this data: (1) to allow, enable, or otherwise support the transmission of mass unsolicited, commercial advertising or solicitations via direct mail, electronic mail, or by telephone; (2) in contravention of any applicable data and privacy protection acts; or (3) to enable high volume, automated, electronic processes that apply to the registry (or its systems). Compilation, repackaging, dissemination, or other use of the WHOIS database in its entirety, or of a substantial portion thereof, is not allowed without NeuStar's prior written permission. NeuStar reserves the right to modify or change these conditions at any time without prior or subsequent notification of any kind. By executing this query, in any manner whatsoever, you agree to abide by these terms. NOTE: FAILURE TO LOCATE A RECORD IN THE WHOIS DATABASE IS NOT INDICATIVE OF THE AVAILABILITY OF A DOMAIN NAME. "
	],
	contacts: {
		tech: {
			city: "Jacksonville",
			name: "Perfect Privacy, LLC",
			email: "c872446c4kx@networksolutionsprivateregistration.com",
			phone: "+1.5707088780",
			state: "Florida",
			handle: "43438924V",
			street: "12808 Gran Bay Parkway West care of Network Solutions",
			country: "United States",
			postalcode: "32258"
		},
		admin: {
			city: "Jacksonville",
			name: "Perfect Privacy, LLC",
			email: "c872446c4kx@networksolutionsprivateregistration.com",
			phone: "+1.5707088780",
			state: "Florida",
			handle: "43438924V",
			street: "12808 Gran Bay Parkway West care of Network Solutions",
			country: "United States",
			postalcode: "32258"
		},
		billing: {
			city: "Jacksonville",
			name: "Perfect Privacy, LLC",
			email: "ha9729574er@networksolutionsprivateregistration.com",
			phone: "+1.5707088780",
			state: "Florida",
			handle: "41251237V",
			street: "12808 Gran Bay Parkway West care of Network Solutions",
			country: "United States",
			postalcode: "32258",
			organization: "Taqa"
		},
		registrant: {
			city: "Jacksonville",
			name: "Perfect Privacy, LLC",
			email: "c872446c4kx@networksolutionsprivateregistration.com",
			phone: "+1.5707088780",
			state: "Florida",
			handle: "43438924V",
			street: "12808 Gran Bay Parkway West care of Network Solutions",
			country: "United States",
			postalcode: "32258"
		}
	},
	emails: [
		"c872446c4kx@networksolutionsprivateregistration.com",
		"ha9729574er@networksolutionsprivateregistration.com"
	],
	added: "2014-12-21T19:57:53.548847+01:00"
}
```
The `data` field contains a base64 encoded JSON object that is taken from [python-whois](https://github.com/joepie91/python-whois) output.

| Method | Endpoint | Return | Description |
| -------- | ------ | ------- | --------- |
| GET | /whois  | Array | Return an array of Whois objects. Searchable. |
| POST | /whois  | Object | Create a Whois instance. Post accepts JSON parameters: `{"domain": "<uuid>"}` or `{"query": "<query>"}`. Return a Whois object.|
| GET | /whois/{id} | Object | Return an instance of a Whois by ID. |
| POST | /whois/query/ | Array | Perform a Whois query. See below for more information. |

### Query
You can perform queries on the `/whois/query/` endpoint, the endpoint takes a JSON object where you define your query:

```
{
	"email": "<email>"
}
```
