# Code structure

This application is written in Golang for the backend and TypeScript with React + Material-UI for the frontend, which are in the `server` and `app` folders respectively.

The React SPA is built with `yarn build` into static files, which are then copied to the server through a postinstall script.

The Go server serves both the static files of the SPA and the backend API at `/` and `/api/` respectively.

The application is [hosted on Heroku](https://virtualfs.herokuapp.com).

Decisions made as well as the general code structure for both backend and frontend are detailed below

## Backend

### Database

The database used is PostgreSQL, using `pgx` to interface with the database from go and `pgxpool` to create a connection pool to ensure concurrency.

The database is initialized with the db.pgsql script, which creates the single table `filedirs`, as well as some triggers to update the size of files and directories when they are added/removed/updated.

The hierarchical path information is preserved using PostgreSQL's [`ltree` extension](https://www.postgresql.org/docs/9.1/ltree.html), which is also indexed to speed up search speed.

The database is setup in `server/db/db.go`.

### Models

The models layer export a `FileDir` type struct, which represents a row in the DB.

In `server/api/models/filedirs.go`, the `FileDirsTable` type holds reference to the connection pool to the database, which it can then use for its type methods to interact with the database.

`FileDirsTable`'s type methods are main layer of logic for the backend server. They are where the queries are constructed from incoming data and executed.

`FileDirsTable` is packaged as a singleton in `server/api/models/common.go`, and is used by the controllers to manipulate database data.

### Controllers

The backend has one single controller in `server/api/controllers/controller.go`, which defines the handlers for the REST API routes, checks and sanitizes incoming data, then call the models layer.

### Routing

Routing is done using `go-chi` in `server/routes`.

### Testing

Due to time constraints, only the most crucial and complicated part of the application (`models`) currently have unit tests, and not extensively.

## Frontend

The frontend project skeleton was created with `create-react-app --template typescript`.

### Components

The frontend mainly uses 2 components for its 2 pages: the terminal in `MainTerminal` (created using `react-console-emulator`), and the help page in `Help`. All components are under `app/src/components`.

### Terminal

The terminal UI and command recognition is done using `react-console-emulator`, and the functions that are called for each command are in `app/src/components/MainTerminal/functions`.

### Help Page

The help page was built with `material-ui` components. The files are in `app/src/components/Help`.

### Routing

Routing within the SPA is done with a `BrowserRouter`, as well as the `useHistory()` hook.

### API calls

API calls are done using `axios`, and are wrapped with methods under `app/src/utils/axios.ts`

## Implementation logic

The application implemented all of the logic required, wtih the following modifications/clarifications:

- `mv SOURCE... DEST` treats `DEST` as the parent directory and not the actual file name to move to, as that would involve a rename as well
- Due to the previous point, updating files are not actually supported yet. The code is present in the backend (`Update` function in `FileDirsTable`) but buggy and not reachable since a route hasn't been made for it. Based on the requirements, it does not seem that updating files/directories is needed.
- The `*` operator is only usable in `mv` and `rm`, as they are the only two logical
- Added an additional `cat` command to print file contents to the terminal.

## REST API Reference

Response object format:

```typescript
RESPONSE_DATA: {
    filedirs: {
        id: int,
        name: string,
        size: int,
        data: string,
        created: string,
        updated: string,
        path: string,
        filetype: string,
    }
}
```

```typescript
DELETE_DATA: {
    numrows: int // number of rows deleted
}
```

Request object format:

```typescript
CREATE_DATA: {
    filedirs: { // Object to create
        name: string,
        data: string,
        filetype: string,
    }
    path: string // path to create on
}
```

```typescript
MOVE_DATA: {
    source: string // source file location
    dest: string // destination parent dir
}
```

### Get the data for a specific file/folder

- URL: `/api/:path/exact`
- Method: `GET`
- Constraints: `path` must be encoded in Base64URL
- Response: single `RESPONSE_DATA`

### Get the data for all files under a directory

- URL: `/api/:path`
- Method: `GET`
- Constraints: `path` must be encoded in Base64URL
- Response: List of `RESPONSE_DATA`

### Create new file/folder

- URL: `/api/`
- Method: `POST`
- Request data: `CREATE_DATA`
- Headers: `Content-Type: application/json`
- Response: single `RESPONSE_DATA` of the created object

### Move file/folder 

- URL: `/api/`
- Method: `PUT`
- Request data: `MOVE_DATA`
- Headers: `Content-Type: application/json`
- Response: List of `RESPONSE_DATA` of moved objects

### Delete a file/folder

- URL: `/api/:path`
- Method: `DELETE`
- Constraints: `path` must be encoded in Base64URL
- Response: Single `DELETE_DATA`

# Potential issues
Security:
- There is no authentication required to interact with the file system, so anyone can drop the entire file system by sending `DELETE`requests

Concurrency:
- Most SQL queries used are done in a single query, including those that update multiple nodes at the same time (`mv`), and `UPDATE`s queries are atomic.  
- The biggest concurrency bottleneck on the server are the size calculation triggers that are called every time a new file is inserted/updated/deleted. 
- In order to alleviate this bottleneck, we could lower the update frequency by running updates only when the data is accessed, handle the size updates ourselves in the application code, or remove it altogether and do dynamic size calculations upon request (which results in a lot more READs).

Efficiency:
- Since we are using a data structure built to represent paths that is indexable, queries are very efficient.
- On the final hosted server, our performance bottleneck are the API calls between frontend and backend, as well as the bootup time of the Heroku dyno.
# Future improvements

- More unit tests
- Fix and expose the `Update` function in backend (`PUT /api/`)
- Setup a CI/CD pipeline for deployment after having tests
- More UI/UX improvements in the frontend
- Add a new update file contents command in the frontend
