--------------------------------------------------
# TODO
--------------------------------------------------

## ------FOCUSED ITEMS------

[] display last-scraped data value for each competitor
[] YOU WERE IMPLEMENTING THE AUTOMATIC SCRAPING SYSTEM?!?!?!??!?!??!?!?!? FIGURE OUT HOW TO HANDLE THIS IN THE BACKGROUND
[] TEST SCRAPING DATA UPDATING ON THE DASHBOARD WITH NEW SCRAPING DATA AND UPDATED LASTSCRAPED DATE

## ------FRONTEND------

[] update/modify dashboard functionality
    [] dates are not populated correctly for last scraped or last change
    [] price data still shows placeholder prices.
    [] email alerts toggle doesn't work.
    [] if no price changes have occured: show default 0%
    [] if price changes have occured, show the difference by a percentage
    [] display competitor name for each competitor
    [] sort buttons don't work
[] make the website responsive

## ------BACKEND------

[] when a user's token expires and they navigate to the homepage, the token should be wiped to have login state change
    [] set middleware to check on homepage for token.
        [] if a token is in localstorage, check if it's valid
            [] if it's not valid, delete it from localstorage and refresh homepage.
            [] if it IS valid, the user should be logged in.
    [] test token expiry
        [] set token expiry to 5 seconds
        [] login
        [] navigate to homepage
[] internal/auth/models.go: move models to internal/models/models.go (SHOULD IT STAY?)


--------------------------------------------------
## ------COMPLETED------
--------------------------------------------------
[x] modify frontend/src/components/Header
    ***NAV-BAR ITEMS SHOULD NOT CALL THE API DIRECTLY!!***
    ***NEVERMIND! THEY DON'T AS LONG AS THEY'RE NOT ROUTED ON THE BACKEND!!!! (so it's all good bro)***
    [x] change the login route to redirect to Login.tsx page
        [x] Create Login Page with option to signup
[x] ADD LOGOUT FUNCTIONALITY
    [x] Add responsive login/logout buttons.
[x] CREATE COMPONENTS!!!!!!
    [x] Create a Login Portal component
    [x] Create a Sign Up Portal component
[x] design a layout for the dashboard on canva
[x] organize all models into one models.go file
    [x] Double check every file for models
    [x] internal/auth/auth.go: move Users model into internal/models/models.go
    [x] internal/storage/postgres.go: move Price model into internal/models/models.go
    [x] internal/ai/openai.go: move model
[x] style dashboard to match canva mockup
    [x] competitors table with data
    [x] side panel with name, email and email alerts toggle 
        - see: https://getbootstrap.com/docs/5.3/forms/checks-radios/ 
[x] style homepage to match canva mockup
[x] add a logo to the header that brings you back to the homepage
    [x] design a logo
    [x] generate png/svg/jpg whatever
    [x] add to assets
    [x] link to header
[x] find a better font
[x] user name on side panel shows placeholder name
    [x] create firstname and lastname fields for User model
    [x] create firstname and lastname fields for User table
    [x] create firstname and lastname form inputs on frontend when signing up.
    [x] modify signup handler to accept a firstname and lastname.
    [x] modify database function query to insert firstname and lastname
    [x] retrieve user data from the database to populate dashboard with name.
        [x] create a handler for retrieving user data
        [x] create database function to retreive user firstname and lastname
[x] change placeholder email to populate with user's actual email
[x] there is no option to enter a competitor name when adding a competitor
    [x] add name to competitor data model so that you can send it to the frontend
