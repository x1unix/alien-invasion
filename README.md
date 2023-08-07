

## Assumptions

- City name is unique and case-sensitive.
- Number of aliens cannot exceed the number of cities.
- This solution assumes that each city should contain correct connections on both sides
    - City A should contain road to city B if city B contains direction to city A.
    - Each direction from and to a city is one-to-one only. Many to many is unsupported.
- At this solution, collision check performed right after the move.
  - This leads to case when only 2 aliens can end up at the same city. 
