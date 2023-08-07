# Alien Invasion

Alien invasion simulation game.

## Assumptions

- City name is unique and case-sensitive.
- Number of aliens cannot exceed the number of cities.
- This solution assumes that each city should contain correct connections on both sides
  - City A should connect to city B and city B connects to city A.
  - Each direction from and to a city is one-to-one only. Many to many is unsupported.
- At this solution, collision check performed right after the move.
  - This leads to case when only 2 aliens max can end up in the same city.

## Usage

```shell
# Run game with 20 aliens and 'us-cities' map file.

go run ./cmd/invasion -c 20 -f ./dataset/us-cities.txt
```

Datasets are available in [datasets](dataset) directory.

Use `-v` flag for verbose debug output. Debug log will be written to stderr.

Use `--help` to get brief about all flags.


## Tests

```shell
make test

# With coverage
make cover
```
 
