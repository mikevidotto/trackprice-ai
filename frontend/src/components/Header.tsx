import "../css/Header.css"
export default function Header() {
    console.log(localStorage.getItem("token"));
    if ((localStorage.getItem("token") == "") || (localStorage.getItem("token") == null)) {
        return (
            <>
                <header>
                    <div id="logo">
                        <a className="logobtn" href="http://localhost:5173/">TrackPriceAi</a>
                    </div>
                    <a href="http://localhost:5173/">Home</a>
                    <a href="http://localhost:5173/signup">Join</a>
                    <a href="http://localhost:5173/login">Login</a>
                    <a href="http://localhost:5173/">Contact</a>
                </header>
            </>
        );
    }
        return (
            <>
                <header>
                    <div id="logo">
                        <a href="http://localhost:5173/">TrackPriceAi</a>
                    </div>
                    <a href="http://localhost:5173/">Home</a>
                    <a href="http://localhost:5173/track">Dashboard</a>
                    <a href="http://localhost:5173/">Contact</a>
                    <a href="http://localhost:5173/logout">Logout</a>

                </header>
            </>
        );
}
