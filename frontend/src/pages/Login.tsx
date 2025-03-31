import { useState } from "react";
import Header from "../components/Header";
import { useNavigate } from "react-router-dom";
import API from "../utils/api";
import { AuthResponse } from "../utils/types";

export default function Login() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const res = await API.post<AuthResponse>("/login", { email, password });
      localStorage.setItem("token", res.data.token);
      console.log(localStorage.getItem("token"))
      navigate("/track");
    } catch (err) {
      alert("Login failed!");
    }
  };

  return (
    <>
    <Header />
    <form onSubmit={handleSubmit}>
      <input
        type="email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
      />
      <br></br>
      <input
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
      /> 
      <br></br>
      <button type="submit">Login</button>
    </form>
    </>
  );
}
