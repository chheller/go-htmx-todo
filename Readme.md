# Structure

Main acts as the IoC container, and wires up most things (namely, database configurations, SMTP connections, etc)

router inits handlers and wires up their services
handlers contain the actual endpoints
web module has gohtml template files and corresponding go files for parsing and executing those templates

# Running
- Install [air](https://github.com/cosmtrek/air)
- Run `air` 
# Todo

- Add tests for things like 404s, 405s, 500s
- Test for mocking email verification
- Test for full user flow
- Complete user flow
  - Create a post signup page with instructions to check email
  - Create sign-in page via Email OTP
- Create Todo page
  - View List of Todos
  - Add Todo
  - Complete Todo
  - Delete Todo
- Add a generic application error type and extend all app errors from it
- Enhance application logging
