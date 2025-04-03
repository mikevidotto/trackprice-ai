# TODO

[] CREATE COMPONENTS!!!!!!
    [x] Create a Login Portal component
    [x] Create a Sign Up Portal component
    [] Modify Track component to look better.
[] modify/split frontend/pages/Track.tsx
    [] create Submit.tsx
        [] Submits a competitor via react component form that makes post call using axios.post("/track")
        [] the Track() function
    [] create Display.tsx
        [] Displays user's competitors in a list react component using axios.get("/track")
        [] the getCompetitors() function within the Track() function.

[] create a page layout that you can put your react components into (use App.tsx??)
    [] add layout components
        [] header
        [] nav-bar
        [] body
        [] footer
    [] make it responsive

[x] modify frontend/src/components/Header
    ***NAV-BAR ITEMS SHOULD NOT CALL THE API DIRECTLY!!***
    ***NEVERMIND! THEY DON'T AS LONG AS THEY'RE NOT ROUTED ON THE BACKEND!!!! (so it's all good bro)***
    [x] change the login route to redirect to Login.tsx page
        [x] Create Login Page with option to signup
        
