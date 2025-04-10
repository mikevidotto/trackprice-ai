import Header from "../components/Header";
import Track from "../pages/Track"

export default function Dashboard() {
    try {
    } catch(error) {
        console.log("HELLLLLOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOO")
        console.log(error)
    }
  return (
        <>
            <Header />
            <Track />
        </>
    );
}
