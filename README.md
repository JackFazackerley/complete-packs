# Complete Packs
A home project for an ordering system

## What it does
Complete packs is an API that will give the best solution for pack sizes for a given target, where target is number of items.

It has three rules that that will decide the pack sizes and amount you will get.

1. Only whole packs can be sent. Packs cannot be broken open.
2. Within the constraints of Rule 1 above, send out no more items than necessary to fulfil the
order.
3. Within the constraints of Rules 1 & 2 above, send out as few packs as possible to fulfil each
order.

An example of this could be, given the pack sizes: `250,500,1000,2000,5000` and target items as `1`. The best solution would be `1x250` - which shows the use of rule 2.

Given the same pack sizes as above, an example of rule 3 would be given the target `251` the best solution is `1x500` as `2x250` would be giving more packs than necessary.


## Orders API 
The orders API has two endpoints `/order/best` and `/order/fast`. The difference between these two is the way they compute the answer. 

## Packs API
The packs API has three endpoints `packs/read`, `packs/write`, and `packs/delete`.

## Calulcation methods

### Best
[Best pkg](pkg/order/best)

Done in a breadth first manner to guarantee minimal packs are used for each subproblem

Best will compute the answer by starting from the target and subtracting each size from it. If the result of `target - size` is `> 0` and it's the first time we've seen this result we add the result to a queue and store the size that got it to the result. 

The reason we only care about storing the result on the first occurrence is because the first occurrence is the fastest way to get to that result, so if we see it again we ignore it.

We keep adding to the queue until all result are negative.

#### Issues
Best will always give the best answer for the target and even though it will only go for as long as needed if the number is big enough it will use a LOT of memory. So if wanting a big number expect to use a lot of memory.

### Fast
[Fast pkg](pkg/order/fast)

Done by checking the amount of times target goes into size `target/sizes[i]`. It then checks to see if `target > sizes[i+1]` so we know if to use mod or max. We then check to see if the divided result is greater than previous and if the count is lesser than the previous count.   

Fast will compute the answer by starting at the target value and finding the first number which divides in the most amount of times with the lowest count. The remainder is then returned and the new remainder is the passed back to the same function and is called again until we reach 0. 

Fast needs to be called twice this is because the first pass will find the total amount of items but not necessarily the best amount of packs. So we pass the total amount of items back and this time it will give the best result.

#### Issues
Fast is well fast, the only issue with it is that each size needs to be a multiplication of each other. If this isn't the case it will sometimes give the wrong answer. A good example is by having sizes `[499,500]` and the target `998` the best answer is `2x499` but it gives `1x499,1x500`. This is because of the 499 not being a multiple of 500.

## SQLite
The database of choice to store the sizes in. It doesn't need the bloat of a full database. This is a proof of concept and doesn't need to be scaled, so sqlite is a fitting choice.

To setup the database run the following command:
```
sqlite3 packs.db < setup.sql 
```
 
Make sure that the `config.yaml` file is pointing to the correct path for the db.

## Running locally
To the run API server simply run these commands: 

```
go mod vendor
go run -mod vendor cmd/api/main.go
``` 

To run the API:
```
cd cmd/ui
npm install
npm run serve
```