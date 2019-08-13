# Golang booking API

<img src="https://cdn-images-1.medium.com/max/1600/1*yh90bW8jL4f8pOTZTvbzqw.png" width="150px"/>

A small API I'm working on to learn Golang.


`POST /register`

Register an account. Responds with auth token if successful.

`POST /login`

Log in with your username and password. Responds with auth token if successful.

`GET /bookings`

Returns all bookings created by the user

`POST /bookings`

Create a new booking. Currently accepts the fields "name", "date" and "location".

`GET /bookings/:id`

Returns a particular booking made by the user

`PUT /bookings/:id`

Update a booking with the associated booking ID

`DELETE /bookings/:id`

Delete a booking with the associated booking ID

