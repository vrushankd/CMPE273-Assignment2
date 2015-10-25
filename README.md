# CMPE273-Assignment2
Author : Vrushank J Doshi.

RESTful services (CRUD operation) using GO and MongoDB

There are 2 files:
1. server.go - main is in here and imports "controllers" package.
2. controllers.go - Contains all CRUD operations.

MongoLab:
There are 2 collections:
1.cmpe273Assgn2 - Keeps track of all documents.
2.counter - Keeps track of "_id" for autoincrement.

How to run ?
1.Create a controllers directory at location where all your Go packages reside.
2.Save "controllers.go" file in your controllers directory.
3.Run server.go. 
   go run server.go
