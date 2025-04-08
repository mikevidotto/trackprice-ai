# TODO
--------------------------------------------------

## ------FOCUSED ITEMS------



## ------FRONTEND------

[] create a page layout that you can put your react components into (use App.tsx??)
    [] add layout components
        [] header
        [] nav-bar
        [] body
        [] footer
    [] make it responsive
[] add a logo to the header that brings you back to the homepage
[] style dashboard to match canva mockup
[] style homepage to match canva mockup


## ------BACKEND------
[] organize all models into one models.go file
    [] Double check every file for models
    [] internal/auth/models.go: move models to internal/models/models.go
    [] internal/auth/auth.go: move Users model into internal/models/models.go
    [] internal/storage/postgres.go: move Price model into internal/models/models.go
    [] internal/ai/openai.go: move model

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
