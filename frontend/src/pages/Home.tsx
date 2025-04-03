import Header from "../components/Header"
import Track from "../pages/Track"

export default function Home() {
    return (
        <section>
        <Header />
        <div className="home">
            <h1>Welcome to the Home Page!</h1>
            <p>This is the home page of our application.</p>
            <p>Feel free to explore the features and functionalities we offer.</p>
        </div>
        </section>
    );
}
