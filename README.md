# Domain Info API

API for a service that allows us to retrieve information about a domain and know if the settings have changed, based on the [SSL Labs API](https://github.com/ssllabs/ssllabs-scan). 

## Getting started

These instructions will get you a copy of the project up and running on your local machine for development. 

### Prerequisites

These are the things we need to install before cloning the project. 

#### CockroachDB

Follow the instructions found in the [official docs](https://www.cockroachlabs.com/docs/stable/install-cockroachdb-linux.html) to install CockroachDB on your machine. After that, run the following command to get a [single-node cluster](https://www.cockroachlabs.com/docs/stable/cockroach-start-single-node.html) up and running on your machine: 

```
cockroach start-single-node --insecure --listen-addr=localhost --background
```

Then, we need to enter the SQL shell to create our database: 

```
cockroach sql --insecure
```

Finally, run the finally command to create the database: 

```
CREATE DATABASE domain_info_api;
```

#### Whois XML API

As of now, in order to get the data we need from the **whois** command, we need to use the following service. Register to [Whois XML API](https://main.whoisxmlapi.com/login) to get access to their [WHOIS API](https://whois.whoisxmlapi.com/) service. However, it only offers a limited amount of requests per month :(.

**NOTE**: After setting up these two services, make sure to add the connection string and API key into a `.env` file and put it at the root of the project.  

## Installation

Once you have [installed go](https://golang.org/doc/install), run this command to get a copy of the project: 

```
git clone https://github.com/aledeltoro/domain-info-api.git
```

## Built With

* [Fasthttprouter](https://github.com/buaazp/fasthttprouter) - HTTP router used
* [Pq](https://github.com/lib/pq) - Pure Go Postgres driver (required to use CockroachDB)
* [Godotenv](https://github.com/joho/godotenv) - Loads environment variables from `.env`
* [Goquery](https://github.com/PuerkitoBio/goquery) - Allows web scraping for the parts we want from a page
* [Govalidator](https://github.com/asaskevich/govalidator) - A package of validators and sanitizers based on [validator.js](https://github.com/validatorjs/validator.js)

## License

MIT.

