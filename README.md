# Noughts and Crosses

## How to run the application?

```bash
go build && ./noughts-crosses
```


## Endpoints explained

#### Create a Game

```bash
http://localhost:8080/game
```

The following endpoint takes a POST Request with an empty body to create a game where IDs are generated (for the game and first player) and persisted to a database. 

Response:

```json
{
    "game": {
        "id": 5555555555555,
        "first_player": {
            "id": 444555352252,
        },
        "second_player": null,
        "noughts": null, 
        "crosses": null, 
    }
}

```


#### Join a player

```bash
http://localhost:8080/game/{id}/join
```

The following endpoint takes a POST request with an empty body to generate an ID for the second player and save it to the game with the given ID in the url.


#### Get a Game

```bash
http://localhost:8080/game/{id}
```

The following endpoint takes a GET request and retrieves the status of a game with the given ID in the url.


response:

```json
{
    "game": {
        "id": 5555555555555,
        "first_player": {
            "id": 444555352252,
        },
        "second_player":{
            "id": 5325252525,
        },
        "noughts": [
            {"x": 1, "y":2, "player": {"id": 5254242524}},
            {"x": 1, "y":2, "player": {"id": 5254242524}},
        ],
        "crosses": [
            {"x": 1, "y":2, "player": {"id": 5254242524}},
            {"x": 1, "y":2, "player": {"id": 5254242524}},
        ],
    }
}

```

#### Set a nought or a cross in the Game

```bash
http://localhost:8080/game/{id}/nought
http://localhost:8080/game/{id}/cross
```

The following endpoint takes a POST request with a json in the body:

```json
{
    "nought": {
        "x": 1,
        "y": 1,
        "player": {
            "id": 5432525252,
        },
    },
}
```

The x and y coordinates are represented as follows:

```bash
|_3_|___|___|
|_2_|___|___|
|_1_|_2_|_3_|
```

So if a player's move that consist of putting a nought in the dead center box, would look like this:

```json
{
    "nought": {
        "x": 2,
        "y": 2,
        "player": {
            "id": 5432525252,
        },
    },
}
```


