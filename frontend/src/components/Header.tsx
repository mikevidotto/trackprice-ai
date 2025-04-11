import "../css/Header.css"
import logo from "../assets/logo.png"
export default function Header() {
    console.log(localStorage.getItem("token"));
    if ((localStorage.getItem("token") == "") || (localStorage.getItem("token") == null)) {
        return (
            <>
                <header>
                    <div id="logo">
                            <a id="logobtn" href="http://localhost:5173/">
                                <img id="logo-image" src={logo}></img>
                                <p>TrackPriceAi</p>
                            </a>
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
