# VirtualFS

VirtualFS is a web application with a Golang backend and a React frontend to emulate a file system using a database (PostgreSQL).

## Installation
You'll need the following: [`yarn`](https://classic.yarnpkg.com/en/docs/install/), [`Node.js`](https://nodejs.org/en/download/), [`go`](https://golang.org/doc/install), and [`postgresql`](https://www.postgresql.org/download/) installed. 

To build the project:

```bash
$ make build
```

Then create a database in PostgreSQL and provide these environment variables (values given here are the defaults):

Either:
- `DATABASE_URL=''` : the URL used to locate the database of the form `postgresql://user:password@host:port/database_name` 

Or all of the following:
- `DB_USER='postgres'` : User used to login to PSQL
- `DB_PASSWORD='postgres'` : Password of above user
- `DB_HOST='localhost'` : Hostname of the database
- `DB_PORT=5432` : Port to connect to database to
- `DB_NAME='virtualfs_test'` : Name of the database
## Usage
Run the server and frontend with:
```Bash
$ make run
```

Or:
```Bash
$ make build && server/bin/server
```

Then the server can be accessed at `localhost:8080`, and the backend API is at `localhost:8080/api`.
## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
[MIT](./LICENSE)