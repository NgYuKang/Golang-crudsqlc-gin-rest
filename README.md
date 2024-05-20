## Go SQLC Gin Self Learning Project
This is a project where I used Go's Gin module to create a rest API, using pgx to call postgresql, and sqlc to generate db calling code and handle migrations<br>
The List route is not efficient because it is using offset to paginate, whereas we should use unique value sorting and filtering.
<br>Say last page value was ID 4, and it was sorted by the column of price descending with price 40, we would filter by an or filter of <br>
price < 40 OR price == 40 and ID > 4
<br>
<br>
The generated SQLC code is not able to handle dynamic query, so in my opinion it is better to just write your own service code and sql.