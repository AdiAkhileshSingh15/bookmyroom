# Bookings and Reservations

This is the repository for my bookings and reservations project.

- Built in Go version 1.20.x
  
Dependencies:

- [chi router](https://github.com/go-chi/chi)
- [alex edwards SCS](https://github.com/alexedwards/scs/v2) session management
- [nosurf](https://github.com/justinas/nosurf)
- [pgx](https://github.com/jackc/pgx/v4)
- [simple mail](https://github.com/xhit/go-simple-mail/v2)
- [Go validator](https://github.com/asaskevich/govalidator)

To check the mailing functionality, you can use the [MailHog](https://github.com/mailhog/MailHog) SMTP testing package.
And run it on a different host on your server, using:
```
~/go/bin/MailHog
```

In order to build and run this application, it is necessary to 
install Soda (go install github.com/gobuffalo/pop/... ), create
a postgres database, fill in the correct values in database.yml, 
and then run soda migrate.

To build and run the application, from the root level of the project,
execute this command:
```
make build && ./bookmyroom -dbname=<database name> -dbuser=<database user> -dbpass=<database password>
```
where you have the correct entires for your database name (dbName) and database user (dbUser)

For the complete list of command flags, run ./bookmyroom -h