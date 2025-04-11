import Header from "../components/Header"
import logo from "../assets/logo.png"
import "../css/Home.css"

export default function Home() {
    try {
    } catch (error) {
        console.log(error)
    }
    return (
        <section>
            <Header />
            <div id="homepage">
                <div id="left-side">
                    <h1 id="home-header">TrackPrice<span style={{ color: "green" }}>AI</span></h1>
                    <p>
                        TrackPrice AI tracks your competitors’ pricing so you don’t have to. Add URLs, and we’ll notify you when prices change—automatically. Stay informed, react faster, and build smarter pricing strategies with zero manual monitoring.
                    </p>
                </div>
                <div id="right-side">
                    <img src={logo}></img>
                </div>
            </div>
        </section>
    );
}
