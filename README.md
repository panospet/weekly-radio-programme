# Weekly Radio Programme API

Simple mechanism for keeping track of radio shows and avoid conflicts.

## Endpoints

- `GET /health` to just see if the application is up and running
- `GET /show` to get all shows
- `GET /show/{id}` to get a show with this specific integer id
- `POST /show` to create a show
- `PUT /show/{id}` to update a show
- `DELETE /show/{id}` to delete a show

`POST` and `PUT` request body for a show should be like:

```shell
{
    "title": "another radio show",
    "weekday": "Sun",
    "timeslot": "14:00-16:00",
    "description": "best radio show ever baby"
}
```

- Timeslot should have `hh:mm-hh:mm` format.
- Weekday can have one of `Mon`, `Tue`, `Wed`, `Thu`, `Fri`, `Sat`, `Sun` values.
- Title cannot have more than `100` characters.

## How to run

### Locally

#### Start Db

To start db, run:

```shell
make db-start
```

Please add a `.env` file which should have `DB_PATH` and `PORT` variables. If no `PORT` is chosen, `6000` is used by
default.

`.env` file should look like

```shell
DB_PATH=postgres://127.0.0.1/weeklyprogrammedb?sslmode=disable&user=admin&password=password
PORT=6000
```

obviously the above vars change according to your configuration.

#### Docker is a prerequisite

I assume that you have Docker installed. If not, check [here](https://docs.docker.com/engine/install/ubuntu/).

#### Run using Go

I assume that you have Go installed. If not, check [here](https://golang.org/doc/install). Then run

```shell
go mod tidy
go run main.go
```