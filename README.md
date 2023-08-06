

## Assumptions

- City name is unique
- Number of aliens cannot exceed the number of cities.
- At this solution, collision check performed right after the move.
  - This leads to case when only 2 aliens can end up at the same city. 
- This solution assumes that each city should contain correct connections on both sides
    - City A should contain road to city B if city B contains direction to city A.