# Convert csv to sql

## Description
This repo to semi-automate converting any csv file to sql file

## How to use
1. Add the csv file into repo
2. Change `csvFileName` varibale with the csv file name.
3. Change database credentials `host`, `user`, `password` and `dbname`
4. Update `updateQueryTemplate` and Query logic section with what you want.
5. Run the command `$ make run`
6. The file `queries.sql` has been created with all query that you what.
