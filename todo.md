# TODO

[] CREATE COMPONENTS!!!!!!
    [] Create a Login Portal component
    [] Create a Sign Up Portal component
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

[] modify frontend/src/components/Header
    ***NAV-BAR ITEMS SHOULD NOT CALL THE API DIRECTLY!!***
    [] change the login route to redirect to Login.tsx page
        [] Create Login Page with option to signup
        
