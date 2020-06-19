# GetBuckets Backend

This is the API server for [GetBuckets](https://github.com/rtravitz/getbuckets). This is a joke with a friend that probably went too far. He goes on long runs, and it turns out he knows where all of the [Honey Buckets](https://www.honeybucket.com/) are located for miles around. This server collects the coordinates of each "bucket" to display on a map. It also collects data about whether a given bucket is locked and it's cleanliness on a 1-5 scale, where 5 is :chef-kiss: incredible.

## Running
Have Postgres listening on port 5432 with a database created by the `postgres` named `getbuckets`. 

With Go installed, run `make migrate` to run database migrations. Then use `make run` or `migrate watch` to start the server.