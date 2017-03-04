## uoft-ace-api

> An API for retrieving real-time data for [Academic and Campus Events](http://www.ace.utoronto.ca/) at the University of Toronto (St. George).

Data source: http://www.ace.utoronto.ca/bookings.

### Installation

  - Requirements:

    - Go should be [installed](https://golang.org/doc/install) and [configured](https://golang.org/doc/install#testing).
    - [Redis](https://redis.io/)

  - Install with [go get](https://golang.org/cmd/go/):
    
    ```sh
    $ go get github.com/kshvmdn/uoft-ace-api
    ```

  - Build from source:

    ```sh
    $ git clone https://github.com/kshvmdn/uoft-ace-api $GOPATH/src/github.com/kshvmdn/uoft-ace-api
    $ cd $_
    $ make
    ```

### Usage

  - Start Redis, Use `--daemonize yes` to run in background.

    ```sh
    $ redis-server [--daemonize yes]
    ```

  - Start the server. The environment variables are optional, the defaults are shown below.

    ```sh
    $ PORT=8080 REDIS_PORT=6379 REDIS_PASSWORD="" REDIS_DB=0 ./uoft-ace-api
    ```

  - **Endpoints**:
    
    - `/calendar/{building:[a-zA-Z]+}`

      - Retrieve a schedule for the current week for all the rooms of the building provided.
      - Output format:

        ```js
        [{
          building_code: String,
          name: String,
          rooms: [{
            room_number: String,
            schedule: [{
              date: String,
              bookings: [{
                time: String,
                description: String
              }]
            }]
          }]
        }]
        ```

      - Example: [`/calendar/BA`](http://localhost:8080/calendar/ba)

    - `/calendar/{building:[a-zA-Z]+}/{room:[0-9a-zA-Z]+}`

      - Retrieve a schedule for the current week for the room provided.
      - Output format:

        ```js
        {
          room_number: String,
          schedule: [{
             date: String,
             bookings: [{
               time: String,
               description: String
             }]
           }]
        }
        ```

      - Example: [`/calendar/BA/1130`](http://localhost:8080/calendar/ba/1130)

### Contribute

This project is completely open source, feel free to [open an issue](https://github.com/kshvmdn/issues) for bugs/requests or [submit a pull request](https://github.com/kshvmdn/pulls)!

##### TODO:

  - [ ] Implement Redis to cache response data (boilerplate is partially implemented).
  - [ ] Add URL parameters to specify dates (possibly date ranges as well).
  - [ ] Determine an efficient way to sort room codes for a given building (right now the order is completely random).
