import './App.css'
import { BrowserRouter, Routes, Route } from "react-router-dom";
import Login from "./pages/Login";
import Signup from "./pages/Signup";
import Track from "./pages/Track";
import Home from "./pages/Home";

function App() {
  return (
    <> 
      <BrowserRouter>
    <Routes>
      <Route path="/" element={<Home />} />
      <Route path="/login" element={<Login />} />
      <Route path="/signup" element={<Signup />} />
      <Route path="/track" element={<Track />} />
    </Routes>
  </BrowserRouter>
    </>
  )
}

export default App
