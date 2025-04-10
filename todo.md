--------------------------------------------------
# TODO
--------------------------------------------------

## ------FOCUSED ITEMS------
[] style dashboard to match canva mockup
    [x] competitors table with data
    [] side panel with name, email and email alerts button


## ------FRONTEND------
[] dashboard functionality
    [] if no price changes have occured, show that
    [] if price changes have occured, show the difference by a percentage
    [] display last-scraped data value for each competitor
    [] display competitor name for each competitor
[] make the website responsive
[] add a logo to the header that brings you back to the homepage
[] style homepage to match canva mockup


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
[] add name to competitor data model so that you can send it to the frontend
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
